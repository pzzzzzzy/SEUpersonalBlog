-- 为articles表添加like_count列
ALTER TABLE articles ADD COLUMN like_count INTEGER DEFAULT 0;

-- 为现有文章设置初始点赞数
UPDATE articles SET like_count = 0;
