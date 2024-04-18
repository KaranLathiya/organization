package dal

import (
	"database/sql"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"

	error_handling "organization/error"
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
