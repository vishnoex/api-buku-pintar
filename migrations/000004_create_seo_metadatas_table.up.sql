CREATE TABLE IF NOT EXISTS `seo_metadatas` (
  `id` VARCHAR(36) PRIMARY KEY,
  `title` VARCHAR(255) NOT NULL,
  `description` TEXT,
  `keywords` TEXT,
  `entity` ENUM('banner', 'category', 'ebook', 'article', 'inspiration') NOT NULL DEFAULT 'banner',
  `entity_id` VARCHAR(36) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_entity_id` (`entity_id`),
  INDEX `idx_entity_entity_id` (`entity`, `entity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
