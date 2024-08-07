package cmd

import (
	"fmt"
	"github.com/huiming23344/kv-raft/db/engines"
	"github.com/huiming23344/kv-raft/network"
)

type Config struct {
	opt    string
	config string
}

var _ Command = (*Config)(nil)

func parseConfigFrame(p *network.Parse) (Command, error) {
	fmt.Println("parseConfigFrame")
	fmt.Printf("p: %v\n", p)
	opt, err := p.NextString()
	if err != nil {
		return nil, err
	}
	config, err := p.NextString()
	if err != nil {
		return nil, err
	}
	cmd := &Config{opt, config}
	return cmd, nil
}

func (c *Config) Apply(db engines.KvsEngine) *network.Frame {
	rsp := new(network.Frame)
	rsp.Ftype = network.Array
	rsp.Value = []*network.Frame{}
	return rsp
}

func (c *Config) IntoFrame() *network.Frame {
	array := []*network.Frame{
		{
			Ftype: network.Bulk,
			Value: DELETE,
		},
		{
			Ftype: network.Bulk,
			Value: "aaa",
		},
	}
	return &network.Frame{
		Ftype: network.Array,
		Value: array,
	}
}

func (c *Config) Name() string {
	return CONFIG
}
