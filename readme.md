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

`-i` is the identity you are signing in with. `-to` is the idenity you wish to start a conversation with.

```
./client -i @nate.errorcode.io -to @kevpatt.errorcode.io
```

## Protocol

### Login

```
LOGIN Identity=@nate.errorcode.io Password=abc123
```

### Send

```
SEND From=@nate.errorcode.io To=@nate.errorcode.io;@kevpatt.errorcode.io Length=12
Hello, World!
```

### OK

```
OK
```