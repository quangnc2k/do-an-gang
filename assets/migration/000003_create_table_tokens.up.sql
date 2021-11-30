CREATE TABLE tokens
(
    id         varchar(40) PRIMARY KEY NOT NULL,
    user_id    varchar(40)             NOT NULL,
    expired_at timestamptz             NOT NULL,
    revoked    boolean,
    user_agent text,
    ip_address varchar(100),
    created_at timestamptz             NOT NULL,
    CONSTRAINT fk_token_user FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX ON tokens (ip_address);