package cmd

import (
	"errors"
	"github.com/huiming23344/kv-raft/db/engines"
	kvsError "github.com/huiming23344/kv-raft/errors"
	"github.com/huiming23344/kv-raft/network"
	"log"
)

type Get struct {
	// the lockup key
	key string
}

func NewGet(key string) Command {
	return &Get{
		key,
	}
}

// 从接收的Frame中解析一个 Get 命令
func parseGetFrame(parse *network.Parse) (Command, error) {
	key, err := parse.NextString()
	if err != nil {
		return nil, err
	}
	cmd := &Get{
		key,
	}
	return cmd, nil
}

func (c *Get) Apply(db engines.KvsEngine) *network.Frame {
	log.Printf("Get the value of `%s` from node value\n", c.key)
	rsp := new(network.Frame)
	value, err := db.Get(c.key)
	if err != nil {
		if errors.Is(err, kvsError.KeyNotFound) {
			rsp.Ftype = network.Null
		} else {
			rsp.Ftype = network.Error
			rsp.Value = err.Error()
		}
	} else {
		rsp.Ftype = network.Bulk
		rsp.Value = value
	}
	return rsp
}

func (c *Get) IntoFrame() *network.Frame {
	array := []*network.Frame{
		{
			Ftype: network.Bulk,
			Value: GET,
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

func (c *Get) Name() string {
	return GET
}
