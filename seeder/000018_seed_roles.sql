-- Seed default roles for RBAC
-- This should be run after the roles table migration (000018_create_roles_table)

INSERT IGNORE INTO `roles` (`id`, `name`, `description`) VALUES
(UUID(), 'admin', 'Administrator with full system access and management capabilities'),
(UUID(), 'editor', 'Content editor who can create, update, and publish content'),
(UUID(), 'reader', 'Basic user who can read and access free content'),
(UUID(), 'premium', 'Premium user with access to exclusive content and features');
