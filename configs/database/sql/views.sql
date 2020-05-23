create or replace view user_profile_view(id, login, email, name, avatar_path, created_at) as
SELECT users.id,
       users.login,
       users.email,
       users.name,
       users.avatar_path,
       users.created_at
FROM users;


create or replace view git_repository_user_view(id, owner_id, name, description, is_fork, is_public, stars, forks, merge_requests_open, parent_id, created_at, user_login, user_email, user_name, user_avatar_path, user_created_at) as
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
FROM git_repositories gr
         JOIN user_profile_view upv ON gr.owner_id = upv.id;


create or replace view git_repository_parent_user_view(id, owner_id, name, description, is_fork, is_public, stars, forks, merge_requests_open, parent_id, created_at, user_login, user_email, user_name, user_avatar_path, user_created_at, parent_owner_id, parent_name, parent_description, parent_is_fork, parent_is_public, parent_stars, parent_forks, parent_merge_requests_open, parent_parent_id, parent_created_at, parent_user_login, parent_user_email, parent_user_name, parent_user_avatar_path, parent_user_created_at) as
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
FROM git_repository_user_view gr
         LEFT JOIN git_repository_user_view grparent ON gr.parent_id = grparent.id;


create or replace view users_git_repositories_view(repository_id, user_id, role, created_at, git_repository_owner_id, git_repository_name, git_repository_description, git_repository_is_fork, git_repository_is_public, git_repository_stars, git_repository_forks, git_repository_merge_requests_open, git_repository_created_at, git_repository_parent_id, user_login, user_email, user_name, user_avatar_path, user_created_at) as
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
FROM users_git_repositories ugr
         JOIN git_repositories gr ON ugr.repository_id = gr.id
         JOIN user_profile_view upv ON ugr.user_id = upv.id;


create or replace view git_repository_user_stars_view(repository_id, author_id, created_at, user_id, user_login, user_email, user_name, user_avatar_path, user_created_at, git_repository_owner_id, git_repository_name, git_repository_description, git_repository_is_fork, git_repository_is_public, git_repository_stars, git_repository_forks, git_repository_merge_requests_open, git_repository_created_at, git_repository_parent_id) as
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
FROM git_repository_user_stars grus
         JOIN user_profile_view upv ON grus.author_id = upv.id
         JOIN git_repositories gr ON grus.repository_id = gr.id;


create or replace view issues_users_view(id, author_id, repository_id, title, message, label, is_closed, created_at, user_id, user_login, user_email, user_name, user_avatar_path, user_created_at) as
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
FROM issues i
         JOIN user_profile_view upv ON i.author_id = upv.id;


create or replace view news_users_view(id, author_id, repository_id, message, label, created_at, user_id, user_login, user_email, user_name, user_avatar_path, user_created_at) as
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
FROM news n
         JOIN user_profile_view upv ON n.author_id = upv.id;


create or replace view merge_requests_view(id, author_id, from_repository_id, to_repository_id, from_repository_branch, to_repository_branch, title, message, label, status, diff, is_closed, is_accepted, created_at, user_login, user_email, user_name, user_avatar_path, user_created_at, git_repository_from_owner_id, git_repository_from_name, git_repository_from_description, git_repository_from_is_fork, git_repository_from_is_public, git_repository_from_stars, git_repository_from_forks, git_repository_from_merge_requests_open, git_repository_from_parent_id, git_repository_from_created_at, git_repository_to_owner_id, git_repository_to_name, git_repository_to_description, git_repository_to_is_fork, git_repository_to_is_public, git_repository_to_stars, git_repository_to_forks, git_repository_to_merge_requests_open, git_repository_to_parent_id, git_repository_to_created_at) as
SELECT mr.id,
       mr.author_id,
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
FROM merge_requests mr
         JOIN user_profile_view upv ON mr.author_id = upv.id
         JOIN git_repositories grfrom ON mr.from_repository_id = grfrom.id
         LEFT JOIN git_repositories grto ON mr.to_repository_id = grto.id;
