package main

import (
	"./handlers"
	"./routes"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	var appConfig routes.AppConfig

	isDebug, exist := os.LookupEnv("GB_DEBUG")
	if !exist {
		fmt.Println(`GB_DEBUG is not set, default is false`)
		appConfig.IsDebug = false
	} else {
		var parseErr error
		if appConfig.IsDebug, parseErr = strconv.ParseBool(isDebug); parseErr != nil {
			fmt.Println(`GB_DEBUG value invalid, shutdown`)
			return
		}
	}

	CSRFTAuthKey, exist := os.LookupEnv("GB_CSCRF_AUTH_KEY")
	if !exist {
		fmt.Println(`No CSCRF_AUTH_KEY value in environment, shutdown`)
		return
	}
	appConfig.CSRFTAuthKey = []byte(CSRFTAuthKey)
	if len(appConfig.CSRFTAuthKey) != 32 {
		fmt.Println(`GB_CSCRF_AUTH_KEY must be 32bytes length`)
		return
	}


	modelsContext := new(handlers.StoresContext)

	r:= routes.NewRouter(&appConfig, modelsContext)
	http.Handle("/", r)

	_ = http.ListenAndServe(":8080", r)
}
