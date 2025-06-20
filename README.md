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

- `GET /api/users` - List all users
- `POST /api/users` - Create new user
- `GET /api/users/{id}` - Get user by ID
- `PUT /api/users/{id}` - Update user
- `DELETE /api/users/{id}` - Delete (deactivate) user
- `POST /api/users/{id}/activate` - Activate user
- `POST /api/users/{id}/deactivate` - Deactivate user

### Groups

- `GET /api/groups` - List all groups
- `POST /api/groups` - Create new group
- `GET /api/groups/{id}` - Get group by ID
- `PUT /api/groups/{id}` - Update group
- `DELETE /api/groups/{id}` - Delete group
- `POST /api/groups/{groupId}/users/{userId}` - Add user to group
- `DELETE /api/groups/{groupId}/users/{userId}` - Remove user from group
