create table if not exists users(
    id serial primary key,
    name varchar(255) not null,
    email varchar(255) unique,
    age int,
    created_at timestamp not null default now()
);
insert into users (name, email, age) values ('john doe', 'john@gmail.com', 32);