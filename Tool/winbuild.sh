#!/bin/bash

# cd ../Common
# COMMONDIR=`pwd`
# export GOPATH=$GOPATH':'$COMMONDIR
# cd ../Tool

# CURDIR=`pwd`

# export GOPATH=$GOPATH':'$COMMONDIR':'$CURDIR
# echo $GOPATH

go build ./src/GxTool
# go build ./src/Test

rm -rf bin/log
mkdir -p bin/log

mv GxTool ./bin/gxtool
# mv Test ./bin/gxTest

echo '==========>build Tool ok'