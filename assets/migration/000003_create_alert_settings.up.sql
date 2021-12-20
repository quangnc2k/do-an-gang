CREATE TABLE alert_configs
(
    id           VARCHAR(50)  NOT NULL,
    name         VARCHAR(100) NOT NULL,
    created_at   TIMESTAMPTZ  NOT NULL,
    severity     VARCHAR(10)  NOT NULL,
    confidence   NUMERIC      NOT NULL,
    recipients    VARCHAR(320)[] NOT NULL,
    suppress_for INTERVAL     NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS alert_configs_unique_id ON alert_configs (id);

CREATE TRIGGER alert_configs_changes
    AFTER INSERT OR UPDATE OR DELETE ON alert_configs
    FOR EACH ROW EXECUTE PROCEDURE onChanges('alert_configs');