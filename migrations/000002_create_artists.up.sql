CREATE TABLE IF NOT EXISTS artist_insights.artists (
    artist_id   String,
    artist_name String
) ENGINE = ReplacingMergeTree()
ORDER BY (artist_id);
