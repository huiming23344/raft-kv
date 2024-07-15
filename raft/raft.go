package raft

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"
	kvscli "github.com/luo/kv-raft/client"
	"github.com/luo/kv-raft/cmd"
	kvscfg "github.com/luo/kv-raft/config"
	"github.com/luo/kv-raft/engines"
	"github.com/luo/kv-raft/network"
	"net"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Node struct {
	raft     *raft.Raft
	serverID raft.ServerID
}

func NewRaftNode(engine engines.KvsEngine) (*Node, error) {
	cfg := kvscfg.GlobalConfig()
	dataDir := cfg.Server.DataDir
	// serverID 对应 kvs 的地址
	serverID := cfg.Server.Addr

	raftConfig := raft.DefaultConfig()
	raftConfig.ProtocolVersion = raft.ProtocolVersionMax
	raftConfig.LocalID = raft.ServerID(serverID)
	leaderNotifyCh := make(chan bool, 1)
	raftConfig.NotifyCh = leaderNotifyCh

	fsm := NewFSM(engine)
	transport, err := newRaftTransport(cfg.Raft.Addr)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}
	// 忽略快照
	snapshotStore := raft.NewDiscardSnapshotStore()
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft-log.bolt"))
	if err != nil {
		return nil, err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft-stable.bolt"))
	if err != nil {
		return nil, err
	}
	raftNode, err := raft.NewRaft(raftConfig, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, err
	}
	if cfg.Raft.Bootstrap {
		cfg := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      raft.ServerID(serverID),
					Address: transport.LocalAddr(),
				},
			},
		}
		raftNode.BootstrapCluster(cfg)
	}
	return &Node{
		raft:     raftNode,
		serverID: raft.ServerID(serverID),
	}, nil
}

func newRaftTransport(addr string) (*raft.NetworkTransport, error) {
	address, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	transport, err := raft.NewTCPTransport(address.String(), address, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}
	return transport, nil
}

// Apply FSM 状态机
func (r *Node) Apply(frame *network.Frame) *network.Frame {
	if !r.isLeader() {
		// 代理调用 Leader
		rspFrame, err := r.proxyInvoke(r.leader(), frame)
		if err != nil {
			rspFrame = &network.Frame{
				Ftype: network.Error,
				Value: err.Error(),
			}
		}
		return rspFrame
	}
	cmdBytes, _ := frame.Bytes()
	ret := r.raft.Apply(cmdBytes, 5*time.Second)
	if ret.Error() != nil {
		return &network.Frame{
			Ftype: network.Error,
			Value: ret.Error().Error(),
		}
	}
	return ret.Response().(*network.Frame)
}

// Member 集群成员
func (r *Node) Member(cm *cmd.Member) *network.Frame {
	rspFrame := &network.Frame{
		Ftype: network.Simple,
		Value: "OK",
	}
	switch cm.Opt() {
	case cmd.MemberAdd:
		if err := r.raft.AddVoter(raft.ServerID(cm.ServerID()), raft.ServerAddress(cm.Address()), 0, 0).Error(); err != nil {
			return &network.Frame{
				Ftype: network.Error,
				Value: err.Error(),
			}
		}
	case cmd.MemberRemove:
		if err := r.raft.RemoveServer(raft.ServerID(cm.ServerID()), 0, 0).Error(); err != nil {
			return &network.Frame{
				Ftype: network.Error,
				Value: err.Error(),
			}
		}
	case cmd.MemberList:
		var buf bytes.Buffer
		for _, s := range r.raft.GetConfiguration().Configuration().Servers {
			buf.WriteString(fmt.Sprintf("id=%s address=%s suffrage=%d isLeader=%t\n", s.ID, s.Address, s.Suffrage, r.checkIsLeader(s.ID)))
		}
		rspFrame.Value = strings.TrimRight(buf.String(), "\n")
	default:
		rspFrame.Value = "unknown member subcommand " + cm.Opt()
	}
	return rspFrame
}

func (r *Node) isLeader() bool {
	_, leaderId := r.raft.LeaderWithID()
	return leaderId == r.serverID
}

func (r *Node) leader() raft.ServerAddress {
	_, serverID := r.raft.LeaderWithID()
	return raft.ServerAddress(serverID)
}

func (r *Node) checkIsLeader(serverID raft.ServerID) bool {
	_, leaderId := r.raft.LeaderWithID()
	return leaderId == serverID
}

func (r *Node) proxyInvoke(addr raft.ServerAddress, frame *network.Frame) (*network.Frame, error) {
	client, err := kvscli.NewClient(string(addr))
	if err != nil {
		return nil, err
	}
	return client.Invoke(frame)
}
