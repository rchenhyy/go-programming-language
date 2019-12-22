## Notes

Go语言提供了一种机制，在编译时不知道类型的情况下，可更新变量、在运行时查看值、调用方法以及直接对它们的布局进行操作，这种机制成为**反射**（reflection）。反射也让我们可以**把类型当作头等值**。

#### 为什么使用反射

有时候我们需要写一个有能力统一处理各种值类型的函数，而这些类型可能*无法共享同一个接口*，也有可能*布局未知*，也有可能这个类型在我们设计函数时*还不存在*，甚至这个类型会同时存在上面三种问题。

#### 反射的基础：reflect.Type和reflect.Value

反射功能由reflect包提供，它定义了两个重要的类型：Type和Value。

1. **reflect.Type**

	reflect.Type接口只有一个实现，即**类型描述符**，接口值中的动态类型也是类型描述符。<br/>
	reflect.TypeOf函数接受任何的interface{}参数，并且把接口中的动态类型以reflect.Type形式返回。<br/>
	因为reflect.Type返回一个接口值对应的动态类型，所以它返回的总是具体类型（而不是接口类型）。

	> 把一个具体值赋给一个接口类型时会发生一个隐式类型转换，转换会生成一个包含两部分内容的**接口值**：**动态类型**部分是操作数的类型，**动态值**部分是操作数的值（interface value = (value, type)）。<br/>

	```
	t := reflect.TypeOf(3) // a reflect.Type
	fmt.Print(t)           // "int"
	```
	
	尽管有无限种类型（Type），但类型的分类（Kind）只有少数几种：
	
	* **基础类型**：bool、string以及各种数字类型
	* **聚合类型**：array和struct
	* **引用类型**：chan、func、ptr、slice和map
	* **接口类型**：interface
	* **无类型**：Invalid类型（reflect.Value的零值）

2. **reflect.Value**

	reflect.Value可以包含一个任意类型的值。<br/>
	reflect.ValueOf函数接受任意的interface{}并将接口的动态值以reflect.Value的形式返回。
		
	```
	v := reflect.ValueOf(3) // a reflect.Value
	fmt.Println(v)          // "3"
	fmt.Printf("%v\n", v)   // "3"
	fmt.Println(v.String()) // 注意："<int Value>"

	t := v.Type()  // a reflect.Type
	fmt.Println(t) // "int"

	x := v.Interface()    // an interface{}
	i := x.(int)          // an int
	fmt.Printf("%d\n", i) // "3"
	```
	
	reflect.Value和interface{}都可以包含任意的值。二者的区别是**空接口**（interface{}）隐藏了值的布局信息、内置操作和相关方法，所以除非我们知道它的动态类型，并用一个类型断言类渗透进去，否则我们对所包含值能做的事情很少。<br/>
	作为对比，Value有很多方法可以用来分析所包含的值，而不用知道它的类型。<br/>
	
	尽管reflect.Value有很多方法，对于每个值，只有少量的方法可以安全调用。<br/>
	
	复杂类型的元素访问：
	
	* Slice和Array：`Len()`和`Index(i)`
	* Map：`MapKeys()`和`MapIndex(key)`
	* Struct：`NumField()`和`Field(i)`
	* Ptr：`Elem()`和`IsNil()`
	* Interface：`Elem`和`IsNil()`

	注意，结构体的字段列表包括了从匿名字段中做了类型提升的字段，并且包括未导出字段！
	
#### 代码示例：模拟fmt的实现

		
	```
	@see format.go
	@see display.go
	```
	
	
#### 值的寻址和设置

一个**变量**是一个**可寻址的存储区域**，其中包含了一个值，并且它的值可以通过这个地址来更新。<br/>

1. **可寻址性**

	事实上，通过`reflect.ValueOf(x)`返回的reflect.Value都是不可寻址的。<br/>
	可以通过`reflect.Value(&x).Elem()`来获得任意变量x的可寻址的Value值。<br/>
	判断一个Value是否为可寻址的，可以调用其`CanAddr()`方法。

	```
	x := 2
	a := reflect.ValueOf(2)  // 2 int
	b := reflect.ValueOf(x)  // 2 int
	c := reflect.ValueOf(&x) // &x *int
	d := c.Elem()            // 2 int (addressable by &x)

	fmt.Println(a.CanAddr()) // false
	fmt.Println(b.CanAddr()) // false
	fmt.Println(c.CanAddr()) // false
	fmt.Println(d.CanAddr()) // true
	```
	
	从一个可寻址的reflect.Value获取变量需要三步：

	1. 首先，调用`Addr()`，返回一个Value，其中包含一个只想变量的指针
	2. 接着，在这个Value上效用`Interface()`，会返回一个包含这个指针的interface{}值
	3. 最后，如果我们知道变量的类型，我们可以使用**类型断言**来把接口内容转换为一个普通指针

	```
	x := 2
	d := reflect.ValueOf(&x).Elem()   // d代表变量x
	px := d.Addr().Interface().(*int) // px := &x
	*px = 2333                        // x= 2333
	fmt.Println(x)
	```
	
2. **可设置性**

	可以直接通过可寻址的reflect.Value来更新变量，不用通过指针，而是直接调用`Set`方法。<br/>
	在这种情况下，会在运行时由Set方法来检查变量的可赋值条件。<br/>
	对于**未导出字段**，调用Set更新值会失败（panic）。<br/>
	判断一个Value是否为可设置的，可以调用其`CanSet()`方法。可设置的必然是可寻址的，反之则不。
	
	```
	x := 1
	rx := reflect.ValueOf(&x).Elem()
	rx.SetInt(98)
	rx.Set(reflect.ValueOf(99))
	// rx.Set(reflect.ValueOf("abc")) // panic...
	fmt.Println(x)

	stdout := reflect.ValueOf(os.Stdout).Elem() // os.Stdout: *os.File
	fmt.Println(stdout.Type())                  // os.File
	// fd.SetInt(2)	// panic...
	```

#### 代码示例：模拟encoding/...的实现

		
	```
	@see sexpr.go
	```

#### 结构体字段标签

略

## Practice

略

## Reference
> https://blog.golang.org/laws-of-reflection
> 
> * Reflection goes from interface value to reflection object.
> * Reflection goes from reflection object to interface value.
> * To modify a reflection object, the value must be settable.