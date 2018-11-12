#!/bin/bash
# basic install script for ricochet-group

set -e

if [ ! -e ricochet-group ]; then
    echo -e "Please build ricochet-group first:\n  $ go build" >&2
    exit 1
fi

# create ricochet-group user if there isn't one already
grep -q ^ricochet-group: /etc/passwd || useradd --system --home-dir /var/lib/ricochet-group/ ricochet-group

# create files and directories
cp ricochet-group /usr/local/bin/ricochet-group.new
mv /usr/local/bin/ricochet-group.new /usr/local/bin/ricochet-group
mkdir -p /etc/ricochet-group/
mkdir -p /var/lib/ricochet-group/
chown ricochet-group /var/lib/ricochet-group/ || echo -e "\n\n*** You'll need to change the owner of /var/lib/ricochet-group to whatever user you'll be running ricochet-group as (ricochet-group will need access to the tor control cookie)\n\n------------------\n"

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
    cp config.yaml.install /etc/ricochet-group/config.yaml
fi

echo -e "Installed!\nEdit /etc/ricochet-group/config.yaml to get the config you desire, and then run:\n\n  $ sudo systemctl start ricochet-group\n\nto start ricochet-group, and:\n\n  $ sudo systemctl enable ricochet-group\n\nto have it start automatically at boot.\n"
