-- Seed default permissions for RBAC
-- This should be run after the permissions table migration (000019_create_permissions_table)

-- User permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'user:create', 'user', 'create', 'Create new users'),
(UUID(), 'user:read', 'user', 'read', 'Read user information'),
(UUID(), 'user:update', 'user', 'update', 'Update user information'),
(UUID(), 'user:delete', 'user', 'delete', 'Delete users'),
(UUID(), 'user:list', 'user', 'list', 'List all users'),
(UUID(), 'user:manage', 'user', 'manage', 'Full user management access');

-- Role permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'role:create', 'role', 'create', 'Create new roles'),
(UUID(), 'role:read', 'role', 'read', 'Read role information'),
(UUID(), 'role:update', 'role', 'update', 'Update role information'),
(UUID(), 'role:delete', 'role', 'delete', 'Delete roles'),
(UUID(), 'role:list', 'role', 'list', 'List all roles'),
(UUID(), 'role:manage', 'role', 'manage', 'Full role management access');

-- Permission permissions (meta-permissions)
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'permission:create', 'permission', 'create', 'Create new permissions'),
(UUID(), 'permission:read', 'permission', 'read', 'Read permission information'),
(UUID(), 'permission:update', 'permission', 'update', 'Update permission information'),
(UUID(), 'permission:delete', 'permission', 'delete', 'Delete permissions'),
(UUID(), 'permission:list', 'permission', 'list', 'List all permissions'),
(UUID(), 'permission:manage', 'permission', 'manage', 'Full permission management access');

-- Category permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'category:create', 'category', 'create', 'Create new categories'),
(UUID(), 'category:read', 'category', 'read', 'Read category information'),
(UUID(), 'category:update', 'category', 'update', 'Update category information'),
(UUID(), 'category:delete', 'category', 'delete', 'Delete categories'),
(UUID(), 'category:list', 'category', 'list', 'List all categories'),
(UUID(), 'category:manage', 'category', 'manage', 'Full category management access');

-- Banner permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'banner:create', 'banner', 'create', 'Create new banners'),
(UUID(), 'banner:read', 'banner', 'read', 'Read banner information'),
(UUID(), 'banner:update', 'banner', 'update', 'Update banner information'),
(UUID(), 'banner:delete', 'banner', 'delete', 'Delete banners'),
(UUID(), 'banner:list', 'banner', 'list', 'List all banners'),
(UUID(), 'banner:manage', 'banner', 'manage', 'Full banner management access');

-- Ebook permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'ebook:create', 'ebook', 'create', 'Create new ebooks'),
(UUID(), 'ebook:read', 'ebook', 'read', 'Read ebook content'),
(UUID(), 'ebook:update', 'ebook', 'update', 'Update ebook information'),
(UUID(), 'ebook:delete', 'ebook', 'delete', 'Delete ebooks'),
(UUID(), 'ebook:list', 'ebook', 'list', 'List all ebooks'),
(UUID(), 'ebook:manage', 'ebook', 'manage', 'Full ebook management access');

-- Summary permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'summary:create', 'summary', 'create', 'Create new summaries'),
(UUID(), 'summary:read', 'summary', 'read', 'Read summary content'),
(UUID(), 'summary:update', 'summary', 'update', 'Update summary information'),
(UUID(), 'summary:delete', 'summary', 'delete', 'Delete summaries'),
(UUID(), 'summary:list', 'summary', 'list', 'List all summaries'),
(UUID(), 'summary:manage', 'summary', 'manage', 'Full summary management access');

-- Article permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'article:create', 'article', 'create', 'Create new articles'),
(UUID(), 'article:read', 'article', 'read', 'Read article content'),
(UUID(), 'article:update', 'article', 'update', 'Update article information'),
(UUID(), 'article:delete', 'article', 'delete', 'Delete articles'),
(UUID(), 'article:list', 'article', 'list', 'List all articles'),
(UUID(), 'article:manage', 'article', 'manage', 'Full article management access');

-- Inspiration permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'inspiration:create', 'inspiration', 'create', 'Create new inspirations'),
(UUID(), 'inspiration:read', 'inspiration', 'read', 'Read inspiration content'),
(UUID(), 'inspiration:update', 'inspiration', 'update', 'Update inspiration information'),
(UUID(), 'inspiration:delete', 'inspiration', 'delete', 'Delete inspirations'),
(UUID(), 'inspiration:list', 'inspiration', 'list', 'List all inspirations'),
(UUID(), 'inspiration:manage', 'inspiration', 'manage', 'Full inspiration management access');

-- Author permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'author:create', 'author', 'create', 'Create new authors'),
(UUID(), 'author:read', 'author', 'read', 'Read author information'),
(UUID(), 'author:update', 'author', 'update', 'Update author information'),
(UUID(), 'author:delete', 'author', 'delete', 'Delete authors'),
(UUID(), 'author:list', 'author', 'list', 'List all authors'),
(UUID(), 'author:manage', 'author', 'manage', 'Full author management access');

-- Payment permissions
INSERT IGNORE INTO `permissions` (`id`, `name`, `resource`, `action`, `description`) VALUES
(UUID(), 'payment:create', 'payment', 'create', 'Create new payments'),
(UUID(), 'payment:read', 'payment', 'read', 'Read payment information'),
(UUID(), 'payment:update', 'payment', 'update', 'Update payment information'),
(UUID(), 'payment:delete', 'payment', 'delete', 'Delete payments'),
(UUID(), 'payment:list', 'payment', 'list', 'List all payments'),
(UUID(), 'payment:manage', 'payment', 'manage', 'Full payment management access');
