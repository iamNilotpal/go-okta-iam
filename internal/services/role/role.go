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
