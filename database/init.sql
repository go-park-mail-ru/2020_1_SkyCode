drop table if exists users;
drop table if exists sessions;
drop table if exists restaurants;
drop table if exists products;
drop table if exists orders;
drop table if exists orderproducts;
drop table if exists reviews;

create table users
(
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

create table sessions
(
    id         serial      not null primary key,
    userId     int         not null,
    token      varchar(50) not null,
    expiration timestamp   not null default current_timestamp,
    foreign key (userId) references users (id)
);

create table restaurants
(
    id          serial      not null primary key,
    moderId     int         not null,
    name        varchar(30) not null,
    description text        not null,
    rating      real        not null,
    image       varchar(50) not null,
    foreign key (moderId) references users (id)
);

create table products
(
    id      serial      not null primary key,
    rest_id int         not null,
    name    varchar(30) not null,
    price   real        not null,
    image   varchar(50) not null,
    foreign key (rest_id) references restaurants (id)
);

create table orders
(
    id        serial       not null primary key,
    userId    int          not null,
    restId    int          not null,
    address   varchar(255) not null,
    price     real         not null,
    phone     varchar(15)  not null,
    comment   varchar(255),
    personNum int          not null,
    datetime  timestamp    not null default current_timestamp,
    role      varchar(30)  not null
        constraint checkRoleInsert CHECK (role IN ('Accepted', 'Delivering', 'Done')),
    foreign key (userId) references users (id),
    foreign key (restId) references restaurants (id)
);

create table orderProducts
(
    id        serial not null primary key,
    orderId   int    not null,
    productId int    not null,
    count     int    not null,
    foreign key (orderId) references orders (id) on delete cascade,
    foreign key (productId) references products (id) on delete cascade
);

create table reviews
(
    id              serial      not null primary key,
    restId          int         not null,
    userId          int         not null,
    message         text        not null,
    creationDate    timestamp   not null,
    rate            real        not null,
    foreign key (restId) references restaurants (id) on delete cascade,
    foreign key (userId) references users (id) on delete cascade,
        constraint uq_rest_user_reviews unique (restId, userId)
);

create or replace function calculate_rating()
returns trigger as $calculate_rating$
begin
    update restaurants
    set rating =   (select avg(rate) from reviews
                    where reviews.restid = restaurants.id)
    where restaurants.id = coalesce (new.restid, old.restid);
    return new;
end;
$calculate_rating$ LANGUAGE plpgsql;

drop trigger if exists calc_rating on reviews;

create trigger calc_rating
after insert or delete or update of rate on reviews
for each row
execute procedure calculate_rating();
