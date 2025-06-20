package models

import (
	"time"

	"github.com/okta/okta-sdk-golang/v5/okta"
)

const (
	RoleTypeSystem string = "SYSTEM"
	RoleTypeCustom string = "CUSTOM"
)

// Role represents a set of permissions that can be assigned to users or groups.
// Examples: "Admin", "ReadOnly", "UserManager", "BillingViewer".
type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"lastUpdated"`
}

// CreateRoleRequest represents the data needed to create a new role.
type CreateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateRoleRequest represents the data that can be updated for a role.
type UpdateRoleRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// UserRoleAssignment represents assigning a role to a user.
type UserRoleAssignment struct {
	UserID string `json:"userId"`
	RoleID string `json:"roleId"`
}

func ConvertOktaIamRoleToModel(oktaRole *okta.IamRole) *Role {
	role := &Role{
		ID:          *oktaRole.Id,
		Type:        RoleTypeCustom,
		Name:        oktaRole.GetLabel(),
		Description: oktaRole.GetDescription(),
		Created:     oktaRole.GetCreated(),
		LastUpdated: oktaRole.GetLastUpdated(),
	}
	return role
}

func ConvertOktaRoleToModel(assignment *okta.Role) *Role {
	role := &Role{
		ID:          assignment.GetId(),
		Type:        RoleTypeCustom,
		Name:        assignment.GetLabel(),
		Description: assignment.GetDescription(),
		Created:     assignment.GetCreated(),
		LastUpdated: assignment.GetLastUpdated(),
	}
	return role
}
