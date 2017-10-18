# discoirc
A text-based IRC client.

## What's the name from?
It was very briefly called 'discourirc', as in the French
[discourir](https://en.wiktionary.org/wiki/discourir), 'to discus'.

But that was too many letters and less memorable for English speakers (like the
author).


## Principles

* Leave window management to window managers.
  * Separate sessions / connections from window sessions.
  * Identify how to create a new window in supported environments (xterm, tmux,
    screen) and leave window management to those.
* Keep it simple.
  * Starting off without plugins; if it matters enough to write a plugin for,
    it's probably worth having in the client itself.
  * Starting off with a subset of IRC features. We have better alternatives for
    sending files than DCC nowadays. (And, does CTCP need to be implemented?)

## Why not irssi?
Because I'm not going to learn Perl, and the config file is untyped and
undocumented.

## 1.0 Roadmap

These are the features I'd like to land to get up to 1.0:

### 0.1: Min UI
UI panels in roughly working order, implemented against a stub backend.
Allow switching between them, but not launching windows. Statically
configure channels and connections on the stub server.

* Channel: new messages incoming, send messages, list connection+channel,
  nick.
* Session: Current connections and channels. Defer editing / starting new
  connections.

### 0.2: IRC
Implement the backend interface as a connection to IRC servers. Still
statically configure connections and channels. via code.

### 0.3: Configuration
Load configuration of connections from file. Write to a file. Support connecting
to new servers and channels in the UI. This may take the form of a new
"connection" master view, kind of in between sesion and channel.

Allow "autoload" (watch for updates) and "autowrite" (persist settings)
both as options. Counterintuitively, if you unset 'autowrite', it doesn't write
it... right? Apply-then-persist.

### 0.4: Multiprocess
Split the backend (IRC connections) and the terminal interface; make the backend a separate
process. Allow detaching and reattaching of interface processes.

Launching new processes is still the responsibility of the end user; a client
won't automatically launch the backend.

(I'm actually thinking that there only be one binary, just launched with
different args; that allows the client to always use `os.Args[0]` rather than
doing some amount of `PATH`ing.)

### 0.5: Multiwindow
Allow the client to automatically launch the backend in certain circumstances.
Support the client launching new windows in the supported WMs:

* tmux
* screen
* some number of terminal emulators, e.g. xterm.

### 0.6: Fit & Finish
Support more IRC operations. Use `discoirc` daily. Clean TODOs, HACKs, etc.

Maybe: Add UI for "channel meta", e.g. user list, and "connection meta" inasmuch
as that's not already covered by whatever the connection editor looks like.

### 0.7: Private release
Get some folks interested in it to try it out. Get feedback. Incorporate
feedback.

Jennifer (coworker) has volunteered.

### 0.8: Packaging and OSSing
Open-source, and package for some platforms (at least Debian, Ubuntu).

### 0.9: Real Good Docs
Try to do this as we go along, but take a specific point release to review and
polish docs.

### 1.0: OK!
Call it.


## Useful libraries

* Config & scripting Interfaces
  * https://github.com/golang/protobuf for internal data structures & config.
  * https://github.com/Shopify/go-lua for (Lua) scripting
  * Or maybe other things at https://github.com/avelino/awesome-go#configuration
* IRC
  * https://github.com/fluffle/goirc for IRC connections.
* UI
  * https://github.com/jroimartin/gocui for discoirc UI.
  * https://github.com/gizak/termui isn't particuarly useful, but it looks cool.
* Build
  * https://glide.sh/ maybe? I do like having dependencies at known versions.
* Model after
  * https://github.com/cantino/huginn for plugins (maybe).

## Difficult things to do
* Unicode support, overall.
* RTL support in particular.
