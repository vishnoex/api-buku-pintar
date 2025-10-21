package mysql

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupPermissionRepoMock(t *testing.T) (*permissionRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewPermissionRepository(db).(*permissionRepository)
	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestPermissionRepository_Create(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	desc := "Create ebook permission"
	permission := &entity.Permission{
		ID:          "perm-123",
		Name:        "ebook:create",
		Description: &desc,
		Resource:    "ebook",
		Action:      "create",
		CreatedAt:   now,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO permissions").
			WithArgs(permission.ID, permission.Name, permission.Description, permission.Resource, permission.Action, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, permission)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate name error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO permissions").
			WithArgs(permission.ID, permission.Name, permission.Description, permission.Resource, permission.Action, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		err := repo.Create(ctx, permission)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("nil permission", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission is nil")
	})
}

func TestPermissionRepository_GetByID(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	permissionID := "perm-123"
	now := time.Now()
	desc := "Create ebook permission"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "created_at", "updated_at"}).
			AddRow(permissionID, "ebook:create", desc, "ebook", "create", now, now)

		mock.ExpectQuery("SELECT (.+) FROM permissions WHERE id = ?").
			WithArgs(permissionID).
			WillReturnRows(rows)

		permission, err := repo.GetByID(ctx, permissionID)
		assert.NoError(t, err)
		assert.NotNil(t, permission)
		assert.Equal(t, permissionID, permission.ID)
		assert.Equal(t, "ebook:create", permission.Name)
		assert.Equal(t, "ebook", permission.Resource)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM permissions WHERE id = ?").
			WithArgs(permissionID).
			WillReturnError(sql.ErrNoRows)

		permission, err := repo.GetByID(ctx, permissionID)
		assert.Error(t, err)
		assert.Nil(t, permission)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPermissionRepository_GetByName(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	permissionName := "ebook:create"
	now := time.Now()
	desc := "Create ebook permission"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "created_at", "updated_at"}).
			AddRow("perm-123", permissionName, desc, "ebook", "create", now, now)

		mock.ExpectQuery("SELECT (.+) FROM permissions WHERE name = ?").
			WithArgs(permissionName).
			WillReturnRows(rows)

		permission, err := repo.GetByName(ctx, permissionName)
		assert.NoError(t, err)
		assert.NotNil(t, permission)
		assert.Equal(t, permissionName, permission.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPermissionRepository_ListByResource(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	resource := "ebook"
	now := time.Now()
	desc := "Permission"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "created_at", "updated_at"}).
			AddRow("perm-1", "ebook:create", desc, resource, "create", now, now).
			AddRow("perm-2", "ebook:read", desc, resource, "read", now, now)

		mock.ExpectQuery("SELECT (.+) FROM permissions WHERE resource = ?").
			WithArgs(resource, 10, 0).
			WillReturnRows(rows)

		permissions, err := repo.ListByResource(ctx, resource, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, permissions, 2)
		assert.Equal(t, "ebook:create", permissions[0].Name)
		assert.Equal(t, "ebook:read", permissions[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPermissionRepository_HasPermission(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	userID := "user-123"
	permissionName := "ebook:create"

	t.Run("has permission", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users u").
			WithArgs(userID, permissionName).
			WillReturnRows(rows)

		hasPermission, err := repo.HasPermission(ctx, userID, permissionName)
		assert.NoError(t, err)
		assert.True(t, hasPermission)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("does not have permission", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users u").
			WithArgs(userID, permissionName).
			WillReturnRows(rows)

		hasPermission, err := repo.HasPermission(ctx, userID, permissionName)
		assert.NoError(t, err)
		assert.False(t, hasPermission)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPermissionRepository_GetPermissionsByUserID(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	userID := "user-123"
	now := time.Now()
	desc := "Permission"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "created_at", "updated_at"}).
			AddRow("perm-1", "ebook:create", desc, "ebook", "create", now, now).
			AddRow("perm-2", "article:read", desc, "article", "read", now, now)

		mock.ExpectQuery("SELECT DISTINCT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at FROM permissions p").
			WithArgs(userID).
			WillReturnRows(rows)

		permissions, err := repo.GetPermissionsByUserID(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, permissions, 2)
		assert.Equal(t, "ebook:create", permissions[0].Name)
		assert.Equal(t, "article:read", permissions[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no permissions", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "created_at", "updated_at"})

		mock.ExpectQuery("SELECT DISTINCT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at FROM permissions p").
			WithArgs(userID).
			WillReturnRows(rows)

		permissions, err := repo.GetPermissionsByUserID(ctx, userID)
		assert.NoError(t, err)
		assert.Len(t, permissions, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPermissionRepository_Update(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	desc := "Updated description"
	permission := &entity.Permission{
		ID:          "perm-123",
		Name:        "ebook:create",
		Description: &desc,
		Resource:    "ebook",
		Action:      "create",
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE permissions SET").
			WithArgs(permission.Name, permission.Description, permission.Resource, permission.Action, permission.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(ctx, permission)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE permissions SET").
			WithArgs(permission.Name, permission.Description, permission.Resource, permission.Action, permission.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(ctx, permission)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPermissionRepository_Delete(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	permissionID := "perm-123"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM permissions WHERE id = ?").
			WithArgs(permissionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(ctx, permissionID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM permissions WHERE id = ?").
			WithArgs(permissionID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(ctx, permissionID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestPermissionRepository_CreateBulk(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	desc1 := "Create permission"
	desc2 := "Read permission"
	permissions := []*entity.Permission{
		{
			ID:          "perm-1",
			Name:        "ebook:create",
			Description: &desc1,
			Resource:    "ebook",
			Action:      "create",
			CreatedAt:   now,
		},
		{
			ID:          "perm-2",
			Name:        "ebook:read",
			Description: &desc2,
			Resource:    "ebook",
			Action:      "read",
			CreatedAt:   now,
		},
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		
		for _, perm := range permissions {
			mock.ExpectExec("INSERT INTO permissions").
				WithArgs(perm.ID, perm.Name, perm.Description, perm.Resource, perm.Action, sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1))
		}
		
		mock.ExpectCommit()

		err := repo.CreateBulk(ctx, permissions)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rollback on error", func(t *testing.T) {
		mock.ExpectBegin()
		
		mock.ExpectExec("INSERT INTO permissions").
			WithArgs(permissions[0].ID, permissions[0].Name, permissions[0].Description, permissions[0].Resource, permissions[0].Action, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		
		mock.ExpectExec("INSERT INTO permissions").
			WithArgs(permissions[1].ID, permissions[1].Name, permissions[1].Description, permissions[1].Resource, permissions[1].Action, sqlmock.AnyArg()).
			WillReturnError(sql.ErrConnDone)
		
		mock.ExpectRollback()

		err := repo.CreateBulk(ctx, permissions)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty list", func(t *testing.T) {
		err := repo.CreateBulk(ctx, []*entity.Permission{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "empty")
	})
}

func TestPermissionRepository_GetByNames(t *testing.T) {
	repo, mock, cleanup := setupPermissionRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	names := []string{"ebook:create", "ebook:read"}
	now := time.Now()
	desc := "Permission"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "created_at", "updated_at"}).
			AddRow("perm-1", "ebook:create", desc, "ebook", "create", now, now).
			AddRow("perm-2", "ebook:read", desc, "ebook", "read", now, now)

		mock.ExpectQuery("SELECT (.+) FROM permissions WHERE name IN \\(\\?,\\?\\)").
			WithArgs(names[0], names[1]).
			WillReturnRows(rows)

		permissions, err := repo.GetByNames(ctx, names)
		assert.NoError(t, err)
		assert.Len(t, permissions, 2)
		assert.Equal(t, "ebook:create", permissions[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty names", func(t *testing.T) {
		permissions, err := repo.GetByNames(ctx, []string{})
		assert.Error(t, err)
		assert.Nil(t, permissions)
		assert.Contains(t, err.Error(), "empty")
	})
}
