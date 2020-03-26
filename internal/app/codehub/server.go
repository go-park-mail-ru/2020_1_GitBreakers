package codehub

import (
	"flag"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	sessDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery"
	sessRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/repository"
	sessUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/usecase"
	userDeliv "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery"
	userRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/repository"
	userUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/usecase"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func StartNew() {
	//TODO получать из конфига+создать схему бд
	connStr := "user=andrey password=167839 dbname=codehub"

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to start db: " + err.Error())
	}
	defer db.Close()

	db.SetMaxOpenConns(10)

	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://89.208.198.186:8080", "http://89.208.198.186:80"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		Debug:            false,
	})

	//TODO чекнуть про редис(как лучше коннектить)
	redisAddr := flag.String("addr", "redis://user:@localhost:6379/0", "redis addr")
	redisConn, err := redis.DialURL(*redisAddr)
	if err != nil {
		log.Fatalf("can't connect to redis")
		return
	}

	userSetHandler, m := initNewHandler(db, &redisConn)

	r.HandleFunc("/signup", userSetHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/login", userSetHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", userSetHandler.Logout).Methods(http.MethodPost)
	r.HandleFunc("/profile", userSetHandler.GetInfo).Methods(http.MethodGet)

	if err = http.ListenAndServe(":8080", c.Handler(m.AuthMiddleware(r))); err != nil {
		log.Println(err)
		return
	}
}

func initNewHandler(db *sqlx.DB, redis *redis.Conn) (*userDeliv.UserHttp, *middleware.Middleware) {
	sessRepos := sessRepo.SessionRedis{redis}
	userRepo := userRepo.DBWork{db}
	sessUC := sessUC.SessionUCWork{&sessRepos}
	sessDeliv := sessDeliv.SessionHttpWork{&sessUC}
	userUCase := userUC.UCUserWork{&userRepo}
	//TODO Интерфейсы поправить, чтобы сжирал реальные данные
	userDeliv := userDeliv.UserHttp{
		SessHttp: sessDeliv,
		UserUC:   userUCase,
	}
	m := middleware.Middleware{
		SessDeliv: sessDeliv,
		UCUser:    userUCase,
	}
	return &userDeliv, &m
}
