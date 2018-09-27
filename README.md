This is a very early-stage work-in-progress towards a group chat server for Ricochet.

See https://ricochet.im/

Basic functionality is working.

Still to implement:
 - nicknames
 - private groups (currently any client can connect to any group)
   - this should support both a configurable list of allowable peers, and a password that will let in anyone who knows it
 - admin rights (e.g. to kick people)
 - config file
 - peer tracking, so it can automatically connect to peers on startup
 - automatically generate a private key
 - different welcome message based on whether it's a new peer or not? (i.e. to tell them to set a nick with /nick)
 - write to stderr some indication of when the bot is ready to connect to, instead of just waiting?
 - stop replies to commands from appearing in the GUI *before* the commands themselves - not sure why this is happening, maybe add 1 second to our message timestamps?

Commands we might still want to implement:
 /whois jes    - given Ricochet ID of given nickname
 /nick jes     - set nickname for client
 /welcome      - show the current welcome message

admin commands:
 /invite $id   - invite given ricochet id to the chat (i.e. group chat should connect to that id and send a contactrequest)
 /welcome $msg - update the welcome message
 /kick jes     - kick the given nick or ricochet id from the group
 /ban jes      - ban the given nick or ricochet id from the group
