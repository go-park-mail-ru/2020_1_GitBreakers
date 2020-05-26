package gitserver

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	codehubMetrics "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/monitoring"
	mergeRepository "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/repository/postgres/merge"
	gitRepository "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/repository"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/gitserver/delivery"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/gitserver/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	middlewareCommon "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/middleware"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"github.com/sosedoff/gitkit"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func StartNew() {
	conf := config.New()
	prometheus.MustRegister(codehubMetrics.Hits, codehubMetrics.RequestDuration, codehubMetrics.DBQueryDuration)

	f, err := os.OpenFile(conf.GIT_SERVER_LOGFILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Fatalln("Failed to open gitserver logfile:", err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			logrus.Errorln("Failed to close logfile:", err)
		}
	}()

	customLogger := logger.NewTextFormatSimpleLogger(f, 1)

	if _, err = fmt.Fprintf(f, ">>>>>>>>>>>>%v<<<<<<<<<<<<\n", time.Now()); err != nil {
		msg := fmt.Sprintln("Failed to write gitserver start timestamp in log output:", err)
		customLogger.Error(msg)
		log.Fatal(msg)
	}

	userClient, err := clients.NewUserClient()
	if err != nil {
		log.Fatal(err, "not connect to auth server")
	}

	//берутся из .env файла
	connStr := "user=" + conf.POSTGRES_USER + " password=" +
		conf.POSTGRES_PASS + " dbname=" + conf.POSTGRES_DBNAME

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		customLogger.Fatalln("failed to start db:", err)
	} else {
		customLogger.Println("connected to postgres:", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			customLogger.Println("failed to close db:", err)
		}
	}()

	db.SetMaxOpenConns(int(conf.MAX_DB_OPEN_CONN)) //10 по дефолту

	absGitRepoDir, pathErr := filepath.Abs(filepath.Clean(conf.GIT_USER_REPOS_DIR))
	if pathErr != nil {
		customLogger.Fatalln("bad git repositories directory path:", err)
		return
	}

	absPullsDir, pathErr := filepath.Abs(filepath.Clean(conf.GIT_USER_PULLRQ_DIR))
	if pathErr != nil {
		log.Fatalln("bad git pull requests directory path:", err)
	}

	gitRepos := gitRepository.NewRepository(db, absGitRepoDir)
	mergeRepos := mergeRepository.NewPullRequestRepository(db, gitRepos, absPullsDir)

	gitkitConfig := gitkit.Config{
		Dir:        conf.GIT_USER_REPOS_DIR,
		AutoHooks:  false, // Do not initialise hooks, because we use nested directories
		AutoCreate: false, // Do not create repository if it not exist
		Auth:       false, // We use own authentication based on middleware
	}

	gitkitServer := gitkit.New(gitkitConfig)

	if err := gitkitServer.Setup(); err != nil {
		customLogger.Fatalln("cannot start gitserver:", err)
	}

	panicMiddleware := middleware.CreatePanicMiddleware(customLogger)
	accessLogMiddleware := middlewareCommon.CreateAccessLogMiddleware(customLogger)

	routerTemplate := fmt.Sprintf("/{%s}/{%s}.git",
		delivery.OwnerLoginMuxParameter, delivery.RepositoryNameMuxParameter)

	mainRouter := mux.NewRouter()

	mainRouter.Use(panicMiddleware)

	metricsRouter := mainRouter.PathPrefix("/metrics").Subrouter()
	gitRouter := mainRouter.PathPrefix(routerTemplate).Subrouter()

	gitRouter.Use(middleware.PrometheusMetricsMiddleware, accessLogMiddleware)

	gitServerDelivery := delivery.GitServerDelivery{
		UseCase: usecase.NewUseCase(gitRepos, mergeRepos, userClient),
		Logger:  customLogger,
	}

	metricsRouter.Handle("", promhttp.Handler()).Methods(http.MethodGet)

	gitInfoRefsHandler := delivery.CreateGitIfoRefsMiddleware(gitServerDelivery)(gitkitServer)
	gitUploadPackHandler := delivery.CreateGitUploadPackMiddleware(gitServerDelivery)(gitkitServer)
	gitReceivePackHandler := delivery.CreateGitReceivePackMiddleware(gitServerDelivery)(gitkitServer)

	gitRouter.Handle("/info/refs", gitInfoRefsHandler).Methods(http.MethodGet)
	gitRouter.Handle("/"+delivery.GitUploadPackService, gitUploadPackHandler).Methods(http.MethodPost)
	gitRouter.Handle("/"+delivery.GitReceivePackService, gitReceivePackHandler).Methods(http.MethodPost)

	customLogger.Printf("starting git server with GIT_SERVER_ENDPOINT=%v\n", conf.GIT_SERVER_ENDPOINT)
	customLogger.Printf("starting git server with GIT_USER_REPOS_DIR=%v\n", absGitRepoDir)

	if err := http.ListenAndServe(conf.GIT_SERVER_ENDPOINT, mainRouter); err != nil {
		customLogger.Fatalln("cannot start http git server:", err)
	}
}
