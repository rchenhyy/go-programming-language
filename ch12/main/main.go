package main

import (
	"fmt"
	"github.com/rchenhyy/go-programming-language/ch12"
	"os"
	"reflect"
)

func main() {
	// demo
	ch12.Types()
	ch12.Values()
	ch12.CanAddr()
	ch12.Addr()
	ch12.Addr2()
	ch12.CanSet()

	// format
	fmt.Println(ch12.Any("hello world"))
	fmt.Println(fmt.Sprintf("hello world"))

	var invalid reflect.Value
	fmt.Println(invalid.Kind())

	// display
	type person struct {
		name string
		age  uint
		sex  string
	}
	type example struct {
		x     person
		y     *int
		num   int
		itf   interface{}
		value reflect.Value
		arr   [3]string
		m     map[string]interface{}
	}
	p := person{"陈帅", 18, "男"}
	e := example{
		x:     p,
		y:     nil,
		num:   9,
		itf:   nil,
		value: reflect.ValueOf(os.Stderr),
		arr:   [3]string{"1", "2", "-"},
		m:     map[string]interface{}{"score": 99.0, "comment": "good", "extra": nil},
	}
	ch12.Display("example", e)

	// encode
	ps := []person{p, p, p}
	bytes, err := ch12.Marshal(ps)
	fmt.Printf("marshal: obj=%v, bytes=%s\n", ps, bytes)
	fmt.Printf("marshal: err=%v\n", err)

	// decode
	// ...
	err = ch12.Unmarshal(bytes, &ps)
	fmt.Printf("unmarshal: bytes=%s, obj=%v\n", bytes, ps)
	fmt.Printf("unmarshal: err=%v\n", err)
}
