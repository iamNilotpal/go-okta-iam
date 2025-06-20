# Flexera IAM Platform

## Overview

This IAM platform serves as a management interface for Flexera's microservices
architecture, providing centralized user, group, role, and permission management
through Okta integration.

## What This System Does

- **User Management**: Create, update, activate/deactivate users
- **Group Management**: Organize users into logical groups
- **Role-Based Access Control**: Define and assign roles with specific
  permissions

## API Endpoints

### Users

- `GET /api/v1/users` - List all users
- `POST /api/v1/users` - Create new user
- `GET /api/v1/users/{userID}` - Get user by ID
- `PUT /api/v1/users/{userID}` - Update user
- `DELETE /api/v1/users/{userID}` - Delete user
- `POST /api/v1/users/{userID}/activate` - Activate user
- `POST /api/v1/users/{userID}/deactivate` - Deactivate user
- `POST /api/v1/users/{userID}/suspend` - Suspend user
- `POST /api/v1/users/{userID}/unsuspend` - Unsuspend user

### Groups

- `GET /api/v1/groups` - List all groups
- `POST /api/v1/groups` - Create new group
- `GET /api/v1/groups/{groupID}` - Get group by ID
- `PUT /api/v1/groups/{groupID}` - Update group
- `DELETE /api/v1/groups/{groupID}` - Delete group
- `POST /api/v1/groups/{groupID}/users/{userID}` - Add user to group
- `DELETE /api/v1/groups/{groupID}/users/{userID}` - Remove user from group
