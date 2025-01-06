#!/bin/bash

echo "checking $VERSION for $DIR/oapi-codegen"

$DIR/oapi-codegen --version | grep $VERSION

if [ $? -eq 0 ]; then
    exit 0
fi

echo "installing $VERSION for $DIR/oapi-codegen"

GOBIN=$DIR go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@$VERSION
