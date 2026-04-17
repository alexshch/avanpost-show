CREATE TABLE IF NOT EXISTS users
(
    id              uuid PRIMARY KEY,
    username        varchar(256) NOT NULL,
    email           varchar(256) NOT NULL,
    is_active       bool         NOT NULL DEFAULT TRUE,
    locked_at       timestamptz  NULL,
    created_at      timestamptz  NOT NULL,
    updated_at      timestamptz  NOT NULL,
    firstname       varchar(100) NOT NULL,
    lastname        varchar(100) NOT NULL,
    middlename      varchar(100) NULL
);
