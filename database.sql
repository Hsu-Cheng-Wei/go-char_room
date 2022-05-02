CREATE DATABASE chat_room;
use chat_room;
CREATE TABLE users(
    id varchar(100) not null,
    password varchar(200) not null,
    name varchar(100) not null,
    email varchar(50),
    PRIMARY KEY(id)
)