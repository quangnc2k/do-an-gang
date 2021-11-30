CREATE MATERIALIZED VIEW stats_threat_by_severity
            WITH (timescaledb.continuous)
AS
SELECT
    time_bucket('1h', seen_at) as _bucket,
    severity,
    COUNT(*) count
FROM threats
GROUP BY time_bucket('1h', seen_at), severity;