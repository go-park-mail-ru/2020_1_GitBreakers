package main

import (
	session "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery/grpc"
	redisRepo "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/repository/redis"
	sessUC "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/usecase"
	"github.com/go-redis/redis/v7"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default db
	})
	sessRepos := redisRepo.NewSessionRedis(redisClient, "codehub/session/")
	sessUCase := sessUC.SessionUC{RepoSession: &sessRepos}
	srv := &session.GRPCServer{&sessUCase}

	session.RegisterSessionServer(s, srv)

	l, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
