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
	roleOfUser, err := c.repo.CheckRoleOfMember(userID, invitationToOrganization.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if roleOfUser == constant.ORGANIZATION_ROLE_ADMIN || roleOfUser == constant.ORGANIZATION_ROLE_OWNER {
		isMemberOfOrganization, err := c.repo.IsMemberOfOrganization(invitationToOrganization.Invitee, invitationToOrganization.OrganizationID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		if isMemberOfOrganization {
			error_handling.ErrorMessageResponse(w, error_handling.AlreadyMember)
			return
		}
		err = c.repo.InvitationToOrganization(invitationToOrganization, userID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.INVITED})
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
	utils.SuccessMessageResponse(w, http.StatusOK, invitationDetailsList)
}

func (c *UserController) RespondToInvitation(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	var respondToInvitation request.RespondToInvitation
	err := utils.BodyReadAndValidate(r.Body, &respondToInvitation, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if *respondToInvitation.InvitationAccept {
		err = c.repo.AcceptInvitationAndJoinTheOrganization(userID, respondToInvitation.OrganizationID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.INVITATION_ACCEPTED})
		return
	}
	err = c.repo.RejectInvitation(userID, respondToInvitation.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.INVITATION_REJECTED})
}
