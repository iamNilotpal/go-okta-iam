package group_service

import (
	"context"
	"fmt"

	"github.com/iamNilotpal/iam/internal/models"
	"github.com/okta/okta-sdk-golang/v5/okta"
	"go.uber.org/zap"
)

type Service struct {
	client *okta.APIClient
	log    *zap.SugaredLogger
}

func New(log *zap.SugaredLogger, client *okta.APIClient) *Service {
	return &Service{log: log, client: client}
}

func (s *Service) CreateGroup(ctx context.Context, req *models.CreateGroupRequest) (*models.Group, error) {
	s.log.Infow("Creating group in Okta", "name", req.Name)

	profile := okta.GroupProfile{
		Name:        &req.Name,
		Description: &req.Description,
	}

	if len(req.Profile) > 0 {
		profile.AdditionalProperties = req.Profile
	}

	group, response, err := s.client.GroupAPI.
		CreateGroup(ctx).Group(okta.Group{Profile: &profile}).Execute()

	if err != nil {
		s.log.Infow("Failed to create group in Okta", zap.Error(err),
			"name", req.Name,
			"statusCode", response.StatusCode,
		)
		return nil, fmt.Errorf("failed to create group in Okta: %w", err)
	}

	s.log.Infow("Group created successfully in Okta", "groupId", *group.Id, "name", req.Name)
	return models.ConvertOktaGroupToModel(group), nil
}

func (s *Service) GetGroup(ctx context.Context, groupID string) (*models.Group, error) {
	s.log.Infow("Getting group from Okta", "groupId", groupID)

	group, response, err := s.client.GroupAPI.GetGroup(ctx, groupID).Execute()
	if err != nil {
		s.log.Infow("Failed to get group from Okta", zap.Error(err),
			"groupId", groupID,
			"statusCode", response.StatusCode,
		)
		return nil, fmt.Errorf("failed to get group from Okta: %w", err)
	}

	return models.ConvertOktaGroupToModel(group), nil
}

func (s *Service) GetGroups(ctx context.Context) ([]*models.Group, error) {
	s.log.Infow("Getting groups from Okta")

	groups, _, err := s.client.GroupAPI.ListGroups(ctx).Execute()
	if err != nil {
		s.log.Infow("Failed to get groups from Okta", zap.Error(err))
		return nil, fmt.Errorf("failed to get groups from Okta: %w", err)
	}

	result := make([]*models.Group, len(groups))
	for i := range groups {
		result[i] = models.ConvertOktaGroupToModel(&groups[i])
	}

	s.log.Infow("Groups retrieved successfully from Okta", "count", len(result))
	return result, nil
}

func (s *Service) UpdateGroup(ctx context.Context, groupID string, req *models.UpdateGroupRequest) (*models.Group, error) {
	s.log.Infow("Updating group in Okta", zap.String("groupId", groupID))

	var updateNeeded bool
	var profile okta.GroupProfile

	if req.Name != "" {
		updateNeeded = true
		profile.SetName(req.Name)
	}

	if req.Description != "" {
		updateNeeded = true
		profile.SetDescription(req.Description)
	}

	if len(req.Profile) > 0 {
		updateNeeded = true
		profile.AdditionalProperties = req.Profile
	}

	if !updateNeeded {
		return s.GetGroup(ctx, groupID)
	}

	updatedGroup, response, err := s.client.GroupAPI.
		ReplaceGroup(ctx, groupID).Group(okta.Group{Profile: &profile}).Execute()
	if err != nil {
		s.log.Infow("Failed to update group in Okta", zap.Error(err),
			"groupId", groupID,
			"statusCode", response.StatusCode,
		)
		return nil, fmt.Errorf("failed to update group in Okta: %w", err)
	}

	s.log.Info("Group updated successfully in Okta", "groupId", groupID)
	return models.ConvertOktaGroupToModel(updatedGroup), nil
}

func (s *Service) DeleteGroup(ctx context.Context, groupID string) error {
	s.log.Infow("Deleting group from Okta", "groupId", groupID)

	response, err := s.client.GroupAPI.DeleteGroup(ctx, groupID).Execute()
	if err != nil {
		s.log.Infow("Failed to delete group from Okta", zap.Error(err),
			"groupId", groupID,
			"statusCode", response.StatusCode,
		)
		return fmt.Errorf("failed to delete group from Okta: %w", err)
	}

	s.log.Infow("Group deleted successfully from Okta", "groupId", groupID)
	return nil
}
