#!/bin/bash
killall -q main || echo "No running processes"
sleep 3
nohup go run cmd/auth/main.go >>authservice_nohup.out 2>&1 &
sleep 1
nohup go run cmd/user/main.go >>userservice_nohup.out 2>&1 &
sleep 1
nohup go run cmd/news/main.go >>newsservice_nohup.out 2>&1 &
sleep 1
nohup go run cmd/gitserver/main.go >>gitservice_logfile.log 2>&1 &
sleep 1
nohup go run cmd/server/main.go >>codehubserivce_nohup.out 2>&1 &
