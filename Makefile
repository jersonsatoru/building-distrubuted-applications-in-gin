.PHONY: app/run
app/run:
	go run ./cmd/http

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

.PHONY: db/redis
db/redis:
	docker stop redis || true
	docker rm redis || true
	docker run -d \
		--name=redis \
		-p 6379:6379 \
		-v ${PWD}/redis.conf:/usr/local/etc/redis \
		redis:7.0-rc2-alpine3.15

.PHONY: auth/oauth2
auth/oauth2:
	curl --request POST \
		--url https://dev-a9-nztvy.us.auth0.com/oauth/token \
		--header 'content-type: application/json' \
		--data '{"client_id":"ZyWQelqGRcVXGyeez3aOJMrxEKMrxesg","client_secret":"FEqckbHFKU6WoFGKtfZrrWgsE4qskpVIY1xk06G7ljSCAvgUhIkVaubl5MDZNEBW","audience":"https://dev-a9-nztvy.us.auth0.com/api/v2/","grant_type":"client_credentials"}'

.PHONY: swagger/gen
swagger/gen:
	PORT=8080 swagger generate spec -o ./swagger.yaml --scan-models

.PHONY: swagger/serve
swagger/serve:
	PORT=8080 swagger serve --flavor swagger --no-open ./swagger.yaml