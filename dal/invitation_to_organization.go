package dal

import (
	"database/sql"
	error_handling "organization/error"
	"organization/model/request"
	"organization/model/response"

	"github.com/lib/pq"
)

func InvitationToOrganization(db *sql.DB, invitationToOrganization request.InvitationToOrganization, userID string) error {
	var id string
	err := db.QueryRow("INSERT INTO public.invitation (invitee, organization_id, role, invited_by) VALUES ($1, $2, $3, $4) returning id", invitationToOrganization.Invitee, invitationToOrganization.OrganizationID, invitationToOrganization.Role, userID).Scan(&id)
	if dbErr, ok := err.(*pq.Error); ok {
		errCode := dbErr.Code
		switch errCode {
		case "23503":
			// foreign key violation
			return error_handling.OrganizationDoesNotExist

		case "23505":
			// unique constraint violation
			return error_handling.AlreadyInvited

		}
		return error_handling.InternalServerError
	}
	return nil
}

func TrackAllInvitations(db *sql.DB, userID string) ([]response.InvitationDetails, error) {
	rows, err := db.Query("SELECT id, role, organization_id, invited_by, invited_at FROM public.invitation WHERE invitee = $1 ", userID)
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
	return invitationDetailsList, nil
}

func AcceptInvitation(tx *sql.Tx, userID string, organizationID string) (string, string, error) {
	var invitedBy, invitedRole string
	err := tx.QueryRow("DELETE FROM public.invitation WHERE invitee = $1 AND organization_id = $2 returning role, invited_by", userID, organizationID).Scan(&invitedRole, &invitedBy)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", "", error_handling.InvalidDetails
		}
		return "", "", error_handling.InternalServerError
	}
	return invitedRole, invitedBy, nil
}

func RejectInvitation(db *sql.DB, userID string, organizationID string) error {
	result, err := db.Exec("DELETE FROM public.invitation WHERE invitee = $1 AND organization_id = $2", userID, organizationID)
	if err != nil {
		return error_handling.InternalServerError
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return error_handling.InternalServerError
	}
	if rowsAffected == 0 {
		return error_handling.InvalidDetails
	}
	return nil
}

func WithdrawSentInvitations(tx *sql.Tx, userID string, organizationID string) error {
	_, err := tx.Exec("DELETE FROM public.invitation WHERE invited_by = $1 AND organization_id = $2", userID, organizationID)
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}
