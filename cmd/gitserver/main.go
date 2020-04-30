package main

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/gitserver"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	govalidator.SetFieldsRequiredByDefault(true)
}

func main() {
	gitserver.StartNew()
}
