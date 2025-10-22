CREATE TABLE `oauth_tokens` (
    `id` CHAR(36) PRIMARY KEY,
    `user_id` CHAR(36) NOT NULL,
    `provider` VARCHAR(20) NOT NULL,
    `access_token` TEXT NOT NULL,
    `refresh_token` TEXT,
    `token_type` VARCHAR(20),
    `expires_at` TIMESTAMP,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (`user_id`) REFERENCES `users`(`id`) ON DELETE CASCADE,
    INDEX `idx_user_provider` (`user_id`, `provider`)
);
