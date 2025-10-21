-- Seed role-permission assignments for RBAC
-- This should be run after both roles and permissions are seeded
-- Uses subqueries to lookup role and permission IDs by name

-- ============================================================================
-- ADMIN ROLE - Full system access
-- ============================================================================
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'admin';

-- ============================================================================
-- EDITOR ROLE - Content creation and management
-- ============================================================================
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'editor'
AND p.name IN (
    -- User permissions (read only)
    'user:read',
    
    -- Category permissions (read and list)
    'category:read',
    'category:list',
    
    -- Banner permissions (read and list)
    'banner:read',
    'banner:list',
    
    -- Ebook permissions (create, read, update, list)
    'ebook:create',
    'ebook:read',
    'ebook:update',
    'ebook:list',
    
    -- Summary permissions (create, read, update, list)
    'summary:create',
    'summary:read',
    'summary:update',
    'summary:list',
    
    -- Article permissions (create, read, update, list)
    'article:create',
    'article:read',
    'article:update',
    'article:list',
    
    -- Inspiration permissions (create, read, update, list)
    'inspiration:create',
    'inspiration:read',
    'inspiration:update',
    'inspiration:list',
    
    -- Author permissions (create, read, update, list)
    'author:create',
    'author:read',
    'author:update',
    'author:list'
);

-- ============================================================================
-- READER ROLE - Basic read access to free content
-- ============================================================================
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'reader'
AND p.name IN (
    -- Category permissions (read and list)
    'category:read',
    'category:list',
    
    -- Banner permissions (read and list)
    'banner:read',
    'banner:list',
    
    -- Ebook permissions (read and list - free content only)
    'ebook:read',
    'ebook:list',
    
    -- Summary permissions (read and list - free content only)
    'summary:read',
    'summary:list',
    
    -- Article permissions (read and list)
    'article:read',
    'article:list',
    
    -- Inspiration permissions (read and list)
    'inspiration:read',
    'inspiration:list',
    
    -- Author permissions (read and list)
    'author:read',
    'author:list'
);

-- ============================================================================
-- PREMIUM ROLE - Extended read access including premium content
-- ============================================================================
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`)
SELECT r.id, p.id
FROM roles r
CROSS JOIN permissions p
WHERE r.name = 'premium'
AND p.name IN (
    -- Category permissions (read and list)
    'category:read',
    'category:list',
    
    -- Banner permissions (read and list)
    'banner:read',
    'banner:list',
    
    -- Ebook permissions (read and list - including premium)
    'ebook:read',
    'ebook:list',
    
    -- Summary permissions (read and list - including premium)
    'summary:read',
    'summary:list',
    
    -- Article permissions (read and list - including premium)
    'article:read',
    'article:list',
    
    -- Inspiration permissions (read and list)
    'inspiration:read',
    'inspiration:list',
    
    -- Author permissions (read and list)
    'author:read',
    'author:list',
    
    -- Payment permissions (read own payments)
    'payment:read'
);
-- User permissions (all)
('role-admin-001', 'perm-user-create'),
('role-admin-001', 'perm-user-read'),
('role-admin-001', 'perm-user-update'),
('role-admin-001', 'perm-user-delete'),
('role-admin-001', 'perm-user-list'),
('role-admin-001', 'perm-user-manage'),

-- Role permissions (all)
('role-admin-001', 'perm-role-create'),
('role-admin-001', 'perm-role-read'),
('role-admin-001', 'perm-role-update'),
('role-admin-001', 'perm-role-delete'),
('role-admin-001', 'perm-role-list'),
('role-admin-001', 'perm-role-manage'),

-- Permission permissions (all)
('role-admin-001', 'perm-permission-create'),
('role-admin-001', 'perm-permission-read'),
('role-admin-001', 'perm-permission-update'),
('role-admin-001', 'perm-permission-delete'),
('role-admin-001', 'perm-permission-list'),
('role-admin-001', 'perm-permission-manage'),

-- Category permissions (all)
('role-admin-001', 'perm-category-create'),
('role-admin-001', 'perm-category-read'),
('role-admin-001', 'perm-category-update'),
('role-admin-001', 'perm-category-delete'),
('role-admin-001', 'perm-category-list'),
('role-admin-001', 'perm-category-manage'),

-- Banner permissions (all)
('role-admin-001', 'perm-banner-create'),
('role-admin-001', 'perm-banner-read'),
('role-admin-001', 'perm-banner-update'),
('role-admin-001', 'perm-banner-delete'),
('role-admin-001', 'perm-banner-list'),
('role-admin-001', 'perm-banner-manage'),

-- Ebook permissions (all)
('role-admin-001', 'perm-ebook-create'),
('role-admin-001', 'perm-ebook-read'),
('role-admin-001', 'perm-ebook-update'),
('role-admin-001', 'perm-ebook-delete'),
('role-admin-001', 'perm-ebook-list'),
('role-admin-001', 'perm-ebook-manage'),

-- Summary permissions (all)
('role-admin-001', 'perm-summary-create'),
('role-admin-001', 'perm-summary-read'),
('role-admin-001', 'perm-summary-update'),
('role-admin-001', 'perm-summary-delete'),
('role-admin-001', 'perm-summary-list'),
('role-admin-001', 'perm-summary-manage'),

-- Article permissions (all)
('role-admin-001', 'perm-article-create'),
('role-admin-001', 'perm-article-read'),
('role-admin-001', 'perm-article-update'),
('role-admin-001', 'perm-article-delete'),
('role-admin-001', 'perm-article-list'),
('role-admin-001', 'perm-article-manage'),

-- Inspiration permissions (all)
('role-admin-001', 'perm-inspiration-create'),
('role-admin-001', 'perm-inspiration-read'),
('role-admin-001', 'perm-inspiration-update'),
('role-admin-001', 'perm-inspiration-delete'),
('role-admin-001', 'perm-inspiration-list'),
('role-admin-001', 'perm-inspiration-manage'),

-- Author permissions (all)
('role-admin-001', 'perm-author-create'),
('role-admin-001', 'perm-author-read'),
('role-admin-001', 'perm-author-update'),
('role-admin-001', 'perm-author-delete'),
('role-admin-001', 'perm-author-list'),
('role-admin-001', 'perm-author-manage'),

-- Payment permissions (all)
('role-admin-001', 'perm-payment-create'),
('role-admin-001', 'perm-payment-read'),
('role-admin-001', 'perm-payment-update'),
('role-admin-001', 'perm-payment-delete'),
('role-admin-001', 'perm-payment-list'),
('role-admin-001', 'perm-payment-manage');

-- ============================================================================
-- EDITOR ROLE - Content creation and management
-- ============================================================================
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`) VALUES
-- User permissions (read only)
('role-editor-002', 'perm-user-read'),

-- Category permissions (read and list)
('role-editor-002', 'perm-category-read'),
('role-editor-002', 'perm-category-list'),

-- Banner permissions (read and list)
('role-editor-002', 'perm-banner-read'),
('role-editor-002', 'perm-banner-list'),

-- Ebook permissions (create, read, update, list)
('role-editor-002', 'perm-ebook-create'),
('role-editor-002', 'perm-ebook-read'),
('role-editor-002', 'perm-ebook-update'),
('role-editor-002', 'perm-ebook-list'),

-- Summary permissions (create, read, update, list)
('role-editor-002', 'perm-summary-create'),
('role-editor-002', 'perm-summary-read'),
('role-editor-002', 'perm-summary-update'),
('role-editor-002', 'perm-summary-list'),

-- Article permissions (create, read, update, list)
('role-editor-002', 'perm-article-create'),
('role-editor-002', 'perm-article-read'),
('role-editor-002', 'perm-article-update'),
('role-editor-002', 'perm-article-list'),

-- Inspiration permissions (create, read, update, list)
('role-editor-002', 'perm-inspiration-create'),
('role-editor-002', 'perm-inspiration-read'),
('role-editor-002', 'perm-inspiration-update'),
('role-editor-002', 'perm-inspiration-list'),

-- Author permissions (create, read, update, list)
('role-editor-002', 'perm-author-create'),
('role-editor-002', 'perm-author-read'),
('role-editor-002', 'perm-author-update'),
('role-editor-002', 'perm-author-list');

-- ============================================================================
-- READER ROLE - Basic read access to free content
-- ============================================================================
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`) VALUES
-- Category permissions (read and list)
('role-reader-003', 'perm-category-read'),
('role-reader-003', 'perm-category-list'),

-- Banner permissions (read and list)
('role-reader-003', 'perm-banner-read'),
('role-reader-003', 'perm-banner-list'),

-- Ebook permissions (read and list - free content only)
('role-reader-003', 'perm-ebook-read'),
('role-reader-003', 'perm-ebook-list'),

-- Summary permissions (read and list - free content only)
('role-reader-003', 'perm-summary-read'),
('role-reader-003', 'perm-summary-list'),

-- Article permissions (read and list)
('role-reader-003', 'perm-article-read'),
('role-reader-003', 'perm-article-list'),

-- Inspiration permissions (read and list)
('role-reader-003', 'perm-inspiration-read'),
('role-reader-003', 'perm-inspiration-list'),

-- Author permissions (read and list)
('role-reader-003', 'perm-author-read'),
('role-reader-003', 'perm-author-list');

-- ============================================================================
-- PREMIUM ROLE - Extended read access including premium content
-- ============================================================================
INSERT IGNORE INTO `role_permissions` (`role_id`, `permission_id`) VALUES
-- Category permissions (read and list)
('role-premium-004', 'perm-category-read'),
('role-premium-004', 'perm-category-list'),

-- Banner permissions (read and list)
('role-premium-004', 'perm-banner-read'),
('role-premium-004', 'perm-banner-list'),

-- Ebook permissions (read and list - including premium)
('role-premium-004', 'perm-ebook-read'),
('role-premium-004', 'perm-ebook-list'),

-- Summary permissions (read and list - including premium)
('role-premium-004', 'perm-summary-read'),
('role-premium-004', 'perm-summary-list'),

-- Article permissions (read and list - including premium)
('role-premium-004', 'perm-article-read'),
('role-premium-004', 'perm-article-list'),

-- Inspiration permissions (read and list)
('role-premium-004', 'perm-inspiration-read'),
('role-premium-004', 'perm-inspiration-list'),

-- Author permissions (read and list)
('role-premium-004', 'perm-author-read'),
('role-premium-004', 'perm-author-list'),

-- Payment permissions (read own payments)
('role-premium-004', 'perm-payment-read');
