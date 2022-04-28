Forked version... with changes not likely to be merged into original.

Added features to make q.Q() a useful debugging tool:

* A kind of "leveled" logging via q.P which can enable logs in only certain functions/pkgs that match regexp pattern.
* q.O can be set to direct output to a particular file, stdout or stderr
* By default nothing printed via q.Q() because q.P == "".  q.P should be set to regexp to match package/function names
* q.O = "xyz" creates /tmp/xyz file to store q.Q messages
* q.O = "stderr" prints q.Q messages to stderr
* q.O = "stdout" prints q.Q messages to stdout
* q.P = "*" turns on all q.Q messages
* q.P = "xyz.*" turns on q.Q messages for functions or packages that contain string "xyz"
* q.P can be set to any regexp to match the package or function names. That way only certain q.Q() messages in certain package or function will be printed
* Unlike glog, init() does not add to flags to enable q.O or q.P. They are set manually or can be mapped to flags using specific packages such as cobra, flag, urfave/cli, etc. This way you can avoid situations where cobra is used for flags and flag package init() done is nullified.

Example for this forked version usage at qqtest/main.go

# q
[![Build Status](https://travis-ci.org/y0ssar1an/q.svg?branch=develop)](https://travis-ci.org/y0ssar1an/q)
[![GoDoc](https://godoc.org/github.com/y0ssar1an/q?status.svg)](https://godoc.org/github.com/y0ssar1an/q)
[![Go Report Card](https://goreportcard.com/badge/github.com/y0ssar1an/q)](https://goreportcard.com/report/github.com/y0ssar1an/q)

q is a better way to do print statement debugging.

Type `q.Q` instead of `fmt.Printf` and your variables will be printed like this:

![q output examples](https://i.imgur.com/OFmm7pb.png)

## Why is this better than `fmt.Printf`?

* Faster to type
* Pretty-printed vars and expressions
* Easier to see inside structs
* Doesn't go to noisy-ass stdout. It goes to `$TMPDIR/q`.
* Pretty colors!

## Basic Usage

```go
import "github.com/y0ssar1an/q"
...
q.Q(a, b, c)
```
```go
// Alternatively, use the . import and you can omit the package name.
import . "github.com/y0ssar1an/q"
...
Q(a, b, c)
```

For best results, dedicate a terminal to tailing `$TMPDIR/q` while you work.

## Install
```sh
go get -u github.com/y0ssar1an/q
```

Put these functions in your shell config. Typing `qq` or `rmqq` will then start
tailing `$TMPDIR/q`.
```sh
qq() {
    clear
    local gpath="${GOPATH:-$HOME/go}"
    "${gpath%%:*}/src/github.com/y0ssar1an/q/q.sh" "$@"
}
rmqq() {
    if [[ -f "$TMPDIR/q" ]]; then
        rm "$TMPDIR/q"
    fi
    qq
}
```

## Editor Integration

#### VS Code
`Preferences > User Snippets > Go`
```json
"qq": {
    "prefix": "qq",
    "body": "q.Q($1) // DEBUG",
    "description": "Pretty-print to $TMPDIR/q"
}
```

#### Sublime Text
`Tools > Developer > New Snippet`
```xml
<snippet>
    <content><![CDATA[
q.Q($1) // DEBUG
]]></content>
    <tabTrigger>qq</tabTrigger>
    <scope>source.go</scope>
</snippet>
```

#### Atom
`Atom > Open Your Snippets`
```coffee
'.source.go':
    'qq':
        'prefix': 'qq'
        'body': 'q.Q($1) // DEBUG'
```

#### vim/Emacs
TBD Send me a PR, please :)

## Haven't I seen this somewhere before?
Python programmers will recognize this as a Golang port of the
[`q` module by zestyping](https://github.com/zestyping/q).

Ping does a great job of explaining `q` in his awesome lightning talk from
PyCon 2013. Watch it! It's funny :)

[![ping's PyCon 2013 lightning talk](https://i.imgur.com/7KmWvtG.jpg)](https://youtu.be/OL3De8BAhME?t=25m14s)

## FAQ

### Why `q.Q`?
It's quick to type and unlikely to cause naming collisions.

### Is `q.Q()` safe for concurrent use?
Yes.
