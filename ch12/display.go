package ch12

import (
	"fmt"
	"reflect"
)

func Display(path string, x interface{}) {
	display(path, reflect.ValueOf(x))
}

func display(path string, v reflect.Value) {
	/*
		注意：
		1）reflect的中间操作都是返回Type/Value类型，除非显示的调用Value.Interface()方法
		2）在Java中，反射操作在单独定义的Class/Field这些工具类中，对象obj作为参数传递，因此看起来和Go不太一样，其实只是封装方式的区别
	*/
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		// 显示所有元素：Len() & Index(i)
		for i := 0; i < v.Len(); i++ {
			display(fmt.Sprintf("%s[%d]", path, i), v.Index(i))
		}
	case reflect.Map:
		// 显示所有键值：MapKeys() & MapIndex(key)
		for _, key := range v.MapKeys() {
			display(fmt.Sprintf("%s[%s]", path, formatAtom(key)), v.MapIndex(key))
		}
	case reflect.Struct:
		// 显示所有字段：NumField() & Field(i) & Type().Field(i).Name
		for i := 0; i < v.NumField(); i++ {
			display(fmt.Sprintf("%s.%s", path, v.Type().Field(i).Name), v.Field(i))
		}
	case reflect.Ptr:
		// 区分nil和非nil，获取有效元素：IsNil()和Elem()
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			display(fmt.Sprintf("(*%s)", path), v.Elem())
		}
	case reflect.Interface:
		// 区分nil和非nil，获取有效元素：IsNil()和Elem()；同时关注类型和值（type, value）
		if v.IsNil() {
			fmt.Printf("%s = nil\n", path)
		} else {
			fmt.Printf("%s.type = %s\n", path, v.Elem().Type())
			display(fmt.Sprintf("%s.value", path), v.Elem())
		}
	default:
		// 基本类型、通道、函数
		fmt.Printf("%s = %s\n", path, formatAtom(v))
	}
}
