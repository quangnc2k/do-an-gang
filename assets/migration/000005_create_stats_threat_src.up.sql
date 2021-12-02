CREATE MATERIALIZED VIEW stats_threat_by_affected
            WITH (timescaledb.continuous)
AS
SELECT
    time_bucket('1h', seen_at) as _bucket,
    affected_host,
    COUNT(*) count
FROM threats
GROUP BY time_bucket('1h', seen_at), affected_host;

