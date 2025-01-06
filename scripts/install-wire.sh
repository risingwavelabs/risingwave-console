#!/bin/bash

echo "checking $VERSION for $DIR/wire"

stat $DIR/wire

if [ $? -eq 0 ]; then
    exit 0
fi

echo "installing $VERSION for $DIR/wire"

GOBIN=$DIR go install github.com/google/wire/cmd/wire@$VERSION
