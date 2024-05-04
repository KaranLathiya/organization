package controller

import (
	"encoding/json"
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
	utils.SuccessMessageResponse(w, http.StatusOK, response.OrganizationID{OrganizationID: organizationId})
}

func (c *UserController) UpdateOrganization(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	var updateOrganizationDetails request.UpdateOrganizationDetails
	err := utils.BodyReadAndValidate(r.Body, &updateOrganizationDetails, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	roleOfUser, err := c.repo.CheckRoleOfMember(userID, updateOrganizationDetails.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if roleOfUser == constant.ORGANIZATION_ROLE_ADMIN || roleOfUser == constant.ORGANIZATION_ROLE_OWNER {
		err := c.repo.UpdateOrganizationDetails(userID, updateOrganizationDetails)
		if err != nil {
			error_handling.ErrorMessageResponse(w, err)
			return
		}
		utils.SuccessMessageResponse(w, http.StatusOK, response.SuccessResponse{Message: constant.ORGANIZATION_DETAILS_UPDATED})
		return
	}
	error_handling.ErrorMessageResponse(w, error_handling.NoAccessRights)
}

func (c *UserController) OTPForDeleteOrganization(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	var organizationID request.OrganizationID
	err := utils.BodyReadAndValidate(r.Body, &organizationID, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	roleOfUser, err := c.repo.CheckRoleOfMember(userID, organizationID.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	if !(roleOfUser == constant.ORGANIZATION_ROLE_OWNER) {
		error_handling.ErrorMessageResponse(w, error_handling.OwnerAccessRights)
		return
	}
	organizationName, err := c.repo.GetOrganizationNameByOrganizationID(organizationID.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	jwtToken, err := utils.CreateJWT("Organization", "Organization")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	deleteOrganization := request.DeleteOrganization{
		OrganizationID: organizationID.OrganizationID,
		OwnerID:        userID,
		Name:           organizationName,
	}
	deleteOrganizationByte, err := json.MarshalIndent(deleteOrganization, "", "  ")
	if err != nil {
		error_handling.ErrorMessageResponse(w, error_handling.MarshalError)
		return
	}
	body, err := utils.CallAnotherService(jwtToken, constant.USER_SERVICE_BASE_URL+"internal/user/otp", deleteOrganizationByte, http.MethodPost)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	var successResponse response.SuccessResponse
	err = json.Unmarshal(body, &successResponse)
	if err != nil {
		error_handling.ErrorMessageResponse(w, error_handling.UnmarshalError)
		return
	}
	utils.SuccessMessageResponse(w, http.StatusOK, successResponse)
}

