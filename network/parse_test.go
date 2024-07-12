package network

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_NextString(t *testing.T) {
	Convey("test parse next", t, func() {
		array := []*Frame{
			{
				Ftype: Bulk,
				Value: "get",
			},
			{
				Ftype: Bulk,
				Value: "name",
			},
		}
		frame := Frame{
			Ftype: Array,
			Value: array,
		}
		parse, err := NewParse(&frame)
		if err != nil {
			t.Fatal(err)
		}
		nextString, err := parse.NextString()
		if err != nil {
			t.Fatal(err)
		}
		So(nextString, ShouldEqual, "get")

		nextString, err = parse.NextString()
		if err != nil {
			t.Fatal(err)
		}
		So(nextString, ShouldEqual, "name")
	})
}

func Test_Finish(t *testing.T) {
	Convey("test parse finish", t, func() {
		array := []*Frame{
			{
				Ftype: Bulk,
				Value: "get",
			},
		}
		frame := Frame{
			Ftype: Array,
			Value: array,
		}
		parse, err := NewParse(&frame)
		if err != nil {
			t.Fatal(err)
		}
		_, err = parse.NextString()
		if err != nil {
			t.Fatal(err)
		}
		err = parse.Finish()
		So(err, ShouldBeNil)
	})
}
