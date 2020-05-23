package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	news "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/delivery/grpc"
	news2 "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/repository/postgres/news"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/usecase"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/repository"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/repository/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
	"path/filepath"
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
		log.Println("cannot connect to database", err)
		return
	} else {
		log.Println("Database started success")
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Println("cannot close database connection", err)
		}
	}()

	db.SetMaxOpenConns(int(conf.MAX_DB_OPEN_CONN)) //10 по дефолту
	newsRepos := news2.NewRepoNews(db)

	absGitRepoDir, pathErr := filepath.Abs(filepath.Clean(conf.GIT_USER_REPOS_DIR))
	if pathErr != nil {
		log.Println("bad git directory path:", err)
		return
	}

	gitRepos := repository.NewRepository(db, absGitRepoDir)
	userRepos := postgres.NewUserRepo(db, conf.DEFAULT_USER_AVATAR_NAME, "/static/image/avatar/", conf.PATH_PREFIX)

	UCCodeHub := usecase.UCCodeHub{
		RepoNews: &newsRepos,
		GitRepo:  gitRepos,
		UserRepo: userRepos,
	}

	news.NewNewsServer(s, &UCCodeHub)

	l, err := net.Listen("tcp", ":8083")
	if err != nil {
		log.Println("cannot start service:", err)
		return
	}

	if err := s.Serve(l); err != nil {
		log.Println("cannot start grpc service:", err)
	}
}
