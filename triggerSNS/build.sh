#!/bin/sh
SCRIPT_DIR=$(cd $(dirname $0); pwd)
cd ${SCRIPT_DIR}
GOOS=linux go build -o main main.go
zip main.zip main