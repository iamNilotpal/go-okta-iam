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
- `DELETE /api/v1/users/{userID}` - Delete (deactivate) user
- `POST /api/v1/users/{userID}/activate` - Activate user
- `POST /api/v1/users/{userID}/deactivate` - Deactivate user
- `POST /api/v1/users/{userID}/suspend` - Suspend user
- `POST /api/v1/users/{userID}/unsuspend` - Unsuspend user

### Groups

- `GET /api/v1/groups` - List all groups
- `POST /api/v1/groups` - Create new group
- `GET /api/v1/groups/{id}` - Get group by ID
- `PUT /api/v1/groups/{id}` - Update group
- `DELETE /api/v1/groups/{id}` - Delete group
- `POST /api/v1/groups/{groupId}/users/{userId}` - Add user to group
- `DELETE /api/v1/groups/{groupId}/users/{userId}` - Remove user from group
