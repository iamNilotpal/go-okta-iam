package models

import "time"

const (
	RoleTypeSystem string = "SYSTEM"
	RoleTypeCustom string = "CUSTOM"
)

// Role represents a set of permissions that can be assigned to users or groups.
// Examples: "Admin", "ReadOnly", "UserManager", "BillingViewer".
type Role struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	Created     time.Time    `json:"created"`
	LastUpdated time.Time    `json:"lastUpdated"`
	Permissions []Permission `json:"permissions,omitempty"`
	Users       []User       `json:"users,omitempty"`
	Groups      []Group      `json:"groups,omitempty"`
}

// CreateRoleRequest represents the data needed to create a new role.
type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions,omitempty"`
}

// UpdateRoleRequest represents the data that can be updated for a role.
type UpdateRoleRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// UserRoleAssignment represents assigning a role to a user.
type UserRoleAssignment struct {
	UserID string `json:"userId"`
	RoleID string `json:"roleId"`
}

// RolePermissionAssignment represents assigning a permission to a role.
type RolePermissionAssignment struct {
	RoleID       string `json:"roleId"`
	PermissionID string `json:"permissionId"`
}
