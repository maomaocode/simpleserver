#!/bin/bash

# 远程服务器信息
REMOTE_HOST="45.77.33.136"
REMOTE_USER="root"
REMOTE_PATH="/data/"

echo "Building binary..."
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o formserver .

echo "Stopping binary on remote server..."
ssh $REMOTE_USER@$REMOTE_HOST "pkill formserver"

echo "Uploading new binary to remote server..."
scp formserver $REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH

echo "Starting new binary on remote server..."
ssh $REMOTE_USER@$REMOTE_HOST "cd $REMOTE_PATH && /bin/bash -c 'nohup ./formserver > ./nohup.log 2>&1 &'"

echo "Upload success"