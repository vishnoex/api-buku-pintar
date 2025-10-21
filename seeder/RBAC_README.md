# RBAC Seeder

This directory contains SQL seeders for Role-Based Access Control (RBAC) initialization.

## Seed Files

### 000018_seed_roles.sql
Seeds the four default roles in the system:
- **Admin**: Full system access
- **Editor**: Content creation and management
- **Reader**: Basic read access to free content
- **Premium**: Extended access including premium content

**Note**: Uses `UUID()` to generate unique IDs automatically.

### 000019_seed_permissions.sql
Seeds all permissions (66+ permissions) across 11 resource types:
- User permissions (6)
- Role permissions (6)
- Permission permissions (6)
- Category permissions (6)
- Banner permissions (6)
- Ebook permissions (6)
- Summary permissions (6)
- Article permissions (6)
- Inspiration permissions (6)
- Author permissions (6)
- Payment permissions (6)

**Note**: Uses `UUID()` to generate unique IDs automatically.

### 000020_seed_role_permissions.sql
Maps permissions to roles based on access levels using **subquery lookups**:
- Queries roles and permissions by name
- Uses `CROSS JOIN` with `WHERE` clauses
- No hardcoded IDs - references are dynamic

#### Admin Role
- **Access**: Full system access (all 66 permissions)
- **Can**: Manage users, roles, permissions, and all content
- **Use Case**: System administrators

#### Editor Role
- **Access**: Content management (27 permissions)
- **Can**: Create/update articles, ebooks, summaries, inspirations, authors
- **Cannot**: Delete content, manage users/roles/permissions, manage payments
- **Use Case**: Content creators and editors

#### Reader Role
- **Access**: Basic read access (14 permissions)
- **Can**: Read and list free content, categories, banners, authors
- **Cannot**: Create, update, or delete anything
- **Use Case**: Free users browsing the platform

#### Premium Role
- **Access**: Extended read access (15 permissions)
- **Can**: Read and list all content (including premium), view own payments
- **Cannot**: Create, update, or delete content
- **Use Case**: Paid subscribers with premium content access

## Running Seeders

### Prerequisites
1. Run migrations first to create the tables:
   ```bash
   migrate -path migrations -database "mysql://user:pass@tcp(host:port)/dbname" up
   ```

2. Ensure these migrations are applied:
   - `000018_create_roles_table`
   - `000019_create_permissions_table`
   - `000020_create_role_permissions_table`

### Execute Seeders

Run the seeders in order:

```bash
# 1. Seed roles
mysql -u username -p database_name < seeder/000018_seed_roles.sql

# 2. Seed permissions
mysql -u username -p database_name < seeder/000019_seed_permissions.sql

# 3. Seed role-permission assignments
mysql -u username -p database_name < seeder/000020_seed_role_permissions.sql
```

Or run all at once:
```bash
cat seeder/000018_seed_roles.sql \
    seeder/000019_seed_permissions.sql \
    seeder/000020_seed_role_permissions.sql | \
    mysql -u username -p database_name
```

### Verification

Check that roles were created:
```sql
SELECT * FROM roles;
```

Check that permissions were created:
```sql
SELECT COUNT(*) as total_permissions FROM permissions;
-- Should return 66
```

Check role-permission assignments:
```sql
SELECT r.name, COUNT(rp.permission_id) as permission_count
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
GROUP BY r.id, r.name;
```

Expected output:
```
+----------+------------------+
| name     | permission_count |
+----------+------------------+
| admin    | 66               |
| editor   | 27               |
| reader   | 14               |
| premium  | 15               |
+----------+------------------+
```

## Assigning Roles to Users

After seeding, assign roles to users by role name:

```sql
-- Assign admin role to a user
UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'admin') 
WHERE id = 'user-id-here';

-- Assign editor role to a user
UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'editor') 
WHERE id = 'user-id-here';

-- Assign reader role to a user (default for new users)
UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'reader') 
WHERE id = 'user-id-here';

-- Assign premium role to a user (after payment)
UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'premium') 
WHERE id = 'user-id-here';
```

## Permission Format

All permissions follow the `resource:action` naming pattern:
- Format: `{resource}:{action}`
- Example: `ebook:create`, `article:read`, `user:manage`

### Actions
- `create`: Create new resources
- `read`: Read resource information
- `update`: Update resource information
- `delete`: Delete resources
- `list`: List all resources
- `manage`: Full management access (includes all actions)

### Resources
- user, role, permission, category, banner
- ebook, summary, article, inspiration, author, payment

## Customization

To add new permissions:

1. Add to `000019_seed_permissions.sql`:
   ```sql
   INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
   (UUID(), 'resource:action', 'resource', 'action', 'Description');
   ```

2. Assign to appropriate roles in `000020_seed_role_permissions.sql`:
   ```sql
   -- Add to the IN clause for the desired role
   INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
   SELECT r.id, p.id
   FROM roles r
   CROSS JOIN permissions p
   WHERE r.name = 'admin'
   AND p.name IN (
       'resource:action',
       -- ... other permissions
   );
   ```

3. Update the permission constants in `internal/domain/entity/permission.go`

## Notes

- Uses `INSERT IGNORE` to prevent duplicate entries on re-run
- Uses `UUID()` to generate unique IDs automatically
- Role-permission assignments use subqueries to lookup IDs by name
- No hardcoded IDs - all references are dynamic and portable
- Admin role gets all permissions (66 total)
- Reader and Premium roles differ only in premium content access
- Editor role cannot delete content (safety measure)
- Payment permissions restricted to admin and premium users

## Related Files

- Entities: `internal/domain/entity/role.go`, `internal/domain/entity/permission.go`
- Repositories: `internal/repository/mysql/role_repository.go`, `internal/repository/mysql/permission_repository.go`
- Services: `internal/service/role_service_impl.go`, `internal/service/permission_service_impl.go`
- Middleware: `internal/delivery/http/middleware/role.go`
- Helper: `internal/helper/permission_checker.go`
