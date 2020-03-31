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
	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

func StartNew() {
	conf := config.New()
	//берутся из .env файла
	connStr := "user=" + conf.POSTGRES_USER + " password=" +
		conf.POSTGRES_PASS + " dbname=" + conf.POSTGRES_DBNAME

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to start db: " + err.Error())
		return
	} else {
		fmt.Println("Connected to postgres ", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			return
		}
	}()

	db.SetMaxOpenConns(conf.MAX_DB_OPEN_CONN) //10 по дефолту

	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   conf.ALLOWED_ORIGINS,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Debug:            false,
	})

	redisConn := redis.NewClient(&redis.Options{
		Addr:     conf.REDIS_ADDR, // use default Addr
		Password: conf.REDIS_PASS, // no password set
		DB:       0,               // use default db
	}).Conn()
	res, err := redisConn.Ping().Result()
	if res != "PONG" {
		log.Fatal("error with redis")
		return
	} else {
		fmt.Println("Connected to redis ", res, err)
	}

	userSetHandler, m := initNewHandler(db, redisConn)

	r.HandleFunc("/signup", userSetHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/login", userSetHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", userSetHandler.Logout).Methods(http.MethodGet)
	r.HandleFunc("/profile", userSetHandler.GetInfo).Methods(http.MethodGet)
	r.HandleFunc("/avatar", userSetHandler.UploadAvatar).Methods(http.MethodPut)

	staticHandler := http.FileServer(http.Dir("./static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static", staticHandler))
	panicMiddleware := middleware.PanicMiddleware(m.AuthMiddleware(r))

	if err = http.ListenAndServe(conf.MAIN_LISTEN_PORT, c.Handler(panicMiddleware)); err != nil {
		log.Println(err)
		return
	}
}

func initNewHandler(db *sqlx.DB, redis *redis.Conn) (*userDeliv.UserHttp, *middleware.Middleware) {
	sessRepos := sessRepo.NewSessionRedis(redis)
	userRepos := userRepo.NewUserRepo(db, "default.jpg", "./static/image/avatar/")
	sessUCase := sessUC.SessionUC{&sessRepos}
	sessDelivery := sessDeliv.SessionHttp{&sessUCase, 48 * time.Hour}
	userUCase := userUC.UCUser{&userRepos}
	userDelivery := userDeliv.UserHttp{
		SessHttp: &sessDelivery,
		UserUC:   &userUCase,
	}
	m := middleware.Middleware{
		SessDeliv: &sessDelivery,
		UCUser:    &userUCase,
	}
	return &userDelivery, &m
}
