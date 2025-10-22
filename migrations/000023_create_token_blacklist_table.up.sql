CREATE TABLE `token_blacklist` (
    `id` CHAR(36) PRIMARY KEY,
    `token_hash` VARCHAR(64) UNIQUE NOT NULL,
    `user_id` CHAR(36),
    `blacklisted_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `expires_at` TIMESTAMP NOT NULL,
    `reason` VARCHAR(100),
    INDEX `idx_token_hash` (`token_hash`),
    INDEX `idx_expires_at` (`expires_at`)
);
