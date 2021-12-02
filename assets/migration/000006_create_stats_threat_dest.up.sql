CREATE MATERIALIZED VIEW stats_threat_by_attacker
            WITH (timescaledb.continuous)
AS
SELECT
    time_bucket('1h', seen_at) as _bucket,
    attacker_host,
    COUNT(*) count
FROM threats
GROUP BY time_bucket('1h', seen_at), attacker_host;
