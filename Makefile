up:
	docker compose build --no-cache && docker-compose up 
down:
	docker compose down
restart:
	docker compose restart
logs:
	docker compose logs -f
login:
	docker exec -it db mysql -u root -p


.PHONY: sql generate

# SQL CLIを起動
sql:
	mysql -u user -p -h 127.0.0.1 -P 53306 app

# スキーマのSQLを実行し、SQLCでコード生成
generate:
	mysql -u user -p -h 127.0.0.1 -P 53306 app < src/db/schema.sql && sqlc generate
