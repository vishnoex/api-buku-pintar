package middleware

import (
	"buku-pintar/internal/domain/entity"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock PermissionService for testing
type mockPermissionService struct {
	hasPermissionFunc       func(ctx context.Context, userID, permissionName string) (bool, error)
	hasPermissionsFunc      func(ctx context.Context, userID string, permissionNames []string) (bool, error)
	hasAnyPermissionFunc    func(ctx context.Context, userID string, permissionNames []string) (bool, error)
	canUserPerformActionFunc func(ctx context.Context, userID, resource, action string) (bool, error)
	getUserPermissionsFunc  func(ctx context.Context, userID string) ([]*entity.Permission, error)
}

func (m *mockPermissionService) HasPermission(ctx context.Context, userID, permissionName string) (bool, error) {
	if m.hasPermissionFunc != nil {
		return m.hasPermissionFunc(ctx, userID, permissionName)
	}
	return false, nil
}

func (m *mockPermissionService) HasPermissions(ctx context.Context, userID string, permissionNames []string) (bool, error) {
	if m.hasPermissionsFunc != nil {
		return m.hasPermissionsFunc(ctx, userID, permissionNames)
	}
	return false, nil
}

func (m *mockPermissionService) HasAnyPermission(ctx context.Context, userID string, permissionNames []string) (bool, error) {
	if m.hasAnyPermissionFunc != nil {
		return m.hasAnyPermissionFunc(ctx, userID, permissionNames)
	}
	return false, nil
}

func (m *mockPermissionService) CanUserPerformAction(ctx context.Context, userID, resource, action string) (bool, error) {
	if m.canUserPerformActionFunc != nil {
		return m.canUserPerformActionFunc(ctx, userID, resource, action)
	}
	return false, nil
}

func (m *mockPermissionService) GetUserPermissions(ctx context.Context, userID string) ([]*entity.Permission, error) {
	if m.getUserPermissionsFunc != nil {
		return m.getUserPermissionsFunc(ctx, userID)
	}
	return nil, nil
}

// Implement remaining interface methods as stubs
func (m *mockPermissionService) CreatePermission(ctx context.Context, permission *entity.Permission) error {
	return nil
}

func (m *mockPermissionService) GetPermissionByID(ctx context.Context, id string) (*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionService) GetPermissionByName(ctx context.Context, name string) (*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionService) ListPermissions(ctx context.Context, page, pageSize int) ([]*entity.Permission, int, error) {
	return nil, 0, nil
}

func (m *mockPermissionService) UpdatePermission(ctx context.Context, permission *entity.Permission) error {
	return nil
}

func (m *mockPermissionService) DeletePermission(ctx context.Context, id string) error {
	return nil
}

func (m *mockPermissionService) GetPermissionsByResource(ctx context.Context, resource string, limit, offset int) ([]*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionService) GetPermissionsByAction(ctx context.Context, action string, limit, offset int) ([]*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	return nil
}

func (m *mockPermissionService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	return nil
}

func (m *mockPermissionService) BulkAssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

func (m *mockPermissionService) BulkRemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

func (m *mockPermissionService) GeneratePermissionName(resource, action string) string {
	return ""
}

func (m *mockPermissionService) ParsePermissionName(permissionName string) (resource string, action string, err error) {
	return "", "", nil
}

func (m *mockPermissionService) GetUserRole(ctx context.Context, userID string) (*entity.Role, error) {
	return nil, nil
}

func (m *mockPermissionService) BulkCreatePermissions(ctx context.Context, permissions []*entity.Permission) error {
	return nil
}

func (m *mockPermissionService) CountPermissions(ctx context.Context) (int, error) {
	return 0, nil
}

func (m *mockPermissionService) CountPermissionsByResource(ctx context.Context, resource string) (int, error) {
	return 0, nil
}

func (m *mockPermissionService) ExistsPermission(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (m *mockPermissionService) GetPermissionList(ctx context.Context, limit, offset int) ([]*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionService) GetPermissionCount(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockPermissionService) GetPermissionsByResourceAndAction(ctx context.Context, resource, action string) ([]*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionService) GetRolesByPermissionID(ctx context.Context, permissionID string, limit, offset int) ([]*entity.Role, error) {
	return nil, nil
}

func (m *mockPermissionService) CountRolesByPermissionID(ctx context.Context, permissionID string) (int64, error) {
	return 0, nil
}

func (m *mockPermissionService) GetPermissionWithRoles(ctx context.Context, permissionID string) (*entity.PermissionWithRoles, error) {
	return nil, nil
}

func (m *mockPermissionService) GetUserPermissionNames(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}

func (m *mockPermissionService) GetUserPermissionsForResource(ctx context.Context, userID, resource string) ([]*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionService) CreatePermissionsBulk(ctx context.Context, permissions []*entity.Permission) error {
	return nil
}

func (m *mockPermissionService) GetPermissionsByNames(ctx context.Context, names []string) ([]*entity.Permission, error) {
	return nil, nil
}

func (m *mockPermissionService) ValidatePermissionName(ctx context.Context, name string) error {
	return nil
}

func (m *mockPermissionService) IsPermissionNameUnique(ctx context.Context, name string, excludeID *string) (bool, error) {
	return false, nil
}

func (m *mockPermissionService) CanDeletePermission(ctx context.Context, permissionID string) error {
	return nil
}

// Mock RoleService for testing
type mockRoleService struct {
	getRoleByIDFunc func(ctx context.Context, id string) (*entity.Role, error)
}

func (m *mockRoleService) GetRoleByID(ctx context.Context, id string) (*entity.Role, error) {
	if m.getRoleByIDFunc != nil {
		return m.getRoleByIDFunc(ctx, id)
	}
	return nil, nil
}

// Implement remaining interface methods as stubs
func (m *mockRoleService) CreateRole(ctx context.Context, role *entity.Role) error {
	return nil
}

func (m *mockRoleService) GetRoleByName(ctx context.Context, name string) (*entity.Role, error) {
	return nil, nil
}

func (m *mockRoleService) ListRoles(ctx context.Context, page, pageSize int) ([]*entity.Role, int, error) {
	return nil, 0, nil
}

func (m *mockRoleService) UpdateRole(ctx context.Context, role *entity.Role) error {
	return nil
}

func (m *mockRoleService) DeleteRole(ctx context.Context, id string) error {
	return nil
}

func (m *mockRoleService) GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error) {
	return nil, nil
}

func (m *mockRoleService) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	return nil
}

func (m *mockRoleService) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	return nil
}

func (m *mockRoleService) GetUsersByRoleID(ctx context.Context, roleID string, limit, offset int) ([]*entity.User, error) {
	return nil, nil
}

func (m *mockRoleService) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	return nil
}

func (m *mockRoleService) RemoveRoleFromUser(ctx context.Context, userID string) error {
	return nil
}

func (m *mockRoleService) BulkAssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

func (m *mockRoleService) BulkRemovePermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

func (m *mockRoleService) CountRoles(ctx context.Context) (int, error) {
	return 0, nil
}

func (m *mockRoleService) ExistsRole(ctx context.Context, name string) (bool, error) {
	return false, nil
}

func (m *mockRoleService) GetDefaultRole(ctx context.Context) (*entity.Role, error) {
	return nil, nil
}

func (m *mockRoleService) GetRoleList(ctx context.Context, limit, offset int) ([]*entity.Role, error) {
	return nil, nil
}

func (m *mockRoleService) GetRoleCount(ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *mockRoleService) AssignPermissionsToRole(ctx context.Context, roleID string, permissionIDs []string) error {
	return nil
}

func (m *mockRoleService) RemoveAllPermissionsFromRole(ctx context.Context, roleID string) error {
	return nil
}

func (m *mockRoleService) GetRoleWithPermissions(ctx context.Context, roleID string) (*entity.RoleWithPermissions, error) {
	return nil, nil
}

func (m *mockRoleService) CountUsersByRoleID(ctx context.Context, roleID string) (int64, error) {
	return 0, nil
}

func (m *mockRoleService) ValidateRoleName(ctx context.Context, name string) error {
	return nil
}

func (m *mockRoleService) IsRoleNameUnique(ctx context.Context, name string, excludeID *string) (bool, error) {
	return false, nil
}

func (m *mockRoleService) CanDeleteRole(ctx context.Context, roleID string) error {
	return nil
}

// TestCheckPermission_Success tests successful permission check
func TestCheckPermission_Success(t *testing.T) {
	mockPermSvc := &mockPermissionService{
		hasPermissionFunc: func(ctx context.Context, userID, permissionName string) (bool, error) {
			return true, nil
		},
	}
	mockRoleSvc := &mockRoleService{}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, nil)

	handler := middleware.CheckPermission(entity.PermissionCategoryCreate)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	req := httptest.NewRequest("POST", "/categories/create", nil)
	
	// Add user to context
	user := &entity.User{ID: "user-123", Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), UserContextKey, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}

	if rr.Body.String() != "success" {
		t.Errorf("Expected 'success', got %s", rr.Body.String())
	}
}

// TestCheckPermission_Denied tests permission denial
func TestCheckPermission_Denied(t *testing.T) {
	mockPermSvc := &mockPermissionService{
		hasPermissionFunc: func(ctx context.Context, userID, permissionName string) (bool, error) {
			return false, nil
		},
	}
	mockRoleSvc := &mockRoleService{}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, nil)

	handler := middleware.CheckPermission(entity.PermissionCategoryCreate)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/categories/create", nil)
	
	user := &entity.User{ID: "user-123", Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), UserContextKey, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("Expected status Forbidden, got %d", rr.Code)
	}

	var response PermissionErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Status != http.StatusForbidden {
		t.Errorf("Expected status 403 in response, got %d", response.Status)
	}
}

// TestCheckPermission_NoUser tests when no user is in context
func TestCheckPermission_NoUser(t *testing.T) {
	mockPermSvc := &mockPermissionService{}
	mockRoleSvc := &mockRoleService{}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, nil)

	handler := middleware.CheckPermission(entity.PermissionCategoryCreate)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/categories/create", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status Unauthorized, got %d", rr.Code)
	}
}

// TestCheckRole_Success tests successful role check
func TestCheckRole_Success(t *testing.T) {
	roleID := "role-123"
	mockPermSvc := &mockPermissionService{}
	mockRoleSvc := &mockRoleService{
		getRoleByIDFunc: func(ctx context.Context, id string) (*entity.Role, error) {
			return &entity.Role{ID: id, Name: "admin"}, nil
		},
	}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, nil)

	handler := middleware.CheckRole("admin")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	
	user := &entity.User{ID: "user-123", Email: "admin@example.com", RoleID: &roleID}
	ctx := context.WithValue(req.Context(), UserContextKey, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

// TestCheckAllPermissions_Success tests AND permission logic
func TestCheckAllPermissions_Success(t *testing.T) {
	mockPermSvc := &mockPermissionService{
		hasPermissionsFunc: func(ctx context.Context, userID string, permissionNames []string) (bool, error) {
			return true, nil
		},
	}
	mockRoleSvc := &mockRoleService{}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, nil)

	handler := middleware.CheckAllPermissions(
		entity.PermissionCategoryCreate,
		entity.PermissionCategoryUpdate,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/categories/advanced", nil)
	
	user := &entity.User{ID: "user-123", Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), UserContextKey, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

// TestCheckAnyPermission_Success tests OR permission logic
func TestCheckAnyPermission_Success(t *testing.T) {
	mockPermSvc := &mockPermissionService{
		hasAnyPermissionFunc: func(ctx context.Context, userID string, permissionNames []string) (bool, error) {
			return true, nil
		},
	}
	mockRoleSvc := &mockRoleService{}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, nil)

	handler := middleware.CheckAnyPermission(
		entity.PermissionCategoryCreate,
		entity.PermissionCategoryUpdate,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/categories/flexible", nil)
	
	user := &entity.User{ID: "user-123", Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), UserContextKey, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

// TestCheckResourceAction_Success tests resource:action permission pattern
func TestCheckResourceAction_Success(t *testing.T) {
	mockPermSvc := &mockPermissionService{
		canUserPerformActionFunc: func(ctx context.Context, userID, resource, action string) (bool, error) {
			return true, nil
		},
	}
	mockRoleSvc := &mockRoleService{}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, nil)

	handler := middleware.CheckResourceAction(
		entity.ResourceCategory,
		entity.ActionCreate,
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/categories", nil)
	
	user := &entity.User{ID: "user-123", Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), UserContextKey, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

// TestInjectPermissions tests permission injection middleware
func TestInjectPermissions(t *testing.T) {
	permissions := []*entity.Permission{
		{ID: "1", Name: entity.PermissionCategoryRead},
		{ID: "2", Name: entity.PermissionCategoryCreate},
	}

	mockPermSvc := &mockPermissionService{
		getUserPermissionsFunc: func(ctx context.Context, userID string) ([]*entity.Permission, error) {
			return permissions, nil
		},
	}
	mockRoleSvc := &mockRoleService{}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, nil)

	handler := middleware.InjectPermissions(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		perms, ok := GetPermissionsFromContext(r.Context())
		if !ok {
			t.Error("Expected permissions in context")
		}
		if len(perms) != 2 {
			t.Errorf("Expected 2 permissions, got %d", len(perms))
		}
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	
	user := &entity.User{ID: "user-123", Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), UserContextKey, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", rr.Code)
	}
}

// TestAuditLogger tests audit logging functionality
func TestAuditLogger(t *testing.T) {
	auditLogger := NewPermissionAuditLogger()

	mockPermSvc := &mockPermissionService{
		hasPermissionFunc: func(ctx context.Context, userID, permissionName string) (bool, error) {
			return true, nil
		},
	}
	mockRoleSvc := &mockRoleService{}

	config := &PermissionMiddlewareConfig{
		EnableAuditLog: true,
		EnableDebug:    false,
		AuditLogger:    auditLogger,
	}

	middleware := NewPermissionMiddleware(mockPermSvc, mockRoleSvc, config)

	handler := middleware.CheckPermission(entity.PermissionCategoryCreate)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("POST", "/categories", nil)
	
	user := &entity.User{ID: "user-123", Email: "test@example.com"}
	ctx := context.WithValue(req.Context(), UserContextKey, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	entries := auditLogger.GetEntries()
	if len(entries) == 0 {
		t.Error("Expected audit log entries")
	}

	entry := entries[0]
	if entry.UserID != "user-123" {
		t.Errorf("Expected user ID user-123, got %s", entry.UserID)
	}
	if !entry.Granted {
		t.Error("Expected permission to be granted")
	}
}
