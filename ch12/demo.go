package ch12

import (
	"fmt"
	"os"
	"reflect"
)

func Types() {
	t := reflect.TypeOf(3) // a reflect.Type
	fmt.Print(t)           // "int"
}

func Values() {
	v := reflect.ValueOf(3) // a reflect.Value
	fmt.Println(v)          // "3"
	fmt.Printf("%v\n", v)   // "3"
	fmt.Println(v.String()) // 注意："<int Value>"

	t := v.Type()  // a reflect.Type
	fmt.Println(t) // "int"

	x := v.Interface()    // an interface{}
	i := x.(int)          // an int
	fmt.Printf("%d\n", i) // "3"
}

func CanAddr() {
	x := 2
	a := reflect.ValueOf(2)  // 2 int
	b := reflect.ValueOf(x)  // 2 int
	c := reflect.ValueOf(&x) // &x *int
	d := c.Elem()            // 2 int (addressable by &x)

	fmt.Println(a.CanAddr()) // false
	fmt.Println(b.CanAddr()) // false
	fmt.Println(c.CanAddr()) // false
	fmt.Println(d.CanAddr()) // true
}

func Addr() {
	x := 2
	d := reflect.ValueOf(&x).Elem()   // d代表变量x
	px := d.Addr().Interface().(*int) // px := &x
	*px = 2333                        // x= 2333
	fmt.Println(x)
}

func Addr2() {
	x := 2
	d := reflect.ValueOf(&x)   // d代表变量&x
	px := d.Interface().(*int) // px := &x
	*px = 2333                 // x= 2333
	fmt.Println(x)
}

func CanSet() {
	x := 1
	rx := reflect.ValueOf(&x).Elem()
	rx.SetInt(98)
	rx.Set(reflect.ValueOf(99))
	// rx.Set(reflect.ValueOf("abc")) // panic...
	fmt.Println(x)

	stdout := reflect.ValueOf(os.Stdout).Elem() // os.Stdout: *os.File
	fmt.Println(stdout.Type())                  // os.File
	// fd := stdout.FieldByName("fd")
	// fmt.Println(fd.Int())
	// fd.SetInt(2)	// panic...
}
