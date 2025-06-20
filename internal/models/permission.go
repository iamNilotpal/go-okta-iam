package models

import "time"

const (
	ActionRead   string = "read"
	ActionWrite  string = "write"
	ActionAdmin  string = "admin"
	ActionDelete string = "delete"
)

// Permission represents a specific action that can be performed on a resource.
// Examples: "read:users", "write:billing", "admin:system".
type Permission struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"lastUpdated"`
	Scope       string    `json:"scope,omitempty"`
}

// CreatePermissionRequest represents the data needed to create a new permission.
type CreatePermissionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Resource    string `json:"resource"`
	Action      string `json:"action"`
	Scope       string `json:"scope,omitempty"`
}

// UpdatePermissionRequest represents the data that can be updated for a permission.
type UpdatePermissionRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Resource    string `json:"resource,omitempty"`
	Action      string `json:"action,omitempty"`
	Scope       string `json:"scope,omitempty"`
}
