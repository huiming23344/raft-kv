package engines

import (
	errs "github.com/huiming23344/kv-raft/errors"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func Test_LsmStore(t *testing.T) {
	Convey("test KvsStore feature", t, func() {
		path, _ := os.Getwd()
		engine := NewLsmEngine(path)
		// 1.日志首次写入
		err := engine.Set("name", "mars")
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
