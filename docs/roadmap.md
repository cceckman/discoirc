# Roadmap

This describes the milestones for `discoirc`.

Each of these should also be Github issues, but writing them up here because
it's usable offline.

## 0.X milestones

### 0.1: Min UI
UI panels are in roughly working order, implemented against a stub backend.
Allow switching between them, but not launching windows. Statically
configure channels and connections on the stub server.

- [ ] Channel
  - [ ] New messages incoming
  - [ ] Send messages
  - [ ] List network + channel
  - [ ] List nick
- [ ] Session
  - [ ] Networks
  - [ ] Channels on networks
  - [ ] Connection status on networks

Deferred:

- Adding / editing / closing connections

### 0.2: IRC
Implement the backend interface as a connection to IRC servers. Still
statically configure connections and channels. via code.

- [ ] IRC connection
  - [ ] Join channel
  - [ ] Send messages
  - [ ] Receive channel messages
  - [ ] Receive notices
- [ ] UI integration
  - [ ] Connect / disconnect notification

Many things deferred here; see milestones below.

### 0.3: Configuration
Reconfigure `discoirc` on the fly.

- [ ] File interface
  - [ ] Load initial connections, channels from a file.
  - [ ] Watch file for updates; validate and load.
  - [ ] Save file with current configuration.
  - [ ] Automatically save file when configuration is updated.
- [ ] IRC management
  - [ ] Create new network entries in UI.
  - [ ] Manage connection state in UI.
  - [ ] Manage channel state in UI.

### 0.4: Multiprocess
Split the backend (IRC connections) and the terminal interface; make the backend a separate
process.

- [ ] Process management
  - [ ] Establish socket convention between UI and daemon.
  - [ ] Add process lifecycle management: start process with socket arg, watch
    for connection-or-death with a timeout.
- [ ] Interface
  - [ ] Create RPC interface between UI and daemon.
  - [ ] Pass messages (of various sorts) across it.

Deferred:

- Launching new UI windows (see next milestone)

Note that this doesn't involve creating a new binary; using the same binary
means that installation is just "use this".

### 0.5: Multiwindow
Support the UI launching new windows.


- [ ] Launch handling
  - [ ] Establish argument conventions: jumping straight to a non-default view.
- [ ] Support terminal-WM launches
  - [ ] `tmux`
  - [ ] `screen`

Deferred:

- GUI terminal emulators, e.g. `xterm`

### 0.6: Fit & Finish

- [ ] Support more IRC operations
  - [ ] Autocommand on startup (e.g. NickServ)
  - [ ] Mode rendering
- [ ] Add useful views
  - [ ] Channel meta: user list and modes
- [ ] Update window title
  - [ ] `tmux` escapes
  - [ ] `screen` escapes
  - [ ] `xterm` escapes
- [ ] Use daily.
- [ ] No `TODO`s or `HACK`s in code.

### 0.7: Private release
Get some folks interested in it to try it out. Get feedback. Incorporate
feedback.

- [ ] Seek volunteers
  - [ ] Jennifer (coworker)
  - [ ] @danderson


### 0.8: Packaging and OSSing

- [ ] Open-source
- [ ] Add packages
  - [ ] Debian
  - [ ] Ubuntu (inasmuch as it's different)
  - [ ] Arch

Some distros may be deferred; depends on interest.


### 0.9: Real Good Docs
Try to do this as we go along, but take a specific point release to review and
polish docs.

- [ ] Installation
- [ ] Getting started
- [ ] Full documentation of configuration options

## 1.X Milestones
Once it's usable for day-to-day use, there's a bunch of FRs and dnice things to
have.

These aren't ordered; I'm open to influence on priorities.

### 1.A: Logging
Persist logs to disk. Automaticaly page in logs while scrolling.

(@danderson is particularly interested in this.)

- [ ] Configure log writing
  - [ ] Provide configurable disk / time limits.
- [ ] Log paging

### 1.B: Keybindings
Revise the keybindings. Make the physics of IRC behave like your favorite
editor.

- [ ] Select initial mode via configuration
- [ ] Legacy mode: default mode, initial implementation. All commands are
  `/`-prefixed messages.
- [ ] Vim modes (per @cceckman's preference)
  - [ ] Insert mode: "send", the default mode.
  - [ ] Normal mode: JK scrolling, paging, commands via `:`
- [ ] Emacs mode (per @danderson's request)

### 1.C: Search
Provide searching and filtering functionality.

- [ ] Search: highlight matches, jump to previous / next.
- [ ] Filtering: show only messages that match a pattern
  - [ ] ...as a filter in a given channel
  - [ ] ...as its own view, cross-channel / network

