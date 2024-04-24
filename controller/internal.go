package controller

import (
	"net/http"
	"organization/constant"
	error_handling "organization/error"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"

	"github.com/go-chi/chi"
)

func (c *UserController) GetJWT(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := utils.CreateJWT("Organization", "Organization list of users")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, response.JWTToken{JWTToken: jwtToken})
}

func (c *UserController) FetchOragnizationListOfUsers(w http.ResponseWriter, r *http.Request) {
	var userIDs request.UserIDs
	err := utils.VerifyJWT(r.Header.Get("Authorization"), "Organization", "Organization list of users")
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
	utils.SuccessMessageResponse(w, http.StatusOK, organizationListOfUsers)
}

func (c *UserController) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	err := utils.VerifyJWT(r.Header.Get("Authorization"), "Organization", "Delete the organization")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	organizationID := request.OrganizationID{
		OrganizationID: chi.URLParam(r, "organization"),
	}
	err = utils.ValidateStruct(&organizationID, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	err = c.repo.DeleteOrganization(organizationID.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.ORGANIZATION_DELETED})
}
