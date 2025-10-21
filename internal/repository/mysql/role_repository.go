package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"time"
)

type roleRepository struct {
	db *sql.DB
}

// NewRoleRepository creates a new instance of RoleRepository
func NewRoleRepository(db *sql.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

// Create creates a new role
func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	query := `INSERT INTO roles (id, name, description, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	role.CreatedAt = now
	role.UpdatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		role.ID,
		role.Name,
		role.Description,
		role.CreatedAt,
		role.UpdatedAt,
	)
	return err
}

// GetByID retrieves a role by its ID
func (r *roleRepository) GetByID(ctx context.Context, id string) (*entity.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at 
		FROM roles WHERE id = ?`

	role := &entity.Role{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// GetByName retrieves a role by its unique name
func (r *roleRepository) GetByName(ctx context.Context, name string) (*entity.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at 
		FROM roles WHERE name = ?`

	role := &entity.Role{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return role, nil
}

// Update updates an existing role
func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	query := `UPDATE roles 
		SET name = ?, description = ?, updated_at = ? 
		WHERE id = ?`

	role.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query,
		role.Name,
		role.Description,
		role.UpdatedAt,
		role.ID,
	)
	return err
}

// Delete deletes a role by its ID
func (r *roleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM roles WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves all roles with pagination
func (r *roleRepository) List(ctx context.Context, limit, offset int) ([]*entity.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at 
		FROM roles ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []*entity.Role
	for rows.Next() {
		role := &entity.Role{}
		err = rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// Count returns the total number of roles
func (r *roleRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM roles`
	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

// GetPermissionsByRoleID retrieves all permissions for a specific role
func (r *roleRepository) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error) {
	query := `SELECT p.id, p.name, p.resource, p.action, p.description, p.created_at 
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.resource ASC, p.action ASC`

	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*entity.Permission
	for rows.Next() {
		permission := &entity.Permission{}
		err = rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Resource,
			&permission.Action,
			&permission.Description,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// AssignPermissionToRole assigns a single permission to a role
func (r *roleRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE role_id = role_id` // Ignore if already exists

	_, err := r.db.ExecContext(ctx, query, roleID, permissionID)
	return err
}

// RemovePermissionFromRole removes a single permission from a role
func (r *roleRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?`

	_, err := r.db.ExecContext(ctx, query, roleID, permissionID)
	return err
}

// AssignPermissionsToRole assigns multiple permissions to a role (bulk operation)
func (r *roleRepository) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	if len(permissionIDs) == 0 {
		return nil
	}

	// Start transaction for bulk insert
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)
		ON DUPLICATE KEY UPDATE role_id = role_id`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, permissionID := range permissionIDs {
		_, err = stmt.ExecContext(ctx, roleID, permissionID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// RemoveAllPermissionsFromRole removes all permissions from a role
func (r *roleRepository) RemoveAllPermissionsFromRole(ctx context.Context, roleID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = ?`

	_, err := r.db.ExecContext(ctx, query, roleID)
	return err
}

// GetUsersByRoleID retrieves all users with a specific role
func (r *roleRepository) GetUsersByRoleID(ctx context.Context, roleID string, limit, offset int) ([]*entity.User, error) {
	query := `SELECT id, name, email, role_id, password, role, avatar, status, created_at, updated_at 
		FROM users WHERE role_id = ? 
		ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, roleID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		user := &entity.User{}
		err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.RoleID,
			&user.Password,
			&user.Role,
			&user.Avatar,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// CountUsersByRoleID counts the number of users with a specific role
func (r *roleRepository) CountUsersByRoleID(ctx context.Context, roleID string) (int64, error) {
	query := `SELECT COUNT(*) FROM users WHERE role_id = ?`
	var count int64
	err := r.db.QueryRowContext(ctx, query, roleID).Scan(&count)
	return count, err
}
