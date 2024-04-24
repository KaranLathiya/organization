package main

import (
	"fmt"
	"net/http"
	"organization/config"
	"organization/controller"
	"organization/db"
	"organization/repository"
	"organization/routes"
)

func main() {
	err := config.LoadConfig("../config")
	if err != nil {
		panic(fmt.Sprintf("cannot load config: %v", err))
	}
	db := db.Connect()
	defer db.Close()
	repos := repository.InitRepositories(db)
	controllers := controller.InitControllers(repos)
	router := routes.InitializeRouter(controllers)
	fmt.Println("server started")
	http.ListenAndServe(":"+config.ConfigVal.Port, router)
}
