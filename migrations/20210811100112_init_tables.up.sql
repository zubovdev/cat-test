-- user table
create table "user"
(
    id            serial,
    email         varchar(255) not null,
    first_name    varchar(255),
    last_name     varchar(255),
    type          int2 default 1,
    password_hash varchar(255) not null,
    auth_token    varchar(64),
    primary key (id)
);

create unique index user__email__uindex on "user" (email);
create unique index user__auth_token__uindex on "user" (auth_token);

insert into "user"(email, first_name, last_name, type, password_hash)
values ('admin@mail.com', 'Admin', 'Admin', 2, '$2a$10$VqXKAOF9c.DVrKxLhLxDQ.GHj6DQE1r14Jh7pc/MUeymBMZLE03s6');

-- task table
create table task
(
    id             bigserial,
    title          varchar(255) not null,
    description    text,
    estimated_time bigint,
    user_id        bigint,
    status         int2    default 1,
    due_date       bigint,
    mail_sent      boolean default false,
    primary key (id),
    constraint fk_task_user
        foreign key (user_id)
            references "user" (id)
            on update cascade on delete set null
);