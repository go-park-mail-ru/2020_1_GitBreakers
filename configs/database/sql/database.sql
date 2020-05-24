-- Tables

CREATE TABLE IF NOT EXISTS users
(
    id             BIGSERIAL PRIMARY KEY                              NOT NULL UNIQUE,
    login          VARCHAR(128)                                       NOT NULL UNIQUE CHECK ( login <> '' ),
    email          VARCHAR(256)                                       NOT NULL UNIQUE CHECK ( email <> '' ),
    password       VARCHAR(256)                                       NOT NULL CHECK ( password <> '' ),
    name           VARCHAR(256)                                       NOT NULL DEFAULT '',
    avatar_path    VARCHAR(1024)                                      NOT NULL,
    email_verified BOOLEAN                  DEFAULT FALSE             NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS git_repositories
(
    id                  BIGSERIAL PRIMARY KEY                              NOT NULL UNIQUE,
    owner_id            BIGINT                                             NOT NULL,
    name                VARCHAR(512)                                       NOT NULL CHECK ( name <> '' ),
    description         VARCHAR(2048)            DEFAULT ''                NOT NULL,
    is_fork             BOOLEAN                                            NOT NULL,
    is_public           BOOLEAN                                            NOT NULL,
    stars               BIGINT                   DEFAULT 0                 NOT NULL,
    forks               BIGINT                   DEFAULT 0                 NOT NULL,
    merge_requests_open BIGINT                   DEFAULT 0                 NOT NULL,
    parent_id           BIGINT                   DEFAULT NULL,
    created_at          TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES git_repositories (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    UNIQUE (owner_id, name)
);

CREATE TABLE IF NOT EXISTS users_git_repositories
(
    user_id       BIGINT                                             NOT NULL,
    repository_id BIGINT                                             NOT NULL,
    role          VARCHAR(128)                                       NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (repository_id) REFERENCES git_repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    UNIQUE (user_id, repository_id),
    CONSTRAINT users_git_repositories_pk PRIMARY KEY (user_id, repository_id)
);

CREATE TABLE IF NOT EXISTS git_repository_user_stars
(
    repository_id BIGINT                                             NOT NULL,
    author_id     BIGINT                                             NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (repository_id) REFERENCES git_repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (author_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,

    UNIQUE (repository_id, author_id),
    CONSTRAINT git_repository_user_star_pk PRIMARY KEY (repository_id, author_id)
);

CREATE TABLE IF NOT EXISTS issues
(
    id            BIGSERIAL PRIMARY KEY                              NOT NULL UNIQUE,
    author_id     BIGINT,
    repository_id BIGINT                                             NOT NULL,
    title         VARCHAR(256)                                       NOT NULL CHECK ( title <> '' ),
    message       VARCHAR(2048)                                      NOT NULL,
    label         VARCHAR(256)             DEFAULT ''                NOT NULL,
    is_closed     BOOLEAN                  DEFAULT FALSE             NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,
    FOREIGN KEY (repository_id) REFERENCES git_repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS news
(
    id            BIGSERIAL PRIMARY KEY                              NOT NULL UNIQUE,
    author_id     BIGINT,
    repository_id BIGINT                                             NOT NULL,
    message       VARCHAR(2048)                                      NOT NULL CHECK ( message <> '' ),
    label         VARCHAR(256)             DEFAULT ''                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,
    FOREIGN KEY (repository_id) REFERENCES git_repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS merge_requests
(
    id                     BIGSERIAL PRIMARY KEY UNIQUE                       NOT NULL,
    author_id              BIGINT,
    closer_user_id         BIGINT                   DEFAULT NULL,
    from_repository_id     BIGINT,
    to_repository_id       BIGINT                                             NOT NULL,
    from_repository_branch VARCHAR(256)                                       NOT NULL CHECK ( from_repository_branch <> '' ),
    to_repository_branch   VARCHAR(256)                                       NOT NULL CHECK ( to_repository_branch <> '' ),
    title                  VARCHAR(256)                                       NOT NULL CHECK ( title <> '' ),
    message                VARCHAR(2048)                                      NOT NULL,
    label                  VARCHAR(256)             DEFAULT ''                NOT NULL,
    status                 VARCHAR(128)             DEFAULT ''                NOT NULL,
    diff                   BYTEA                    DEFAULT ''                NOT NULL,
    is_closed              BOOLEAN                  DEFAULT FALSE             NOT NULL,
    is_accepted            BOOLEAN                  DEFAULT FALSE             NOT NULL,
    created_at             TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users (id)
        ON DELETE SET NULL ON UPDATE CASCADE,
    FOREIGN KEY (closer_user_id) REFERENCES users (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,
    FOREIGN KEY (from_repository_id) REFERENCES git_repositories (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,
    FOREIGN KEY (to_repository_id) REFERENCES git_repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
