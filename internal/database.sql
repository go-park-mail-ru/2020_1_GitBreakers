-- Main

-- DROP TABLE IF EXISTS git_repository_users;
-- DROP TABLE IF EXISTS git_repository;
-- DROP TABLE IF EXISTS users;

CREATE TABLE IF NOT EXISTS users
(
    id             SERIAL PRIMARY KEY                                 NOT NULL UNIQUE,
    login          VARCHAR(128)                                       NOT NULL UNIQUE,
    email          VARCHAR(128)                                       NOT NULL UNIQUE,
    password       VARCHAR(256)                                       NOT NULL,
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
    is_public   BOOLEAN                                            NOT NULL,
    is_fork     BOOLEAN                                            NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES users (id),
    UNIQUE (id, owner_id)
);


CREATE TABLE IF NOT EXISTS git_repository_users
(
    repository_id INTEGER                                            NOT NULL,
    user_id       INTEGER                                            NOT NULL,
    role          VARCHAR(64)              DEFAULT ''                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (repository_id) REFERENCES git_repository (id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    UNIQUE (repository_id, user_id),
    CONSTRAINT git_repository_users_pk PRIMARY KEY (repository_id, user_id)
);
