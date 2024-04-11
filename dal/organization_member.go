package dal

import (
	"database/sql"
	error_handling "organization/error"
	"organization/utils"

	"github.com/lib/pq"
)

func AddMembersToOrganization(tx *sql.Tx, organizationID string, memberID string, role string, invitedBy *string) error {
	var id string
	err := tx.QueryRow("INSERT INTO public.member (organization_id, member_id, role, joined_at, invited_by) VALUES ($1, $2, $3, $4, $5) returning id", organizationID, memberID, role, utils.CurrentUTCTime(0), invitedBy).Scan(&id)
	if dbErr, ok := err.(*pq.Error); ok {
		errCode := dbErr.Code
		switch errCode {
		case "23505":
			// unique constraint violation
			return error_handling.OrganizationDoesNotExist
		}
		return error_handling.InternalServerError
	}
	return nil
}

func CheckRole(db *sql.DB, memberID string, organizationID string) (string,error) {
	var role string
	err := db.QueryRow("SELECT role from public.member WHERE member_id = $1 AND organization_id = $2", memberID, organizationID).Scan(&role)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", error_handling.InvalidDetails
		}
		return "", error_handling.InternalServerError
	}
	return role, nil
}
