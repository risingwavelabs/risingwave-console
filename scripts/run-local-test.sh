#!/bin/bash 
set +e

CONTAINER_NAME=wavekit-testpg

docker stop $CONTAINER_NAME
docker run --rm -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres --name $CONTAINER_NAME postgres:14.1-alpine

export WK_PG_DSN=postgres://postgres:postgres@localhost:5432/postgres
export WK_PG_MIGRATION=../sql/migrations
export WK_ROOT_PASSWORD=123456

exit_code=1
go clean -testcache 
if [ -z "$1" ]; then
    go test -v -timeout 300s -v $TEST_DIR/...
    exit_code=$?
else
    go test -v -timeout 120s -run ^$1$ github.com/risingwavelabs/wavekit/e2e
    exit_code=$?
fi

handle_sigint() {
    echo "Received SIGINT (Ctrl+C). Exiting..."
    echo "shutting down database, please wait..."
    docker stop $CONTAINER_NAME
    exit $exit_code
}

if [ -z "$HOLD" ]; then
    echo "shutting down database, please wait..."
    docker stop $CONTAINER_NAME
else
    trap handle_sigint SIGINT
    echo "Hold for debug. Press Ctrl+C to stop."
    while true; do
        sleep 1
    done
fi
