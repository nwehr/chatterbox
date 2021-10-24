## Usage

### Install Dependencies

```
$ go get ./...
```

### Build Client

```
$ make client
```

### Run Client

`-i` is the identity you are signing in with. `-r` is a semicolon-delimited list of recipients.

```
./client -i @nate.errorcode.io -r @kevpatt.errorcode.io;@nate.errorcode.io
```

## Protocol

### MSG

The `MSG` type sends a text message to the recipients. Both clients and servers send `MSG` types.

`MSG` Args:

* `Uuid` required. The client-generated uuid for the message.
* `From` required. The identity of the originator of the message.
* `Recipients` required. The list of identities who will receive the message (including the originator).
* `Encoding` required. The encoding of the datatype.
* `Length` required. The size of the data in bytes.
* `ConversationId` optional. Sha256 hash of the recipients.

Example:

```
[C|S]: MSG From=@nate.errorcode.io Recipients=@kevpatt.errorcode.io;@nate.errorcode.io Encoding=text/plain Length=12 Uuid=cf52b8b0-32a9-11ec-aed9-db66117c16da
       Hello, World!
```

### LCONV

The `LCONV` type requests a list of conversations from the server. The server response with a `CONV` type. 

`LCONV` Args:

* `Recipient` required. The identity of the recipient of conversations.

`CONV` Args:

* `Id` required. The conversation id which is a sha256 hash of the recipients.
* `Recipients` required. The list of identities who are recipients of the conversation.
* `Unread` required. The number of unread messages for the recipient identified in the request.

Example:

```
[C]: LCONV Recipient=@nate.errorcode.io

[S]: CONV Id=76359067c82a3f8e1dadc0e93570f2113c4f6b2b66994ce5e115d9be6d983de6 Recipients=@nate.errorcode.io Unread=1
[S]: CONV Id=a83087325cec029b2c39464aa0d2ea94f9282531cb53ad41143ea8bcea78205a Recipients=@kevpatt.errorcode.io;@nate.errorcode.io Unread=3
```

### LMSG

The `LMSG` type requests a list of messages from the server identified by the conversation id.

`LMSG` Args:

* `ConversationId` required. 

Example:



```
[C]: LMSG ConversationId=a83087325cec029b2c39464aa0d2ea94f9282531cb53ad41143ea8bcea78205a

[S]: MSG From=@nate.errorcode.io Recipients=@nate.errorcode.io;@kevpatt.errorcode.io Encoding=text/plain Length=12 Uuid=cf52b8b0-32a9-11ec-aed9-db66117c16da
     Hello, World!
```