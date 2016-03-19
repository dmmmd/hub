# TCP message hub written in Go

Just playing around to learn the language.

Not for production.

## Usage

Start the server:

    ./hub -port=1234

Now connect to it:

    nc 127.0.0.1 1234
    
Commands:
    
To get your ID:

    identity
    
To get the list of connected clients:

    list
    
Relay message format:

    relay // Type
    42,100500,9001 // Receivers
    foo bar // From here till the end - body
    umad?
    
    still writing
    
To try:

    echo "relay\n42,100500,9001\nfoo\nbar\numad?" | nc 127.0.0.1 1234