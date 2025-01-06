#!/bin/bash

echo "checking $VERSION for $DIR/sqlc"

$DIR/sqlc version | grep $VERSION

if [ $? -eq 0 ]; then
    exit 0
fi

echo "installing $VERSION for $DIR/sqlc"

GOBIN=$DIR go install github.com/sqlc-dev/sqlc/cmd/sqlc@$VERSION
