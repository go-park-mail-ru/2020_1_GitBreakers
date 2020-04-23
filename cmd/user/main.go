package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	usergrpc "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery/grpc"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/repository/postgres"
	userUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/usecase"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"log"
	"net"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	govalidator.SetFieldsRequiredByDefault(true)
}

func main() {
	s := grpc.NewServer()
	conf := config.New()

	connStr := "user=" + conf.POSTGRES_USER + " password=" +
		conf.POSTGRES_PASS + " dbname=" + conf.POSTGRES_DBNAME

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Database started success")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	db.SetMaxOpenConns(conf.MAX_DB_OPEN_CONN) //10 по дефолту
	userRepos := postgres.NewUserRepo(db, "default.jpg", "/static/image/avatar/", conf.HOST_TO_SAVE)
	userUCase := userUC.UCUser{RepUser: &userRepos}

	usergrpc.NewUserServer(s, &userUCase)

	l, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
