CREATE TABLE IF NOT EXISTS artist_insights.users (
    user_id String,
    country String
) ENGINE = ReplacingMergeTree()
ORDER BY (user_id);
