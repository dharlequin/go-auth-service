package api

import (
	"os"

	"github.com/dharlequin/go-auth-service/api/controllers"
)

var server = controllers.Server{}

//Run runs all server parts
func Run() {
	server.Initialize(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	server.Run(":8181")
}
