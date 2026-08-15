CREATE TABLE users
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) not null,
    username varchar(255) not null unique,
    password_hash varchar(255) not null
);


CREATE TYPE chat_type AS ENUM ('chat', 'group', 'chain');

CREATE TABLE chats
(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name varchar(255) not null,
    type chat_type not null
);

CREATE TABLE chat_users
(
    chat_id UUID REFERENCES chats (id) ON DELETE CASCADE,
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (chat_id, user_id)
);