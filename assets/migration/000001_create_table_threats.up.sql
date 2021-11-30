CREATE TABLE threats
(
    id         VARCHAR(50),
    created_at TIMESTAMPTZ NOT NULL,
    seen_at    TIMESTAMPTZ,
    src_host   VARCHAR(25),
    dst_host   VARCHAR(25),
    conn_id    VARCHAR(30),
    confidence NUMERIC,
    severity   VARCHAR(10),
    phase      VARCHAR(100),
    metadata   JSON
);
SELECT create_hypertable('threats', 'seen_at');

CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE EXTENSION IF NOT EXISTS pg_trgm;

CREATE UNIQUE INDEX IF NOT EXISTS threats_id_unique_idx ON threats (seen_at, id);

CREATE INDEX IF NOT EXISTS threats_id_idx ON threats (id);
CREATE INDEX IF NOT EXISTS threats_conn_id_idx ON threats (conn_id);
CREATE INDEX IF NOT EXISTS threats_src_host_idx ON threats (src_host);
CREATE INDEX IF NOT EXISTS threats_dst_host_idx ON threats (dst_host);
CREATE INDEX IF NOT EXISTS threats_created_at_idx ON threats (created_at);
CREATE INDEX IF NOT EXISTS threats_resource_idx ON threats (phase);
