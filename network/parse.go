package network

import (
	"errors"
	"fmt"
)

type Parse struct {
	parts []*Frame
	index int
}

func NewParse(frame *Frame) (*Parse, error) {
	if frame.Ftype != Array {
		return nil, errors.New(fmt.Sprintf("protocol error; expected array, got %d", frame.Ftype))
	}
	parse := &Parse{
		parts: frame.Value.([]*Frame),
		index: 0,
	}
	return parse, nil
}

func (p *Parse) next() *Frame {
	l := len(p.parts)
	index := p.index
	if index >= l {
		return nil
	}
	frame := p.parts[index]
	p.index += 1
	return frame
}

func (p *Parse) NextString() (string, error) {
	frame := p.next()
	if frame == nil {
		return "", errors.New("end of frame")
	}
	switch frame.Ftype {
	case Simple:
		return frame.Value.(string), nil
	case Bulk:
		return frame.Value.(string), nil
	default:
		return "", errors.New(fmt.Sprintf("protocol error; expected simple frame or bulk frame, got %d", frame.Ftype))
	}
}

func (p *Parse) Finish() error {
	if p.next() == nil {
		return nil
	} else {
		return errors.New("protocol error; expected end of frame, but there was more")
	}
}
