#!/usr/bin/env bash

sqlc generate

mkdir -p dist
go build -o dist/racl main.go

