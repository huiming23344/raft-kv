package cmd

import (
	"fmt"
	"github.com/luo/kv-raft/db/engines"
	"github.com/luo/kv-raft/network"
)

const (
	SET    = "set"
	GET    = "get"
	DELETE = "del"
	MEMBER = "member"
)

const (
	MemberAdd    = "add"
	MemberRemove = "remove"
	MemberList   = "list"
)

type Command interface {
	Apply(db engines.KvsEngine) *network.Frame
	IntoFrame() *network.Frame
	Name() string
}

// FromFrame 从接收的Frame中解析出command
func FromFrame(frame *network.Frame) (Command, error) {
	parse, err := network.NewParse(frame)
	if err != nil {
		return nil, err
	}
	var cmd Command
	commandName, err := parse.NextString()
	switch commandName {
	case SET:
		cmd, err = parseSetFrame(parse)
	case GET:
		cmd, err = parseGetFrame(parse)
	case DELETE:
		cmd, err = parseDeleteFrame(parse)
	case MEMBER:
		cmd, err = parseMemberFrame(parse)
	default:
		err = fmt.Errorf("unknown command %s", commandName)
	}
	if err == nil {
		err = parse.Finish()
	}
	return cmd, err
}
