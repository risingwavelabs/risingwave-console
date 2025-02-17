#!/bin/bash

echo "checking $VERSION for $DIR/$BIN"

$DIR/$BIN --version | grep $VERSION

if [ $? -eq 0 ]; then
    exit 0
fi

echo "installing $VERSION for $DIR/$BIN"

GOBIN=$DIR go install $REPO@$VERSION
