package controller

import (
	"net/http"
	"organization/constant"
	error_handling "organization/error"
	"organization/middleware"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"
)

func (c *UserController) InvitationToOrganization(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	var invitationToOrganization request.InvitationToOrganization
	err := utils.BodyReadAndValidate(r.Body, &invitationToOrganization, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	role, err := c.repo.CheckRole(userID, invitationToOrganization.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if role == "admin" || role == "owner" {
		done, err := c.repo.InvitationToOrganization(invitationToOrganization, userID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		if !done {
			utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.ALREADY_INVITED})
			return
		}
		utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.INVITED})
		return
	}
	error_handling.ErrorMessageResponse(w, error_handling.NoAccessRights)
}

func (c *UserController) TrackAllInvitations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)

	invitationDetailsList, err := c.repo.TrackAllInvitations(userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, 200, invitationDetailsList)
}

func (c *UserController) RespondToInvitation(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	respondToInvitation := request.RespondToInvitation{
		OrganizationID: r.FormValue("organization"),
		Respond:        r.FormValue("respond"),
	}
	err := utils.ValidateStruct(respondToInvitation, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if respondToInvitation.Respond == "accept" {
		err = c.repo.AcceptInvitationAndJoinTheOrganization(userID, respondToInvitation.OrganizationID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.INVITATION_ACCEPTED})
		return
	}
	err = c.repo.RejectInvitation(userID, respondToInvitation.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.INVITATION_REJECTED})
}
