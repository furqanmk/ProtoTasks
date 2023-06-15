PROJECT_NAME=protobuf
DB_FILE=./deploy/db.yml
LIQUIBASE_FILE=./deploy/liquibase.yml
ALL=${DB_FILE}

export COMPOSE_IGNORE_ORPHANS=True

start: db app grpcui

db:
	export COMPOSE_PROJECT_NAME=${PROJECT_NAME}; export COMPOSE_FILE=${DB_FILE}; docker-compose up db -d

dbtool:
	export COMPOSE_PROJECT_NAME=${PROJECT_NAME}; export COMPOSE_FILE=${DB_FILE}; docker-compose up cloudbeaver -d

stop:
	export COMPOSE_PROJECT_NAME=${PROJECT_NAME}; export COMPOSE_FILE=${ALL}; docker-compose down

migrate:
	export COMPOSE_PROJECT_NAME=${PROJECT_NAME}; export COMPOSE_FILE=${LIQUIBASE_FILE}; docker-compose up --force-recreate liquibase

local-build:
	go build -o ./deploy/${PROJECT_NAME} .

app: local-build
	./deploy/${PROJECT_NAME} > ./deploy/${PROJECT_NAME}.log 2>&1 & echo $$! > ./deploy/${PROJECT_NAME}.pid;
	@until [ -s ./deploy/${PROJECT_NAME}.log ]; do sleep 0.1; done
	@tail -n1 ./deploy/${PROJECT_NAME}.log

grpcui:
	 grpcui -plaintext localhost:8008
