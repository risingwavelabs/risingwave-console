#!/bin/bash

echo "checking $VERSION for $DIR/mockgen"

$DIR/mockgen --version | grep $VERSION

if [ $? -eq 0 ]; then
    exit 0
fi

echo "installing $VERSION for $DIR/mockgen"

GOBIN=$DIR go install github.com/golang/mock/mockgen@v$VERSION
