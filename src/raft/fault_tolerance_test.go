package raft

import (
	"testing"
)

func TestCheckResponseTermIgnoresStaleReply(t *testing.T) {
	rf := &Raft{
		currentTerm: 5,
		state:       LEADER,
	}
	args := &RequestVoteArgs{BaseRPC: BaseRPC{Term: 3}}
	reply := &RequestVoteReply{BaseRPC: BaseRPC{Term: 4}}

	if rf.checkResponseTerm(args, reply, true) {
		t.Fatal("accepted a response from an earlier term")
	}
	if rf.currentTerm != 5 || rf.state != LEADER {
		t.Fatalf("stale response changed state: term=%d state=%d", rf.currentTerm, rf.state)
	}
}

func TestCheckResponseTermStepsDownForHigherTerm(t *testing.T) {
	rf := &Raft{
		currentTerm: 5,
		state:       LEADER,
	}
	args := &AppendEntriesArgs{BaseRPC: BaseRPC{Term: 5}}
	reply := &AppendEntriesReply{BaseRPC: BaseRPC{Term: 6}}

	if rf.checkResponseTerm(args, reply, false) {
		t.Fatal("accepted a response reporting a newer term")
	}
	if rf.currentTerm != 6 || rf.state != FOLLOWER {
		t.Fatalf("did not step down for newer term: term=%d state=%d", rf.currentTerm, rf.state)
	}
}

func TestInstallSnapshotRetainsMatchingSuffix(t *testing.T) {
	rf := &Raft{
		currentTerm: 3,
		state:       FOLLOWER,
		votedFor:    -1,
		logs: []Entry{
			{Index: 0, Term: 0},
			{Index: 1, Term: 2, Command: "one"},
			{Index: 2, Term: 2, Command: "two"},
			{Index: 3, Term: 3, Command: "three"},
		},
		persister: MakePersister(),
	}
	args := &InstallSnapshotArgs{
		BaseRPC:           BaseRPC{Term: 3},
		LastIncludedIndex: 1,
		LastIncludedTerm:  2,
		Data:              []byte("snapshot"),
		Done:              true,
	}

	rf.InstallSnapshot(args, &InstallSnapshotReply{})

	if len(rf.logs) != 3 || rf.logs[0].Index != 1 || rf.logs[1].Index != 2 || rf.logs[2].Index != 3 {
		t.Fatalf("matching log suffix was not retained: %+v", rf.logs)
	}
	if rf.logs[0].Term != 2 || rf.logs[0].Command != nil {
		t.Fatalf("snapshot boundary entry is incorrect: %+v", rf.logs[0])
	}
}

func TestInstallSnapshotDiscardsConflictingSuffix(t *testing.T) {
	rf := &Raft{
		currentTerm: 3,
		state:       FOLLOWER,
		votedFor:    -1,
		logs: []Entry{
			{Index: 0, Term: 0},
			{Index: 1, Term: 1, Command: "conflicting"},
			{Index: 2, Term: 2, Command: "stale"},
		},
		persister: MakePersister(),
	}
	args := &InstallSnapshotArgs{
		BaseRPC:           BaseRPC{Term: 3},
		LastIncludedIndex: 1,
		LastIncludedTerm:  2,
		Data:              []byte("snapshot"),
		Done:              true,
	}

	rf.InstallSnapshot(args, &InstallSnapshotReply{})

	if len(rf.logs) != 1 || rf.logs[0].Index != 1 || rf.logs[0].Term != 2 {
		t.Fatalf("conflicting suffix was not discarded: %+v", rf.logs)
	}
}
