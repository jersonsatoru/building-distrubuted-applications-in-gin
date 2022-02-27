.PHONY: db/mongodb
db/mongodb:
	docker stop mongodb || true
	docker rm mongodb || true
	docker run -d \
		--name mongodb \
		-p 27017:27017 \
		-e MONGO_INITDB_ROOT_USERNAME=satoru \
		-e MONGO_INITDB_ROOT_PASSWORD=satoru \
		-v /mongodb-data:/data/db \
		mongo:5.0.6
.PHONY: db/pg
db/pg:
	docker stop postgres || true
	docker rm postgres || true
	docker run -d \
		--name postgres \
		-p 5433:5432 \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_USER=postgres \
		-e POSTGRES_DB=app_gin \
		postgres:14.1-alpine3.15
