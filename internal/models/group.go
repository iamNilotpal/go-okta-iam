package models

import (
	"time"

	"github.com/okta/okta-sdk-golang/v5/okta"
)

const (
	GroupTypeBuiltIn string = "BUILT_IN"
	GroupTypeApp     string = "APP_GROUP"
	GroupTypeOkta    string = "OKTA_GROUP"
)

// Group represents a collection of users with similar access needs.
// For example: "Engineering", "Sales", "Managers", "contractors".
type Group struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Type        string         `json:"type"`
	Created     time.Time      `json:"created"`
	LastUpdated time.Time      `json:"lastUpdated"`
	Profile     map[string]any `json:"profile,omitempty"`
	Members     []User         `json:"members,omitempty"`
	Roles       []Role         `json:"roles,omitempty"`
}

// CreateGroupRequest represents the data needed to create a new group.
type CreateGroupRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Profile     map[string]any `json:"profile,omitempty"`
}

// UpdateGroupRequest represents the data that can be updated for a group.
type UpdateGroupRequest struct {
	Name        string         `json:"name,omitempty"`
	Description string         `json:"description,omitempty"`
	Profile     map[string]any `json:"profile,omitempty"`
}

// GroupRoleAssignment represents assigning a role to a group
type GroupRoleAssignment struct {
	RoleID  string `json:"roleId"`
	GroupID string `json:"groupId"`
}

func ConvertOktaGroupToModel(oktaGroup *okta.Group) *Group {
	group := &Group{
		ID:   *oktaGroup.Id,
		Type: *oktaGroup.Type,
	}

	if oktaGroup.Profile != nil {
		group.Name = oktaGroup.Profile.GetName()
		group.Description = oktaGroup.Profile.GetDescription()
		group.Profile = oktaGroup.Profile.AdditionalProperties
	}

	if oktaGroup.Created != nil {
		group.Created = oktaGroup.GetCreated()
	}

	if oktaGroup.LastUpdated != nil {
		group.LastUpdated = oktaGroup.GetLastUpdated()
	}

	return group
}
