package client

import (
	"errors"
	"github.com/huiming23344/kv-raft/cmd"
	"github.com/huiming23344/kv-raft/network"
	"net"
	"strconv"
)

type Client struct {
	connnection network.Connection
}

func NewClient(addr string) (*Client, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	client := &Client{
		connnection: network.NewConnection(conn),
	}
	return client, nil
}

func (c *Client) Member(opt, serverID, address string) (string, error) {
	frame := cmd.NewMember(opt, serverID, address).IntoFrame()
	rsp, err := c.Invoke(frame)
	if err != nil {
		return "", err
	}
	switch rsp.Ftype {
	case network.Simple:
		return rsp.Value.(string), nil
	case network.Error:
		return rsp.Value.(string), nil
	default:
		return "", errors.New("protocol error; expected simple frame or error frame")
	}
}

func (c *Client) Set(key, value string) (string, error) {
	frame := cmd.NewSet(key, value).IntoFrame()
	rsp, err := c.Invoke(frame)
	if err != nil {
		return "", err
	}
	switch rsp.Ftype {
	case network.Simple:
		return rsp.Value.(string), nil
	case network.Error:
		return rsp.Value.(string), nil
	default:
		return "", errors.New("protocol error; expected simple frame or error frame")
	}
}

func (c *Client) Get(key string) (string, error) {
	frame := cmd.NewGet(key).IntoFrame()
	rsp, err := c.Invoke(frame)
	if err != nil {
		return "", err
	}
	switch rsp.Ftype {
	case network.Bulk:
		return rsp.Value.(string), nil
	case network.Null:
		return "null", nil
	case network.Error:
		return rsp.Value.(string), nil
	default:
		return "", errors.New("protocol error; expected simple frame or error frame")
	}
}

func (c *Client) Del(key string) (string, error) {
	frame := cmd.NewDelete(key).IntoFrame()
	rsp, err := c.Invoke(frame)
	if err != nil {
		return "", err
	}
	switch rsp.Ftype {
	case network.Integer:
		return strconv.FormatInt(int64(rsp.Value.(int)), 10), nil
	case network.Error:
		return rsp.Value.(string), nil
	default:
		return "", errors.New("protocol error; expected simple frame or error frame")
	}
}

func (c *Client) readResponse() (*network.Frame, error) {
	frame, err := c.connnection.ReadFrame()
	if err != nil {
		return nil, err
	}
	return frame, nil
}

func (c *Client) Invoke(frame *network.Frame) (*network.Frame, error) {
	if err := c.connnection.WriteFrame(frame); err != nil {
		return nil, err
	}
	return c.readResponse()
}
