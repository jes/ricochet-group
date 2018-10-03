ricochet-group
==============

This is a pretty rudimentary group chat system for Ricochet IM.

See https://ricochet.im/

`ricochet-group` is compatible with the standard Ricochet GUI. It works by creating a central ricochet peer that runs
the group, which relays messages from the sending user to all of the other users. It essentially implements an IRC
channel over the Ricochet protocol.

It hasn't been security-tested and is only a hobby project. Rely on it at your own risk.

Current status
--------------

Public and private group chats work.

There are certainly bugs.
It is certainly possible to evade kicking and banning.
It is probably possible to connect to a private chat without permission.

See file `TODO` for things that still need to be done.

Installation
------------

### Basic usage

Clone the repo:

    $ git clone https://github.com/jes/ricochet-group

If you don't already have go and tor, you'll need to install them.

    $ sudo apt install golang-go tor # on Ubuntu
    $ sudo yum install golang tor    # on CentOS

Fetch the dependencies:

    $ go get -d

(This may take a while and produce no output). Build ricochet-group:

    $ go build

After that, you can get started by running:

    $ ./ricochet-group

The first line will tell you the Ricochet ID of your group chat:

    ricochet-group coming up at ricochet:3yah8ol5a6ub3rto ...

Edit `config.yaml` to customise the configuration.

### Permanent installation

The included script `install.sh` should install `ricochet-group` on Ubuntu and CentOS systems, and hopefully others.

    $ sudo ./install.sh

Having installed `ricochet-group` with `install.sh`, you should edit the config in `/etc/ricochet-group/config.yaml`
and then start it with systemd:

    $ sudo systemctl start ricochet-group

Examine log output with journalctl:

    $ journalctl -u ricochet-group

If all went well, the first line should tell you your group chat's ricochet id:

    Sep 30 01:09:36 localhost ricochet-group[26754]: ricochet-group coming up at ricochet:3yah8ol5a6ub3rto ...

(If it didn't start properly, you'll want to see the error message. In my experience, error messages when the program
exits immediately are not shown by `journalctl -u ricochet-group`, but you should be able to find them with
`journalctl | grep ricochet-group | tail`. I don't know why this is.)

You can then connect to the group chat and verify that it works. It might take 30 seconds or more for the hidden
service to become connectable, so don't be alarmed if it doesn't work immediately.

You can edit the configuration in `/etc/ricochet-group/config.yaml`. Remember to restart `ricochetgroup` whenever you change it.

Once you're satisfied, you can make it start at boot:

    $ sudo systemctl enable ricochet-group

It is possible to run multiple different group chats on one machine by simply running multiple instances of `ricochet-group`.
You'll have to sort out your own systemd configuration for this however.

Configuration
-------------

`ricochet-group` looks for configuration in `./config.yaml` and `/etc/ricochet-group/config.yaml`, with configuration in
`./config.yaml` taking priority.

All of the available configuration options are either used or described in the example config file.

If you want your group chat to use a specific ricochet id (e.g. a vanity address), you can copy the corresponding private key into
`/var/run/ricochet-group/private_key` instead of using the auto-generated key.

### Private groups

Private groups do appear to work, but I can't warrant that there aren't bugs that would allow non-allowed users to read
a private chat, so you rely on it at your own risk.

Contact me
----------

`ricochet-group` is written by James Stanley. You can read my blog at https://incoherency.co.uk/ , email me at
james@incoherency.co.uk, or message me on ricochet at ricochet:it2j3z6t6ksumpzd
