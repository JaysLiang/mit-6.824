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

### leader election

- followers keep follower state on condition receiving valid rpc message from candidate or leader
- followers change it follwer state on condition election timeout
- new election, increment current term, follower change it state to candidate, and vote for itself and start RequestVote RPC
- candidate win election become a new leader
- other candidate win election
- timeout no leader elected


### log replication

- leader begin to handle client request. 
- each client request contains a command.
- leader append the command as a new entry
- leader apply new entry to state machine and issue to followers in parallel
- as picture show, each entry contains a term id and unique log index, leader's log always the newest

![image](../../img/logs.png)