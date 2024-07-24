#!/bin/bash

set -e

BUILD_DIR="map_lan_ip_amd64"

if [ ! -d "$BUILD_DIR" ]
then
    mkdir -p "$BUILD_DIR"
fi

echo -e "\033[32mBuild Start... \033[0m"

CGO_ENABLE=0  GOOS=linux  GOARCH=amd64 go build -a -o $BUILD_DIR/map_lan_ip

cp -a map_lan_ip.yml map_lan_ip.service install.sh $BUILD_DIR/

tar czvf map_lan_ip_amd64.tar.gz $BUILD_DIR/

echo -e "\033[32mBuild successfully! \033[0m"
echo -e "\033[32mFile: ./map_lan_ip_amd64.tar.gz \033[0m"

