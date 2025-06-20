# Flexera IAM Platform

## Overview

This IAM platform serves as a management interface for Flexera's microservices
architecture, providing centralized user, group, role, and permission management
through Okta integration.

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
- `GET /api/v1/users/{userID}/roles` - Get roles of a user
- `PUT /api/v1/users/{userID}/roles/{roleID}` - Assign a role to a user
- `DELETE /api/v1/users/{userID}/roles/{roleID}` - Unassign a role from a user

### Groups

- `GET /api/v1/groups` - List all groups
- `POST /api/v1/groups` - Create new group
- `GET /api/v1/groups/{groupID}` - Get group by ID
- `PUT /api/v1/groups/{groupID}` - Update group
- `DELETE /api/v1/groups/{groupID}` - Delete group
- `GET /api/v1/groups/{groupID}/members` - Get group members by ID
- `POST /api/v1/groups/{groupID}/users/{userID}` - Add user to group
- `DELETE /api/v1/groups/{groupID}/users/{userID}` - Remove user from group
- `GET /api/v1/groups/{groupID}/roles` - Get roles of a group
- `PUT /api/v1/groups/{groupID}/roles/{roleID}` - Assign a role to a group
- `DELETE /api/v1/groups/{groupID}/roles/{roleID}` - Unassign a role from a
  group

### Roles

- `GET /api/v1/roles` - List all roles
- `POST /api/v1/roles` - Create new role
- `GET /api/v1/roles/{roleID}` - Get role by ID
- `PUT /api/v1/roles/{roleID}` - Update role
- `DELETE /api/v1/roles/{roleID}` - Delete role
