This is a very early-stage work-in-progress towards a group chat server for Ricochet.

See https://ricochet.im/

Basic functionality is working.

Still to implement:
 - nicknames
 - private groups (currently any client can connect to any group)
 - admin rights (e.g. to kick people)
 - config file
 - peer tracking, so it can automatically connect to peers on startup
 - automatically generate a private key
 - some "command" system, so you can go "/whois jes" to find out the Ricochet ID belonging to the user with nickname "jes", or "/who" to get a list of connected peers
