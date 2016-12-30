# discoirc
A text-based IRC client.

## What's the name from?
It was very briefly called 'discourirc', as in the French
[discourir](https://en.wiktionary.org/wiki/discourir), 'to discus'.

But that was too many letters and less memorable for English speakers (like the
author).

## Why not irssi?
Because I understand neither the configuration language nor Perl, and wanted to
make make something.

In particular, to make something well-documented.

## Components

* discod is the main component. It's a daemon that holds on to the server
  connection, listens for events, runs user scripts (Lua),
  fills in the buffers of windows, etc.
* discoirc, a terminal frontend for discod.

## Dependencies

* Config & scripting Interfaces
  * https://github.com/golang/protobuf for internal data structures & config.
  * https://github.com/Shopify/go-lua for (Lua) scripting
* IRC
  * https://github.com/fluffle/goirc for IRC connections.
* UI
  * https://github.com/jroimartin/gocui for discoirc UI.
  * https://github.com/gizak/termui isn't particuarly useful, but it looks cool.

