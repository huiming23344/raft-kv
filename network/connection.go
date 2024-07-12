package network

import (
	"bufio"
	"errors"
	"net"
)

type Connection struct {
	conn   *net.TCPConn
	reader Buffer
	writer *bufio.Writer
}

func NewConnection(conn *net.TCPConn) Connection {
	return Connection{
		conn:   conn,
		reader: newBuffer(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (c *Connection) ReadFrame() (*Frame, error) {
	for {
		// 1.当缓存区有足够正常的数据，则解析一个Frame返回
		frame, ok, err := c.parseFrame()
		if err != nil {
			return nil, err
		}
		if ok {
			return frame, nil
		}
		// 2.读满缓存区
		err = c.reader.readFromReader()
		if err != nil {
			return nil, err
		}
	}
}

func (c *Connection) parseFrame() (*Frame, bool, error) {
	// 1.检查缓冲区数据能够读取一个Frame
	cursor := newCursor(c.reader.chunk())
	err := check(&cursor)
	if err != nil {
		if err == Incomplete {
			// 缓冲区数据不够
			return nil, false, nil
		} else {
			// 编码的Frame无效，终止连接
			return nil, false, errors.New("frame is invalid")
		}
	}
	// 2.读取一个Frame
	length := cursor.position()
	cursor.setPosition(0)
	frame, err := parse(&cursor)
	if err != nil {
		// 编码的Frame无效，终止连接
		return nil, false, errors.New("frame is invalid")
	}
	// 3.移动缓冲区的读取位置
	_ = c.reader.advance(length)
	return frame, true, nil
}

// WriteFrame 写入一个Frame
func (c *Connection) WriteFrame(frame *Frame) error {
	frameBytes, err := frame.Bytes()
	if err != nil {
		return err
	}
	if _, err := c.writer.Write(frameBytes); err != nil {
		return err
	}
	return c.writer.Flush()
}
