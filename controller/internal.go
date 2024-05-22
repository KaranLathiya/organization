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

// jwt for organization service example
//
// @tags PublicAPI
//	@Summary		jwt for organization service
//	@Description	jwt for organization service call
//	@ID				jwt-for-organization
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	response.JWTToken "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/internal/jwt [get]
func (c *UserController) GetJWTForOragnizationService(w http.ResponseWriter, r *http.Request) {
	jwtToken, err := utils.CreateJWT("Organization", "Organization")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, response.JWTToken{JWTToken: jwtToken})
}

// fetch organization list of users example
//
// @tags PublicAPI
// @Security JWTAuth
//	@Summary		fetch organization list of users
//	@Description	fetch organization list of users 
//	@ID				organization-list-of-users
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]response.OrganizationListOfUser "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/internal/organizations [post]
func (c *UserController) FetchOragnizationListOfUsers(w http.ResponseWriter, r *http.Request) {
	var userIDs request.UserIDs
	err := utils.VerifyJWT(r.Header.Get("Authorization"), "Organization", "Organization")
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

// delete organization example
//
// @tags PublicAPI
// @Security JWTAuth
//	@Summary		delete organization
//	@Description	delete organization after verifying otp
//	@ID				organization-delete
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]response.OrganizationListOfUser "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/internal/organization/{organization-id} [delete]
func (c *UserController) DeleteOrganization(w http.ResponseWriter, r *http.Request) {
	err := utils.VerifyJWT(r.Header.Get("Authorization"), "Organization", "Organization")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	organizationID := request.OrganizationID{
		OrganizationID: chi.URLParam(r, "organization-id"),
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
