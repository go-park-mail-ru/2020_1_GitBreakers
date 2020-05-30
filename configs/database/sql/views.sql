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
       gr.forks,
       gr.merge_requests_open,
       gr.parent_id,
       gr.created_at,
       upv.login       AS user_login,
       upv.email       AS user_email,
       upv.name        AS user_name,
       upv.avatar_path AS user_avatar_path,
       upv.created_at  AS user_created_at
FROM git_repositories AS gr
         JOIN user_profile_view upv ON gr.owner_id = upv.id;


DROP VIEW IF EXISTS git_repository_parent_user_view CASCADE;
CREATE VIEW git_repository_parent_user_view AS
SELECT gr.id,
       gr.owner_id,
       gr.name,
       gr.description,
       gr.is_fork,
       gr.is_public,
       gr.stars,
       gr.forks,
       gr.merge_requests_open,
       gr.parent_id,
       gr.created_at,
       gr.user_login,
       gr.user_email,
       gr.user_name,
       gr.user_avatar_path,
       gr.user_created_at,
       grparent.owner_id            AS parent_owner_id,
       grparent.name                AS parent_name,
       grparent.description         AS parent_description,
       grparent.is_fork             AS parent_is_fork,
       grparent.is_public           AS parent_is_public,
       grparent.stars               AS parent_stars,
       grparent.forks               AS parent_forks,
       grparent.merge_requests_open AS parent_merge_requests_open,
       grparent.parent_id           AS parent_parent_id,
       grparent.created_at          AS parent_created_at,
       grparent.user_login          AS parent_user_login,
       grparent.user_email          AS parent_user_email,
       grparent.user_name           AS parent_user_name,
       grparent.user_avatar_path    AS parent_user_avatar_path,
       grparent.user_created_at     AS parent_user_created_at
FROM git_repository_user_view AS gr
         LEFT JOIN git_repository_user_view grparent ON gr.parent_id = grparent.id;


DROP VIEW IF EXISTS users_git_repositories_view CASCADE;
CREATE VIEW users_git_repositories_view AS
SELECT ugr.repository_id,
       ugr.user_id,
       ugr.role,
       ugr.created_at,
       gr.owner_id            AS git_repository_owner_id,
       gr.name                AS git_repository_name,
       gr.description         AS git_repository_description,
       gr.is_fork             AS git_repository_is_fork,
       gr.is_public           AS git_repository_is_public,
       gr.stars               AS git_repository_stars,
       gr.forks               AS git_repository_forks,
       gr.merge_requests_open AS git_repository_merge_requests_open,
       gr.created_at          AS git_repository_created_at,
       gr.parent_id           AS git_repository_parent_id,
       upv.login              AS user_login,
       upv.email              AS user_email,
       upv.name               AS user_name,
       upv.avatar_path        AS user_avatar_path,
       upv.created_at         AS user_created_at
FROM users_git_repositories AS ugr
         JOIN git_repositories AS gr ON ugr.repository_id = gr.id
         JOIN user_profile_view AS upv ON ugr.user_id = upv.id;


DROP VIEW IF EXISTS git_repository_user_stars_view CASCADE;
CREATE VIEW git_repository_user_stars_view AS
SELECT grus.repository_id,
       grus.author_id,
       grus.created_at,
       upv.id                 AS user_id,
       upv.login              AS user_login,
       upv.email              AS user_email,
       upv.name               AS user_name,
       upv.avatar_path        AS user_avatar_path,
       upv.created_at         AS user_created_at,
       gr.owner_id            AS git_repository_owner_id,
       gr.name                AS git_repository_name,
       gr.description         AS git_repository_description,
       gr.is_fork             AS git_repository_is_fork,
       gr.is_public           AS git_repository_is_public,
       gr.stars               AS git_repository_stars,
       gr.forks               AS git_repository_forks,
       gr.merge_requests_open AS git_repository_merge_requests_open,
       gr.created_at          AS git_repository_created_at,
       gr.parent_id           AS git_repository_parent_id
FROM git_repository_user_stars AS grus
         JOIN user_profile_view AS upv ON grus.author_id = upv.id
         JOIN git_repositories AS gr ON grus.repository_id = gr.id;


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

DROP VIEW IF EXISTS merge_requests_view CASCADE;
CREATE VIEW merge_requests_view AS
SELECT mr.id,
       mr.author_id,
       mr.closer_user_id,
       mr.from_repository_id,
       mr.to_repository_id,
       mr.from_repository_branch,
       mr.to_repository_branch,
       mr.title,
       mr.message,
       mr.label,
       mr.status,
       mr.diff,
       mr.is_closed,
       mr.is_accepted,
       mr.created_at,
       upv.login                  AS user_login,
       upv.email                  AS user_email,
       upv.name                   AS user_name,
       upv.avatar_path            AS user_avatar_path,
       upv.created_at             AS user_created_at,
       upv_closer.login           AS closer_user_login,
       upv_closer.email           AS closer_user_email,
       upv_closer.name            AS closer_user_name,
       upv_closer.avatar_path     AS closer_user_avatar_path,
       upv_closer.created_at      AS closer_user_created_at,
       grfrom.owner_id            AS git_repository_from_owner_id,
       grfrom.name                AS git_repository_from_name,
       grfrom.description         AS git_repository_from_description,
       grfrom.is_fork             AS git_repository_from_is_fork,
       grfrom.is_public           AS git_repository_from_is_public,
       grfrom.stars               AS git_repository_from_stars,
       grfrom.forks               AS git_repository_from_forks,
       grfrom.merge_requests_open AS git_repository_from_merge_requests_open,
       grfrom.parent_id           AS git_repository_from_parent_id,
       grfrom.created_at          AS git_repository_from_created_at,
       grto.owner_id              AS git_repository_to_owner_id,
       grto.name                  AS git_repository_to_name,
       grto.description           AS git_repository_to_description,
       grto.is_fork               AS git_repository_to_is_fork,
       grto.is_public             AS git_repository_to_is_public,
       grto.stars                 AS git_repository_to_stars,
       grto.forks                 AS git_repository_to_forks,
       grto.merge_requests_open   AS git_repository_to_merge_requests_open,
       grto.parent_id             AS git_repository_to_parent_id,
       grto.created_at            AS git_repository_to_created_at
FROM merge_requests AS mr
         LEFT JOIN user_profile_view AS upv ON mr.author_id = upv.id
         LEFT JOIN user_profile_view AS upv_closer ON mr.closer_user_id = upv_closer.id
         LEFT JOIN git_repositories AS grfrom ON mr.from_repository_id = grfrom.id
         JOIN git_repositories AS grto ON mr.to_repository_id = grto.id;

