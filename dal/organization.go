package dal

import (
	"database/sql"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"

	error_handling "organization/error"

	"github.com/lib/pq"
)

func CreateOrganization(tx *sql.Tx, createOrganization request.CreateOrganization, ownerID string) (string, error) {
	var organizationID string
	err := tx.QueryRow("INSERT INTO public.organization (name, owner_id, created_at, privacy) VALUES ($1, $2, $3, $4) returning id", createOrganization.Name, ownerID, utils.CurrentUTCTime(0), createOrganization.Privacy).Scan(&organizationID)
	if err != nil {
		return "", error_handling.InternalServerError
	}
	return organizationID, nil
}

func UpdateOrganizationDetails(db *sql.DB, memberID string, updateOrganizationDetails request.UpdateOrganizationDetails) error {
	_, err := db.Exec("UPDATE public.organization SET privacy=$1,name=$2 WHERE id=$3 ;", updateOrganizationDetails.Privacy, updateOrganizationDetails.Name, updateOrganizationDetails.OrganizationID)
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func FetchAllOrganizationDetailsOfUser(db *sql.DB, userID string) (response.AllOrganizationDetailsOfUser, []string, error) {
	var allMemberIDs []string
	allOrganizationDetailsOfUser := response.AllOrganizationDetailsOfUser{
		UserID: userID,
	}
	rows, err := db.Query("SELECT member.organization_id, organization.name, organization.owner_id, organization.created_at, organization.privacy,  organization.updated_by, organization.updated_at, org_member.role, org_member.user_id FROM member left join organization ON organization.id = member.organization_id LEFT JOIN member AS org_member ON org_member.organization_id = member.organization_id WHERE member.user_id = $1;", userID)
	if err != nil {
		return allOrganizationDetailsOfUser, nil, error_handling.InternalServerError
	}
	allOrganizationDetailsOfUserMap := make(map[string]response.Organization)
	allMemberIDsMap := make(map[string]bool)
	for rows.Next() {
		var organization response.Organization
		organizationMember := &response.OrganizationMember{}
		err = rows.Scan(&organization.OrganizationID, &organization.Name, &organization.OwnerID, &organization.CreatedAt, &organization.Privacy, &organization.UpdatedBy, &organization.UpdatedAt, &organizationMember.Role, &organizationMember.UserID)
		if err != nil {
			return allOrganizationDetailsOfUser, nil, err
		}
		val, ok := allOrganizationDetailsOfUserMap[organization.OrganizationID]
		if !ok {
			organization.OrganizationMembers = &[]*response.OrganizationMember{}
			*organization.OrganizationMembers = append(*organization.OrganizationMembers, organizationMember)
			allOrganizationDetailsOfUserMap[organization.OrganizationID] = organization
			allOrganizationDetailsOfUser.Organizations = append(allOrganizationDetailsOfUser.Organizations, allOrganizationDetailsOfUserMap[organization.OrganizationID])
		} else {
			*val.OrganizationMembers = append(*val.OrganizationMembers, organizationMember)
		}
		_, ok = allMemberIDsMap[organizationMember.UserID]
		if !ok {
			allMemberIDsMap[organizationMember.UserID] = true
			allMemberIDs = append(allMemberIDs, organizationMember.UserID)
		}
	}
	return allOrganizationDetailsOfUser, allMemberIDs, nil
}

func FetchOnlyOrganizationDetailsOfCurrentUser(db *sql.DB, userID string, organizationID string) (response.OrganizationDetailsOfUser, []string, error) {
	var allMemberIDs []string
	organizationDetailsOfUser := response.OrganizationDetailsOfUser{
		UserID: userID,
	}
	var organization response.Organization
	err := db.QueryRow("SELECT id, name, owner_id FROM public.organization WHERE id = $1;", organizationID).Scan(&organization.OrganizationID, &organization.Name, &organization.OwnerID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return organizationDetailsOfUser, nil, error_handling.OrganizationDoesNotExist
		}
		return organizationDetailsOfUser, nil, error_handling.InternalServerError
	}
	allMemberIDs = append(allMemberIDs, organization.OwnerID)
	organization.OrganizationMembers = &[]*response.OrganizationMember{}
	organizationMember := &response.OrganizationMember{
		UserID: organization.OwnerID,
		Role:   "owner",
	}
	*organization.OrganizationMembers = append(*organization.OrganizationMembers, organizationMember)
	organizationDetailsOfUser.Organization = organization
	return organizationDetailsOfUser, allMemberIDs, nil
}

func FetchOrganizationDetailsOfCurrentUser(db *sql.DB, userID string, organizationID string) (response.OrganizationDetailsOfUser, []string, error) {
	var allMemberIDs []string
	organizationDetailsOfUser := response.OrganizationDetailsOfUser{
		UserID: userID,
	}
	rows, err := db.Query("SELECT org_member.organization_id, organization.name, organization.owner_id, organization.created_at, organization.privacy, organization.updated_by, organization.updated_at, org_member.role, org_member.user_id FROM member AS org_member LEFT JOIN organization ON organization.id = org_member.organization_id WHERE organization.id = $1;", organizationID)
	if err != nil {
		return organizationDetailsOfUser, nil, error_handling.InternalServerError
	}
	var organization response.Organization
	organization.OrganizationMembers = &[]*response.OrganizationMember{}
	for rows.Next() {
		organizationMember := &response.OrganizationMember{}
		err = rows.Scan(&organization.OrganizationID, &organization.Name, &organization.OwnerID, &organization.CreatedAt, &organization.Privacy, &organization.UpdatedBy, &organization.UpdatedAt, &organizationMember.Role, &organizationMember.UserID)
		if err != nil {
			return organizationDetailsOfUser, nil, err
		}
		allMemberIDs = append(allMemberIDs, organizationMember.UserID)
		*organization.OrganizationMembers = append(*organization.OrganizationMembers, organizationMember)
	}
	organizationDetailsOfUser.Organization = organization
	return organizationDetailsOfUser, allMemberIDs, nil
}

func FetchOragnizationListOfUsers(db *sql.DB, userIDs []string) ([]response.OrganizationListOfUser, error) {
	rows, err := db.Query("SELECT user_id,name,role,organization_id FROM public.organization INNER JOIN public.member AS org_member on org_member.organization_id = organization.id WHERE org_member.user_id = ANY($1) ;", pq.Array(userIDs))
	if err != nil {
		return nil, error_handling.InternalServerError
	}
	userMap := make(map[string]response.OrganizationListOfUser)
	var organizationListOfUsers []response.OrganizationListOfUser
	for rows.Next() {
		var organizationListOfUser response.OrganizationListOfUser
		var organizationInfoOfUser response.OrganizationInfoOfUser
		err = rows.Scan(&organizationListOfUser.UserID, &organizationInfoOfUser.Name, &organizationInfoOfUser.Role, &organizationInfoOfUser.OrganizationID)
		if err != nil {
			return nil, err
		}
		organizationListOfUserVal, ok := userMap[organizationListOfUser.UserID]
		if !ok {
			organizationListOfUser.Organizations = &[]response.OrganizationInfoOfUser{}
			*organizationListOfUser.Organizations = append(*organizationListOfUser.Organizations, organizationInfoOfUser)
			organizationListOfUsers = append(organizationListOfUsers, organizationListOfUser)
			userMap[organizationListOfUser.UserID] = organizationListOfUser
		} else {
			*organizationListOfUserVal.Organizations = append(*organizationListOfUserVal.Organizations, organizationInfoOfUser)
		}
	}
	return organizationListOfUsers, nil
}

func GetOrganizationNameByOrganizationID(db *sql.DB, organizationID string) (string, error) {
	var organizationName string
	err := db.QueryRow("SELECT name FROM public.organization WHERE id = $1;", organizationID).Scan(&organizationName)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", error_handling.OrganizationDoesNotExist
		}
		return "", error_handling.InternalServerError
	}
	return organizationName,nil
}
