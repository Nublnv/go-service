up:
	docker-compose up -d --build
down:
	docker-compose down -v
debug:
	docker-compose exec -it go-service /bin/bash
click:
	docker-compose exec -it clickhouse-server clickhouse-client