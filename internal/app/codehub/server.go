package codehub

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/csrf"
	gitDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/delivery"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/repository"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	http2 "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery/http"
	redisRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/repository/redis"
	sessUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/usecase"
	userDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/repository/postgres"
	userUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	middlewareCommon "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"time"
)

func StartNew() {
	conf := config.New()
	customLogger := logger.SimpleLogger{}

	f, err := os.OpenFile(conf.LOGFILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Errorln("Failed to open logfile:", err)
		f = os.Stdout
	}

	customLogger = logger.NewTextFormatSimpleLogger(f)

	defer func() {
		if f != os.Stdout {
			if err := f.Close(); err != nil {
				logrus.Errorln("Failed to close logfile:", err)
			}
		}
	}()

	if _, err = fmt.Fprintf(f, ">>>>>>>>>>>>%v<<<<<<<<<<<<\n", time.Now()); err != nil {
		msg := fmt.Sprintln("Failed to write server start timestamp in log output:", err)
		customLogger.Error(msg)
		log.Fatal(msg)
	}

	//берутся из .env файла
	connStr := "user=" + conf.POSTGRES_USER + " password=" +
		conf.POSTGRES_PASS + " dbname=" + conf.POSTGRES_DBNAME

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		msg := fmt.Sprintln("Failed to start db:", err)
		customLogger.Error(msg)
		log.Fatal(msg)
	} else {
		customLogger.Println("Connected to postgres:", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			customLogger.Infoln("Failed to close db:", err)
		}
	}()

	db.SetMaxOpenConns(conf.MAX_DB_OPEN_CONN) //10 по дефолту

	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   conf.ALLOWED_ORIGINS,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		Debug:            false,
		AllowedHeaders: []string{"Content-Type", "User-Agent",
			"Cache-Control", "Accept", "X-Requested-With", "If-Modified-Since", "Origin", "X-CSRF-Token"},
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.REDIS_ADDR, // use default Addr
		Password: conf.REDIS_PASS, // no password set
		DB:       0,               // use default db
	})

	res, err := redisClient.Ping().Result()
	if err != nil {
		msg := fmt.Sprintln("error with redis:", err)
		customLogger.Error(msg)
		log.Fatal(msg)
	} else {
		customLogger.Println("Connected to redis:", res)
	}

	r.Use(middleware.JsonContentTypeMiddleware, middleware.ProtectHeadersMiddleware)

	csrfMiddleware := middleware.CreateCsrfMiddleware(
		[]byte(conf.CSRF_SECRET_KEY),
		conf.ALLOWED_ORIGINS,
		false,
		conf.COOKIE_EXPIRE_HOURS*3600)

	withCsrfRouter := r.PathPrefix("").Subrouter()
	withCsrfRouter.Use(csrfMiddleware)

	userSetHandler, m, repoHandler := initNewHandler(db, redisClient, customLogger, conf)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(csrfMiddleware)
	api.HandleFunc("/csrftoken", csrf.GetNewCsrfToken).Methods(http.MethodGet)

	r.HandleFunc("/session", userSetHandler.Login).Methods(http.MethodPost)
	withCsrfRouter.HandleFunc("/session", userSetHandler.Logout).Methods(http.MethodDelete)
	r.HandleFunc("/user/profile", userSetHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/user/profile", userSetHandler.GetInfo).Methods(http.MethodGet)
	withCsrfRouter.HandleFunc("/users/profile", userSetHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/user/profile/{login}", userSetHandler.GetInfoByLogin).Methods(http.MethodGet)
	withCsrfRouter.HandleFunc("/user/avatar", userSetHandler.UploadAvatar).Methods(http.MethodPut)
	r.HandleFunc("/user/repo/{username}", repoHandler.GetRepoList).Methods(http.MethodGet)
	r.HandleFunc("/user/repo", repoHandler.GetRepoList).Methods(http.MethodGet)
	withCsrfRouter.HandleFunc("/user/repo", repoHandler.CreateRepo).Methods(http.MethodPost)

	r.HandleFunc("/repo/{username}/{reponame}", repoHandler.GetRepo).Methods(http.MethodGet)
	r.HandleFunc("/repo/{username}/{reponame}/branches", repoHandler.GetBranchList).Methods(http.MethodGet)
	r.HandleFunc("/repo/{username}/{reponame}/commits/hash/{hash}", repoHandler.GetCommitsList).Methods(http.MethodGet)
	r.HandleFunc("/repo/{username}/{reponame}/files/{hashcommits}", repoHandler.ShowFiles).Methods(http.MethodGet)
	r.HandleFunc("/repo/{username}/{reponame}/commits/branch/{branchname}", repoHandler.GetCommitsByBranchName).Methods(http.MethodGet)

	//r.HandleFunc("/repo/issues", nil).Methods(http.MethodPost)
	//r.HandleFunc("/repo/issues", nil).Methods(http.MethodPut)
	//r.HandleFunc("/repo/issues", nil).Methods(http.MethodGet)
	//r.HandleFunc("/repo/issues", nil).Methods(http.MethodDelete)
	//
	//r.HandleFunc("/repo/stars", nil).Methods(http.MethodPost)
	//r.HandleFunc("/repo/stars", nil).Methods(http.MethodGet)
	//r.HandleFunc("/repo/stars", nil).Methods(http.MethodDelete)
	//
	//r.HandleFunc("/repo/news", nil).Methods(http.MethodGet)

	staticHandler := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticHandler))

	panicMiddleware := middleware.CreatePanicMiddleware(customLogger)(m.AuthMiddleware(r))
	loggerMWare := middlewareCommon.CreateAccessLogMiddleware(1, customLogger)

	if err = http.ListenAndServe(conf.MAIN_LISTEN_PORT, c.Handler(loggerMWare(panicMiddleware))); err != nil {
		log.Fatal(err)
	}
}

func initNewHandler(db *sqlx.DB, redis *redis.Client, logger logger.SimpleLogger, conf *config.Config) (*userDeliv.UserHttp, *middleware.Middleware, *gitDeliv.GitDelivery) {
	sessRepos := redisRepo.NewSessionRedis(redis, "codehub/session/")
	userRepos := postgres.NewUserRepo(db, "default.jpg", "/static/image/avatar/", conf.HOST_TO_SAVE)
	sessUCase := sessUC.SessionUC{RepoSession: &sessRepos}
	sessDelivery := http2.SessionHttp{
		SessUC:     &sessUCase,
		ExpireTime: time.Duration(conf.COOKIE_EXPIRE_HOURS) * time.Hour,
	}

	userUCase := userUC.UCUser{RepUser: &userRepos}

	userDelivery := userDeliv.UserHttp{
		SessHttp: &sessDelivery,
		UserUC:   &userUCase,
		Logger:   &logger,
	}

	repogit := repository.NewRepository(db, conf.GIT_USER_REPOS_DIR)

	gitUseCase := usecase.GitUseCase{Repo: &repogit}

	gitDelivery := gitDeliv.GitDelivery{
		UC:     &gitUseCase,
		Logger: &logger,
		UserUC: &userUCase,
	}

	m := middleware.Middleware{
		SessDeliv: &sessDelivery,
		UCUser:    &userUCase,
	}

	return &userDelivery, &m, &gitDelivery
}
