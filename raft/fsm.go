package raft

import (
	"github.com/hashicorp/raft"
	"github.com/luo/kv-raft/cmd"
	"github.com/luo/kv-raft/engines"
	"github.com/luo/kv-raft/network"
	"io"
)

type FSM struct {
	db engines.KvsEngine
}

func NewFSM(engine engines.KvsEngine) raft.FSM {
	return &FSM{
		db: engine,
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
