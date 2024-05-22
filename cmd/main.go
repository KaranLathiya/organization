package main

import (
	"fmt"
	"net/http"
	"organization/config"
	"organization/controller"
	"organization/db"
	"organization/internal/cronjob"
	"organization/repository"
	"organization/routes"
	
	_ "organization/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

//	@title			Organization-Service API
//	@version		1.0
//	@description	Organization service for create/update/delete organization for user. It allows to invite/remove other users, assign them role and update their roles.
//	@host		localhost:9000/

// @schemes http
// @tag.name MicrosoftAuth
// @tag.description for microsoft login,tokens
// @tag.name Organization
// @tag.description Organization create, update, delete 
// @tag.name OrganizationMember
// @tag.description Organization member role update, leave organization, remove member, transfer ownership
// @tag.name OrganizationInvitation
// @tag.description sent, respond, track of invitation
// @tag.name UserOrganizationData
// @tag.description get users organizations details 
// @tag.name PublicAPI
// @tag.description inter service apis
// @securitydefinitions.apikey UserIDAuth
// @in header
// @name Auth-user

// @securitydefinitions.apikey jwtAuth
// @in header
// @name Authorization
func main() {
	err := config.LoadConfig("../config")
	if err != nil {
		panic(fmt.Sprintf("cannot load config: %v", err))
	}
	db := db.Connect()
	defer db.Close()
	repos := repository.InitRepositories(db)
	go cronjob.InitializeCronjob(repos)
	controllers := controller.InitControllers(repos)
	router := routes.InitializeRouter(controllers)
	router.Mount("/swagger/", httpSwagger.WrapHandler)
	fmt.Println("server started")
	http.ListenAndServe(":"+config.ConfigVal.Port, router)
}
