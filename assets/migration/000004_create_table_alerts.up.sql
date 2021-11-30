CREATE TABLE alerts (
    id VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    details JSON NOT NULL,
    resolved BOOLEAN NOT NULL,
    resolved_at TIMESTAMPTZ,
    resolved_by VARCHAR(40),

    CONSTRAINT fk_alerts_users FOREIGN KEY (resolved_by) REFERENCES users(id) ON DELETE SET NULL
);

SELECT create_hypertable('alerts', 'created_at');

CREATE INDEX ON alerts (created_at DESC, id);