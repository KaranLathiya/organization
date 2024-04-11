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
	CheckRole(memberID string, organizationID string) (string, error)
	InvitationToOrganization(invitationToOrganization request.InvitationToOrganization, memberID string) (bool, error)
	TrackAllInvitations(memberID string) ([]response.InvitationDetails, error)
	UpdateOrganizationDetails(memberID string, updateOrganizationDetails request.UpdateOrganizationDetails) error
	AcceptInvitation(userID string, organizationID string) error
	RejectInvitation(userID string, organizationID string) error
	UpdateMemberRole(userID string, role string, organizationID string, memberID string) error
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
		return "", err
	}
	err = dal.AddMemberToOrganization(tx, organizationID, ownerID, "owner", nil)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", error_handling.InternalServerError
	}
	return organizationID, nil
}

func (r *Repositories) CheckRole(memberID string, organizationID string) (string, error) {
	return dal.CheckRole(r.db, memberID, organizationID)
}

func (r *Repositories) InvitationToOrganization(invitationToOrganization request.InvitationToOrganization, memberID string) (bool, error) {
	isMemberOfOrganization, err := dal.IsMemberOfOrganization(r.db, invitationToOrganization.Invitee, invitationToOrganization.OrganizationID)
	if err != nil {
		return false,err
	}
	if isMemberOfOrganization {
		return false,error_handling.AlreadyMember
	}
	return dal.InvitationToOrganization(r.db, invitationToOrganization, memberID)
}

func (r *Repositories) TrackAllInvitations(memberID string) ([]response.InvitationDetails, error) {
	return dal.TrackAllInvitations(r.db, memberID)
}

func (r *Repositories) UpdateOrganizationDetails(memberID string, updateOrganizationDetails request.UpdateOrganizationDetails) error {
	return dal.UpdateOrganizationDetails(r.db, memberID, updateOrganizationDetails)
}

func (r *Repositories) AcceptInvitation(userID string, organizationID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	defer tx.Rollback()
	invitedRole, invitedBy, err := dal.AcceptInvitation(tx, userID, organizationID)
	if err != nil {
		return err
	}
	err = dal.AddMemberToOrganization(tx, organizationID, userID, invitedRole, &invitedBy)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func (r *Repositories) RejectInvitation(userID string, organizationID string) error {
	return dal.RejectInvitation(r.db, userID, organizationID)
}

func (r *Repositories) UpdateMemberRole(userID string, role string, organizationID string, memberID string) error {
	role, err := dal.CheckRole(r.db, userID, organizationID)
	if err != nil {
		return err
	}
	if role == "owner" {
		return error_handling.NoAccessRights
	}
	return dal.UpdateMemberRole(r.db, userID, role, organizationID, memberID)
}
