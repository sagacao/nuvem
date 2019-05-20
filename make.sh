#!/bin/bash

if [ $# != 1 ]; then
	echo "Usage:"
	echo "	./make.sh [game/gate/gmtools/websvr]"
	exit 0
fi

model=$1

echo "start make ..."
go build -tags=jsoniter app/${model}/${model}.go 
if [ $? -ne 0 ]; then
    echo "make error ..."
    exit 1
fi

echo "make success"
