# Learning
- Named pipes
  - Single send-and-receive relation. Not very complicated.
  - File can't be opened by writer until receiver is ready.
  - Re-opening by reader is iffy. Doesn't seem to receive until writer exits  
    (and flushes?)
  - Re-opening by writer seems to be OK.
  - Flushing is also iffy. fsync doesn't work (f.Sync()), cancelling the writer
    seems to flush. Not sure where the buffering is.
- Domain sockets
  - Like network sockets, they're connection-oriented; you can know about
    connects and disconnects.
  - Probably a better model overall; also leads to remote-head operation.
  - Can domain sockets be wrapped as a pair of Go channels- one for "send", one
    for "disconnected"? Maybe that isn't the right way anyway.

