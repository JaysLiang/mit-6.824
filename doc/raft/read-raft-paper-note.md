# the raft consensus algorithm

this doc is my note to reading the paper [raft-extend.pdf](./raft-extended.pdf)

## main step
- select leader
- leader accepts log entries from clients
- leader decide and tell other replicates to apply to log entries to state machine

## basic conception

### server state
- servers have only three state: leader, follower, candidate
- leader decide apply log entry, candidate redirect request to leader to decide to apply log entry.
- candidate, is used to elect a new leader.

![image](../../img/server-states.png)

### term

- time divided into terms, raft ensure only one leader at in a given term

![image](../../img/terms.png)

- t3 no leader is elected.
- stale term will be rejected

### communicate

- raft server communicate using RPC

####  RPC types

- Request Vote 
- Append Entries
- transfer snapshots
