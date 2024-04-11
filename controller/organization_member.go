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

func (c *UserController) UpdateMemberRole(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	var updateMemberRole request.UpdateMemberRole
	err := utils.BodyReadAndValidate(r.Body, &updateMemberRole, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
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
