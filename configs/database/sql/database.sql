-- Main

CREATE TABLE IF NOT EXISTS users
(
    id             BIGSERIAL PRIMARY KEY                              NOT NULL UNIQUE,
    login          VARCHAR(128)                                       NOT NULL UNIQUE CHECK ( login <> '' ),
    email          VARCHAR(128)                                       NOT NULL UNIQUE CHECK ( email <> '' ),
    password       VARCHAR(256)                                       NOT NULL CHECK ( password <> '' ),
    name           VARCHAR(256)                                       NOT NULL DEFAULT '',
    avatar_path    VARCHAR(1024)                                      NOT NULL,
    email_verified BOOLEAN                  DEFAULT FALSE             NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS git_repositories
(
    id          BIGSERIAL PRIMARY KEY                              NOT NULL UNIQUE,
    owner_id    INTEGER                                            NOT NULL,
    name        VARCHAR(512)                                       NOT NULL,
    description VARCHAR(2048)                                      NOT NULL,
    is_public   BOOLEAN                                            NOT NULL,
    is_fork     BOOLEAN                                            NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    UNIQUE (id, owner_id)
);


CREATE TABLE IF NOT EXISTS users_git_repositories
(
    user_id       BIGINT                                             NOT NULL,
    repository_id BIGINT                                             NOT NULL,
    role          VARCHAR(64)              DEFAULT ''                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (repository_id) REFERENCES git_repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    UNIQUE (repository_id, user_id),
    CONSTRAINT users_git_repositories_pk PRIMARY KEY (repository_id, user_id)
);

alter table git_repositories
    add if not exists stars integer default 0 not null;

CREATE TABLE IF NOT EXISTS git_repository_user_star
(
    repository_id BIGINT                                             NOT NULL,
    user_id       BIGINT                                             NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (repository_id) REFERENCES git_repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    UNIQUE (repository_id, user_id),
    CONSTRAINT git_repository_user_star_pk PRIMARY KEY (repository_id, user_id)
)
