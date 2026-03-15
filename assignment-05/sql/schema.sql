drop table if exists users;
create table users (
    id serial primary key,
    name varchar(100) not null,
    email varchar(120) not null unique,
    gender varchar(20) not null,
    birth_date date not null
);

create table user_friends (
    user_id integer not null REFERENCES users(id) on delete CASCADE,
    friend_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    CONSTRAINT no_self_friend CHECK (user_id <> friend_id)
);