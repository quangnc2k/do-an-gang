CREATE MATERIALIZED VIEW stats_threat_by_dst
            WITH (timescaledb.continuous)
AS
SELECT
    time_bucket('1h', seen_at) as _bucket,
    dst_host,
    COUNT(*) count
FROM threats
GROUP BY time_bucket('1h', seen_at), dst_host;
