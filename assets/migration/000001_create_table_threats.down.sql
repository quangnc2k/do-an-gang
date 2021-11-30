DROP INDEX IF EXISTS threats_id_idx;
DROP INDEX IF EXISTS threats_conn_id_idx;
DROP INDEX IF EXISTS threats_src_host_idx;
DROP INDEX IF EXISTS threats_dst_host_idx;
DROP INDEX IF EXISTS threats_sensor_id_idx;
DROP INDEX IF EXISTS threats_created_at_idx;
DROP INDEX IF EXISTS threats_resource_idx;
DROP INDEX IF EXISTS threats_tactics_gin_idx;
DROP INDEX IF EXISTS threats_techniques_gin_idx;
DROP INDEX IF EXISTS threats_origin_idx;
DROP INDEX IF EXISTS threats_channel_idx;

DROP INDEX IF EXISTS threats_resource_gin_idx;
DROP INDEX IF EXISTS threats_tactics_gin_idx;
DROP INDEX IF EXISTS threats_techniques_gin_idx;

DROP EXTENSION IF EXISTS btree_gin;
DROP EXTENSION IF EXISTS pg_trgm;

DROP TABLE IF EXISTS threats;