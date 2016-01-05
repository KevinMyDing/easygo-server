#!/bin/bash

cd Common
./winbuild.sh

cd ../Login
./winbuild.sh

cd ../Gate
./winbuild.sh

cd ../Center
./winbuild.sh

cd ../Client
./winbuild.sh

cd ../Tool
./winbuild.sh

cd ..