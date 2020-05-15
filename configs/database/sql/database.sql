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
    owner_id    BIGINT                                             NOT NULL,
    name        VARCHAR(512)                                       NOT NULL CHECK ( name <> '' ),
    description VARCHAR(2048)            DEFAULT ''                NOT NULL,
    is_fork     BOOLEAN                                            NOT NULL,
    is_public   BOOLEAN                                            NOT NULL,
    stars       BIGINT                   DEFAULT 0                 NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    parent_id   BIGINT                   DEFAULT NULL,

    FOREIGN KEY (owner_id) REFERENCES users (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES git_repositories (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,

    UNIQUE (id, owner_id)
);

CREATE TABLE IF NOT EXISTS users_git_repositories
(
    user_id       BIGINT                                             NOT NULL,
    repository_id BIGINT                                             NOT NULL,
    role          VARCHAR(64)                                        NOT NULL,
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
    label         VARCHAR(64)              DEFAULT ''                NOT NULL,
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
    label         VARCHAR(64)              DEFAULT ''                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users (id)
        ON DELETE SET NULL
        ON UPDATE CASCADE,
    FOREIGN KEY (repository_id) REFERENCES git_repositories (id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- Views

DROP VIEW IF EXISTS user_profile_view CASCADE;
CREATE VIEW user_profile_view AS
SELECT id,
       login,
       email,
       name,
       avatar_path,
       created_at
FROM users;


DROP VIEW IF EXISTS git_repository_user_view CASCADE;
CREATE VIEW git_repository_user_view AS
SELECT gr.id,
       gr.owner_id,
       gr.name,
       gr.description,
       gr.is_fork,
       gr.is_public,
       gr.stars,
       gr.created_at,
       upv.id          AS user_id,
       upv.login       AS user_login,
       upv.email       AS user_email,
       upv.name        AS user_name,
       upv.avatar_path AS user_avatar_path,
       upv.created_at  AS user_created_at
FROM git_repositories AS gr
         JOIN user_profile_view upv ON gr.owner_id = upv.id;


DROP VIEW IF EXISTS git_repository_user_stars_view CASCADE;
CREATE VIEW git_repository_user_stars_view AS
SELECT grus.repository_id,
       grus.author_id,
       grus.created_at,
       upv.id          AS user_id,
       upv.login       AS user_login,
       upv.email       AS user_email,
       upv.name        AS user_name,
       upv.avatar_path AS user_avatar_path,
       upv.created_at  AS user_created_at
FROM git_repository_user_stars AS grus
         JOIN user_profile_view AS upv ON grus.author_id = upv.id;


DROP VIEW IF EXISTS issues_users_view CASCADE;
CREATE VIEW issues_users_view AS
SELECT i.id,
       i.author_id,
       i.repository_id,
       i.title,
       i.message,
       i.label,
       i.is_closed,
       i.created_at,
       upv.id          AS user_id,
       upv.login       AS user_login,
       upv.email       AS user_email,
       upv.name        AS user_name,
       upv.avatar_path AS user_avatar_path,
       upv.created_at  AS user_created_at
FROM issues AS i
         JOIN user_profile_view AS upv on i.author_id = upv.id;


DROP VIEW IF EXISTS news_users_view CASCADE;
CREATE VIEW news_users_view AS
SELECT n.id,
       n.author_id,
       n.repository_id,
       n.message,
       n.label,
       n.created_at,
       upv.id          AS user_id,
       upv.login       AS user_login,
       upv.email       AS user_email,
       upv.name        AS user_name,
       upv.avatar_path AS user_avatar_path,
       upv.created_at  AS user_created_at
FROM news AS n
         JOIN user_profile_view AS upv ON n.author_id = upv.id;
