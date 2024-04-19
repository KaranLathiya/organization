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
	role, err := c.repo.CheckRole(userID, updateMemberRole.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if role == "admin" || role == "owner" {
		err := c.repo.UpdateMemberRole(userID, updateMemberRole.Role, updateMemberRole.OrganizationID, updateMemberRole.MemberID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.ROLE_UPDATED})
		return
	}
	error_handling.ErrorMessageResponse(w, error_handling.NoAccessRights)
}

func (c *UserController) LeaveOrganization(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	organizationID := request.OrganizationID{
		OrganizationID: chi.URLParam(r, "organization"),
	}
	err := utils.ValidateStruct(&organizationID, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	role, err := c.repo.CheckRole(userID, organizationID.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if role == "owner" {
		error_handling.ErrorMessageResponse(w, error_handling.OwnerLeaveRestriction)
		return
	}
	err = c.repo.DeleteSentInvitationsAndLeaveOrganization(organizationID.OrganizationID, userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.LEAVE_ORGANIZATION})
}

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
	role, err := c.repo.CheckRole(userID, removeMemberFromOrganization.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if role == "admin" || role == "owner" {
		err := c.repo.DeleteSentInvitationsAndRemoveMemberFromOrganization(removeMemberFromOrganization.OrganizationID, removeMemberFromOrganization.MemberID)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.REMOVED_FROM_ORGANIZATION})
		return
	}
	error_handling.ErrorMessageResponse(w, error_handling.NoAccessRights)
}

func (c *UserController) TransferOwnership(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	var transferOwnership request.TransferOwnership
	err := utils.BodyReadAndValidate(r.Body, &transferOwnership, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if userID == transferOwnership.MemberID {
		error_handling.ErrorMessageResponse(w, error_handling.OwnRoleChangeRestriction)
		return
	}
	err = c.repo.TransferOwnership(transferOwnership.OrganizationID, transferOwnership.MemberID, userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.OWNERSHIP_TRANSFERRED})
}
