CREATE MATERIALIZED VIEW stats_threat_by_phase
            WITH (timescaledb.continuous)
AS
SELECT
    time_bucket('1h', seen_at) as _bucket,
    phase,
    COUNT(*) count
FROM threats
GROUP BY time_bucket('1h', seen_at), phase;