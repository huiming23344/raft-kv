package engines

import (
	"fmt"
	errs "github.com/huiming23344/kv-raft/errors"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SortedGenList(t *testing.T) {
	path, _ := os.Getwd()
	genList, err := sortedGenList(path)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(genList)
}

func Test_KvsStore(t *testing.T) {
	Convey("test KvsStore feature", t, func() {
		path, _ := os.Getwd()
		engine, err := NewKvsStore(path)
		if err != nil {
			t.Fatal(err)
		}
		// 1.日志首次写入
		err = engine.Set("name", "mars")
		if err != nil {
			t.Fatal(err)
		}
		val, err := engine.Get("name")
		So(val, ShouldEqual, "mars")
		// 2.日志追加
		err = engine.Set("age", "25")
		if err != nil {
			t.Fatal(err)
		}
		val, err = engine.Get("age")
		if err != nil {
			t.Fatal(err)
		}
		So(val, ShouldEqual, "25")
		// 3.删除索引，日志不动
		err = engine.Remove("name")
		if err != nil {
			t.Fatal(err)
		}
		_, err = engine.Get("name")
		So(err, ShouldResemble, errs.KeyNotFound)
	})
}

func Test_Compact(t *testing.T) {
	Convey("test KvsStore compact log", t, func() {
		path, _ := os.Getwd()
		engine, err := NewKvsStore(path)
		if err != nil {
			t.Fatal(err)
		}
		err = engine.(*KvsStore).compact()
		if err != nil {
			t.Fatal(err)
		}
	})
}
