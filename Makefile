PROJECT_NAME=protobuf
DB_FILE=./deploy/db.yml
LIQUIBASE_FILE=./deploy/liquibase.yml
ALL=${DB_FILE}

export COMPOSE_IGNORE_ORPHANS=True

db:
	export COMPOSE_PROJECT_NAME=${PROJECT_NAME}; export COMPOSE_FILE=${DB_FILE}; docker-compose up db -d

dbtool:
	export COMPOSE_PROJECT_NAME=${PROJECT_NAME}; export COMPOSE_FILE=${DB_FILE}; docker-compose up cloudbeaver -d

stop:
	export COMPOSE_PROJECT_NAME=${PROJECT_NAME}; export COMPOSE_FILE=${ALL}; docker-compose down

migrate:
	export COMPOSE_PROJECT_NAME=${PROJECT_NAME}; export COMPOSE_FILE=${LIQUIBASE_FILE}; docker-compose up --force-recreate liquibase
