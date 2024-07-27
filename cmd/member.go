package cmd

import (
	"github.com/luo/kv-raft/db/engines"
	"github.com/luo/kv-raft/network"
)

type Member struct {
	// the cluster operate command
	opt string

	// the member serverID
	serverID string

	// the member address
	address string
}

func NewMember(opt, serverID, address string) Command {
	return &Member{
		opt, serverID, address,
	}
}

// 将接收到的 Frame 解析为一个 Member 命令
func parseMemberFrame(p *network.Parse) (Command, error) {
	opt, err := p.NextString()
	if err != nil {
		return nil, err
	}
	serverID, err := p.NextString()
	if err != nil {
		return nil, err
	}
	address, err := p.NextString()
	if err != nil {
		return nil, err
	}
	cmd := &Member{
		opt, serverID, address,
	}
	return cmd, nil
}

func (m *Member) Apply(engines.KvsEngine) *network.Frame {
	return &network.Frame{
		Ftype: network.Simple,
		Value: "OK",
	}
}

func (m *Member) IntoFrame() *network.Frame {
	array := []*network.Frame{
		{
			Ftype: network.Bulk,
			Value: MEMBER,
		},
		{
			Ftype: network.Bulk,
			Value: m.opt,
		},
		{
			Ftype: network.Bulk,
			Value: m.serverID,
		},
		{
			Ftype: network.Bulk,
			Value: m.address,
		},
	}
	return &network.Frame{
		Ftype: network.Array,
		Value: array,
	}
}

func (m *Member) Name() string {
	return MEMBER
}

func (m *Member) Opt() string {
	return m.opt
}

func (m *Member) ServerID() string {
	return m.serverID
}

func (m *Member) Address() string {
	return m.address
}
