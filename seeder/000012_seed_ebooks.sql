-- Seed data for ebooks table

-- Make sure authors and categories exist first
INSERT IGNORE INTO `authors` (`id`, `name`, `avatar`) VALUES
('a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'John Doe', 'https://example.com/avatars/johndoe.jpg'),
('b2c3d4e5-f6a7-8901-bcde-f23456789012', 'Jane Smith', 'https://example.com/avatars/janesmith.jpg'),
('c3d4e5f6-a7b8-9012-cdef-345678901234', 'Robert Johnson', 'https://example.com/avatars/robertjohnson.jpg');

INSERT IGNORE INTO `categories` (`id`, `name`, `description`, `icon_link`, `is_active`, `order_number`) VALUES
('d4e5f6a7-b8c9-0123-defg-456789012345', 'Fiction', 'Fiction books and novels', 'https://example.com/icons/fiction.png', TRUE, 1),
('e5f6a7b8-c9d0-1234-efgh-567890123456', 'Science', 'Science and educational books', 'https://example.com/icons/science.png', TRUE, 2),
('f6a7b8c9-d0e1-2345-fghi-678901234567', 'Business', 'Business and entrepreneurship books', 'https://example.com/icons/business.png', TRUE, 3);

INSERT IGNORE INTO `content_statuses` (`id`, `name`) VALUES
('a7b8c9d0-e1f2-3456-ghij-789012345678', 'draft'),
('b8c9d0e1-f2a3-4567-hijk-890123456789', 'published'),
('c9d0e1f2-a3b4-5678-ijkl-901234567890', 'archived');

-- Insert ebooks data
INSERT INTO `ebooks` (
  `id`, 
  `author_id`, 
  `title`, 
  `synopsis`, 
  `slug`, 
  `cover_image`, 
  `category_id`, 
  `content_status_id`, 
  `price`, 
  `language`, 
  `duration`, 
  `filesize`, 
  `format`, 
  `page_count`, 
  `preview_page`, 
  `url`, 
  `published_at`
) VALUES 
(
  'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 
  'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 
  'The Art of Programming', 
  'A comprehensive guide to programming fundamentals and best practices.', 
  'the-art-of-programming', 
  'https://example.com/covers/art-of-programming.jpg', 
  'e5f6a7b8-c9d0-1234-efgh-567890123456', 
  'b8c9d0e1-f2a3-4567-hijk-890123456789', 
  25000, 
  'English', 
  180, 
  15728640, 
  'pdf', 
  350, 
  20, 
  'https://example.com/ebooks/the-art-of-programming.pdf', 
  '2023-01-15 10:00:00'
),
(
  'f2a3b4c5-d6e7-8901-bcde-f23456789012', 
  'b2c3d4e5-f6a7-8901-bcde-f23456789012', 
  'Business Strategy 101', 
  'Learn the fundamentals of business strategy and how to apply them to your organization.', 
  'business-strategy-101', 
  'https://example.com/covers/business-strategy.jpg', 
  'f6a7b8c9-d0e1-2345-fghi-678901234567', 
  'b8c9d0e1-f2a3-4567-hijk-890123456789', 
  30000, 
  'English', 
  210, 
  20971520, 
  'epub', 
  420, 
  25, 
  'https://example.com/ebooks/business-strategy-101.epub', 
  '2023-02-20 14:30:00'
),
(
  'a3b4c5d6-e7f8-9012-cdef-345678901234', 
  'c3d4e5f6-a7b8-9012-cdef-345678901234', 
  'The Mystery of the Lost Key', 
  'A thrilling mystery novel that will keep you on the edge of your seat.', 
  'mystery-lost-key', 
  'https://example.com/covers/mystery-lost-key.jpg', 
  'd4e5f6a7-b8c9-0123-defg-456789012345', 
  'b8c9d0e1-f2a3-4567-hijk-890123456789', 
  15000, 
  'English', 
  150, 
  10485760, 
  'pdf', 
  280, 
  15, 
  'https://example.com/ebooks/mystery-lost-key.pdf', 
  '2023-03-10 09:15:00'
),
(
  'b4c5d6e7-f8a9-0123-defg-456789012345', 
  'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 
  'Advanced Data Science', 
  'Explore advanced concepts in data science, machine learning, and artificial intelligence.', 
  'advanced-data-science', 
  'https://example.com/covers/advanced-data-science.jpg', 
  'e5f6a7b8-c9d0-1234-efgh-567890123456', 
  'b8c9d0e1-f2a3-4567-hijk-890123456789', 
  35000, 
  'English', 
  240, 
  25165824, 
  'pdf', 
  500, 
  30, 
  'https://example.com/ebooks/advanced-data-science.pdf', 
  '2023-04-05 11:45:00'
),
(
  'c5d6e7f8-a9b0-1234-efgh-567890123456', 
  'b2c3d4e5-f6a7-8901-bcde-f23456789012', 
  'Marketing in the Digital Age', 
  'A guide to modern marketing strategies and techniques in the digital era.', 
  'marketing-digital-age', 
  'https://example.com/covers/marketing-digital.jpg', 
  'f6a7b8c9-d0e1-2345-fghi-678901234567', 
  'a7b8c9d0-e1f2-3456-ghij-789012345678', 
  28000, 
  'English', 
  190, 
  18874368, 
  'mobi', 
  380, 
  22, 
  'https://example.com/ebooks/marketing-digital-age.mobi', 
  NULL
);

-- Add some table of contents entries for the ebooks
INSERT INTO `table_of_contents` (`id`, `ebook_id`, `title`, `page_number`) VALUES
(UUID(), 'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 'Introduction to Programming', 1),
(UUID(), 'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 'Basic Syntax', 15),
(UUID(), 'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 'Data Structures', 45),
(UUID(), 'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 'Algorithms', 120),
(UUID(), 'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 'Best Practices', 250),

(UUID(), 'f2a3b4c5-d6e7-8901-bcde-f23456789012', 'What is Business Strategy?', 1),
(UUID(), 'f2a3b4c5-d6e7-8901-bcde-f23456789012', 'Market Analysis', 30),
(UUID(), 'f2a3b4c5-d6e7-8901-bcde-f23456789012', 'Competitive Advantage', 85),
(UUID(), 'f2a3b4c5-d6e7-8901-bcde-f23456789012', 'Implementation', 200),
(UUID(), 'f2a3b4c5-d6e7-8901-bcde-f23456789012', 'Case Studies', 300);

-- Add some ebook discounts
INSERT INTO `ebook_discounts` (`id`, `ebook_id`, `discount_price`, `started_at`, `ended_at`) VALUES
(UUID(), 'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 20000, '2023-05-01 00:00:00', '2023-06-01 23:59:59'),
(UUID(), 'f2a3b4c5-d6e7-8901-bcde-f23456789012', 25000, '2023-05-15 00:00:00', '2023-06-15 23:59:59');

-- Add some ebook summaries
INSERT INTO `ebook_summaries` (`id`, `ebook_id`, `description`, `url`, `audio_url`) VALUES
(UUID(), 'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 'A brief overview of programming concepts', 'https://example.com/summaries/art-of-programming.pdf', 'https://example.com/audio/art-of-programming.mp3'),
(UUID(), 'f2a3b4c5-d6e7-8901-bcde-f23456789012', 'Key business strategy concepts summarized', 'https://example.com/summaries/business-strategy.pdf', 'https://example.com/audio/business-strategy.mp3'),
(UUID(), 'a3b4c5d6-e7f8-9012-cdef-345678901234', 'Plot summary and character guide', 'https://example.com/summaries/mystery-lost-key.pdf', 'https://example.com/audio/mystery-lost-key.mp3');

-- Add some premium summaries
INSERT INTO `ebook_premium_summaries` (`id`, `ebook_id`, `description`, `url`, `audio_url`) VALUES
(UUID(), 'e1f2a3b4-c5d6-7890-abcd-ef1234567890', 'Detailed programming guide with examples', 'https://example.com/premium/art-of-programming.pdf', 'https://example.com/premium-audio/art-of-programming.mp3'),
(UUID(), 'f2a3b4c5-d6e7-8901-bcde-f23456789012', 'In-depth business strategy analysis', 'https://example.com/premium/business-strategy.pdf', 'https://example.com/premium-audio/business-strategy.mp3');