This is a very early-stage work-in-progress towards a group chat server for Ricochet.

See https://ricochet.im/

Basic functionality is working.

Still required for "launch" (i.e. convenient to install and run a public group):
 - automatically connect to known peers at startup
 - automatically generate a private key
 - how should persistent state be tracked (peer list, etc.)?
   - would be handy to be in config file, but it's a bit dirty to be programmatically rewriting the config file
 - easy first-run experience & documentation about how to set it up

Nice to have later:
 - password-protected private groups
   - getting the password correct should add you to allowedusers
   - /invite should also add an id to allowedusers
 - different welcome message based on whether it's a new peer or not? (i.e. to tell them to set a nick with /nick)
 - write to stderr some indication of when the bot is ready to connect to, instead of just waiting?
 - stop replies to commands from appearing in the GUI *before* the commands themselves
   - not sure why this is happening, and fudging the message timestamps doesn't fix it
 - reload config on sighup
   - make sure to disconnect from non-allowed users, banned users, etc.
 - accept config from either yaml file or command line (think viper has this built-in)

admin commands still to implement:

    /invite $id   - invite given ricochet id to the chat (i.e. group chat should connect to that id and send a contactrequest)
    /welcome $msg - update the welcome message
    /kick jes     - kick the given nick or ricochet id from the group
    /ban jes      - ban the given nick or ricochet id from the group
    /admins       - list admins
    /allowedusers - list allowed users

Persistent state we will need to store:
 - welcome message
 - list of allowedusers
 - list of peers to connect to, if different from allowedusers (e.g. on a public group?)
 - list of admins
 - list of banned users
 - onion2Nick

Installation
------------

### Basic usage

As a minimum, you will need to install and run `tor`, and enable its control port. You'll want something in your `torrc`
(`/etc/tor/torrc` on Ubuntu) like:

    ControlPort 9051
    CookieAuthentication 1

After that, you can get started by simply editing `config.yaml` in this directory and then:

    $ sudo -u debian-tor ./ricochet-group

(You need to run it as whatever user has access to the `tor` authentication cookie).

### Permanent installation

The included script `install.sh` should install `ricochet-group` on Ubuntu systems. If you're running something else
it probably won't work. Specifically, the systemd unit in `ricochet-group.service` assumes that the tor user is called
`debian-tor`.

Having installed `ricochet-group` with `install.sh`, you should edit the config in `/etc/ricochet-group/config.yaml`
and then start it with systemd:

    $ sudo systemctl start ricochet-group

Examine log output with journalctl:

    $ journalctl -u ricochet-group

If all went well, the first line should tell you your group chat's ricochet id:

    Sep 30 01:09:36 localhost ricochet-group[26754]: ricochet-group coming up at ricochet:3yah8ol5a6ub3rto ...

You can then connect to the group chat and verify that it works.

Once satisfied, make it start at boot:

    $ sudo systemctl enable ricochet-group

Configuration
-------------

`ricochet-group` looks for configuration in `./config.yaml` and `/etc/ricochet-group/config.yaml`, with configuration in
`./config.yaml` taking priority.

All of the available configuration options are either used or described in the example config file.

If you want your group chat to use a specific ricochet id (e.g. a vanity address), you can copy the corresponding private key into
`/var/run/ricochet-group/private_key` instead of using the auto-generated key.
