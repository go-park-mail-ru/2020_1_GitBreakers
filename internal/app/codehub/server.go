package codehub

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/csrf"
	gitDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/delivery"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/repository"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	sessDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/redis/delivery"
	sessRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/redis/repository"
	sessUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/redis/usecase"
	userDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/database/delivery"
	userRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/database/repository"
	userUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/database/usecase"
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
		logrus.Error("Failed to open logfile:", err)
		f = os.Stdout
	}
	customLogger = logger.NewTextFormatSimpleLogger(f)
	defer func() {
		if f != os.Stdout {
			if err := f.Close(); err != nil {
				log.Println("Failed to close logfile: " + err.Error())
			}
		}
	}()
	if _, err = fmt.Fprintf(f, ">>>>>>>>>>>>%v<<<<<<<<<<<<\n", time.Now()); err != nil {
		msg := "Failed to write server start timestamp in log output: " + err.Error()
		customLogger.Error(msg)
		log.Fatal(msg)
	}

	//берутся из .env файла
	connStr := "user=" + conf.POSTGRES_USER + " password=" +
		conf.POSTGRES_PASS + " dbname=" + conf.POSTGRES_DBNAME

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		msg := "Failed to start db: " + err.Error()
		customLogger.Error(msg)
		log.Fatal(msg)
	} else {
		customLogger.Println("Connected to postgres ", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			customLogger.Info("Failed to close db: " + err.Error())
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
			"Cache-Control", "Accept", "X-Requested-With", "If-Modified-Since", "Origin"},
	})

	redisConn := redis.NewClient(&redis.Options{
		Addr:     conf.REDIS_ADDR, // use default Addr
		Password: conf.REDIS_PASS, // no password set
		DB:       0,               // use default db
	}).Conn()

	res, err := redisConn.Ping().Result()
	if err != nil {
		msg := "error with redis: " + err.Error()
		customLogger.Error(msg)
		log.Fatal(msg)
	} else {
		customLogger.Println("Connected to redis: ", res)
	}

	csrfMiddleware := middleware.CreateCsrfMiddleware(
		[]byte(conf.CSRF_SECRET_KEY),
		conf.ALLOWED_ORIGINS,
		false,
		conf.COOKIE_EXPIRE_HOURS*3600)

	userSetHandler, m, repoHandler := initNewHandler(db, redisConn, customLogger, conf)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/csrftoken", csrf.GetNewCsrfToken).Methods(http.MethodGet)

	r.HandleFunc("/signup", userSetHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/login", userSetHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", userSetHandler.Logout).Methods(http.MethodGet)
	r.HandleFunc("/whoami", userSetHandler.GetInfo).Methods(http.MethodGet)
	r.HandleFunc("/profile", userSetHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/profile/{login}", userSetHandler.GetInfoByLogin).Methods(http.MethodGet)
	r.HandleFunc("/avatar", userSetHandler.UploadAvatar).Methods(http.MethodPut)

	r.HandleFunc("/repo", repoHandler.CreateRepo).Methods(http.MethodPost)
	r.HandleFunc("/{username}/{reponame}", repoHandler.GetRepo).Methods(http.MethodGet)
	r.HandleFunc("/repolist", repoHandler.GetRepoList).Methods(http.MethodGet)
	r.HandleFunc("/{username}", repoHandler.GetRepoList).Methods(http.MethodGet)
	r.HandleFunc("/{username}/{reponame}/branches", repoHandler.GetBranchList).Methods(http.MethodGet)
	r.HandleFunc("/{username}/{reponame}/commits/{branchname}", repoHandler.GetCommitsList).Methods(http.MethodGet)
	r.HandleFunc("/{username}/{reponame}/files/{hashcommits}", repoHandler.ShowFiles).Methods(http.MethodGet)
	r.HandleFunc("/{username}/{reponame}/{branchname}/commits", repoHandler.GetCommitsByBranchName).Methods(http.MethodGet)

	r.Use(csrfMiddleware)

	staticHandler := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticHandler))

	panicMiddleware := middleware.CreatePanicMiddleware(customLogger)(m.AuthMiddleware(r))
	loggerMWare := middlewareCommon.CreateAccessLogMiddleware(1, customLogger)

	if err = http.ListenAndServe(conf.MAIN_LISTEN_PORT, c.Handler(loggerMWare(panicMiddleware))); err != nil {
		log.Fatal(err)
	}
}

func initNewHandler(db *sqlx.DB, redis *redis.Conn, logger logger.SimpleLogger, conf *config.Config) (*userDeliv.UserHttp, *middleware.Middleware, *gitDeliv.GitDelivery) {
	sessRepos := sessRepo.NewSessionRedis(redis, "codehub/session/")
	userRepos := userRepo.NewUserRepo(db, "default.jpg", "/static/image/avatar/", conf.HOST_TO_SAVE)
	sessUCase := sessUC.SessionUC{RepoSession: &sessRepos}
	sessDelivery := sessDeliv.SessionHttp{
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
