#!/bin/bash
killall main
sleep 3
nohup go run cmd/auth/main.go >> authservice_nohup.log &
sleep 1
nohup go run cmd/user/main.go >> userservice_nohup.log &
sleep 1
nohup go run cmd/news/main.go >> newsservice_nohup.out &
sleep 1
nohup go run cmd/gitserver/main.go >> gitservice_logfile.log &
sleep 1
nohup go run cmd/server/main.go >> codehubserivce_nohup.out &
