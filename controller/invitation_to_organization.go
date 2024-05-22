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


// Invitation to organization example
//
// @tags OrganizationInvitation
// @Security UserIDAuth
//	@Summary		invitaton to organization 
//	@Description	invitaton for join organization
//	@ID				invitation-to-organization
//	@Accept			json
//	@Produce		json
// @Param request body request.InvitationToOrganization true "The input for invite to organization"
//	@Success		200		{object}	response.SuccessResponse "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/invitation/ [post]
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


// Track all invitations example
//
// @tags OrganizationInvitation
// @Security UserIDAuth
//	@Summary		track all invitations 
//	@Description	track all invitations of user 
//	@ID				track-invitations-of-user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	[]response.InvitationDetails "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/invitations/ [get]
func (c *UserController) TrackAllInvitations(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	invitationDetailsList, err := c.repo.TrackAllInvitations(userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, invitationDetailsList)
}

// respond to invitation example
//
// @tags OrganizationInvitation
// @Security UserIDAuth
//	@Summary		respond to invitations 
//	@Description	accept or reject organization invitation
//	@ID				respond-to-invitation
//	@Accept			json
//	@Produce		json
// @Param request body request.RespondToInvitation true "The input for respond to invitation"
//	@Success		200		{object}	response.SuccessResponse "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/invitations/ [post]
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
