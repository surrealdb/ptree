// Copyright © SurrealDB Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ptree

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var s = []string{
	"/some",                 // 0
	"/test",                 // 1
	"/test/one",             // 2
	"/test/one/sub-one",     // 3
	"/test/one/sub-one/1st", // 4
	"/test/one/sub-one/2nd", // 5
	"/test/one/sub-two",     // 6
	"/test/one/sub-two/1st", // 7
	"/test/one/sub-two/2nd", // 8
	"/test/one/sub-zen",     // 9
	"/test/one/sub-zen/1st", // 10 ----------
	"/test/one/sub-zen/2nd", // 11
	"/test/two",             // 12
	"/test/two/sub-one",     // 13
	"/test/two/sub-one/1st", // 14
	"/test/two/sub-one/2nd", // 15
	"/test/two/sub-two",     // 16
	"/test/two/sub-two/1st", // 17
	"/test/two/sub-two/2nd", // 18
	"/test/two/sub-zen",     // 19
	"/test/two/sub-zen/1st", // 20 ----------
	"/test/two/sub-zen/2nd", // 21
	"/test/zen",             // 22
	"/test/zen/sub-one",     // 23
	"/test/zen/sub-one/1st", // 24
	"/test/zen/sub-one/2nd", // 25
	"/test/zen/sub-two",     // 26
	"/test/zen/sub-two/1st", // 27
	"/test/zen/sub-two/2nd", // 28
	"/test/zen/sub-zen",     // 29
	"/test/zen/sub-zen/1st", // 30
	"/test/zen/sub-zen/2nd", // 31
	"/zoo",                  // 32
	"/zoo/some",             // 33
	"/zoo/some/path",        // 34
}

var p = [][]int{
	{0, 0},                   // 0
	{0, 1},                   // 1
	{0, 1, 0, 0},             // 2
	{0, 1, 0, 0, 0, 0},       // 3
	{0, 1, 0, 0, 0, 0, 0, 0}, // 4
	{0, 1, 0, 0, 0, 0, 0, 1}, // 5
	{0, 1, 0, 0, 0, 1},       // 6
	{0, 1, 0, 0, 0, 1, 0, 0}, // 7
	{0, 1, 0, 0, 0, 1, 0, 1}, // 8
	{0, 1, 0, 0, 0, 2},       // 9
	{0, 1, 0, 0, 0, 2, 0, 0}, // 10 ----------
	{0, 1, 0, 0, 0, 2, 0, 1}, // 11
	{0, 1, 0, 1},             // 12
	{0, 1, 0, 1, 0, 0},       // 13
	{0, 1, 0, 1, 0, 0, 0, 0}, // 14
	{0, 1, 0, 1, 0, 0, 0, 1}, // 15
	{0, 1, 0, 1, 0, 1},       // 16
	{0, 1, 0, 1, 0, 1, 0, 0}, // 17
	{0, 1, 0, 1, 0, 1, 0, 1}, // 18
	{0, 1, 0, 1, 0, 2},       // 19
	{0, 1, 0, 1, 0, 2, 0, 0}, // 20 ----------
	{0, 1, 0, 1, 0, 2, 0, 1}, // 21
	{0, 1, 0, 2},             // 22
	{0, 1, 0, 2, 0, 0},       // 23
	{0, 1, 0, 2, 0, 0, 0, 0}, // 24
	{0, 1, 0, 2, 0, 0, 0, 1}, // 25
	{0, 1, 0, 2, 0, 1},       // 26
	{0, 1, 0, 2, 0, 1, 0, 0}, // 27
	{0, 1, 0, 2, 0, 1, 0, 1}, // 28
	{0, 1, 0, 2, 0, 2},       // 29
	{0, 1, 0, 2, 0, 2, 0, 0}, // 30
	{0, 1, 0, 2, 0, 2, 0, 1}, // 31
	{0, 2},       // 32
	{0, 2, 0},    // 33
	{0, 2, 0, 0}, // 34
}

func TestBasic(t *testing.T) {

	p := New()

	c := p.Copy()

	Convey("Get initial size", t, func() {
		So(p.Size(), ShouldEqual, 0)
	})

	Convey("Can insert 1st item", t, func() {
		val := c.Put([]byte("/foo"), []byte("FOO"))
		So(val, ShouldBeNil)
		So(c.Size(), ShouldEqual, 1)
		So(c.Get([]byte("/foo")), ShouldResemble, []byte("FOO"))
	})

	Convey("Can insert 2nd item", t, func() {
		val := c.Put([]byte("/bar"), []byte("BAR"))
		So(val, ShouldBeNil)
		So(c.Size(), ShouldEqual, 2)
		So(c.Get([]byte("/bar")), ShouldResemble, []byte("BAR"))
	})

	Convey("Can get nil item", t, func() {
		val := c.Get([]byte("/"))
		So(val, ShouldEqual, nil)
	})

	Convey("Can delete nil item", t, func() {
		val := c.Del([]byte("/foobar"))
		So(val, ShouldEqual, nil)
		So(c.Size(), ShouldEqual, 2)
		So(c.Get([]byte("/foobar")), ShouldEqual, nil)
	})

	Convey("Can delete 1st item", t, func() {
		val := c.Del([]byte("/foo"))
		So(val, ShouldResemble, []byte("FOO"))
		So(c.Size(), ShouldEqual, 1)
		So(c.Get([]byte("/foo")), ShouldEqual, nil)
	})

	Convey("Can delete 2nd item", t, func() {
		val := c.Del([]byte("/bar"))
		So(val, ShouldResemble, []byte("BAR"))
		So(c.Size(), ShouldEqual, 0)
		So(c.Get([]byte("/bar")), ShouldEqual, nil)
	})

	Convey("Can commit transaction", t, func() {
		n := c.Tree()
		So(n, ShouldNotBeNil)
		So(n.Size(), ShouldEqual, 0)
	})

}

func TestComplex(t *testing.T) {

	p := New()
	c := p.Copy()

	Convey("Can get empty `min`", t, func() {
		r := c.Root()
		k, v := r.Min()
		So(k, ShouldBeNil)
		So(v, ShouldBeNil)
	})

	Convey("Can get empty `max`", t, func() {
		r := c.Root()
		k, v := r.Max()
		So(k, ShouldBeNil)
		So(v, ShouldBeNil)
	})

	Convey("Can insert tree items", t, func() {
		for _, v := range s {
			c.Put([]byte(v), []byte(v))
		}
		So(c.Size(), ShouldEqual, 35)
		for i := len(s) - 1; i > 0; i-- {
			c.Put([]byte(s[i]), []byte(s[i]))
		}
		So(c.Size(), ShouldEqual, 35)
	})

	Convey("Can get proper `min`", t, func() {
		k, v := c.Root().Min()
		So(k, ShouldResemble, []byte("/some"))
		So(v, ShouldResemble, []byte("/some"))
	})

	Convey("Can get proper `max`", t, func() {
		k, v := c.Root().Max()
		So(k, ShouldResemble, []byte("/zoo/some/path"))
		So(v, ShouldResemble, []byte("/zoo/some/path"))
	})

	// ------------------------------------------------------------

	Convey("Can iterate tree items at `nil` with `walk`", t, func() {
		i := 0
		c.Root().Walk(nil, func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 35)
	})

	Convey("Can iterate tree items at `/test/zen/s` with `walk`", t, func() {
		i := 0
		c.Root().Walk([]byte("/test/zen/s"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 9)
	})

	Convey("Can iterate tree items at `/test/zen/sub` with `walk`", t, func() {
		i := 0
		c.Root().Walk([]byte("/test/zen/sub"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 9)
	})

	Convey("Can iterate tree items at `/test/zen/sub-o` with `walk`", t, func() {
		i := 0
		c.Root().Walk([]byte("/test/zen/sub-o"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 3)
	})

	Convey("Can iterate tree items at `/test/zen/sub-one` with `walk`", t, func() {
		i := 0
		c.Root().Walk([]byte("/test/zen/sub-one"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 3)
	})

	Convey("Can iterate tree items at `/test/zen/sub` with `walk` and exit", t, func() {
		i := 0
		c.Root().Walk([]byte("/test/zen/sub"), func(k []byte, v interface{}) (e bool) {
			i++
			return true
		})
		So(i, ShouldEqual, 1)
	})

	// ------------------------------------------------------------

	Convey("Can iterate tree items at `/test/` with `subs`", t, func() {
		i := 0
		c.Root().Subs([]byte("/test/"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 3)
	})

	Convey("Can iterate tree items at `/test/zen/s` with `subs`", t, func() {
		i := 0
		c.Root().Subs([]byte("/test/zen/s"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 3)
	})

	Convey("Can iterate tree items at `/test/zen/sub` with `subs`", t, func() {
		i := 0
		c.Root().Subs([]byte("/test/zen/sub"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 3)
	})

	Convey("Can iterate tree items at `/test/zen/sub-o` with `subs`", t, func() {
		i := 0
		c.Root().Subs([]byte("/test/zen/sub-t"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 1)
	})

	Convey("Can iterate tree items at `/test/zen/sub-one` with `subs`", t, func() {
		i := 0
		c.Root().Subs([]byte("/test/zen/sub-one"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 2)
	})

	Convey("Can iterate tree items at `/test/zen/sub` with `subs` and exit", t, func() {
		i := 0
		c.Root().Subs([]byte("/test/zen/sub"), func(k []byte, v interface{}) (e bool) {
			i++
			return true
		})
		So(i, ShouldEqual, 1)
	})

	// ------------------------------------------------------------

	Convey("Can iterate tree items at `nil` with `path`", t, func() {
		i := 0
		c.Root().Path(nil, func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 0)
	})

	Convey("Can iterate tree items at `/test/zen/s` with `path`", t, func() {
		i := 0
		c.Root().Path([]byte("/test/zen/s"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 2)
	})

	Convey("Can iterate tree items at `/test/zen/sub` with `path`", t, func() {
		i := 0
		c.Root().Path([]byte("/test/zen/sub"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 2)
	})

	Convey("Can iterate tree items at `/test/zen/sub-o` with `path`", t, func() {
		i := 0
		c.Root().Path([]byte("/test/zen/sub-o"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 2)
	})

	Convey("Can iterate tree items at `/test/zen/sub-one` with `path`", t, func() {
		i := 0
		c.Root().Path([]byte("/test/zen/sub-one"), func(k []byte, v interface{}) (e bool) {
			i++
			return
		})
		So(i, ShouldEqual, 3)
	})

	Convey("Can iterate tree items at `/test/zen/sub` with `path` and exit", t, func() {
		i := 0
		c.Root().Path([]byte("/test/zen/sub"), func(k []byte, v interface{}) (e bool) {
			i++
			return true
		})
		So(i, ShouldEqual, 1)
	})

}

func TestIritate(t *testing.T) {

	c := New().Copy()

	i := c.Cursor()

	Convey("Can iterate to the min with no items", t, func() {
		k, v := i.First()
		So(v, ShouldBeNil)
		So(k, ShouldBeNil)
	})

	Convey("Can iterate to the max with no items", t, func() {
		k, v := i.Last()
		So(v, ShouldBeNil)
		So(k, ShouldBeNil)
	})

	Convey("Can seek to a key with no items", t, func() {
		k, v := i.Seek([]byte(""))
		So(v, ShouldBeNil)
		So(k, ShouldBeNil)
	})

	Convey("Can seek to a key with no items", t, func() {
		k, v := i.Seek([]byte("/something"))
		So(v, ShouldBeNil)
		So(k, ShouldBeNil)
	})

}

func TestIterate(t *testing.T) {

	c := New().Copy()

	Convey("Can insert tree items", t, func() {
		for _, v := range s {
			c.Put([]byte(v), []byte(v))
		}
		So(c.Size(), ShouldEqual, 35)
	})

	i := c.Cursor()

	Convey("Can get iterator", t, func() {
		So(i, ShouldNotBeNil)
	})

	Convey("Prev with no seek returns nil", t, func() {
		k, v := i.Prev()
		So(k, ShouldBeNil)
		So(v, ShouldBeNil)
	})

	Convey("Next with no seek returns nil", t, func() {
		k, v := i.Next()
		So(k, ShouldBeNil)
		So(v, ShouldBeNil)
	})

	Convey("Can iterate to the min", t, func() {
		k, v := i.First()
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[0]))
		So(k, ShouldResemble, []byte(s[0]))
		So(v, ShouldResemble, []byte(s[0]))
	})

	Convey("Can iterate using `next`", t, func() {
		for j := 1; j < len(s); j++ {
			k, v := i.Next()
			var t []int
			for _, q := range i.path {
				t = append(t, q.pos)
			}
			So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[j]))
			So(k, ShouldResemble, []byte(s[j]))
			So(v, ShouldResemble, []byte(s[j]))
		}
	})

	Convey("Next item is nil and doesn't change cursor", t, func() {
		k, v := i.Next()
		So(k, ShouldBeNil)
		So(v, ShouldBeNil)
	})

	Convey("Can iterate to the max", t, func() {
		k, v := i.Last()
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[len(p)-1]))
		So(k, ShouldResemble, []byte(s[len(p)-1]))
		So(v, ShouldResemble, []byte(s[len(p)-1]))
	})

	Convey("Can iterate using `prev`", t, func() {
		for j := len(s) - 2; j >= 0; j-- {
			k, v := i.Prev()
			var t []int
			for _, q := range i.path {
				t = append(t, q.pos)
			}
			So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[j]))
			So(k, ShouldResemble, []byte(s[j]))
			So(v, ShouldResemble, []byte(s[j]))
		}
	})

	Convey("Prev item is nil and doesn't change cursor", t, func() {
		k, v := i.Prev()
		So(k, ShouldBeNil)
		So(v, ShouldBeNil)
	})

	Convey("Seek nonexistant nil", t, func() {
		k, v := i.Seek(nil)
		So(k, ShouldResemble, []byte(s[0]))
		So(v, ShouldResemble, []byte(s[0]))
	})

	Convey("Seek nonexistant first byte", t, func() {
		k, v := i.Seek([]byte{0})
		So(k, ShouldResemble, []byte(s[0]))
		So(v, ShouldResemble, []byte(s[0]))
	})

	Convey("Seek nonexistant first item", t, func() {
		k, v := i.Seek([]byte("/aaa"))
		So(k, ShouldResemble, []byte(s[0]))
		So(v, ShouldResemble, []byte(s[0]))
	})

	Convey("Seek just over last item", t, func() {
		k, v := i.Seek([]byte("/zoo/some/path/-"))
		So(v, ShouldBeNil)
		So(k, ShouldBeNil)
	})

	Convey("Seek nonexistant last item", t, func() {
		k, v := i.Seek([]byte("/zzz"))
		So(v, ShouldBeNil)
		So(k, ShouldBeNil)
	})

	Convey("Seek nonexistant last byte", t, func() {
		k, v := i.Seek([]byte{255})
		So(v, ShouldBeNil)
		So(k, ShouldBeNil)
	})

	Convey("Seek half item is correct", t, func() {
		k, v := i.Seek([]byte(s[10][:len(s[10])-3]))
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[10]))
		So(k, ShouldResemble, []byte(s[10]))
		So(v, ShouldResemble, []byte(s[10]))
	})

	Convey("Seek full item is correct", t, func() {
		k, v := i.Seek([]byte(s[10]))
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[10]))
		So(k, ShouldResemble, []byte(s[10]))
		So(v, ShouldResemble, []byte(s[10]))
	})

	Convey("Seek next item is correct", t, func() {
		k, v := i.Seek([]byte("/test/one/sub-zen/0th"))
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[10]))
		So(k, ShouldResemble, []byte(s[10]))
		So(v, ShouldResemble, []byte(s[10]))
	})

	Convey("Seek next item is correct", t, func() {
		k, v := i.Seek([]byte("/test/one/sub-zen/1zz"))
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[11]))
		So(k, ShouldResemble, []byte(s[11]))
		So(v, ShouldResemble, []byte(s[11]))
	})

	Convey("Seek next item is correct", t, func() {
		k, v := i.Seek([]byte("/test/one/sub-zen/2zz"))
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[12]))
		So(k, ShouldResemble, []byte(s[12]))
		So(v, ShouldResemble, []byte(s[12]))
	})

	Convey("Seek jump item is correct", t, func() {
		k, v := i.Seek([]byte("/test/one/sub-zen/3rd"))
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[12]))
		So(k, ShouldResemble, []byte(s[12]))
		So(v, ShouldResemble, []byte(s[12]))
	})

	Convey("Seek finishing item is correct", t, func() {
		k, v := i.Seek([]byte("/test/zzz"))
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[32]))
		So(k, ShouldResemble, []byte(s[32]))
		So(v, ShouldResemble, []byte(s[32]))
	})

	Convey("Seek finalising item is correct", t, func() {
		k, v := i.Seek([]byte("/zoo/some/xxxx"))
		So(v, ShouldBeNil)
		So(k, ShouldBeNil)
	})

	Convey("Prev item after seek is correct", t, func() {
		i.Seek([]byte(s[10]))
		k, v := i.Prev()
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[9]))
		So(k, ShouldResemble, []byte(s[9]))
		So(v, ShouldResemble, []byte(s[9]))
	})

	Convey("Next item after seek is correct", t, func() {
		i.Seek([]byte(s[10]))
		k, v := i.Next()
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[11]))
		So(k, ShouldResemble, []byte(s[11]))
		So(v, ShouldResemble, []byte(s[11]))
	})

	Convey("FINAL", t, func() {
		i.Seek([]byte(s[10]))
		i.Del()
		k, v := i.Next()
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[11]))
		So(k, ShouldResemble, []byte(s[11]))
		So(v, ShouldResemble, []byte(s[11]))
	})

	Convey("FINAL", t, func() {
		var k []byte
		i.Seek([]byte(s[10]))
		i.Del()
		k, _ = i.Next()
		i.Del()
		k, _ = i.Next()
		i.Del()
		k, _ = i.Next()
		i.Del()
		k, v := i.Next()
		var t []int
		for _, q := range i.path {
			t = append(t, q.pos)
		}
		So(fmt.Sprint(t), ShouldEqual, fmt.Sprint(p[15]))
		So(k, ShouldResemble, []byte(s[15]))
		So(v, ShouldResemble, []byte(s[15]))
	})

}

func TestUpdate(t *testing.T) {

	c := New().Copy()

	Convey("Can insert 1st item", t, func() {
		val := c.Put([]byte("/test"), []byte("ONE"))
		So(val, ShouldBeNil)
		So(val, ShouldEqual, nil)
		So(c.Size(), ShouldEqual, 1)
		So(c.Get([]byte("/test")), ShouldResemble, []byte("ONE"))
	})

	Convey("Can insert 2nd item", t, func() {
		val := c.Put([]byte("/test"), []byte("TWO"))
		So(val, ShouldNotBeNil)
		So(val, ShouldResemble, []byte("ONE"))
		So(c.Size(), ShouldEqual, 1)
		So(c.Get([]byte("/test")), ShouldResemble, []byte("TWO"))
	})

	Convey("Can insert 3rd item", t, func() {
		val := c.Put([]byte("/test"), []byte("TRE"))
		So(val, ShouldNotBeNil)
		So(val, ShouldResemble, []byte("TWO"))
		So(c.Size(), ShouldEqual, 1)
		So(c.Get([]byte("/test")), ShouldResemble, []byte("TRE"))
	})

}

func TestDelete(t *testing.T) {

	c := New().Copy()

	Convey("Can insert 1st item", t, func() {
		val := c.Put([]byte("/test"), []byte("TEST"))
		So(val, ShouldBeNil)
		So(val, ShouldEqual, nil)
		So(c.Size(), ShouldEqual, 1)
		So(c.Get([]byte("/test")), ShouldResemble, []byte("TEST"))
	})

	Convey("Can delete 1st item", t, func() {
		val := c.Del([]byte("/test"))
		So(val, ShouldNotBeNil)
		So(val, ShouldResemble, []byte("TEST"))
		So(c.Size(), ShouldEqual, 0)
		So(c.Get([]byte("/test")), ShouldBeNil)
	})

	Convey("Can delete 1st item", t, func() {
		val := c.Del([]byte("/test"))
		So(val, ShouldBeNil)
		So(val, ShouldEqual, nil)
		So(c.Size(), ShouldEqual, 0)
		So(c.Get([]byte("/test")), ShouldBeNil)
	})

}
