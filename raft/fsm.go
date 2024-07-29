package raft

import (
	"github.com/hashicorp/raft"
	"github.com/huiming23344/kv-raft/cmd"
	dbs "github.com/huiming23344/kv-raft/db"
	"github.com/huiming23344/kv-raft/network"
	"io"
)

type FSM struct {
	db dbs.DB
}

func NewFSM(db dbs.DB) raft.FSM {
	return &FSM{
		db: db,
	}
}

func (f *FSM) Apply(logEntry *raft.Log) interface{} {
	frame := new(network.Frame)
	frame, err := network.ParseRESP(logEntry.Data)
	if err != nil {
		return &network.Frame{
			Ftype: network.Error,
			Value: err.Error(),
		}
	}
	command, err := cmd.FromFrame(frame)
	if err != nil {
		return &network.Frame{
			Ftype: network.Error,
			Value: err.Error(),
		}
	}
	return command.Apply(f.db)
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	// nothing to do
	return nil, nil
}

func (f *FSM) Restore(io.ReadCloser) error {
	// nothing to do
	return nil
}
