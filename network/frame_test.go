package network

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Simple(t *testing.T) {
	Convey("test parse simple frame", t, func() {
		cursor := newCursor([]byte("+OK\r\n"))
		err := check(&cursor)
		So(err, ShouldBeNil)

		cursor.setPosition(0)
		frame, err := parse(&cursor)
		So(err, ShouldBeNil)
		So(frame, ShouldNotBeNil)
		fmt.Printf("ftype = %d, value = %s\n", frame.Ftype, frame.Value.(string))
	})
}

func Test_Error(t *testing.T) {
	Convey("test parse error frame", t, func() {
		cursor := newCursor([]byte("-Error message\r\n"))
		err := check(&cursor)
		So(err, ShouldBeNil)

		cursor.setPosition(0)
		frame, err := parse(&cursor)
		So(err, ShouldBeNil)
		So(frame, ShouldNotBeNil)
		fmt.Printf("ftype = %d, value = %s\n", frame.Ftype, frame.Value.(string))
	})
}

func Test_Integer(t *testing.T) {
	Convey("test parse integer frame", t, func() {
		cursor := newCursor([]byte(":1000\r\n"))
		err := check(&cursor)
		So(err, ShouldBeNil)

		cursor.setPosition(0)
		frame, err := parse(&cursor)
		So(err, ShouldBeNil)
		So(frame, ShouldNotBeNil)
		fmt.Printf("ftype = %d, value = %d\n", frame.Ftype, frame.Value.(int))
	})
}

func Test_Bulk(t *testing.T) {
	Convey("test parse bulk frame", t, func() {
		cursor := newCursor([]byte("$6\r\nfoobar\r\n"))
		err := check(&cursor)
		So(err, ShouldBeNil)

		cursor.setPosition(0)
		frame, err := parse(&cursor)
		So(err, ShouldBeNil)
		So(frame, ShouldNotBeNil)
		fmt.Printf("ftype = %d, value = %s\n", frame.Ftype, frame.Value.(string))
	})
}

func Test_Array(t *testing.T) {
	Convey("test parse array frame", t, func() {
		cursor := newCursor([]byte("*3\r\n$3\r\nset\r\n$4\r\nname\r\n$4\r\nmars\r\n"))
		err := check(&cursor)
		So(err, ShouldBeNil)

		cursor.setPosition(0)
		frame, err := parse(&cursor)
		So(err, ShouldBeNil)
		So(frame, ShouldNotBeNil)

		fmt.Printf("ftype = %d\n", frame.Ftype)
		for _, frame := range frame.Value.([]*Frame) {
			fmt.Printf("ftype = %d, value = %s\n", frame.Ftype, frame.Value.(string))
		}
	})
}
