#!/bin/bash
set -e

_os=$(uname -s)

if [ "$_os" != "Linux" ] && [ "$_os" != "Darwin" ]; then
	echo "Unsupported OS: $_os"
	exit 1
fi

_arch=$(uname -m)

if [ "$_arch" != "x86_64" ] && [ "$_arch" != "i386" ] && [ "$_arch" != "arm64" ]; then
	echo "Unsupported arch: $_arch"
	exit 1
fi

if ! [ -x "$(command -v curl)" ]; then
	echo "command 'curl' not found."
	exit 1
fi

# check if the binary exists
download_url="https://risingwave-console.s3.ap-southeast-1.amazonaws.com/$_os/$_arch/risingwave-console"
status_code=$(curl -s -o /dev/null -I -w '%{http_code}' "$download_url")

echo "$download_url"

if [ "$status_code" != 200 ]; then
	echo "Error $status_code: failed to install RisingWave Console."
	exit 1
fi

# download
curl -L -o ./risingwave-console "$download_url"
chmod 755 ./risingwave-console

echo "RisingWave Console installed successfully."
