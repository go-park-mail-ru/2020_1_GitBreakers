package codehub

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	codehubMetrics "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/monitoring"
	http4 "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/delivery/http"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/repository/postgres/issues"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/repository/postgres/merge"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/repository/postgres/news"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/repository/postgres/search"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/repository/postgres/stars"
	usecaseCodeHub "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/csrf"
	gitDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/delivery"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/repository"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	http2 "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery/http"
	http3 "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery/http"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/repository/postgres"
	userUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	middlewareCommon "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/middleware"
	gorCSRF "github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func StartNew() {
	conf := config.New()
	prometheus.MustRegister(codehubMetrics.Hits, codehubMetrics.RequestDuration, codehubMetrics.DBQueryDuration)

	f, err := os.OpenFile(conf.LOGFILE, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		logrus.Fatalln("Failed to open logfile:", err)
	}

	customLogger := logger.NewTextFormatSimpleLogger(f, 1)

	defer func() {
		if err := f.Close(); err != nil {
			logrus.Errorln("Failed to close logfile:", err)
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

	db.SetMaxOpenConns(int(conf.MAX_DB_OPEN_CONN)) //10 по дефолту

	mainRouter := mux.NewRouter()

	metricsRouter := mainRouter.PathPrefix("/metrics").Subrouter() // prometheus /metrics route

	const apiMainRoute = "/api/v1" // all api methods start with /api/v1

	apiRouter := mainRouter.PathPrefix(apiMainRoute).Subrouter()

	staticRouter := apiRouter.PathPrefix("").Subrouter()

	handlersRouter := apiRouter.PathPrefix("").Subrouter()

	CsrfRouter := handlersRouter.PathPrefix("").Subrouter()

	// handlers

	userSetHandler, m, repoHandler, CHubHandler := initNewHandler(db, customLogger, conf)
	csrfHandlers := csrf.NewHandlers(csrf.DefaultTokenHeaderName)

	// middleware

	c := cors.New(cors.Options{
		AllowedOrigins:   conf.ALLOWED_ORIGINS,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		Debug:            false,
		AllowedHeaders: []string{
			"Content-Type",
			"User-Agent",
			"Cache-Control",
			"Accept",
			"X-Requested-With",
			"If-Modified-Since",
			"Origin",
			"X-CSRF-Token",
		},
	})

	panicMiddleware := middleware.CreatePanicMiddleware(customLogger)
	loggerMWare := middlewareCommon.CreateAccessLogMiddleware(customLogger)
	checkAuthMiddleware := middleware.CreateCheckAuthMiddleware(customLogger)

	csrfMiddleware := middleware.CreateCSRFMiddleware(
		[]byte(conf.CSRF_SECRET_KEY),
		conf.ALLOWED_ORIGINS,
		middleware.DefaultCSRFCookieName,
		csrfHandlers.TokenHeaderName,
		false,
		gorCSRF.SameSiteNoneMode,
		conf.COOKIE_EXPIRE_HOURS*3600,
	)

	// set middleware

	mainRouter.Use(panicMiddleware)

	apiRouter.Use(
		loggerMWare,
		middleware.PrometheusMetricsMiddleware,
		middleware.ProtectHeadersMiddleware,
	)

	handlersRouter.Use(
		middleware.JsonContentTypeMiddleware,
		m.AuthMiddleware,
	)

	CsrfRouter.Use(checkAuthMiddleware, csrfMiddleware)

	// Handlers

	metricsRouter.Handle("", promhttp.Handler()).Methods(http.MethodGet)

	CsrfRouter.HandleFunc("/csrftoken", csrfHandlers.GetNewCsrfTokenHandler).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/session", userSetHandler.LoginByLoginOrEmail).Methods(http.MethodPost)
	handlersRouter.HandleFunc("/session", userSetHandler.Logout).Methods(http.MethodDelete)

	handlersRouter.HandleFunc("/user/profile", userSetHandler.Create).Methods(http.MethodPost)
	handlersRouter.HandleFunc("/user/profile", userSetHandler.GetInfo).Methods(http.MethodGet)
	CsrfRouter.HandleFunc("/user/profile", userSetHandler.Update).Methods(http.MethodPut)
	handlersRouter.HandleFunc("/user/profile/{login_or_email}", userSetHandler.GetInfoByLoginOrEmail).Methods(http.MethodGet)
	handlersRouter.HandleFunc("/user/id/{id:[0-9]+}/profile", userSetHandler.GetInfoByID).Methods(http.MethodGet)
	CsrfRouter.HandleFunc("/user/avatar", userSetHandler.UploadAvatar).Methods(http.MethodPut)
	handlersRouter.HandleFunc("/user/repo", repoHandler.GetRepoList).Methods(http.MethodGet)
	handlersRouter.HandleFunc("/user/repo/{username}", repoHandler.GetRepoList).Methods(http.MethodGet)
	handlersRouter.HandleFunc("/user/id/{id:[0-9]+}/repo", repoHandler.GetRepoListByUserID).Methods(http.MethodGet)
	CsrfRouter.HandleFunc("/user/repo", repoHandler.CreateRepo).Methods(http.MethodPost)
	CsrfRouter.HandleFunc("/user/repo", repoHandler.DeleteRepo).Methods(http.MethodDelete)

	handlersRouter.HandleFunc("/user/pullrequests", CHubHandler.GetPLFromUser).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/repo/{username}/{reponame}", repoHandler.GetRepo).Methods(http.MethodGet)
	handlersRouter.HandleFunc("/repo/{username}/{reponame}/head", repoHandler.GetRepoHead).Methods(http.MethodGet)
	handlersRouter.HandleFunc("/repo/{username}/{reponame}/branch/{branchname}",
		repoHandler.GetBranchInfoByNames).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/repo/{username}/{reponame}/branch/{branchname}/tree/{path:.*}",
		repoHandler.GetFileContentByBranch).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/repo/{username}/{reponame}/branches",
		repoHandler.GetBranchList).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/repo/{username}/{reponame}/commits/hash/{hash}",
		repoHandler.GetCommitsList).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/repo/{username}/{reponame}/commit/{hash}/tree/{path:.*}",
		repoHandler.GetFileContentByCommitHash).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/repo/{username}/{reponame}/files/{hashcommits}",
		repoHandler.ShowFiles).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/repo/{username}/{reponame}/commits/branch/{branchname}",
		repoHandler.GetCommitsByBranchName).Methods(http.MethodGet)

	CsrfRouter.HandleFunc("/func/repo/{repoID:[0-9]+}/issues", CHubHandler.NewIssue).Methods(http.MethodPost)
	CsrfRouter.HandleFunc("/func/repo/{repoID:[0-9]+}/issues", CHubHandler.UpdateIssue).Methods(http.MethodPut)
	handlersRouter.HandleFunc("/func/repo/{repoID:[0-9]+}/issues", CHubHandler.GetIssues).Methods(http.MethodGet)
	CsrfRouter.HandleFunc("/func/repo/{repoID:[0-9]+}/issues", CHubHandler.CloseIssue).Methods(http.MethodDelete)
	//
	CsrfRouter.HandleFunc("/func/repo/{repoID:[0-9]+}/stars", CHubHandler.ModifyStar).Methods(http.MethodPut)
	CsrfRouter.HandleFunc("/func/repo/{repoID:[0-9]+}/stars/users", CHubHandler.UserWithStar).Methods(http.MethodGet)
	handlersRouter.HandleFunc("/func/repo/{login}/stars", CHubHandler.StarredRepos).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/func/repo/{repoID:[0-9]+}/news", CHubHandler.GetNews).Methods(http.MethodGet)

	CsrfRouter.HandleFunc("/func/repo/fork", repoHandler.Fork).Methods(http.MethodPost)

	CsrfRouter.HandleFunc("/func/repo/pullrequests", CHubHandler.CreatePullReq).Methods(http.MethodPost)
	handlersRouter.HandleFunc("/func/repo/{repoID:[0-9]+}/pullrequests/{direction}",
		CHubHandler.GetPullReqList).Methods(http.MethodGet)

	CsrfRouter.HandleFunc("/func/repo/pullrequests", CHubHandler.UndoPullReq).Methods(http.MethodDelete)
	CsrfRouter.HandleFunc("/func/repo/pullrequests", CHubHandler.ApproveMerge).Methods(http.MethodPut)

	handlersRouter.HandleFunc("/func/repo/pullrequest/{id:[0-9]+}", CHubHandler.GetMRByID).Methods(http.MethodGet)
	handlersRouter.HandleFunc("/func/repo/pullrequest/{id:[0-9]+}/diff", CHubHandler.GetMRDiffByID).Methods(http.MethodGet)

	handlersRouter.HandleFunc("/func/search/{params}", CHubHandler.Search).Methods(http.MethodGet)

	// static files server
	staticHandler := http.FileServer(http.Dir("./static"))
	staticRouter.PathPrefix("/static").Handler(
		http.StripPrefix(apiMainRoute+"/static", staticHandler),
	)

	// use cors middleware firs and start
	if err = http.ListenAndServe(conf.MAIN_LISTEN_ENDPOINT, c.Handler(mainRouter)); err != nil {
		log.Fatal(err)
	}
}

func initNewHandler(db *sqlx.DB, logger logger.SimpleLogger, conf *config.Config) (*http3.UserHttp, *middleware.Middleware, *gitDeliv.GitDelivery, *http4.HttpCodehub) {
	userRepos := postgres.NewUserRepo(db, conf.DEFAULT_USER_AVATAR_NAME, "/static/image/avatar/", conf.PATH_PREFIX)
	sessClient, err := clients.NewSessClient()
	if err != nil {
		logger.Fatal(err, "not connect to auth server")
	}

	userClient, err := clients.NewUserClient()
	if err != nil {
		logger.Fatal(err, "not connect to auth server")
	}
	newsClient, err := clients.NewNewsClient()
	if err != nil {
		logger.Fatal(err, "not connect to news server")
	}

	userUCase := userUC.UCUser{RepUser: &userRepos}

	absGitRepoDir, pathErr := filepath.Abs(filepath.Clean(conf.GIT_USER_REPOS_DIR))
	if pathErr != nil {
		log.Fatalln("bad git directory path:", err)
	}

	absPullsDir, pathErr := filepath.Abs(filepath.Clean(conf.GIT_USER_PULLRQ_DIR))
	if pathErr != nil {
		log.Fatalln("bad git directory path:", err)
	}

	repogit := repository.NewRepository(db, absGitRepoDir)
	repoCodeHubIssue := issues.NewIssueRepository(db)
	repoCodeHubStar := stars.NewStarRepository(db)
	repoCodeHubNews := news.NewRepoNews(db)
	repoCodeHubSearch := search.NewSearchRepository(db)
	repoMerge := merge.NewPullRequestRepository(db, repogit, absPullsDir)

	gitUseCase := usecase.GitUseCase{Repo: &repogit}

	codeHubUseCase := usecaseCodeHub.UCCodeHub{
		RepoIssue:  &repoCodeHubIssue,
		RepoStar:   &repoCodeHubStar,
		RepoNews:   &repoCodeHubNews,
		GitRepo:    repogit,
		UserRepo:   userRepos,
		SearchRepo: repoCodeHubSearch,
		RepoMerge:  repoMerge,
	}

	codeHubDelivery := http4.HttpCodehub{
		Logger:     &logger,
		CodeHubUC:  &codeHubUseCase,
		NewsClient: &newsClient,
		UserClient: &userClient,
	}

	sessDelivery := http2.SessionHttp{
		CookieName:       "session_id",
		CookieExpireTime: time.Duration(conf.COOKIE_EXPIRE_HOURS) * time.Hour,
		CookieSecure:     false,
		CookieSiteMode:   http.SameSiteNoneMode,
		CookiePath:       "/",
		Client:           &sessClient,
	}

	userDelivery := http3.UserHttp{
		SessHttp: &sessDelivery,
		Logger:   &logger,
		UClient:  &userClient,
		UCUser:   &userUCase,
	}

	gitDelivery := gitDeliv.GitDelivery{
		UC:     &gitUseCase,
		Logger: &logger,
		UserUC: &userUCase,
	}

	m := middleware.Middleware{
		SessDeliv: &sessClient,
	}

	return &userDelivery, &m, &gitDelivery, &codeHubDelivery
}
