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
	"log"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Node struct {
	raft     *raft.Raft
	serverID raft.ServerID
}

func NewRaftNode(engine engines.KvsEngine) (*Node, error) {
	cfg := kvscfg.GlobalConfig()
	dataDir := fmt.Sprintf("./nodes/node0")
	serverID := "0"

	raftConfig := raft.DefaultConfig()
	raftConfig.ProtocolVersion = raft.ProtocolVersionMax
	raftConfig.LocalID = raft.ServerID(serverID)
	leaderNotifyCh := make(chan bool, 1)
	raftConfig.NotifyCh = leaderNotifyCh

	fsm := NewFSM(engine)
	var raftAddr string
	// init raft ip
	if cfg.Raft.UseLoopBack {
		raftAddr = fmt.Sprintf("127.0.0.1:%s", cfg.Raft.Port)
	} else {
		addrs, err := GetHostIPAddresses()
		if err != nil {
			return nil, err
		}
		addr := addrs[len(addrs)-1]
		raftAddr = fmt.Sprintf("%s:%s", addr, cfg.Raft.Port)
	}
	transport, err := newRaftTransport(raftAddr)
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
		timeOut := time.Microsecond * 100
		// add node with loopback addr will get a new raft node
		loRe := regexp.MustCompile(`^127\.0\.0\.1`)
		if loRe.FindString(cm.Address()) != "" {
			dataDir := fmt.Sprintf("./nodes/node%s", cm.ServerID())
			engine, err := engines.NewKvsStore(dataDir)
			var fsm = NewFSM(engine)
			raftConfig := raft.DefaultConfig()
			raftConfig.ProtocolVersion = raft.ProtocolVersionMax
			raftConfig.LocalID = raft.ServerID(cm.ServerID())
			LeaderNotifyCh := make(chan bool, 1)
			raftConfig.NotifyCh = LeaderNotifyCh
			if err != nil {
				log.Fatal(err)
			}
			var transport, _ = newRaftTransport(cm.Address())
			os.MkdirAll(dataDir, 0700)
			// 忽略快照
			snapshotStore := raft.NewDiscardSnapshotStore()
			logStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft-log.bolt"))
			stableStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft-stable.bolt"))

			raftNode, err := raft.NewRaft(raftConfig, fsm, logStore, stableStore, snapshotStore, transport)

			cfg := raft.Configuration{
				Servers: []raft.Server{
					{
						Suffrage: raft.Voter,
						ID:       raft.ServerID(cm.ServerID()),
						Address:  transport.LocalAddr(),
					},
				},
			}
			raftNode.BootstrapCluster(cfg)
		}
		if err := r.raft.AddVoter(raft.ServerID(cm.ServerID()), raft.ServerAddress(cm.Address()), 0, timeOut).Error(); err != nil {
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

func canConnect(address string, timeout time.Duration) (bool, error) {
	// 解析地址
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return false, fmt.Errorf("failed to resolve tcp addr: %v", err)
	}

	// 设置连接超时
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// 尝试建立连接
	conn, err := dialer.Dial("tcp", tcpAddr.String())
	if err != nil {
		// 连接失败
		return false, nil
	}
	// 连接成功，不要忘记关闭连接
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	// 如果到达这里，说明连接成功
	return true, nil
}

func GetHostIPAddresses() ([]string, error) {
	var addresses []string

	// 获取所有网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 遍历网络接口
	for _, i := range interfaces {
		// 获取接口的地址列表
		addrs, err := i.Addrs()
		if err != nil {
			return nil, err
		}

		// 遍历地址列表
		for _, addr := range addrs {
			// 检查是否为IPv4地址
			ip := addr.(*net.IPNet)
			if ip.IP.To4() != nil {
				// 添加IPv4地址到结果列表
				addresses = append(addresses, ip.IP.String())
			}
		}
	}

	return addresses, nil
}
