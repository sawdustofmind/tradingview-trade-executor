PG_DSN?=$(shell stoml configs/config.toml db.dsn)
TEST_PG_DSN?=$(shell stoml tests/configs/config.toml db.dsn)
IMAGES=postgres

STAGE_HOST_USER?=root
STAGE_SERVER_HOST?=5.223.46.44

up:
	docker-compose up -d ${IMAGES}

stop:
	docker-compose stop

ps:
	docker-compose ps

pull:
	docker-compose pull ${IMAGES}

generate-api:
	oapi-codegen --config openapi_codegen/types.yaml api/openapi.yaml
	oapi-codegen --config openapi_codegen/user_server.yaml api/openapi.yaml
	oapi-codegen --config openapi_codegen/admin_server.yaml api/openapi.yaml
	oapi-codegen --config openapi_codegen/webhook_server.yaml api/openapi.yaml

dep:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run

build:
	go build cmd/fren/fren.go

build-linux: dep
	GOOS="linux" GOARCH="amd64" go build cmd/fren/fren.go

deploy:
	docker --debug build --platform linux/amd64 -t fren -f Dockerfile .
	docker save fren -o fren.tar
	scp ./fren.tar ${STAGE_HOST_USER}@${STAGE_SERVER_HOST}:/tmp
	#scp ./dist/docker-compose.yml ${STAGE_HOST_USER}@${STAGE_SERVER_HOST}:/tmp
	#ssh ${STAGE_HOST_USER}@${STAGE_SERVER_HOST} 'cp /tmp/docker-compose.yml .'
	ssh ${STAGE_HOST_USER}@${STAGE_SERVER_HOST} 'cd ~ && cat /tmp/fren.tar | sudo docker load'
	ssh ${STAGE_HOST_USER}@${STAGE_SERVER_HOST} 'sudo docker rm -f fren'
	#ssh ${STAGE_HOST_USER}@${STAGE_SERVER_HOST} 'sudo docker network inspect betting >/dev/null 2>&1 || sudo docker network create --driver bridge betting'
	ssh ${STAGE_HOST_USER}@${STAGE_SERVER_HOST} 'cd ~ && sudo docker-compose up -d fren'
	sleep 5
	ssh ${STAGE_HOST_USER}@${STAGE_SERVER_HOST} 'cd ~ && sudo docker-compose ps'

migrate:
	migrate -database postgresql://fren:password@localhost/fren?sslmode=disable -path migrations up

migrate-create:
	migrate create -ext sql -dir ./migrations $(filter-out $@,$(MAKECMDGOALS))