package mysql

import (
	"buku-pintar/internal/domain/entity"
	"buku-pintar/internal/domain/repository"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type permissionRepository struct {
	db *sql.DB
}

// NewPermissionRepository creates a new instance of PermissionRepository
func NewPermissionRepository(db *sql.DB) repository.PermissionRepository {
	return &permissionRepository{db: db}
}

// Create creates a new permission
func (r *permissionRepository) Create(ctx context.Context, permission *entity.Permission) error {
	query := `INSERT INTO permissions (id, name, description, resource, action, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	permission.CreatedAt = now

	_, err := r.db.ExecContext(ctx, query,
		permission.ID,
		permission.Name,
		permission.Description,
		permission.Resource,
		permission.Action,
		permission.CreatedAt,
	)
	return err
}

// GetByID retrieves a permission by its ID
func (r *permissionRepository) GetByID(ctx context.Context, id string) (*entity.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at 
		FROM permissions WHERE id = ?`

	permission := &entity.Permission{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description,
		&permission.Resource,
		&permission.Action,
		&permission.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// GetByName retrieves a permission by its name
func (r *permissionRepository) GetByName(ctx context.Context, name string) (*entity.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at 
		FROM permissions WHERE name = ?`

	permission := &entity.Permission{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&permission.ID,
		&permission.Name,
		&permission.Description,
		&permission.Resource,
		&permission.Action,
		&permission.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return permission, nil
}

// Update updates a permission
func (r *permissionRepository) Update(ctx context.Context, permission *entity.Permission) error {
	query := `UPDATE permissions 
		SET name = ?, description = ?, resource = ?, action = ? 
		WHERE id = ?`

	_, err := r.db.ExecContext(ctx, query,
		permission.Name,
		permission.Description,
		permission.Resource,
		permission.Action,
		permission.ID,
	)
	return err
}

// Delete deletes a permission by its ID
func (r *permissionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM permissions WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// List retrieves all permissions with pagination
func (r *permissionRepository) List(ctx context.Context, limit, offset int) ([]*entity.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at 
		FROM permissions 
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []*entity.Permission{}
	for rows.Next() {
		permission := &entity.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, rows.Err()
}

// Count returns the total number of permissions
func (r *permissionRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM permissions`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ListByResource retrieves permissions by resource type
func (r *permissionRepository) ListByResource(ctx context.Context, resource string, limit, offset int) ([]*entity.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at 
		FROM permissions 
		WHERE resource = ?
		ORDER BY action ASC 
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, resource, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []*entity.Permission{}
	for rows.Next() {
		permission := &entity.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, rows.Err()
}

// ListByAction retrieves permissions by action type
func (r *permissionRepository) ListByAction(ctx context.Context, action string, limit, offset int) ([]*entity.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at 
		FROM permissions 
		WHERE action = ?
		ORDER BY resource ASC 
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, action, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []*entity.Permission{}
	for rows.Next() {
		permission := &entity.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, rows.Err()
}

// ListByResourceAndAction retrieves permissions by resource and action
func (r *permissionRepository) ListByResourceAndAction(ctx context.Context, resource, action string, limit, offset int) ([]*entity.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at 
		FROM permissions 
		WHERE resource = ? AND action = ?
		ORDER BY created_at DESC 
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, resource, action, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []*entity.Permission{}
	for rows.Next() {
		permission := &entity.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, rows.Err()
}

// GetRolesByPermissionID retrieves all roles that have a specific permission
func (r *permissionRepository) GetRolesByPermissionID(ctx context.Context, permissionID string, limit, offset int) ([]*entity.Role, error) {
	query := `SELECT r.id, r.name, r.description, r.created_at, r.updated_at 
		FROM roles r
		INNER JOIN role_permissions rp ON r.id = rp.role_id
		WHERE rp.permission_id = ?
		ORDER BY r.name ASC
		LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, permissionID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := []*entity.Role{}
	for rows.Next() {
		role := &entity.Role{}
		err := rows.Scan(
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

	return roles, rows.Err()
}

// CountRolesByPermissionID counts the number of roles that have a specific permission
func (r *permissionRepository) CountRolesByPermissionID(ctx context.Context, permissionID string) (int64, error) {
	query := `SELECT COUNT(DISTINCT r.id) 
		FROM roles r
		INNER JOIN role_permissions rp ON r.id = rp.role_id
		WHERE rp.permission_id = ?`

	var count int64
	err := r.db.QueryRowContext(ctx, query, permissionID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetPermissionsByUserID retrieves all permissions for a user through their role
func (r *permissionRepository) GetPermissionsByUserID(ctx context.Context, userID string) ([]*entity.Permission, error) {
	query := `SELECT DISTINCT p.id, p.name, p.description, p.resource, p.action, p.created_at 
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN users u ON u.role_id = r.id
		WHERE u.id = ?
		ORDER BY p.resource ASC, p.action ASC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []*entity.Permission{}
	for rows.Next() {
		permission := &entity.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, rows.Err()
}

// HasPermission checks if a user has a specific permission
func (r *permissionRepository) HasPermission(ctx context.Context, userID, permissionName string) (bool, error) {
	query := `SELECT EXISTS(
		SELECT 1 
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN users u ON u.role_id = r.id
		WHERE u.id = ? AND p.name = ?
	)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, userID, permissionName).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// HasPermissions checks if a user has all specified permissions (AND logic)
func (r *permissionRepository) HasPermissions(ctx context.Context, userID string, permissionNames []string) (bool, error) {
	if len(permissionNames) == 0 {
		return false, nil
	}

	// Build IN clause with placeholders
	placeholders := make([]string, len(permissionNames))
	args := []interface{}{userID}
	for i := range permissionNames {
		placeholders[i] = "?"
		args = append(args, permissionNames[i])
	}

	// Check if user has ALL permissions (count must equal number of requested permissions)
	query := fmt.Sprintf(`SELECT COUNT(DISTINCT p.name) 
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN users u ON u.role_id = r.id
		WHERE u.id = ? AND p.name IN (%s)`, strings.Join(placeholders, ","))

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	// User has all permissions if count equals the number of requested permissions
	return count == len(permissionNames), nil
}

// CreateBulk creates multiple permissions in a single transaction
func (r *permissionRepository) CreateBulk(ctx context.Context, permissions []*entity.Permission) error {
	if len(permissions) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO permissions (id, name, description, resource, action, created_at) 
		VALUES (?, ?, ?, ?, ?, ?)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, permission := range permissions {
		permission.CreatedAt = now

		_, err := stmt.ExecContext(ctx,
			permission.ID,
			permission.Name,
			permission.Description,
			permission.Resource,
			permission.Action,
			permission.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByNames retrieves multiple permissions by their names
func (r *permissionRepository) GetByNames(ctx context.Context, names []string) ([]*entity.Permission, error) {
	if len(names) == 0 {
		return []*entity.Permission{}, nil
	}

	// Build IN clause with placeholders
	placeholders := make([]string, len(names))
	args := make([]interface{}, len(names))
	for i, name := range names {
		placeholders[i] = "?"
		args[i] = name
	}

	query := fmt.Sprintf(`SELECT id, name, description, resource, action, created_at 
		FROM permissions 
		WHERE name IN (%s)
		ORDER BY resource ASC, action ASC`, strings.Join(placeholders, ","))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := []*entity.Permission{}
	for rows.Next() {
		permission := &entity.Permission{}
		err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, rows.Err()
}
