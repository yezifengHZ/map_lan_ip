#!/bin/bash

set -e

INSTALL_DIR="/opt/map_lan_ip"

if [ ! -d "$INSTALL_DIR" ]
then
    mkdir -p "$INSTALL_DIR"
fi

echo -e "\033[34m======== Installing map_lan_ip ======== \033[0m"

cp map_lan_ip map_lan_ip.yml $INSTALL_DIR/

echo -e "\033[34m======== Configuring service ======== \033[0m"

cp map_lan_ip.service /usr/lib/systemd/system/

systemctl daemon-reload
systemctl start map_lan_ip
systemctl enable map_lan_ip

echo -e "\033[32mInstalled successfully! \033[0m"

