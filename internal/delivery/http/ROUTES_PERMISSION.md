# API Routes Permission Matrix

This document outlines all API routes and their required permissions for access control.

## Route Organization

Routes are organized into three security levels:
1. **Public Routes** - No authentication required
2. **Authenticated Routes** - Requires valid authentication token
3. **Role-Based Routes** - Requires authentication + specific role/permission

---

## 📖 Public Routes (No Authentication)

### Categories
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/categories` | List categories with pagination | Public |
| GET | `/categories/all` | List all categories | Public |
| GET | `/categories/view/{id}` | Get category by ID | Public |
| GET | `/categories/parent/{parentID}` | List categories by parent | Public |

### Banners
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/banners` | List banners with pagination | Public |
| GET | `/banners/active` | List active banners only | Public |
| GET | `/banners/view/{id}` | Get banner by ID | Public |

### Ebooks
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/ebooks` | List ebooks with pagination | Public |
| GET | `/ebooks/{id}` | Get ebook by ID | Public |
| GET | `/ebooks/slug/{slug}` | Get ebook by slug | Public |

### Summaries
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/summaries` | List summaries with pagination | Public |
| GET | `/summaries/{id}` | Get summary by ID | Public |
| GET | `/summaries/ebook/{ebookID}` | Get summaries by ebook ID | Public |

### Authentication
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| GET | `/oauth2/login` | Initiate OAuth2 login | Public |
| GET | `/oauth2/callback` | OAuth2 callback handler | Public |
| GET | `/oauth2/providers` | List available OAuth2 providers | Public |
| GET | `/oauth2/{provider}/redirect` | OAuth2 provider redirect | Public |

### Users
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| POST | `/users/register` | Register new user | Public |

### Payments
| Method | Endpoint | Description | Access |
|--------|----------|-------------|--------|
| POST | `/payments/callback` | Xendit payment webhook | Public (Webhook) |

---

## 🔒 Authenticated Routes (Requires Login)

These routes require a valid authentication token but no specific role.

### User Profile
| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|---------------------|
| GET | `/users` | Get current user profile | Authenticated |
| PUT | `/users/update` | Update current user profile | Authenticated |
| DELETE | `/users/delete` | Delete current user account | Authenticated |

### Payments
| Method | Endpoint | Description | Required Permission |
|--------|----------|-------------|---------------------|
| POST | `/payments/initiate` | Initiate payment transaction | Authenticated |

---

## 👑 Admin Only Routes

These routes require authentication + **Admin** role.

### Category Management
| Method | Endpoint | Description | Required Role | Permission |
|--------|----------|-------------|---------------|------------|
| POST | `/categories/create` | Create new category | Admin | `category:create` |
| PUT | `/categories/edit/{id}` | Update existing category | Admin | `category:update` |
| DELETE | `/categories/delete/{id}` | Delete category | Admin | `category:delete` |

### Banner Management
| Method | Endpoint | Description | Required Role | Permission |
|--------|----------|-------------|---------------|------------|
| POST | `/banners/create` | Create new banner | Admin | `banner:create` |
| PUT | `/banners/edit/{id}` | Update existing banner | Admin | `banner:update` |
| DELETE | `/banners/delete/{id}` | Delete banner | Admin | `banner:delete` |

---

## ✏️ Editor+ Routes

These routes require authentication + **Editor** OR **Admin** role.

### Summary Management
| Method | Endpoint | Description | Required Roles | Permission |
|--------|----------|-------------|----------------|------------|
| POST | `/summaries/create` | Create new summary | Editor, Admin | `summary:create` |
| PUT | `/summaries/edit/{id}` | Update existing summary | Editor, Admin | `summary:update` |
| DELETE | `/summaries/delete/{id}` | Delete summary | Admin only | `summary:delete` |

> **Note:** Delete operations require Admin role even for editors.

---

## 🔐 Authentication Flow

### Middleware Chain

Routes are protected using a middleware chain:

```go
// Public route (no middleware)
mux.HandleFunc("/ebooks", ebookHandler.ListEbooks)

// Authenticated route (auth middleware only)
mux.Handle("/users", 
    authMiddleware.Authenticate(
        http.HandlerFunc(userHandler.GetUser)))

// Role-based route (auth + role middleware)
mux.Handle("/categories/create", 
    authMiddleware.Authenticate(
        roleMiddleware.RequireRole("admin")(
            http.HandlerFunc(categoryHandler.CreateCategory))))

// Multi-role route (auth + any role middleware)
mux.Handle("/summaries/create", 
    authMiddleware.Authenticate(
        roleMiddleware.RequireAnyRole("admin", "editor")(
            http.HandlerFunc(summaryHandler.CreateSummary))))
```

### Middleware Order

1. **AuthMiddleware** - Validates authentication token, loads user into context
2. **RoleMiddleware** - Checks user role/permissions against requirements
3. **Handler** - Processes the actual request

---

## 📊 Permission Matrix by Role

### Admin Role
- ✅ All permissions (66 total)
- ✅ Full CRUD on all resources
- ✅ User management
- ✅ Role and permission management

### Editor Role
- ✅ Content creation (ebooks, summaries, articles, etc.)
- ✅ Content updates
- ✅ Content reading
- ❌ Content deletion (Admin only)
- ❌ User/role/permission management

### Reader Role
- ✅ Read public content
- ✅ List resources
- ❌ Create/update/delete operations
- ❌ Premium content access

### Premium Role
- ✅ All Reader permissions
- ✅ Access to premium content
- ✅ View own payment history
- ❌ Create/update/delete operations

---

## 🚀 Adding New Protected Routes

### Step 1: Define the route with appropriate middleware

```go
// Admin only route
mux.Handle("/resource/create", 
    r.authMiddleware.Authenticate(
        r.roleMiddleware.RequireRole("admin")(
            http.HandlerFunc(r.resourceHandler.Create))))

// Editor or Admin route
mux.Handle("/resource/edit/{id}", 
    r.authMiddleware.Authenticate(
        r.roleMiddleware.RequireAnyRole("admin", "editor")(
            http.HandlerFunc(r.resourceHandler.Update))))

// Permission-based route
mux.Handle("/resource/publish/{id}", 
    r.authMiddleware.Authenticate(
        r.roleMiddleware.RequirePermission("resource:publish")(
            http.HandlerFunc(r.resourceHandler.Publish))))
```

### Step 2: Ensure permissions exist in seed data

Add to `seeder/000019_seed_permissions.sql`:
```sql
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'resource:create', 'resource', 'create', 'Create new resources');
```

### Step 3: Assign permissions to roles

Add to `seeder/000020_seed_role_permissions.sql`:
```sql
-- Add to admin role (already gets all permissions)
-- Add to editor role if needed
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'editor'
AND p.name IN (
    'resource:create',
    'resource:update'
);
```

---

## 🛡️ Security Best Practices

### 1. Always Use Middleware Chain
- Never rely on handler-level checks alone
- Use middleware for consistent security enforcement

### 2. Principle of Least Privilege
- Grant minimum required permissions
- Default to Reader role for new users

### 3. Separate Concerns
- Authentication: Verify user identity
- Authorization: Check permissions/roles
- Business Logic: Implement in handlers

### 4. Audit Logging
- Log all permission denied events
- Track role changes
- Monitor suspicious access patterns

### 5. Defense in Depth
- Middleware protection
- Handler-level validation
- Database constraints
- Input sanitization

---

## 📝 Testing Routes

### Test Public Routes
```bash
curl http://localhost:8080/ebooks
# Should return 200 OK
```

### Test Authenticated Routes
```bash
curl -H "Authorization: Bearer <token>" http://localhost:8080/users
# Should return 200 OK with valid token
# Should return 401 Unauthorized without token
```

### Test Role-Protected Routes
```bash
# As reader (should fail)
curl -H "Authorization: Bearer <reader-token>" \
     -X POST http://localhost:8080/categories/create
# Should return 403 Forbidden

# As admin (should succeed)
curl -H "Authorization: Bearer <admin-token>" \
     -X POST http://localhost:8080/categories/create \
     -d '{"name":"New Category"}'
# Should return 201 Created
```

---

## 🔧 Troubleshooting

### 401 Unauthorized
- Check if authentication token is valid
- Verify token is included in `Authorization` header
- Ensure user exists in database

### 403 Forbidden
- Verify user has required role
- Check if permissions are seeded correctly
- Confirm role-permission assignments in database

### 404 Not Found
- Verify route path is correct
- Check if route is registered in router
- Ensure method (GET/POST/PUT/DELETE) matches

---

## 📚 Related Files

- **Middleware**: `internal/delivery/http/middleware/role.go`
- **Permission Checker**: `internal/helper/permission_checker.go`
- **Role Service**: `internal/service/role_service_impl.go`
- **Permission Service**: `internal/service/permission_service_impl.go`
- **Seed Files**: `seeder/000018_seed_roles.sql`, `seeder/000019_seed_permissions.sql`

---

**Last Updated:** October 21, 2025  
**Version:** 1.0
