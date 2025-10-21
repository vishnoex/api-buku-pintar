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

func setupRoleRepoMock(t *testing.T) (*roleRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewRoleRepository(db).(*roleRepository)
	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestRoleRepository_Create(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()
	desc := "Editor role"
	role := &entity.Role{
		ID:          "role-123",
		Name:        "editor",
		Description: &desc,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO roles").
			WithArgs(role.ID, role.Name, role.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(ctx, role)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("duplicate name error", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO roles").
			WithArgs(role.ID, role.Name, role.Description, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		err := repo.Create(ctx, role)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("nil role", func(t *testing.T) {
		err := repo.Create(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role is nil")
	})
}

func TestRoleRepository_GetByID(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	roleID := "role-123"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow(roleID, "editor", "Editor role", now, now)

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE id = ?").
			WithArgs(roleID).
			WillReturnRows(rows)

		role, err := repo.GetByID(ctx, roleID)
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, roleID, role.ID)
		assert.Equal(t, "editor", role.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM roles WHERE id = ?").
			WithArgs(roleID).
			WillReturnError(sql.ErrNoRows)

		role, err := repo.GetByID(ctx, roleID)
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_GetByName(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	roleName := "editor"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow("role-123", roleName, "Editor role", now, now)

		mock.ExpectQuery("SELECT (.+) FROM roles WHERE name = ?").
			WithArgs(roleName).
			WillReturnRows(rows)

		role, err := repo.GetByName(ctx, roleName)
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, roleName, role.Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery("SELECT (.+) FROM roles WHERE name = ?").
			WithArgs(roleName).
			WillReturnError(sql.ErrNoRows)

		role, err := repo.GetByName(ctx, roleName)
		assert.Error(t, err)
		assert.Nil(t, role)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_List(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	now := time.Now()

	t.Run("success with pagination", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"}).
			AddRow("role-1", "admin", "Admin role", now, now).
			AddRow("role-2", "editor", "Editor role", now, now)

		mock.ExpectQuery("SELECT (.+) FROM roles ORDER BY created_at DESC LIMIT ? OFFSET ?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		roles, err := repo.List(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, roles, 2)
		assert.Equal(t, "admin", roles[0].Name)
		assert.Equal(t, "editor", roles[1].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty result", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "created_at", "updated_at"})

		mock.ExpectQuery("SELECT (.+) FROM roles ORDER BY created_at DESC LIMIT ? OFFSET ?").
			WithArgs(10, 0).
			WillReturnRows(rows)

		roles, err := repo.List(ctx, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, roles, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_Update(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	updatedDesc := "Updated description"
	role := &entity.Role{
		ID:          "role-123",
		Name:        "editor",
		Description: &updatedDesc,
	}

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("UPDATE roles SET").
			WithArgs(role.Name, role.Description, sqlmock.AnyArg(), role.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(ctx, role)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("UPDATE roles SET").
			WithArgs(role.Name, role.Description, sqlmock.AnyArg(), role.ID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Update(ctx, role)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("nil role", func(t *testing.T) {
		err := repo.Update(ctx, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "role is nil")
	})
}

func TestRoleRepository_Delete(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	roleID := "role-123"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM roles WHERE id = ?").
			WithArgs(roleID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(ctx, roleID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM roles WHERE id = ?").
			WithArgs(roleID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.Delete(ctx, roleID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_Count(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(5)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM roles").
			WillReturnRows(rows)

		count, err := repo.Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("zero count", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM roles").
			WillReturnRows(rows)

		count, err := repo.Count(ctx)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_GetPermissionsByRoleID(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	roleID := "role-123"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "created_at", "updated_at"}).
			AddRow("perm-1", "ebook:create", "Create ebook", "ebook", "create", now, now).
			AddRow("perm-2", "ebook:read", "Read ebook", "ebook", "read", now, now)

		mock.ExpectQuery("SELECT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at FROM permissions p").
			WithArgs(roleID).
			WillReturnRows(rows)

		permissions, err := repo.GetPermissionsByRoleID(ctx, roleID)
		assert.NoError(t, err)
		assert.Len(t, permissions, 2)
		assert.Equal(t, "ebook:create", permissions[0].Name)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no permissions", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "description", "resource", "action", "created_at", "updated_at"})

		mock.ExpectQuery("SELECT p.id, p.name, p.description, p.resource, p.action, p.created_at, p.updated_at FROM permissions p").
			WithArgs(roleID).
			WillReturnRows(rows)

		permissions, err := repo.GetPermissionsByRoleID(ctx, roleID)
		assert.NoError(t, err)
		assert.Len(t, permissions, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_AssignPermissionToRole(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	roleID := "role-123"
	permissionID := "perm-456"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO role_permissions").
			WithArgs(sqlmock.AnyArg(), roleID, permissionID, sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.AssignPermissionToRole(ctx, roleID, permissionID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("already assigned", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO role_permissions").
			WithArgs(sqlmock.AnyArg(), roleID, permissionID, sqlmock.AnyArg()).
			WillReturnError(sql.ErrNoRows)

		err := repo.AssignPermissionToRole(ctx, roleID, permissionID)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_RemovePermissionFromRole(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	roleID := "role-123"
	permissionID := "perm-456"

	t.Run("success", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM role_permissions WHERE role_id = \\? AND permission_id = \\?").
			WithArgs(roleID, permissionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.RemovePermissionFromRole(ctx, roleID, permissionID)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM role_permissions WHERE role_id = \\? AND permission_id = \\?").
			WithArgs(roleID, permissionID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.RemovePermissionFromRole(ctx, roleID, permissionID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_GetUsersByRoleID(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	roleID := "role-123"
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "name", "picture", "role_id", "created_at", "updated_at"}).
			AddRow("user-1", "user1@example.com", "User 1", "pic1.jpg", roleID, now, now).
			AddRow("user-2", "user2@example.com", "User 2", "pic2.jpg", roleID, now, now)

		mock.ExpectQuery("SELECT (.+) FROM users WHERE role_id = ?").
			WithArgs(roleID, 10, 0).
			WillReturnRows(rows)

		users, err := repo.GetUsersByRoleID(ctx, roleID, 10, 0)
		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "user1@example.com", users[0].Email)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestRoleRepository_CountUsersByRoleID(t *testing.T) {
	repo, mock, cleanup := setupRoleRepoMock(t)
	defer cleanup()

	ctx := context.Background()
	roleID := "role-123"

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"count"}).AddRow(3)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE role_id = ?").
			WithArgs(roleID).
			WillReturnRows(rows)

		count, err := repo.CountUsersByRoleID(ctx, roleID)
		assert.NoError(t, err)
		assert.Equal(t, int64(3), count)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
