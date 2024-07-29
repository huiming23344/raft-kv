package cmd

import (
	"github.com/huiming23344/kv-raft/network"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SetFrame(t *testing.T) {
	Convey("test SET frame", t, func() {
		array := []*network.Frame{
			{
				Ftype: network.Bulk,
				Value: "set",
			},
			{
				Ftype: network.Bulk,
				Value: "name",
			},
			{
				Ftype: network.Bulk,
				Value: "mars",
			},
		}
		frame := network.Frame{
			Ftype: network.Array,
			Value: array,
		}
		command, err := FromFrame(&frame)
		if err != nil {
			t.Fatal(err)
		}
		So(command.Name(), ShouldEqual, SET)
	})
}

func Test_GetFrame(t *testing.T) {
	Convey("test GET frame", t, func() {
		array := []*network.Frame{
			{
				Ftype: network.Bulk,
				Value: "get",
			},
			{
				Ftype: network.Bulk,
				Value: "name",
			},
		}
		frame := network.Frame{
			Ftype: network.Array,
			Value: array,
		}
		command, err := FromFrame(&frame)
		if err != nil {
			t.Fatal(err)
		}
		So(command.Name(), ShouldEqual, GET)
	})
}

func Test_DeleteFrame(t *testing.T) {
	Convey("test DELETE frame", t, func() {
		array := []*network.Frame{
			{
				Ftype: network.Bulk,
				Value: "del",
			},
			{
				Ftype: network.Bulk,
				Value: "name",
			},
		}
		frame := network.Frame{
			Ftype: network.Array,
			Value: array,
		}
		command, err := FromFrame(&frame)
		if err != nil {
			t.Fatal(err)
		}
		So(command.Name(), ShouldEqual, DELETE)
	})
}
