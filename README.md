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
 - some "command" system, so you can go "/whois jes" to find out the Ricochet ID belonging to the user with nickname "jes", or "/who" to get a list of connected peers
 - maybe some welcome message when they first login

Commands we might want:
 /who          - list connected peers
 /whois jes    - given Ricochet ID of given nickname
 /nick jes     - set nickname for client
 /welcome      - show the current welcome message
 /help         - explain what commands are available

admin commands:
 /invite $id   - invite given ricochet id to the chat (i.e. group chat should connect to that id and send a contactrequest)
 /welcome $msg - update the welcome message
 /kick jes     - kick the given nick or ricochet id from the group
 /ban jes      - ban the given nick or ricochet id from the group
