package cmd

import (
	"errors"
	"github.com/huiming23344/kv-raft/db/engines"
	kvsError "github.com/huiming23344/kv-raft/errors"
	"github.com/huiming23344/kv-raft/network"
	"log"
)

type Delete struct {
	key string
}

func NewDelete(key string) Command {
	return &Delete{
		key,
	}
}

// 从接收的Frame中解析一个 Delete 命令
func parseDeleteFrame(parse *network.Parse) (Command, error) {
	key, err := parse.NextString()
	if err != nil {
		return nil, err
	}
	cmd := &Delete{
		key,
	}
	return cmd, nil
}

func (c *Delete) Apply(db engines.KvsEngine) *network.Frame {
	log.Printf("Remove `%s` from node value", c.key)
	rsp := new(network.Frame)
	err := db.Remove(c.key)
	if err != nil {
		if errors.Is(err, kvsError.KeyNotFound) {
			rsp.Ftype = network.Integer
			rsp.Value = 0
		} else {
			rsp.Ftype = network.Error
			rsp.Value = err.Error()
		}
	} else {
		rsp.Ftype = network.Integer
		rsp.Value = 1
	}
	return rsp
}

func (c *Delete) IntoFrame() *network.Frame {
	array := []*network.Frame{
		{
			Ftype: network.Bulk,
			Value: DELETE,
		},
		{
			Ftype: network.Bulk,
			Value: c.key,
		},
	}
	return &network.Frame{
		Ftype: network.Array,
		Value: array,
	}
}

func (c *Delete) Name() string {
	return DELETE
}
