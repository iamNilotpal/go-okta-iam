package role_service

import (
	"context"
	"fmt"

	"github.com/okta/okta-sdk-golang/v5/okta"
	"go.uber.org/zap"

	"github.com/iamNilotpal/iam/internal/models"
)

type Service struct {
	client *okta.APIClient
	log    *zap.SugaredLogger
}

func New(log *zap.SugaredLogger, client *okta.APIClient) *Service {
	return &Service{log: log, client: client}
}

func (s *Service) CreateRole(ctx context.Context, req *models.CreateRoleRequest) (*models.Role, error) {
	s.log.Infow("Creating role in Okta", "name", req.Name)

	createRoleRequest := okta.CreateIamRoleRequest{
		Label:       req.Name,
		Description: req.Description,
	}

	role, response, err := s.client.RoleAPI.CreateRole(ctx).Instance(createRoleRequest).Execute()
	if err != nil {
		s.log.Infow("Failed to create role in Okta", zap.Error(err), "name", req.Name, "statusCode", response.StatusCode)
		return nil, fmt.Errorf("failed to create role in Okta: %w", err)
	}

	s.log.Infow("Role created successfully in Okta", "roleId", *role.Id, "name", req.Name)
	return models.ConvertOktaIamRoleToModel(role), nil
}

func (s *Service) GetRole(ctx context.Context, roleID string) (*models.Role, error) {
	s.log.Infow("Getting role from Okta", "roleId", roleID)

	role, response, err := s.client.RoleAPI.GetRole(ctx, roleID).Execute()
	if err != nil {
		s.log.Infow("Failed to get role from Okta", zap.Error(err), "roleId", roleID, "statusCode", response.StatusCode)
		return nil, fmt.Errorf("failed to get role from Okta: %w", err)
	}

	s.log.Infow("Role retrieved successfully from Okta", "roleId", roleID)
	return models.ConvertOktaIamRoleToModel(role), nil
}

func (s *Service) GetRoles(ctx context.Context) ([]*models.Role, error) {
	s.log.Infow("Getting roles from Okta")

	roles, _, err := s.client.RoleAPI.ListRoles(ctx).Execute()
	if err != nil {
		s.log.Infow("Failed to get roles from Okta", zap.Error(err))
		return nil, fmt.Errorf("failed to get roles from Okta: %w", err)
	}

	result := make([]*models.Role, len(roles.Roles))
	for i := range roles.Roles {
		result[i] = models.ConvertOktaIamRoleToModel(&roles.Roles[i])
	}

	s.log.Infow("Roles retrieved successfully from Okta", "count", len(result))
	return result, nil
}

func (s *Service) UpdateRole(ctx context.Context, roleID string, req *models.UpdateRoleRequest) (*models.Role, error) {
	s.log.Infow("Updating role in Okta", "roleId", roleID)

	updateRoleRequest := okta.UpdateIamRoleRequest{}
	updateNeeded := false

	if req.Name != "" {
		updateNeeded = true
		updateRoleRequest.Label = req.Name
	}

	if req.Description != "" {
		updateNeeded = true
		updateRoleRequest.Description = req.Description
	}

	if !updateNeeded {
		return s.GetRole(ctx, roleID)
	}

	role, response, err := s.client.RoleAPI.ReplaceRole(ctx, roleID).Instance(updateRoleRequest).Execute()
	if err != nil {
		s.log.Infow("Failed to update role in Okta", zap.Error(err), "roleId", roleID, "statusCode", response.StatusCode)
		return nil, fmt.Errorf("failed to update role in Okta: %w", err)
	}

	s.log.Infow("Role updated successfully in Okta", "roleId", roleID)
	return models.ConvertOktaIamRoleToModel(role), nil
}

func (s *Service) DeleteRole(ctx context.Context, roleID string) error {
	s.log.Infow("Deleting role from Okta", "roleId", roleID)

	response, err := s.client.RoleAPI.DeleteRole(ctx, roleID).Execute()
	if err != nil {
		s.log.Infow("Failed to delete role from Okta", zap.Error(err),
			"roleId", roleID,
			"statusCode", response.StatusCode,
		)
		return fmt.Errorf("failed to delete role from Okta: %w", err)
	}

	s.log.Infow("Role deleted successfully from Okta", "roleId", roleID)
	return nil
}
