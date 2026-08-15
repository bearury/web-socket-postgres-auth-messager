create_migrate:
	migrate create -ext sql -dir ./schema -seq create_chats

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5436/postgres?sslmode=disable' up