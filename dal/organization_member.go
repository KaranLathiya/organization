package dal

import (
	"database/sql"
	error_handling "organization/error"
)

func AddMemberToOrganization(tx *sql.Tx, organizationID string, memberID string, role string, invitedBy *string) error {
	var id string
	err := tx.QueryRow("INSERT INTO public.member (organization_id, user_id, role, invited_by) VALUES ($1, $2, $3, $4) returning id", organizationID, memberID, role, invitedBy).Scan(&id)
	if err != nil {
		return error_handling.DatabaseErrorShow(err)
	}
	return nil
}

func RemoveMemberFromOrganization(tx *sql.Tx, memberID string, organizationID string) error {
	return RemoveOrganizationMember(tx, memberID, organizationID)
}

func LeaveOrganization(tx *sql.Tx, userID string, organizationID string) error {
	return RemoveOrganizationMember(tx, userID, organizationID)
}

func RemoveOrganizationMember(tx *sql.Tx, memberID string, organizationID string) error {
	result, err := tx.Exec("DELETE FROM public.member WHERE user_id = $1 AND organization_id = $2", memberID, organizationID)
	if err != nil {
		return error_handling.DatabaseErrorShow(err)
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

func CheckRoleOfMember(db *sql.DB, memberID string, organizationID string) (string, error) {
	var role string
	err := db.QueryRow("SELECT role FROM public.member WHERE user_id = $1 AND organization_id = $2", memberID, organizationID).Scan(&role)
	if err != nil {
		return "", error_handling.DatabaseErrorShow(err)
	}
	return role, nil
}

func IsMemberOfOrganization(db *sql.DB, memberID string, organizationID string) (bool, error) {
	var id string
	err := db.QueryRow("SELECT role FROM public.member WHERE user_id = $1 AND organization_id = $2", memberID, organizationID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, error_handling.DatabaseErrorShow(err)
	}
	return true, nil
}

func IsMemberInvitedByOrganization(db *sql.DB, memberID string, organizationID string) (bool, error) {
	var id string
	err := db.QueryRow("SELECT id FROM public.invitation WHERE invitee = $1 AND organization_id = $2", memberID, organizationID).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, error_handling.DatabaseErrorShow(err)
	}
	return true, nil
}

func UpdateMemberRole(db *sql.DB, userID string, role string, organizationID string, memberID string) error {
	result, err := db.Exec("UPDATE public.member SET role = $1, updated_by = $2, updated_at = current_timestamp() WHERE user_id = $3 AND organization_id = $4;", role, userID, memberID, organizationID)
	if err != nil {
		return error_handling.DatabaseErrorShow(err)
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
	result, err := tx.Exec("UPDATE public.member SET role = $1, updated_by = $2, updated_at = current_timestamp() WHERE user_id = $3 AND organization_id = $4;", role, userID, memberID, organizationID)
	if err != nil {
		return error_handling.DatabaseErrorShow(err)
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
