-- Main

CREATE TABLE IF NOT EXISTS users
(
    id             SERIAL PRIMARY KEY                  NOT NULL UNIQUE,
    login          VARCHAR(128)                        NOT NULL UNIQUE,
    email          VARCHAR(128)                        NOT NULL UNIQUE,
    name           VARCHAR(256)                        NOT NULL DEFAULT '',
    avatar_path    VARCHAR(1024)                       NOT NULL,
    email_verified BOOLEAN   DEFAULT FALSE             NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS git_repository
(
    id          SERIAL PRIMARY KEY NOT NULL UNIQUE,
    ownerid     INTEGER            NOT NULL,
    name        VARCHAR(512)       NOT NULL,
    description VARCHAR(2048)      NOT NULL,

    FOREIGN KEY (ownerid) REFERENCES users (id),
    UNIQUE (id, ownerid)
);
