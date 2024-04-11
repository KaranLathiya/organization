package controller

import (
	"organization/repository"
)

type UserController struct {
	repo repository.Repository
}

func InitControllers(repo repository.Repository) *UserController {
	return &UserController{repo: repo}
}
