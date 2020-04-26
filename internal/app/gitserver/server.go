package gitserver

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/config"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/monitoring"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	middlewareCommon "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/middleware"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sosedoff/gitkit"
	"log"
	"net/http"
	"os"
)

func StartNew() {
	conf := config.New()
	prometheus.MustRegister(monitoring.Hits, monitoring.RequestDuration, monitoring.DBQueryDuration)

	//userClient, err := clients.NewUserClient()
	//if err != nil {
	//	log.Fatal(err, "not connect to auth server")
	//}

	customLogger := logger.NewTextFormatSimpleLogger(os.Stdout)

	//берутся из .env файла
	connStr := "user=" + conf.POSTGRES_USER + " password=" +
		conf.POSTGRES_PASS + " dbname=" + conf.POSTGRES_DBNAME

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		customLogger.Fatalln("failed to start db:", err)
	} else {
		customLogger.Println("connected to postgres:", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Println("failed to close db:", err)
		}
	}()

	db.SetMaxOpenConns(int(conf.MAX_DB_OPEN_CONN)) //10 по дефолту

	// repogit := repository.NewRepository(db, conf.GIT_USER_REPOS_DIR)

	gitkitConfig := gitkit.Config{
		Dir:        conf.GIT_USER_REPOS_DIR,
		AutoCreate: false, // Do not create repository if it not exist
		Auth:       false, // TODO use ow authentication based on middleware
	}

	gitkitServer := gitkit.New(gitkitConfig)

	if err := gitkitServer.Setup(); err != nil {
		customLogger.Fatalln("cannot start gitserver:", err)
	}

	panicMiddleware := middleware.CreatePanicMiddleware(customLogger)
	loggerMWare := middlewareCommon.CreateAccessLogMiddleware(1, customLogger)

	gitServer := middleware.PrometheusMetricsMiddleware(panicMiddleware(loggerMWare(gitkitServer)))

	http.Handle("/", gitServer)

	if err := http.ListenAndServe(conf.GIT_SERVER_PORT, nil); err != nil {
		customLogger.Fatalln("cannot start http git server:", err)
	}
}
