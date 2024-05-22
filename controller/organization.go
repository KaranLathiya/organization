package controller

import (
	"encoding/json"
	"net/http"
	"organization/constant"
	error_handling "organization/error"
	microsoftauth "organization/internal/microsoft-auth"
	"organization/middleware"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"
)

// Create Organization example
//
// @tags Organization
// @Security UserIDAuth
//	@Summary		create organization
//	@Description	create organization 
//	@ID				create-organization
//	@Accept			json
//	@Produce		json
// @Param request body request.CreateOrganization true "The input for create organization"
//	@Success		200		{object}	response.OrganizationID "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/ [post]
func (c *UserController) CreateOrganization(w http.ResponseWriter, r *http.Request) {
	ownerID := r.Context().Value(middleware.UserCtxKey).(string)
	var createOrganization request.CreateOrganization
	err := utils.BodyReadAndValidate(r.Body, &createOrganization, nil)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	organizationID, err := c.repo.CreateOrganization(createOrganization, ownerID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	go func() {
		jwtToken, err := utils.CreateJWT("User", "User")
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		body, err := utils.CallHttpService(jwtToken, constant.USER_SERVICE_BASE_URL+"internal/users/details/"+ownerID, nil, http.MethodGet)
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		var organizationCreatorDetails response.OrganizationCreatorDetails
		err = json.Unmarshal(body, &organizationCreatorDetails)
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		refreshToken,err := c.repo.FetchMicrosoftRefreshToken()
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		microsoftAuthToken, err := microsoftauth.GetAccessTokenUsingRefreshToken(refreshToken)
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		go c.repo.StoreMicrosoftRefreshToken(microsoftAuthToken.RefreshToken)
		var ownerPhoneNumberOrEmail string
		if organizationCreatorDetails.Email != nil {
			ownerPhoneNumberOrEmail = *organizationCreatorDetails.Email
		} else if organizationCreatorDetails.PhoneNumber != nil && organizationCreatorDetails.CountryCode != nil {
			ownerPhoneNumberOrEmail = *organizationCreatorDetails.CountryCode + *organizationCreatorDetails.PhoneNumber
		}
		messgae := utils.OrganizationCreatedMessageTemplate(ownerPhoneNumberOrEmail, createOrganization.Name, organizationCreatorDetails.Fullname)
		err = microsoftauth.SendMessageToChannel(messgae, microsoftAuthToken.AccessToken)
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
	}()
	utils.SuccessMessageResponse(w, http.StatusOK, response.OrganizationID{OrganizationID: organizationID})
}

// Update Organization example
//
// @tags Organization
// @Security UserIDAuth
//	@Summary		update organization
//	@Description	update organization 
//	@ID				update-organization
//	@Accept			json
//	@Produce		json
// @Param request body request.UpdateOrganizationDetails true "The input for update organization"
//	@Success		200		{object}	response.SuccessResponse "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/ [put]
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

// otp for delete organization example
//
// @tags Organization
// @Security UserIDAuth
//	@Summary		otp for delete organization
//	@Description	otp for delete organization on owners's registered mail/phoneNumber
//	@ID				delete-organization
//	@Accept			json
//	@Produce		json
// @Param request body request.OrganizationID true "The input for update organization"
//	@Success		200		{object}	response.SuccessResponse "OK"
//	@Failure		400		{object}	error.CustomError	"Bad Request"
//	@Failure		401		{object}	error.CustomError	"Unauthorized"
//	@Failure		403		{object}	error.CustomError	"Forbidden"
//	@Failure		404		{object}	error.CustomError	"Not Found"
//	@Failure		405		{object}	error.CustomError	"Method Not Allowed"
//	@Failure		409		{object}	error.CustomError	"Conflict"
//	@Failure		500		{object}	error.CustomError	"Internal Server Error"
//	@Router			/organization/delete/otp [post]
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
	organizationName, err := c.repo.FetchOrganizationNameByOrganizationID(organizationID.OrganizationID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	jwtToken, err := utils.CreateJWT("User", "User")
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
	body, err := utils.CallHttpService(jwtToken, constant.USER_SERVICE_BASE_URL+"internal/user/otp", deleteOrganizationByte, http.MethodPost)
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
