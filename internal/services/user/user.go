package user_service

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

func (s *Service) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
	var profile okta.UserProfile
	s.log.Infow("Creating user in Okta", "email", req.Email, "login", req.Login)

	profile.SetEmail(req.Email)
	profile.SetLogin(req.Login)
	profile.SetLastName(req.LastName)
	profile.SetFirstName(req.FirstName)

	if len(req.Profile) > 0 {
		profile.AdditionalProperties = req.Profile
	}

	createUserRequest := okta.CreateUserRequest{Profile: profile}
	if req.Password != "" {
		createUserRequest.Credentials = &okta.UserCredentials{
			Password: &okta.PasswordCredential{
				Value: &req.Password,
			},
		}
	}

	// Execute the API call and handle the response
	user, response, err := s.client.UserAPI.CreateUser(ctx).Body(createUserRequest).Activate(req.Activate).Execute()
	if err != nil {
		s.log.Infow("Failed to create user in Okta", zap.Error(err),
			"email", req.Email, "statusCode", response.StatusCode,
		)
		return nil, fmt.Errorf("failed to create user in Okta: %w", err)
	}

	s.log.Infow("User created successfully in Okta",
		"userId", *user.Id,
		"email", req.Email,
		"statusCode", response.StatusCode,
	)

	return models.ConvertOktaUserToModel(user), nil
}

func (s *Service) GetUser(ctx context.Context, userID string) (*models.User, error) {
	s.log.Infow("Getting user from Okta", "userId", userID)

	user, response, err := s.client.UserAPI.GetUser(ctx, userID).Execute()
	if err != nil {
		s.log.Infow("Failed to get user from Okta", zap.Error(err),
			"userId", userID,
			"statusCode", response.StatusCode,
		)
		return nil, fmt.Errorf("failed to get user from Okta: %w", err)
	}

	s.log.Infow("User retrieved successfully from Okta",
		zap.String("userId", userID))

	return models.ConvertOktaUserToModel(&okta.User{
		Id:                    user.Id,
		Created:               user.Created,
		Activated:             user.Activated,
		LastLogin:             user.LastLogin,
		Credentials:           user.Credentials,
		LastUpdated:           user.LastUpdated,
		PasswordChanged:       user.PasswordChanged,
		Profile:               user.Profile,
		RealmId:               user.RealmId,
		Status:                user.Status,
		StatusChanged:         user.StatusChanged,
		TransitioningToStatus: user.TransitioningToStatus,
		Type:                  user.Type,
		Links:                 user.Links,
		AdditionalProperties:  user.AdditionalProperties,
	}), nil
}

func (s *Service) GetUsers(ctx context.Context) ([]*models.User, error) {
	s.log.Infow("Getting users from Okta")

	users, _, err := s.client.UserAPI.ListUsers(ctx).Execute()
	if err != nil {
		s.log.Infow("Failed to get users from Okta", zap.Error(err))
		return nil, fmt.Errorf("failed to get users from Okta: %w", err)
	}

	result := make([]*models.User, len(users))
	for i := range users {
		result[i] = models.ConvertOktaUserToModel(&users[i])
	}

	s.log.Infow("Users retrieved successfully from Okta", "count", len(result))
	return result, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID string, req *models.UpdateUserRequest) (*models.User, error) {
	s.log.Info("Updating user in Okta", zap.String("userId", userID))

	var profile okta.UserProfile
	updateNeeded := false

	if req.FirstName != "" {
		updateNeeded = true
		profile.SetFirstName(req.FirstName)
	}

	if req.LastName != "" {
		updateNeeded = true
		profile.SetLastName(req.LastName)
	}

	if len(req.Profile) > 0 {
		updateNeeded = true
		profile.AdditionalProperties = req.Profile
	}

	if !updateNeeded {
		return s.GetUser(ctx, userID)
	}

	user, response, err := s.client.UserAPI.
		UpdateUser(ctx, userID).User(okta.UpdateUserRequest{Profile: &profile}).Execute()

	if err != nil {
		s.log.Infow("Failed to update user in Okta",
			zap.Error(err),
			"userId", userID,
			"statusCode", response.StatusCode,
		)
		return nil, fmt.Errorf("failed to update user in Okta: %w", err)
	}

	s.log.Info("User updated successfully in Okta", "userId", userID)
	return models.ConvertOktaUserToModel(user), nil
}

func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	s.log.Info("Deleting user in Okta", "userId", userID)

	response, err := s.client.UserAPI.DeactivateUser(ctx, userID).Execute()
	if err != nil {
		s.log.Infow("Failed to deactivate user in Okta", zap.Error(err),
			"userId", userID,
			"statusCode", response.StatusCode,
		)
		return fmt.Errorf("failed to deactivate user in Okta: %w", err)
	}

	response, err = s.client.UserAPI.DeleteUser(ctx, userID).Execute()
	if err != nil {
		s.log.Infow("Failed to delete user in Okta", zap.Error(err),
			"userId", userID,
			"statusCode", response.StatusCode,
		)
		return fmt.Errorf("failed to delete user in Okta: %w", err)
	}

	s.log.Info("User deleted successfully in Okta", zap.String("userId", userID))
	return nil
}

func (s *Service) ActivateUser(ctx context.Context, userID string) error {
	s.log.Info("Activating user in Okta", "userId", userID)

	_, response, err := s.client.UserAPI.ActivateUser(ctx, userID).Execute()
	if err != nil {
		s.log.Infow("Failed to activate user in Okta", zap.Error(err),
			"userId", userID,
			"statusCode", response.StatusCode,
		)
		return fmt.Errorf("failed to activate user in Okta: %w", err)
	}

	s.log.Info("User activated successfully in Okta", zap.String("userId", userID))
	return nil
}

func (s *Service) DeactivateUser(ctx context.Context, userID string) error {
	s.log.Info("Deactivating user in Okta", "userId", userID)

	response, err := s.client.UserAPI.DeactivateUser(ctx, userID).Execute()
	if err != nil {
		s.log.Infow("Failed to deactivate user in Okta", zap.Error(err),
			"userId", userID,
			"statusCode", response.StatusCode,
		)
		return fmt.Errorf("failed to deactivate user in Okta: %w", err)
	}

	s.log.Info("User deactivated successfully in Okta", zap.String("userId", userID))
	return nil
}

func (s *Service) SetUserPassword(ctx context.Context, userID, newPassword string) error {
	s.log.Infow("Setting user password in Okta", "userId", userID)

	changePasswordRequest := okta.ChangePasswordRequest{
		NewPassword: &okta.PasswordCredential{
			Value: &newPassword,
		},
	}

	_, response, err := s.client.UserAPI.
		ChangePassword(ctx, userID).ChangePasswordRequest(changePasswordRequest).Execute()
	if err != nil {
		s.log.Infow("Failed to set user password in Okta", zap.Error(err),
			"userId", userID,
			"statusCode", response.StatusCode,
		)
		return fmt.Errorf("failed to set user password in Okta: %w", err)
	}

	s.log.Infow("User password set successfully in Okta", "userId", userID)
	return nil
}
