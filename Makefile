migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty12@0.0.0.0:5436/postgres?sslmode=disable' up

dockerdb:
	sudo docker run --name=tasks -e POSTGRES_PASSWORD='qwerty12' -p 5436:5432 -d --rm postgres
