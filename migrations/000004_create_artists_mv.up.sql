CREATE MATERIALIZED VIEW IF NOT EXISTS artist_insights.artists_mv
TO artist_insights.artists
AS
SELECT artist_id, artist_name
FROM artist_insights.listens;
