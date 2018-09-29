This is a very early-stage work-in-progress towards a group chat server for Ricochet.

See https://ricochet.im/

Basic functionality is working.

Required for "launch" (i.e. convenient to install and run a public group):
 - automatically connect to known peers at startup
 - automatically generate a private key
 - systemd unit
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
