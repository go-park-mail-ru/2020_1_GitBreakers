package main

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	session "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery/grpc"
	redisRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/repository/redis"
	sessUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/usecase"
	"github.com/go-redis/redis/v7"
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
	redisClient := redis.NewClient(&redis.Options{
		Addr:     conf.REDIS_ADDR, // use default Addr
		Password: conf.REDIS_PASS, // no password set
		DB:       0,               // use default db
	})
	res, err := redisClient.Ping().Result()

	if err != nil {
		msg := fmt.Sprintln("error with redis:", err)
		log.Fatal(msg)
	} else {
		log.Println(res)
	}

	sessRepos := redisRepo.NewSessionRedis(redisClient, "codehub/session/")
	sessUCase := sessUC.SessionUC{RepoSession: &sessRepos}

	session.NewSessServer(s, &sessUCase)

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
