CREATE TABLE users
(
    id         varchar(40) PRIMARY KEY NOT NULL,
    email      varchar(100) UNIQUE     NOT NULL,
    password   varchar(100)            NOT NULL,
    created_at TIMESTAMPTZ             NOT NULL,
    updated_at TIMESTAMPTZ             NOT NULL
);

CREATE INDEX ON users (email);