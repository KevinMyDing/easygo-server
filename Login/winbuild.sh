#!/bin/bash

go build ./src/Login

rm -rf bin/log
mkdir -p bin/log

mv Login ./bin/gxlogin

echo '==========>build Login ok'

