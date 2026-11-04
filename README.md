# Introduction

Implementations for Labs 1-5 (including no-credit challenges in Lab 5) of MIT 6.5840 Distributed Systems course, passed provided tests >500 times.

For more information about the course, please visit the [MIT Course Page](https://pdos.csail.mit.edu/6.824/general.html).

A special thanks to [OneSizeFitsQuorum](https://github.com/OneSizeFitsQuorum/MIT6.824-2021), whose repository and ideas greatly contributed to my understanding of the course material.

# Disclaimer

**Please Do not submit this code as your own assignment.**

By using this code, you agree to the following:
- You will not submit any part of this code as your own work in any academic or professional setting.
- You will use this code responsibly and only for learning and understanding the concepts.
- You acknowledge that submitting someone else's work as your own is considered academic dishonesty and can have serious consequences.

# Running the distributed service tests

The replicated key/value service is implemented in `src/kvraft` on top of the Raft consensus implementation in `src/raft`. The lab RPC network and tests run the replicas in-process and simulate crashes, partitions, unreliable messages, restarts, and snapshots; this repository does not currently provide a standalone network daemon.

Install Go, then from the repository root run the distributed-service test suites:

```powershell
Set-Location .\src
go test ./raft ./kvraft ./shardctrler ./shardkv
```

To run focused checks for the Raft safety fixes:

```powershell
go test ./raft -run 'TestCheckResponseTerm|TestInstallSnapshot' -count=1
```

The first command covers Raft replication and snapshots, replicated key/value persistence and linearizability, shard-controller configuration, and sharded key/value migration. The focused command covers stale Raft RPC responses and snapshot suffix handling.

`go test ./...` is not a valid repository-wide check for this multi-lab checkout: the unrelated `main` package imports lab packages that are not present, `mrapps` contains mutually exclusive applications with duplicate symbols, and the `kvsrv` memory stress test exceeded its 120-second limit in the Windows test run. The distributed-service package tests above completed successfully.

## Validation results

- `go test ./raft ./kvraft ./shardctrler ./shardkv` — passed.
- Focused Raft regression tests — all four passed.
- Regression output capture: [`screenshots/raft-regression-tests.png`](screenshots/raft-regression-tests.png).
