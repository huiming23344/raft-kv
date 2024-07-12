package network

import (
	"errors"
	"io"
)

type Buffer struct {
	reader io.Reader
	buf    []byte
	start  int
	end    int
}

func newBuffer(reader io.Reader) Buffer {
	return Buffer{
		reader: reader,
		buf:    make([]byte, 4*1024),
		start:  0,
		end:    0,
	}
}

// 返回 start 位置和 end 位置之间的字节数
func (b *Buffer) remaining() int {
	return b.end - b.start
}

// 有效字节前移
func (b *Buffer) grow() {
	if b.start == 0 {
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.end - b.start
	b.start = 0
}

// 从reader中读取字节，如果reader阻塞，发生阻塞
func (b *Buffer) readFromReader() error {
	b.grow()
	n, err := b.reader.Read(b.buf[b.end:])
	if err != nil {
		return err
	}
	b.end += n
	return nil
}

// 将 start 位置向前移动 length 个字节
func (b *Buffer) advance(length int) error {
	l := b.remaining()
	if length > l {
		return errors.New("overflow")
	}
	b.start += length
	return nil
}

// 返回从 start 位置开始的切片，其长度是 remaining() 返回的字节数
func (b *Buffer) chunk() []byte {
	start := b.start
	end := b.end
	if start >= end {
		return make([]byte, 0)
	}
	return b.buf[start:end]
}
