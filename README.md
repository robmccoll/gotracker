tracker for go
==============

Tracker is designed to signal a dynamic number of
go routines to exit cleanly via a channel, boolean, or both.

Why use this over a WaitGroup? Tracker can report its count. Tracker can 
be used easily with go routines that are built around select statements.

Note: since tracker doesn't maintain a mapping of which routine created
each channel, channels are never removed.  Tracker _SHOULD NOT_ be used
with a large number of go routines that will be constantly joining
and leaving.

methods
-------
- Join returns a channel that will fire when the tracker
 wants all go routines to exit. In order for the tracker
 to exit, there must be a matching call to Leave
- Leave tells the tracker that a tracked go routine is exiting
 either because it was asked to do so or simply because it is finished.
- IsRunning returns ture if KillAll() has not been called and
 the tracked goroutines should continue to run.
- Count returns how many go routines are currently being tracked (those
 that have called Join() but not Leave())
- KillAll will signal tracked go routines to quite via the IsRunning() boolean
 and the channel that was returned to the thread.


