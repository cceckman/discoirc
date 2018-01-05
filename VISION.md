# `discoirc` goals

Someone once describe the [Go](https://golang.org) programming language to me as "an opinionated language." Similarly, `discoirc` is an opinionanted IRC client.

These opinions are my own; if you disagree with them, you're welcome to use any other client. :-)

## DO..

### ...let me work.
I use IRC, and the terminal, every day at my day job. That's likely to continue
for the indefinite future.

`discoirc` should let the user do basic IRC stuff:

- Join, leave, manage channels and networks
- Filter, search, and prioritize messages (highlights / notifications)
- Extract logs
- Do so securely to the requirement of the network

### ...not get in the way.
`discoirc` should be no more annoying to use than my current `irssi` setup.
That's the low bar - the higher bar is "easy and intuitive to use."

### ...look to the future.
The [IRCv3](https://ircv3.net) working group is working to better standardize,
modernize, and extend the IRC protocol. `discoirc` should support / pursue /
adhere to that, and should support modern technologies pertinent to IRC:

- Basic (not really secure) auth (e.g. passwords)
- Crypto (SSL, SASL) auth
- 🦄Stretch: TCP handoff for transparent restarts
- Unicode-first

It should also integrate, where reasonable, with other chat systems (e.g. Slack
IRC gateway)

## DON'T...

### ...be a window manager.
Lots of terminal programs incorporate a layer of window management
functionality. `vim` has multiple windows in which multiple files, or selections from the same file, can be open. `irssi` has windows and channels. Both let you lay those windows out in various ways.

But this means that most of the time, there's a bunch of layers of window
management open:

1. A GUI window manager
2. A terminal emulator
3. A terminal window manager
4. Whatever program is open

Each of these offers ways to split their view, and to navigate between adjacent views. That's a lot of overhead to switch between one view and the adjacent one- which keybinding do I use? Do I have to move "left" at one layer of management, and then "right" at the next?

There's some ways to reduce this complexity. For instance, I use `xterm` as my
terminal emulator, without messing with tabs or splits or anything like that.
There's also some neat tricks you can play to get `vim` and `tmux` to use the
same bindings for window-switching.

But, if we're already using a higher-level window manager... why not integrate
with it, rather than doing internal window management? Many GUI applications
already take this approach to some extent, e.g. [Pidgin](https://www.pidgin.im/) and [GIMP](https://www.gimp.org/). It's more difficult to do with terminal window managers... but at least there are fewer to integrate with (`tmux` and `screen`).

### ...couple network and UI lifetimes.
It's a core principle in GUI applications to avoid blocking UI threads (often,
"the UI thread", singular) on long events, e.g. network operations. This
desire to not do window management within the application encourages this
behavior on a larger scale - to separate out the *process* that has an IRC
connection from the UI views, which may come and go.

Many IRC users already do this through a slightly different mechanism - starting
an IRC bouncer that connects to networks, and connecting their IRC client to the
bouncer. But this doesn't necessarily result in persistence; you still have to
handle logs, handle highlights of messages when away, etc.

`discoirc`'s goal is to implement this by default- to have persistent
connections and presence by default.

(As a note, I don't use an IRC bouncer myself- I leave `irssi` open in `tmux`,
and [`attach`](https://github.com/cceckman/Tilde/blob/master/scripts/attach)
to the `tmux` session.)

### ...implement unnecessary features.
The Internet has progressed a lot since IRC's heyday. It's [still in use](https://xkcd.com/1782/) - I use it on an everyday basis - but there's a lot of stuff in the protocol, and in clients, that we don't need.

We don't need to send files over IRC. We don't need to have a 'now-playing'
plugin hook in the client. We don't need e.g. URL rewriting within the client
(if it's useful to do, have a bot that does it for *everyone* in the channel.)
We don't need to imlement non-UTF-8 character sets.
