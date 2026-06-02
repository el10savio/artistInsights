CREATE TABLE IF NOT EXISTS artist_insights.listens (
    artist_id   String,
    track_id    String,
    track_name  String,
    user_id     String,
    ts          DateTime CODEC(Delta, ZSTD(1)),
    artist_name String
) ENGINE = MergeTree()
ORDER BY (artist_id, ts);
