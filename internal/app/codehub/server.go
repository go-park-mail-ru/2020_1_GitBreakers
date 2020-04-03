package codehub

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	sessDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/redis/delivery"
	sessRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/redis/repository"
	sessUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/redis/usecase"
	userDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/database/delivery"
	userRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/database/repository"
	userUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/database/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	middleareCommon "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/middleware"
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
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
				customLogger.Info("Failed to close logger: " + err.Error())
			}
		}
	}()
	if _, err = fmt.Fprintf(f, ">>>>>>>>>>>>%v<<<<<<<<<<<<\n", time.Now()); err != nil {
		customLogger.Error("Failed to write server start timestamp in log output: " + err.Error())
		return
	}

	//берутся из .env файла
	connStr := "user=" + conf.POSTGRES_USER + " password=" +
		conf.POSTGRES_PASS + " dbname=" + conf.POSTGRES_DBNAME

	db, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		customLogger.Error("Failed to start db: " + err.Error())
		return
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
	if res != "PONG" {
		customLogger.Error("error with redis")
		return
	} else {
		customLogger.Println("Connected to redis ", res, err)
	}

	userSetHandler, m := initNewHandler(db, redisConn, customLogger, conf)

	r.HandleFunc("/signup", userSetHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/login", userSetHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", userSetHandler.Logout).Methods(http.MethodGet)
	r.HandleFunc("/whoami", userSetHandler.GetInfo).Methods(http.MethodGet)
	r.HandleFunc("/profile", userSetHandler.Update).Methods(http.MethodPut)
	r.HandleFunc("/profile/{login}", userSetHandler.GetInfoByLogin).Methods(http.MethodGet)
	r.HandleFunc("/avatar", userSetHandler.UploadAvatar).Methods(http.MethodPut)

	staticHandler := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticHandler))
	panicMiddleware := middleware.PanicMiddleware(m.AuthMiddleware(r))
	loggerMWare := middleareCommon.CreateAccessLogMiddleware(1, customLogger)
	if err = http.ListenAndServe(conf.MAIN_LISTEN_PORT, c.Handler(loggerMWare(panicMiddleware))); err != nil {
		log.Fatal(err)
	}
}

func initNewHandler(db *sqlx.DB, redis *redis.Conn, logger logger.SimpleLogger, conf *config.Config) (*userDeliv.UserHttp, *middleware.Middleware) {
	sessRepos := sessRepo.NewSessionRedis(redis, "codehub/session/")
	userRepos := userRepo.NewUserRepo(db, "default.jpg", "/static/image/avatar/", conf.HOST_TO_SAVE)
	sessUCase := sessUC.SessionUC{&sessRepos}
	//todo expiretime в конфиге
	sessDelivery := sessDeliv.SessionHttp{&sessUCase, 48 * time.Hour}
	userUCase := userUC.UCUser{&userRepos}
	userDelivery := userDeliv.UserHttp{
		SessHttp: &sessDelivery,
		UserUC:   &userUCase,
		Logger:   &logger, 
	}
	//todo создать репо для гита
	//gitUseCase := usecase.GitUseCase{}
	//gitDelivery := gitDeliv.GitDelivery{gitUseCase}
	m := middleware.Middleware{
		SessDeliv: &sessDelivery,
		UCUser:    &userUCase,
	}
	return &userDelivery, &m
}
