INSERT IGNORE INTO `categories` (`id`, `name`, `slug`, `description`, `icon_link`, `is_active`, `order_number`) VALUES
(UUID(), 'Semua', 'all', 'Semua', 'star', TRUE, 1),
(UUID(), 'Ringkasan', 'summary', 'Ringkasan Buku', 'book-text', TRUE, 2),
(UUID(), 'Inspirasi', 'inspiration', 'Inspirasi harian', 'lightbulb', TRUE, 3),
(UUID(), 'E-Book', 'ebook', 'Buku Digital', 'book-open-check', TRUE, 4),
(UUID(), 'Artikel', 'article', 'Artikel', 'file-text', TRUE, 5);