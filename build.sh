#!/bin/bash

cd ./${0%/*}
b=${PWD#$HOME/}
export GOPATH=$PWD

echo Working dir $b
GOOS=linux GOARCH=arm GOARM=6 go build -o bin/bbq src/bbqcontroller.go && rsync -av bin/bbq 10.0.0.155:$b/bin/bbq
