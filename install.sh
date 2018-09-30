#!/bin/bash
# basic install script for ricochet-group

set -e

if [ ! -e ricochet-group ]; then
    echo -e "Please build ricochet-group first:\n  $ go build" >&2
    exit 1
fi

cp ricochet-group /usr/local/bin/
mkdir -p /etc/ricochet-group/
mkdir -p /var/run/ricochet-group/

# don't overwrite handwritten systemd unit
if [ -e /etc/systemd/system/ricochet-group.service ]; then
    cmp --quiet ricochet-group.service /etc/systemd/system/ricochet-group.service || echo -e "\n\n*** Your installed systemd unit is different from the default!\n*** You might want to examine changes between /etc/systemd/system/ricochet-group.service and ./ricochet-group.service and decide whether you want to edit your config.\n\n------------------\n"
else
    cp ricochet-group.service /etc/systemd/system/
fi

# don't overwrite handwritten config
if [ -e /etc/ricochet-group/config.yaml ]; then
    cmp --quiet config.yaml.install /etc/ricochet-group/config.yaml || echo -e "\n\n*** Your installed config file is different from the default!\n*** You might want to examine changes between /etc/ricochet-group/config.yaml and ./config.yaml.install and decide whether you want to edit your config.\n\n------------------\n"
else
    cp config.yaml.install /etc/ricochet-group/
fi

echo -e "Installed!\nMake sure you enable:\n\n  ControlPort 9051\n  CookieAuthentication 1\n\nin your torrc (maybe /etc/tor/torrc).\nNext edit /etc/ricochet-group/config.yaml to taste, and then run:\n\n  $ sudo systemctl start ricochet-group\n\nto start ricochet-group, and:\n\n  $ sudo systemctl enable ricochet-group\n\nto have it start automatically at boot.\n\nIf this is not Ubuntu you may have to edit /etc/systemd/system/ricochet-group.service to make it work, specifically the name of the tor user"
