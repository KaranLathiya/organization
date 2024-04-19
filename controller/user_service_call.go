package controller

import (
	"net/http"
	error_handling "organization/error"
	"organization/middleware"
	"organization/utils"

	"github.com/go-chi/chi"
)

func (c *UserController) FetchAllOrganizationDetailsOfCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	allOrganizationDetailsOfUser, allMemberIDs, err := c.repo.FetchAllOrganizationDetailsOfUser(userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	userDetails, err := utils.CreateJWTAndCallAnotherService(allMemberIDs, "User", "User details of organization")
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
	utils.SuccessMessageResponse(w, 200, allOrganizationDetailsOfUser)
}

func (c *UserController) FetchOrganizationDetailsOfCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	organizationID := chi.URLParam(r, "organization")
	organizationDetailsOfUser, memberIDs, err := c.repo.FetchOrganizationDetailsOfCurrentUser(userID, organizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	userDetails, err := utils.CreateJWTAndCallAnotherService(memberIDs, "User", "User details of organization")
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
	utils.SuccessMessageResponse(w, 200, organizationDetailsOfUser)
}
