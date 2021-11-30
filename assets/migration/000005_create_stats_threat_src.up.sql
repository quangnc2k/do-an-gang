CREATE MATERIALIZED VIEW stats_threat_by_src
            WITH (timescaledb.continuous)
AS
SELECT
    time_bucket('1h', seen_at) as _bucket,
    src_host,
    COUNT(*) count
FROM threats
GROUP BY time_bucket('1h', seen_at), src_host;

