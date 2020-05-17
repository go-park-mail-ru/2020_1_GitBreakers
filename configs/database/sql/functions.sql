create or replace function upd_news_with_collaborator() returns trigger
    language plpgsql
as
$$
DECLARE
    mess  varchar;
    login varchar;
BEGIN
    select users.login from users where id = new.user_id into login;

    IF (TG_OP = 'UPDATE') THEN
        mess = concat('Update role for collaborator ', login, ' ', new.role);
    ELSIF (TG_OP = 'INSERT') THEN
        mess = concat('New collaborator ', login);
    END IF;
    INSERT INTO news(author_id, repository_id, message) values (new.user_id, new.repository_id, mess);
    RETURN null;
END;
$$;

alter function upd_news_with_collaborator() owner to codehub_dev;

create function upd_news_with_issues() returns trigger
    language plpgsql
as
$$
DECLARE
    mess varchar;
BEGIN
    IF (TG_OP = 'UPDATE') THEN
        mess = concat('Update issues');
    ELSIF (TG_OP = 'INSERT') THEN
        mess = concat('New issues');
    END IF;
    INSERT INTO news(author_id, repository_id, message) values (new.author_id, new.repository_id, mess);
    RETURN null;
END;
$$;

alter function upd_news_with_issues() owner to codehub_dev;

create function upd_news_with_star() returns trigger
    language plpgsql
as
$$
DECLARE
    mess  varchar;
    login varchar;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        select users.login from users where id = new.author_id into login;
        mess = concat('Star added by ', login);
        INSERT INTO news(author_id, repository_id, message) values (new.author_id, new.repository_id, mess);

    ELSIF (TG_OP = 'DELETE') THEN
        select users.login from users where id = old.author_id into login;
        mess = concat('Star deleted by ', login);
        INSERT INTO news(author_id, repository_id, message) values (old.author_id, old.repository_id, mess);

    END IF;
    RETURN null;
END ;
$$;

alter function upd_news_with_star() owner to codehub_dev;

create function upd_stars_repo() returns trigger
    language plpgsql
as
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        update git_repositories set stars=stars + 1 where id = new.repository_id;
    ELSIF (TG_OP = 'DELETE') THEN
        update git_repositories set stars=stars - 1 where id = old.repository_id;
    END IF;
    RETURN null;
END;
$$;

alter function upd_stars_repo() owner to codehub_dev;

DROP FUNCTION IF EXISTS update_forks_on_git_repository_insert_or_delete CASCADE;
CREATE OR REPLACE FUNCTION update_forks_on_git_repository_insert_or_delete() RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE git_repositories SET forks = forks + 1 WHERE id = new.parent_id;
    ELSEIF (TG_OP = 'DELETE') THEN
        UPDATE git_repositories SET forks = forks - 1 WHERE id = old.parent_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_forks_on_git_repository_insert
    AFTER INSERT
    ON git_repositories
    FOR EACH ROW
EXECUTE PROCEDURE update_forks_on_git_repository_insert_or_delete();

CREATE TRIGGER update_forks_on_git_repository_delete
    AFTER DELETE
    ON git_repositories
    FOR EACH ROW
EXECUTE PROCEDURE update_forks_on_git_repository_insert_or_delete();
