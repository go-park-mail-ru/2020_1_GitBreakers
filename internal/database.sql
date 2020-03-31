create table if not exists users
(
    id serial not null
        constraint users_pk
            primary key,
    email text not null,
    password text not null,
    login text not null,
    name text,
    avatar_path text,
    created_at timestamp default CURRENT_TIMESTAMP not null,
    email_verified boolean default false not null
);


create unique index if not exists users_email_uindex
    on users (email);

create unique index if not exists users_id_uindex
    on users (id);

create unique index if not exists users_login_uindex
    on users (login);

create table if not exists storage
(
    id serial not null
        constraint storage_pk
            primary key,
    ownerid integer not null,
    name text not null,
    description text
);


create unique index if not exists storage_id_uindex
    on storage (id);

create table if not exists branches
(
    id serial not null
        constraint branches_pk
            primary key,
    storageid integer not null,
    commits integer default 0 not null
);


create unique index if not exists branches_id_uindex
    on branches (id);

create table if not exists commits
(
    id serial not null
        constraint commits_pk
            primary key,
    description text not null,
    hash text not null,
    branchid integer not null,
    date timestamp default CURRENT_TIMESTAMP not null
);


create unique index if not exists commits_id_uindex
    on commits (id);

