package dal

import (
	"database/sql"
	"organization/model/request"
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
