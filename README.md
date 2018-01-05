# discoirc
[![Build
Status](https://travis-ci.org/cceckman/discoirc.svg?branch=master)](https://travis-ci.org/cceckman/discoirc)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/cceckman/discoirc)
[![Go Report
Card](https://goreportcard.com/badge/github.com/cceckman/discoirc)](https://goreportcard.com/report/github.com/cceckman/discoirc)

```
             ||
             ||
           <><><>
         <><><><><>
        <><><><><><>
        <><><><><><>
        <><><><><><>
         <><><><><>
           <><><>
```


`discoirc` is a terminal-based IRC client. It's similar in some ways to `irssi`,
but with some specific [goals](docs/goals.md) in mind.

## What's with the name?
It was very briefly called `discourirc` - a pun on the French term
[discourir](https://en.wiktionary.org/wiki/discourir), 'to discuss'.

But that was too many letters and is fairly obscure; so, shortened to
`discoirc`, which also suggests a logo should anyone want to make one.

## Documents

See the [Roadmap](docs/roadmap.md) doc for a summary of planned features.

See the [Goals](docs/goals.md) doc for some principles for design.

## Alternatives
I use the venerable `irssi` on a day-to-day basis. But it is insufficiently
documented - I can never get my config file to actually do what's asked - and
frankly I'd rather write my own IRC client than learn Perl in order to make it
do what I want.

I tried the [Komanda](https://github.com/mephux/komanda-cli) CLI. It doesn't
adhere to the principles above, and has some fairly basic bugs - `/me` not
working, for instance.

## Useful libraries

* UI
  * https://github.com/marcusolsson/tui-go is a widget-based UI for the
    terminal.
* IRC
  * https://github.com/fluffle/goirc for IRC connections.
* Config & scripting Interfaces
  * https://github.com/golang/protobuf for internal data structures & config.
  * Or maybe other things at https://github.com/avelino/awesome-go#configuration
* Build
  * https://glide.sh/ maybe? I do like having dependencies at known versions.
