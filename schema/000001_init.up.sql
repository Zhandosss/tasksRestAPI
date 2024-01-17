CREATE TABLE users
(
    id serial primary key,
    firstname varchar(255) not null,
    lastname varchar(255) not null,
    login varchar(255) not null,
    password_hash varchar(255) not null
);

CREATE TABLE tasks
(
    id serial primary key,
    task text not null,
    date timestamp not null,
    owner_id int references users (id) on delete cascade not null
);

create TABLE tags (
    id serial primary key,
    tag varchar(255) not null
);

create TABLE tags_in_task(
    task_id int references tasks (id) on delete cascade not null,
    tag_id int references tags (id) not null
);