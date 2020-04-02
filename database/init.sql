drop table if exists users;
drop table if exists sessions;

create table users
(
    id        serial      not null primary key,
    firstName varchar(20) not null,
    lastName  varchar(20) not null,
    email     varchar(30),
    phone     varchar(15) not null unique,
    password  varchar(20) not null,
    avatar    varchar(50)
);

create table sessions
(
    id      serial      not null primary key,
    userId  int         not null,
    token   varchar(50) not null,
    expires timestamp   not null,
    foreign key (userId) references users (id)
)
