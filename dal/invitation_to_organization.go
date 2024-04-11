package dal

import (
	"database/sql"
	error_handling "organization/error"
	"organization/model/request"
	"organization/model/response"
	"organization/utils"

	"github.com/lib/pq"
)

func InvitationToOrganization(db *sql.DB, invitationToOrganization request.InvitationToOrganization, memberID string) (bool,error) {
	var id string
	err := db.QueryRow("INSERT INTO public.invitation (invitee, organization_id, role, invited_at, invited_by) VALUES ($1, $2, $3, $4, $5) returning id", invitationToOrganization.Invitee, invitationToOrganization.OrganizationID, invitationToOrganization.Role, utils.CurrentUTCTime(0), memberID).Scan(&id)
	if dbErr, ok := err.(*pq.Error); ok {
		errCode := dbErr.Code
		switch errCode {
		case "23503":
			// foreign key violation
			return false,error_handling.OrganizationDoesNotExist

		case "23505":
			// unique constraint violation
			return false,nil

		}
		return false,error_handling.InternalServerError
	}
	return true,nil
}

func TrackAllInvitations(db *sql.DB, memberID string) ([]response.InvitationDetails,error) {
	rows, err := db.Query("SELECT id, role, organization_id, invited_by, invited_at FROM public.invite WHERE invitee = $1 ", memberID)
	if err != nil {
		return nil, error_handling.InternalServerError
	}
	var invitationDetailsList []response.InvitationDetails
	for rows.Next() {
		var invitationDetails response.InvitationDetails
		err = rows.Scan(&invitationDetails.ID, &invitationDetails.Role, &invitationDetails.OrganizationID, &invitationDetails.InvitedBy, &invitationDetails.InvitedAt)
		if err != nil {
			return nil, err
		}
		invitationDetailsList = append(invitationDetailsList, invitationDetails)
	}
	return invitationDetailsList,nil
}

