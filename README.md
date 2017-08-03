# Server Setup
To set a password, go to the server folder and run `rcserver.exe pwgen [password]`.
To generate the public and private keys, run `rcserver.exe privgen` then `rcserver.exe pubgen`.
The private key is used by the server, and the public key should be sent to those gaining access (it doesn't need to stay in the server folder).
Finally, to run the server, run rcserver.exe.
# Client Setup
To set up a client, first add the public key to the client folder. Then set the password using `rcclient.exe pwgen [password]`.
To run, do one of the following:
 - `rcclient.exe cmd [address] "[command]"` - run a command
 - `rcclient.exe console [address]` - run multiple commands

(note that address should be on port 3924)
