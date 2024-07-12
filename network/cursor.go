package network

import "errors"

type Cursor struct {
	buf []byte
	pos int
}

func newCursor(buf []byte) Cursor {
	return Cursor{
		buf: buf,
		pos: 0,
	}
}

func (c *Cursor) hasRemaining() bool {
	return c.remaining() > 0
}

func (c *Cursor) getByte() (byte, error) {
	if c.remaining() < 1 {
		return 0, errors.New("not enough data")
	}
	ret := c.chunk()[0]
	_ = c.advance(1)
	return ret, nil
}

func (c *Cursor) position() int {
	return c.pos
}

func (c *Cursor) setPosition(pos int) {
	c.pos = pos
}

func (c *Cursor) remaining() int {
	l := len(c.buf)
	pos := c.pos
	if pos >= l {
		return 0
	}
	return l - pos
}

func (c *Cursor) chunk() []byte {
	l := len(c.buf)
	pos := c.position()
	if pos >= l {
		return make([]byte, 0)
	}
	return c.buf[pos:]
}

func (c *Cursor) advance(cnt int) error {
	pos := c.pos
	pos += cnt
	if pos > len(c.buf) {
		return errors.New("overflow")
	}
	c.setPosition(pos)
	return nil
}
