ALTER TABLE `users` 
ADD COLUMN `role_id` CHAR(36) DEFAULT NULL AFTER `email`,
ADD CONSTRAINT `fk_users_role_id` FOREIGN KEY (`role_id`) REFERENCES `roles`(`id`) ON DELETE SET NULL;

CREATE INDEX `idx_users_role_id` ON `users`(`role_id`);
