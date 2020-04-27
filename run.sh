#!/bin/bash
killall main
sleep 3
go run cmd/auth/main.go &
sleep 1
go run cmd/user/main.go &
sleep 1
go run cmd/gitserver/main.go &
sleep 1
go run cmd/server/main.go &
