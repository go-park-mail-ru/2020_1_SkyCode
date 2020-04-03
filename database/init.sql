drop table if exists users;
drop table if exists sessions;
drop table if exists restaurants;
drop table if exists products;

create table users (
    id        serial      not null primary key,
    firstName varchar(20) not null,
    lastName  varchar(20) not null,
    email     varchar(30),
    phone     varchar(15) not null unique,
    password  varchar(20) not null,
    avatar    varchar(50)
);

create table sessions (
    id      serial      not null primary key,
    userId  int         not null,
    token   varchar(50) not null,
    foreign key (userId) references users (id)
);

create table restaurants (
    id          serial      not null primary key,
    name        varchar(30) not null,
    description text        not null,
    rating      real        not null,
    image       varchar(50)
);

create table products (
    id      serial      not null primary key,
    rest_id int         not null,
    name    varchar(30) not null,
    price   real        not null,
    image   varhcar(50),
    foreign key (restaurant_id) references restaurants (id)
)
