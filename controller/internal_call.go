package controller

import (
	"net/http"
	error_handling "organization/error"
	"organization/middleware"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"
)

func (c *UserController) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := middleware.CreateJWT("Organization", "Organization list of users")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, 200, response.JWTToken{JWTToken: jwtToken})
}

func (c *UserController) FetchOragnizationListOfUsers(w http.ResponseWriter, r *http.Request) {
	var userIDs request.UserIDs
	err := middleware.VerifyJWTToken(r.Header.Get("Authorization"), "Organization", "Organization list of users")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	err = utils.BodyReadAndValidate(r.Body, &userIDs, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	organizationListOfUsers, err := c.repo.FetchOragnizationListOfUsers(userIDs.UserIDs)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, 200, organizationListOfUsers)
}
