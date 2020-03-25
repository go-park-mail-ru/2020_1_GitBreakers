package codehub

import (
	"flag"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery"
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
	connStr := "user=andrey password=kekmda dbname=codehub"

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

	redisAddr := flag.String("addr", "redis://user:@localhost:6379/0", "redis addr")
	//TODO чекнуть про редис(как лучше коннектить)
	redisConn := &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.DialURL(*redisAddr)
			if err != nil {
				log.Fatal("fail init redis pool: ", err)
			}
			return conn, err
		},
	}
	defer redisConn.Close()
	//TODO должен вернуть наши обработчики
	initNewHandler(db, redisConn)
	//TODO сжирать реальные коннекты к бд
	m := middleware.Middleware{
		SessDeliv: nil,
		UCUser:    nil,
	}
	//TODO сжирать реальные данные а не nil
	userSetHandler := delivery.UserHttp{
		SessHttp: nil,
		UserUC:   nil,
	}

	r.HandleFunc("/signup", userSetHandler.Create).Methods(http.MethodPost)
	r.HandleFunc("/login", userSetHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/logout", userSetHandler.Logout).Methods(http.MethodPost)
	r.HandleFunc("/profile", userSetHandler.GetInfo).Methods(http.MethodGet)

	if err = http.ListenAndServe(":8080", c.Handler(m.AuthMiddleware(r))); err != nil {
		log.Println(err)
		return
	}
}
func initNewHandler(db *sqlx.DB, redis *redis.Pool) {
	///проброс всех коннектов в структуры
}
