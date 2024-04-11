package repository

import (
	"database/sql"
	"organization/dal"
	error_handling "organization/error"
	"organization/model/request"
	"organization/model/response"
)

type Repository interface {
	CreateOrganization(organizationCreate request.CreateOrganization, ownerID string) (string, error)
	CheckRole(memberID string,organizationID string)  (string, error)
	InvitationToOrganization(invitationToOrganization request.InvitationToOrganization,memberID string) (bool,error)
	TrackAllInvitations(memberID string) ([]response.InvitationDetails,error)
	UpdateOrganizationDetails(memberID string, updateOrganizationDetails request.UpdateOrganizationDetails) error
}

type Repositories struct {
	db *sql.DB
}

// InitRepositories should be called in main.go
func InitRepositories(db *sql.DB) *Repositories {
	return &Repositories{db: db}
}

func (r *Repositories) CreateOrganization(createOrganization request.CreateOrganization, ownerID string) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", error_handling.InternalServerError
	}
	defer tx.Rollback()
	organizationID, err := dal.CreateOrganization(tx, createOrganization, ownerID)
	if err != nil {
		return "", error_handling.InternalServerError
	}
	err = dal.AddMembersToOrganization(tx, organizationID, ownerID, "owner", nil)
	if err != nil {
		return "", error_handling.InternalServerError
	}
	err = tx.Commit()
	if err != nil {
		return "", error_handling.InternalServerError
	}
	return organizationID, nil
}

func (r *Repositories) CheckRole(memberID string,organizationID string)  (string, error) {
	return dal.CheckRole(r.db, memberID, organizationID)
}

func (r *Repositories) InvitationToOrganization(invitationToOrganization request.InvitationToOrganization,memberID string) (bool,error){
	return dal.InvitationToOrganization(r.db, invitationToOrganization, memberID)
}

func (r *Repositories) TrackAllInvitations(memberID string) ([]response.InvitationDetails,error){
	return dal.TrackAllInvitations(r.db, memberID)
}

func (r *Repositories) UpdateOrganizationDetails(memberID string, updateOrganizationDetails request.UpdateOrganizationDetails) error{
	return dal.UpdateOrganizationDetails(r.db, memberID, updateOrganizationDetails)
}