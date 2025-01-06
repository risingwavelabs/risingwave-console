#!/bin/bash 
set +e

docker stop aneith-testpg
docker run --rm -d -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=postgres --name aneith-testpg postgres:14.1-alpine

export AEN_PG_HOST=localhost
export AEN_PG_PORT=5432
export AEN_PG_USER=postgres
export AEN_PG_PASSWORD=postgres
export AEN_PG_DB=postgres
export AEN_PG_MIGRATION=../sql/migrations
export AEN_JWT_SECRET=jwt_secret

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
    docker stop aneith-testpg
    exit $exit_code
}

if [ -z "$HOLD" ]; then
    echo "shutting down database, please wait..."
    docker stop aneith-testpg
else
    trap handle_sigint SIGINT
    echo "Hold for debug. Press Ctrl+C to stop."
    while true; do
        sleep 1
    done
fi
