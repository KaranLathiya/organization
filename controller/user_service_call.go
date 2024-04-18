package controller

import (
	"encoding/json"
	"net/http"
	error_handling "organization/error"
	"organization/internal"
	"organization/middleware"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"
)

func (c *UserController) FetchAllOrganizationDetailsOfCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserCtxKey).(string)
	allOrganizationDetailsOfUser, allMemberIDs, err := c.repo.FetchAllOrganizationDetailsOfUser(userID)
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	jwtToken, err := middleware.CreateJWT()
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}
	userIDs := request.UserIDs{
		UserIDs: allMemberIDs,
	}

	userIDsByte, _ := json.MarshalIndent(userIDs, "", "  ")
	body, err := internal.CallAnotherService(jwtToken, "http://localhost:8000/users/details", userIDsByte, "POST")
	if err != nil {
		error_handling.ErrorMessageResponse(w, err)
		return
	}

	userDetails := make(map[string]response.UserDetails)
	err = json.Unmarshal(body, &userDetails)
	if err != nil {
		error_handling.ErrorMessageResponse(w, error_handling.UnmarshalError)
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
