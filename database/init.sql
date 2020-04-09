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
    avatar    varchar(50),
    role      varchar(30) not null
    constraint checkRoleInsert CHECK (role IN ('Admin', 'User', 'Moderator'))
);

create table sessions (
    id      serial      not null primary key,
    userId  int         not null,
    token   varchar(50) not null,
    foreign key (userId) references users (id)
);

create table restaurants (
    id          serial      not null primary key,
    moderId     int         not null,
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
    image   varchar(50),
    foreign key (rest_id) references restaurants (id)
);

create table orders (
    id        serial      not null primary key,
    userId    int         not null,
    address   varchar(255) not null,
    price     real        not null,
    phone     varchar(15) not null,
    comment   varchar(255),
    personNum int         not null,
    foreign key (userId) references users (id)
);

create table orderProducts (
    id        serial      not null primary key,
    orderId   int         not null,
    productId int         not null,
    count     int         not null,
    foreign   key (orderId) references orders (id) on delete cascade,
    foreign   key (productId) references products (id) on delete cascade
);
