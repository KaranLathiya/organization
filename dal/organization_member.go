package dal

import (
	"database/sql"
	error_handling "organization/error"
	"organization/utils"

	"github.com/lib/pq"
)

func AddMemberToOrganization(tx *sql.Tx, organizationID string, memberID string, role string, invitedBy *string) error {
	var id string
	err := tx.QueryRow("INSERT INTO public.member (organization_id, user_id, role, joined_at, invited_by) VALUES ($1, $2, $3, $4, $5) returning id", organizationID, memberID, role, utils.CurrentUTCTime(0), invitedBy).Scan(&id)
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

func RemoveMemberFromOrganization(tx *sql.Tx, memberID string, organizationID string) error {
	return RemoveOrganizationMember(tx, organizationID, memberID)
}

func LeaveOrganization(tx *sql.Tx, userID string, organizationID string) error {
	return RemoveOrganizationMember(tx, userID, organizationID)
}

func RemoveOrganizationMember(tx *sql.Tx, memberID string, organizationID string) error {
	result, err := tx.Exec("DELETE FROM public.member WHERE user_id = $1 AND organization_id = $2", memberID, organizationID)
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

func CheckRole(db *sql.DB, memberID string, organizationID string) (string, error) {
	var role string
	err := db.QueryRow("SELECT role FROM public.member WHERE user_id = $1 AND organization_id = $2", memberID, organizationID).Scan(&role)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", error_handling.OrganizationDoesNotExist
		}
		return "", error_handling.InternalServerError
	}
	return role, nil
}

func IsMemberOfOrganization(db *sql.DB, memberID string, organizationID string) (bool, error) {
	var id string
	err := db.QueryRow("SELECT role FROM public.member WHERE user_id = $1 AND organization_id = $2", memberID, organizationID).Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, error_handling.InternalServerError
	}
	return true, nil
}

func IsMemberInvitedByOrganization(db *sql.DB, memberID string, organizationID string) (bool, error) {
	var id string
	err := db.QueryRow("SELECT id FROM public.invitation WHERE invitee = $1 AND organization_id = $2", memberID, organizationID).Scan(&id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, error_handling.InternalServerError
	}
	return true, nil
}

func UpdateMemberRole(db *sql.DB, userID string, role string, organizationID string, memberID string) error {
	result, err := db.Exec("UPDATE public.member SET role = $1, updated_by = $2, updated_at = $3 WHERE user_id = $4 AND organization_id = $5;", role, userID, utils.CurrentUTCTime(0), memberID, organizationID)
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

func UpdateMemberRoleWithTransaction(tx *sql.Tx, userID string, role string, organizationID string, memberID string) error {
	result, err := tx.Exec("UPDATE public.member SET role = $1, updated_by = $2, updated_at = $3 WHERE user_id = $4 AND organization_id = $5;", role, userID, utils.CurrentUTCTime(0), memberID, organizationID)
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
