package models

import (
	"time"

	"github.com/okta/okta-sdk-golang/v5/okta"
)

const (
	UserStatusActive          string = "ACTIVE"
	UserStatusRecovery        string = "RECOVERY"
	UserStatusSuspended       string = "SUSPENDED"
	UserStatusLockedOut       string = "LOCKED_OUT"
	UserStatusProvisioned     string = "PROVISIONED"
	UserStatusDeprovisioned   string = "DEPROVISIONED"
	UserStatusPasswordExpired string = "PASSWORD_EXPIRED"
)

// User represents a person who can access your system.
type User struct {
	ID          string         `json:"id"`
	Email       string         `json:"email"`
	FirstName   string         `json:"firstName"`
	LastName    string         `json:"lastName"`
	Login       string         `json:"login"`
	Status      string         `json:"status"`
	Created     time.Time      `json:"created"`
	Activated   *time.Time     `json:"activated,omitempty"`
	LastLogin   *time.Time     `json:"lastLogin,omitempty"`
	LastUpdated *time.Time     `json:"lastUpdated,omitempty"`
	Profile     map[string]any `json:"profile,omitempty"`
	Groups      []Group        `json:"groups,omitempty"`
	Roles       []Role         `json:"roles,omitempty"`
}

// CreateUserRequest represents the data needed to create a new user.
type CreateUserRequest struct {
	Email     string         `json:"email"`
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Login     string         `json:"login"`
	Password  string         `json:"password"`
	Profile   map[string]any `json:"profile"`
	Activate  bool           `json:"activate"`
}

// UpdateUserRequest represents the data that can be updated for a user.
type UpdateUserRequest struct {
	FirstName string         `json:"firstName"`
	LastName  string         `json:"lastName"`
	Profile   map[string]any `json:"profile"`
}

// UserGroupAssignment represents assigning a user to a group.
type UserGroupAssignment struct {
	UserID  string `json:"userId" validate:"required"`
	GroupID string `json:"groupId" validate:"required"`
}

func ConvertOktaUserToModel(oktaUser *okta.User) *User {
	user := &User{
		ID:     *oktaUser.Id,
		Status: *oktaUser.Status,
	}

	if oktaUser.Profile != nil {
		user.Email = oktaUser.Profile.GetEmail()
		user.FirstName = oktaUser.Profile.GetFirstName()
		user.LastName = oktaUser.Profile.GetLastName()
		user.Login = oktaUser.Profile.GetLogin()
		user.Profile = oktaUser.Profile.AdditionalProperties
	}

	if oktaUser.Created != nil {
		user.Created = oktaUser.GetCreated()
	}

	if oktaUser.Activated.Get() != nil {
		user.Activated = oktaUser.Activated.Get()
	}

	if oktaUser.LastLogin.Get() != nil {
		user.LastLogin = oktaUser.LastLogin.Get()
	}

	if oktaUser.LastUpdated != nil {
		user.LastUpdated = oktaUser.LastUpdated
	}

	return user
}
