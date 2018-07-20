#!/bin/bash
set -e

# start and configure testdb Docker container, execute go test, cleanup

TEST_DB_CONTAINER=rgmp-test-db
TEST_DB_USERNAME=postgres
TEST_DB_PASSWORD=restgomuxpq
TEST_DB_NAME=postgres
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_SSL_MODE=disable

trap 'docker rm -f $TEST_DB_CONTAINER' EXIT

if docker run --name $TEST_DB_CONTAINER -d -e POSTGRES_PASSWORD=$TEST_DB_PASSWORD -p $TEST_DB_PORT:5432 postgres:alpine; then
    sleep 1
    export TEST_DB_USERNAME TEST_DB_PASSWORD TEST_DB_NAME TEST_DB_HOST TEST_DB_PORT TEST_DB_SSL_MODE; go test -v
else
    exit 1
fi

exit 0
