package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/config"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	sessDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery"
	sessRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/repository"
	sessUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/usecase"
	userDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery"
	userRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/repository"
	userUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/usecase"
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
	}
	defer db.Close()

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
		Password: "",              // no password set
		DB:       0,               // use default DB
	}).Conn()

	userSetHandler, m := initNewHandler(db, redisConn)

	r.HandleFunc("/signup", userSetHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/login", userSetHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", userSetHandler.Logout).Methods(http.MethodGet)
	r.HandleFunc("/profile", userSetHandler.GetInfo).Methods(http.MethodGet)

	if err = http.ListenAndServe(conf.MAIN_LISTEN_PORT, c.Handler(m.AuthMiddleware(r))); err != nil {
		log.Println(err)
		return
	}
}

func initNewHandler(db *sqlx.DB, redis *redis.Conn) (*userDeliv.UserHttp, *middleware.Middleware) {
	sessRepos := sessRepo.SessionRedis{redis}
	userRepo := userRepo.DBWork{db, "/static/img/avatar/default.jpg"}
	sessUC := sessUC.SessionUCWork{&sessRepos}
	sessDeliv := sessDeliv.SessionHttpWork{&sessUC, 48 * time.Hour}
	userUCase := userUC.UCUserWork{&userRepo}
	userDeliv := userDeliv.UserHttp{
		SessHttp: &sessDeliv,
		UserUC:   &userUCase,
	}
	m := middleware.Middleware{
		SessDeliv: &sessDeliv,
		UCUser:    &userUCase,
	}
	return &userDeliv, &m
}
