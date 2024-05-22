package controller

import (
	"net/http"
	error_handling "organization/error"
	"organization/middleware"
	"organization/utils"

	"github.com/go-chi/chi"
)

// fetch all organizations of current user example
//
// @tags UserOrganizationData
// @Security UserIDAuth
//	@Summary		fetch all organizations 
//	@Description	fetch all organizations of current user
//	@ID				organizations-of-user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	response.AllOrganizationDetailsOfUser "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/user/organizations/ [get]
func (c *UserController) FetchAllOrganizationDetailsOfCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	allOrganizationDetailsOfUser, allMemberIDs, err := c.repo.FetchAllOrganizationDetailsOfUser(userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	userDetails, err := utils.CreateJWTAndCallUserServiceForUserDetails(allMemberIDs, "User", "User")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	for _, organizations := range allOrganizationDetailsOfUser.Organizations {
		for _, organizationMember := range *organizations.OrganizationMembers {
			organizationMember.Firstname = userDetails[organizationMember.UserID].Firstname
			organizationMember.Lastname = userDetails[organizationMember.UserID].Lastname
			organizationMember.Fullname = userDetails[organizationMember.UserID].Fullname
			organizationMember.Username = userDetails[organizationMember.UserID].Username
		}
	}
	utils.SuccessMessageResponse(w, http.StatusOK, allOrganizationDetailsOfUser)
}

// fetch organization details of current user example
//
// @tags UserOrganizationData
// @Security UserIDAuth
//	@Summary		fetch organization details of current user
//	@Description	fetch organization details of current user by organizationID
//	@ID				organization-of-user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	response.OrganizationDetailsOfUser "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/user/organizations/{organization-id} [get]
func (c *UserController) FetchOrganizationDetailsOfCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	organizationID := chi.URLParam(r, "organization-id")
	organizationDetailsOfUser, memberIDs, err := c.repo.FetchOrganizationDetailsOfCurrentUser(userID, organizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	userDetails, err := utils.CreateJWTAndCallUserServiceForUserDetails(memberIDs, "User", "User")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	for _, organizationMember := range *organizationDetailsOfUser.Organization.OrganizationMembers {
		organizationMember.Firstname = userDetails[organizationMember.UserID].Firstname
		organizationMember.Lastname = userDetails[organizationMember.UserID].Lastname
		organizationMember.Fullname = userDetails[organizationMember.UserID].Fullname
		organizationMember.Username = userDetails[organizationMember.UserID].Username
	}
	utils.SuccessMessageResponse(w, http.StatusOK, organizationDetailsOfUser)
}
