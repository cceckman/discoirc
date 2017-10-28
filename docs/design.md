# Design

See also the [Roadmap](roadmap.md).

## Processes

I'd like to think of `discoirc` as "communicating sequential processes". So what
does each of those processes do?

* Connection
  1. Wait for start request (may be autoconnect)
  1. Attempt connection
  2. Notify on connect or timeout
  3. Wait for disconnect (network event) or close (request)
  4. Notify on disconnect
* Channel
  1. Wait for connection
  2. Join channel; await sync
  3. Notify UI if applicable
  4. Wait for
    * Disconnect / kick: store and notify UI
    * Update (e.g. message) received: store and notify UI
    * Layout update: Return stored messages
    * Message send: send along connection
* Configuration: Loop
  1. Wait for update (which may be init), from
    * File
    * UI update
  2. Request validation from all Configurables;
  3. Send commit to all Configurables
  4. Commit to file, if autowrite
* UI: Session
  1. Get initial connection/channel state
  2. Await one of:
    * Connection / channel update: Update view
    * Relayout: Relayout
    * User event: Start connection, reconfigure, etc.
* UI: Channel
  1. Wait for Connection, Channel join
  2. Wait for:
    * Relayout, or user scroll:
      1. Request messages from range
      1. Clear, add all those messages to tail
    * Message receive: Add message to tail
    * Message send: Send to channel

## Common components

### Configuration
Various components will need to be configurable. For instance,
UI components will need to be configured with colors and patterns; the daemon
needs to be reconfigured with new connections or channels.

Configuration should go through a validate-and-config cycle.

### Process management

There's two threads of process management: the daemon, and UI terminals.
There's some common things:

* Finding the `discoirc` binary
* Settings up `discoirc` args
  * Socket / daemon launching in particular
* Lifecycle
  * Launch process
  * Await early exit (uh oh, there's an issue!)
  * Or confirm that it's attached to the process

And some uncommon things:

* What's on the other end of the process

But this does suggest some things about the protocol used against the daemon.

### Debug logging
I definitely want to have lots of logs for debugging.

