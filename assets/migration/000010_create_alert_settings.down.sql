CREATE TABLE alert_configs (
    id VARCHAR(50) NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    severity VARCHAR(10) NOT NULL,
    confidence NUMERIC NOT NULL,
    recipient VARCHAR(320)[] NOT NULL,
)