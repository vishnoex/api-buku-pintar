ALTER TABLE `users` 
DROP FOREIGN KEY `fk_users_role_id`,
DROP INDEX `idx_users_role_id`,
DROP COLUMN `role_id`;
