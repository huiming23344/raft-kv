package kv_raft

import (
	"fmt"
	"github.com/luo/kv-raft/cmd"
	"github.com/luo/kv-raft/config"
	"github.com/luo/kv-raft/engines"
	"github.com/luo/kv-raft/network"
	"io"
	"log"
	"net"
)

type KvsServer struct {
	addr string
	db   engines.KvsEngine
	raft *raft.Node
}

func NewKvsServer() *KvsServer {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("load config fail: " + err.Error())
	}
	config.SetGlobalConfig(cfg)
	engine, err := engines.NewKvsStore(cfg.Server.DataDir)
	if err != nil {
		log.Fatal(err)
	}
	raftNode, err := raft.NewRaftNode(engine)
	if err != nil {
		log.Fatal(err)
	}
	return &KvsServer{
		addr: cfg.Server.Addr,
		db:   engine,
		raft: raftNode,
	}
}

func (s *KvsServer) Serve() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", s.addr)
	if err != nil {
		return err
	}
	l, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			return err
		}
		handler := Handler{
			db:         s.db,
			connection: network.NewConnection(conn),
			raft:       s.raft,
		}
		go handler.run()
	}
}

type Handler struct {
	db         engines.KvsEngine
	connection network.Connection
	raft       *raft.Node
}

func (h *Handler) run() {
	for {
		// 1.读取一个 Frame
		frame, err := h.connection.ReadFrame()
		if err == io.EOF {
			return
		}
		if err != nil {
			// 网络读取 Frame 失败或无效协议无法解析，终止连接
			fmt.Printf("connection terminate %s, read frame error: %v\n", h.connection.RemoteAddr(), err)
			return
		}
		// 2.转换一个 Frame 为 command 结构
		command, err := cmd.FromFrame(frame)
		if err != nil {
			// 解析 Frame 是不支持的命令，终止连接
			fmt.Printf("connection terminate %s, parse command error: %v\n", h.connection.RemoteAddr(), err)
			return
		}
		var rspFrame *network.Frame
		switch command.Name() {
		case cmd.GET:
			rspFrame = command.Apply(h.db)
		case cmd.SET, cmd.DELETE:
			rspFrame = h.raft.Apply(frame)
		case cmd.MEMBER:
			rspFrame = h.raft.Member(command.(*cmd.Member))
		}
		// 3.回包
		if err := h.connection.WriteFrame(rspFrame); err != nil {
			fmt.Printf("connection terminate %s, write frame error: %v\n", h.connection.RemoteAddr(), err)
			return
		}
	}
}
