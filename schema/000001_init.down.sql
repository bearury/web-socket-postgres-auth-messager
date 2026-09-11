DROP TABLE users;
DROP TABLE chats;
DROP TABLE chat_users;
DROP TABLE messages;

-- migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' down