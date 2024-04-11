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

func (c *UserController) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value(middleware.UserCtxKey).(string)
	var createOrganization request.CreateOrganization
	err := utils.BodyReadAndValidate(r.Body, &createOrganization, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	organizationId, err := c.repo.CreateOrganization(createOrganization, ownerID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	utils.SuccessMessageResponse(w, 200, response.OrganizationID{OrganizationID: organizationId})
}

func (c *UserController) UpdateOrganizationDetails(w http.ResponseWriter, r *http.Request) {
	memberID := r.Context().Value(middleware.UserCtxKey).(string)
	var updateOrganizationDetails request.UpdateOrganizationDetails
	err := utils.BodyReadAndValidate(r.Body, &updateOrganizationDetails, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	role, err := c.repo.CheckRole(memberID, updateOrganizationDetails.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if role == "admin" || role == "owner" {
		err := c.repo.UpdateOrganizationDetails(memberID, updateOrganizationDetails)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, 200, response.SuccessResponse{Message: constant.ORGANIZATION_DETAILS_UPDATED})
	}
	error_handling.ErrorMessageResponse(w, error_handling.NoAccessRights)
}
