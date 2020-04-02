-- Main

DROP TABLE IF EXISTS git_repository;
DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users
(
    id             SERIAL PRIMARY KEY                                 NOT NULL UNIQUE,
    login          VARCHAR(128)                                       NOT NULL UNIQUE,
    email          VARCHAR(128)                                       NOT NULL UNIQUE,
    name           VARCHAR(256)                                       NOT NULL DEFAULT '',
    avatar_path    VARCHAR(1024)                                      NOT NULL,
    email_verified BOOLEAN                  DEFAULT FALSE             NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS git_repository
(
    id          SERIAL PRIMARY KEY                                 NOT NULL UNIQUE,
    owner_id    INTEGER                                            NOT NULL,
    name        VARCHAR(512)                                       NOT NULL,
    description VARCHAR(2048)                                      NOT NULL,
    is_fork     BOOLEAN                  DEFAULT FALSE             NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES users (id),
    UNIQUE (id, owner_id)
);
