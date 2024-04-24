package repository

import (
	"database/sql"
	"organization/constant"
	"organization/dal"
	error_handling "organization/error"
	"organization/model/request"
	"organization/model/response"
)

type Repository interface {
	CreateOrganization(organizationCreate request.CreateOrganization, ownerID string) (string, error)
	CheckRole(memberID string, organizationID string) (string, error)
	InvitationToOrganization(invitationToOrganization request.InvitationToOrganization, userID string) (bool, error)
	TrackAllInvitations(userID string) ([]response.InvitationDetails, error)
	UpdateOrganizationDetails(userID string, updateOrganizationDetails request.UpdateOrganizationDetails) error
	AcceptInvitationAndJoinTheOrganization(userID string, organizationID string) error
	RejectInvitation(userID string, organizationID string) error
	UpdateMemberRole(userID string, role string, organizationID string, memberID string) error
	DeleteSentInvitationsAndLeaveOrganization(organizationID string, userID string) error
	DeleteSentInvitationsAndRemoveMemberFromOrganization(organizationID string, memberID string) error
	TransferOwnership(organizationID string, memberID string, userID string) error
	FetchAllOrganizationDetailsOfUser(userID string) (response.AllOrganizationDetailsOfUser, []string, error)
	FetchOrganizationDetailsOfCurrentUser(userID string, organizationID string) (response.OrganizationDetailsOfUser, []string, error)
	FetchOragnizationListOfUsers(userIDs []string) ([]response.OrganizationListOfUser, error)
	GetOrganizationNameByOrganizationID(organizationID string) (string, error)
	DeleteOrganization(organizationID string) error
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
	err = dal.AddMemberToOrganization(tx, organizationID, ownerID, constant.ORGANIZATION_ROLE_OWNER, nil)
	if err != nil {
		return "", err
	}
	err = tx.Commit()
	if err != nil {
		return "", error_handling.InternalServerError
	}
	return organizationID, nil
}

func (r *Repositories) CheckRole(userID string, organizationID string) (string, error) {
	return dal.CheckRole(r.db, userID, organizationID)
}

func (r *Repositories) InvitationToOrganization(invitationToOrganization request.InvitationToOrganization, userID string) (bool, error) {
	isMemberOfOrganization, err := dal.IsMemberOfOrganization(r.db, invitationToOrganization.Invitee, invitationToOrganization.OrganizationID)
	if err != nil {
		return false, err
	}
	if isMemberOfOrganization {
		return false, error_handling.AlreadyMember
	}
	return dal.InvitationToOrganization(r.db, invitationToOrganization, userID)
}

func (r *Repositories) TrackAllInvitations(userID string) ([]response.InvitationDetails, error) {
	return dal.TrackAllInvitations(r.db, userID)
}

func (r *Repositories) UpdateOrganizationDetails(userID string, updateOrganizationDetails request.UpdateOrganizationDetails) error {
	return dal.UpdateOrganizationDetails(r.db, userID, updateOrganizationDetails)
}

func (r *Repositories) AcceptInvitationAndJoinTheOrganization(userID string, organizationID string) error {
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
	roleOfUser, err := dal.CheckRole(r.db, memberID, organizationID)
	if err != nil {
		return err
	}
	if roleOfUser == constant.ORGANIZATION_ROLE_OWNER {
		return error_handling.OwnerRoleChangeRestriction
	}
	return dal.UpdateMemberRole(r.db, userID, role, organizationID, memberID)
}

func (r *Repositories) DeleteSentInvitationsAndLeaveOrganization(organizationID string, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	defer tx.Rollback()
	err = dal.DeleteSentInvitations(tx, userID, organizationID)
	if err != nil {
		return err
	}
	err = dal.LeaveOrganization(tx, userID, organizationID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func (r *Repositories) DeleteSentInvitationsAndRemoveMemberFromOrganization(organizationID string, memberID string) error {
	roleOfMember, err := dal.CheckRole(r.db, memberID, organizationID)
	if err != nil {
		return err
	}
	if roleOfMember == constant.ORGANIZATION_ROLE_OWNER {
		return error_handling.OwnerRemoveRestriction
	}
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	defer tx.Rollback()
	err = dal.DeleteSentInvitations(tx, memberID, organizationID)
	if err != nil {
		return err
	}
	err = dal.RemoveMemberFromOrganization(tx, memberID, organizationID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func (r *Repositories) TransferOwnership(organizationID string, memberID string, userID string) error {
	roleOfMember, err := dal.CheckRole(r.db, memberID, organizationID)
	if err != nil {
		return err
	}
	if roleOfMember != constant.ORGANIZATION_ROLE_OWNER {
		return error_handling.OwnerAccessRights
	}
	tx, err := r.db.Begin()
	if err != nil {
		return error_handling.InternalServerError
	}
	defer tx.Rollback()
	err = dal.UpdateMemberRoleWithTransaction(tx, userID, constant.ORGANIZATION_ROLE_ADMIN, organizationID, userID)
	if err != nil {
		return err
	}
	err = dal.UpdateMemberRoleWithTransaction(tx, userID, constant.ORGANIZATION_ROLE_OWNER, organizationID, memberID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func (r *Repositories) FetchAllOrganizationDetailsOfUser(userID string) (response.AllOrganizationDetailsOfUser, []string, error) {
	return dal.FetchAllOrganizationDetailsOfUser(r.db, userID)
}

func (r *Repositories) FetchOrganizationDetailsOfCurrentUser(userID string, organizationID string) (response.OrganizationDetailsOfUser, []string, error) {
	isMemberOfOrganization, err := dal.IsMemberOfOrganization(r.db, userID, organizationID)
	if err != nil {
		return response.OrganizationDetailsOfUser{}, nil, err
	}
	if !isMemberOfOrganization {
		isMemberInvitedByOrganization, err := dal.IsMemberInvitedByOrganization(r.db, userID, organizationID)
		if err != nil {
			return response.OrganizationDetailsOfUser{}, nil, err
		}
		if !isMemberInvitedByOrganization {
			return response.OrganizationDetailsOfUser{}, nil, error_handling.NotMemberOfOrganization
		}
		return dal.FetchOnlyOrganizationDetailsOfCurrentUser(r.db, userID, organizationID)
	}
	return dal.FetchOrganizationDetailsOfCurrentUser(r.db, userID, organizationID)
}

func (r *Repositories) FetchOragnizationListOfUsers(userIDs []string) ([]response.OrganizationListOfUser, error) {
	return dal.FetchOragnizationListOfUsers(r.db, userIDs)
}

func (r *Repositories) GetOrganizationNameByOrganizationID(organizationID string) (string, error) {
	return dal.GetOrganizationNameByOrganizationID(r.db, organizationID)
}

func (r *Repositories) DeleteOrganization(organizationID string) error {
	return dal.DeleteOrganizationByOrganizationID(r.db, organizationID)
}
