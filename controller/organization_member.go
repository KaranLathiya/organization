package controller

import (
	"net/http"
	"organization/constant"
	error_handling "organization/error"
	"organization/middleware"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"

	"github.com/go-chi/chi"
)

// Update Organization member role example
//
// @tags OrganizationMember
// @Security UserIDAuth
//	@Summary		update organization memeber role
//	@Description	update organization memeber role
//	@ID				update-member-role
//	@Accept			json
//	@Produce		json
// @Param request body request.UpdateMemberRole true "The input for update organization member role"
//	@Success		200		{object}	response.SuccessResponse "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/members/role/ [put]
func (c *UserController) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	var updateMemberRole request.UpdateMemberRole
	err := utils.BodyReadAndValidate(r.Body, &updateMemberRole, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if userID == updateMemberRole.MemberID {
		error_handling.ErrorMessageResponse(w, error_handling.OwnRoleChangeRestriction)
		return
	}
	roleOfUser, err := c.repo.CheckRoleOfMember(userID, updateMemberRole.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if roleOfUser == constant.ORGANIZATION_ROLE_ADMIN || roleOfUser == constant.ORGANIZATION_ROLE_OWNER {
		roleOfMember, err := c.repo.CheckRoleOfMember(updateMemberRole.MemberID, updateMemberRole.OrganizationID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		if roleOfMember == constant.ORGANIZATION_ROLE_OWNER {
			error_handling.ErrorMessageResponse(w, error_handling.OwnerRoleChangeRestriction)
			return
		}
		err = c.repo.UpdateMemberRole(userID, updateMemberRole.Role, updateMemberRole.OrganizationID, updateMemberRole.MemberID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.ROLE_UPDATED})
		return
	}
	error_handling.ErrorMessageResponse(w, error_handling.NoAccessRights)
}

// leave organization example
//
// @tags OrganizationMember
// @Security UserIDAuth
//	@Summary		leave organization
//	@Description	organization member leave the organization
//	@ID				leave-organization
//	@Accept			json
//	@Produce		json
// @Param  organization-id path string true "organizationID"
//	@Success		200		{object}	response.SuccessResponse "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/{organization-id}/member/leave [delete]
func (c *UserController) LeaveOrganization(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	organizationID := request.OrganizationID{
		OrganizationID: chi.URLParam(r, "organization-id"),
	}
	err := utils.ValidateStruct(&organizationID, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	roleOfUser, err := c.repo.CheckRoleOfMember(userID, organizationID.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if roleOfUser == constant.ORGANIZATION_ROLE_OWNER {
		error_handling.ErrorMessageResponse(w, error_handling.OwnerLeaveRestriction)
		return
	}
	err = c.repo.WithdrawSentInvitationsAndLeaveOrganization(organizationID.OrganizationID, userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.LEAVE_ORGANIZATION})
}

// Remove organization member example
//
// @tags OrganizationMember
// @Security UserIDAuth
//	@Summary		remove organization member
//	@Description	remove organization member from the organization
//	@ID				remove-organization-member
//	@Accept			json
//	@Produce		json
// @Param     organization     query     string     true  " "
// @Param     member           query     string     true   " "
//	@Success		200		{object}	response.SuccessResponse "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/members/ [delete]
func (c *UserController) RemoveMemberFromOrganization(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	removeMemberFromOrganization := request.RemoveMemberFromOrganization{
		OrganizationID: r.FormValue("organization"),
		MemberID:       r.FormValue("member"),
	}
	err := utils.ValidateStruct(removeMemberFromOrganization, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if userID == removeMemberFromOrganization.MemberID {
		error_handling.ErrorMessageResponse(w, error_handling.OwnRemoveRestriction)
		return
	}
	roleOfUser, err := c.repo.CheckRoleOfMember(userID, removeMemberFromOrganization.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if roleOfUser == constant.ORGANIZATION_ROLE_ADMIN || roleOfUser == constant.ORGANIZATION_ROLE_OWNER {
		roleOfMember, err := c.repo.CheckRoleOfMember(removeMemberFromOrganization.MemberID, removeMemberFromOrganization.OrganizationID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		if roleOfMember == constant.ORGANIZATION_ROLE_OWNER {
			error_handling.ErrorMessageResponse(w, error_handling.OwnerRoleChangeRestriction)
			return
		}
		err = c.repo.WithdrawSentInvitationsAndRemoveMemberFromOrganization(removeMemberFromOrganization.OrganizationID, removeMemberFromOrganization.MemberID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.REMOVED_FROM_ORGANIZATION})
		return
	}
	error_handling.ErrorMessageResponse(w, error_handling.NoAccessRights)
}

// Transfer organization ownership example
//
// @tags OrganizationMember
// @Security UserIDAuth
//	@Summary		transfer organization ownership
//	@Description	transfer organization ownership by the owner
//	@ID				transfer-organization-ownership
//	@Accept			json
//	@Produce		json
// @Param     organization     query     string     true  " "
// @Param     member           query     string     true   " "
//	@Success		200		{object}	response.SuccessResponse "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/members/transfer-ownership [put]
func (c *UserController) TransferOwnership(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	var transferOwnership request.TransferOwnership
	err := utils.BodyReadAndValidate(r.Body, &transferOwnership, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	roleOfUser, err := c.repo.CheckRoleOfMember(userID, transferOwnership.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if roleOfUser != constant.ORGANIZATION_ROLE_OWNER {
		error_handling.ErrorMessageResponse(w, error_handling.OwnerAccessRights)
		return
	}
	if userID == transferOwnership.MemberID {
		error_handling.ErrorMessageResponse(w, error_handling.OwnRoleChangeRestriction)
		return
	}
	isMemberOfOrganization, err := c.repo.IsMemberOfOrganization(transferOwnership.MemberID, transferOwnership.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if !isMemberOfOrganization {
		error_handling.ErrorMessageResponse(w, error_handling.NotMemberOfOrganization)
		return
	}
	err = c.repo.TransferOwnership(transferOwnership.OrganizationID, transferOwnership.MemberID, userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.OWNERSHIP_TRANSFERRED})
}
