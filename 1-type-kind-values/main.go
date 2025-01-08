package main

import (
	"fmt"
	"reflect"
)

func main() {

	// Type
	{
		type Foo struct{}

		var x int
		xt := reflect.TypeOf(x)
		fmt.Println(xt.Name()) // returns int

		f := Foo{}
		ft := reflect.TypeOf(f)
		fmt.Println(ft.Name()) // returns Foo

		xpt := reflect.TypeOf(&x)
		fmt.Println(xpt.Name()) // returns empty string
		fmt.Println("----")
	}

	// Kind
	{
		var x int
		xpt := reflect.TypeOf(&x)
		fmt.Println(xpt.Name())        // returns empty string
		fmt.Println(xpt.Kind())        // returns reflect.Pointer
		fmt.Println(xpt.Elem().Name()) // returns int
		fmt.Println(xpt.Elem().Kind()) // returns reflect.Int
		fmt.Println("----")
	}

	// Elem on composite types
	{
		type Foo struct {
			A int    `myTag:"value1"`
			B string `myTag:"value2"`
		}

		var f Foo
		ft := reflect.TypeOf(f)

		// NumField returns the number of fields contained in a struct
		for i := 0; i < ft.NumField(); i++ {
			// Field returns a struct type's i'th field.
			field := ft.Field(i)

			fmt.Println(field.Name)             // returns A or B
			fmt.Println(field.Type.Name())      // returns int or string
			fmt.Println(field.Tag.Get("myTag")) // returns value1 or value 2
		}
		fmt.Println("----")
	}

	// Values
	{
		i := 10
		// we need to pass a pointer for being able to modifying the underlying value of the variable i (since Go is a "pass by value" language)
		iv := reflect.ValueOf(&i)
		ivv := iv.Elem()

		ivv.SetInt(500)
		fmt.Println(i) // prints 500
		fmt.Println("----")
	}

	// Make new Values
	{
		//reflect.MakeChan()
		//reflect.MakeSlice()
		//reflect.MakeMap()
		//reflect.MakeMapWithSize()

		// NOTE: you must always start from a value when constructing a [reflect.Type] with [reflect.New] or the Make... functions
		// however, if you do not have a value handy, a trick lets you create a variable to represent a [reflect.Type]
		//
		// (1) using nil to initialize a string pointer
		// (*string)(nil) ... explicit type conversion of nil to *string
		stringType := reflect.TypeOf((*string)(nil)).Elem()
		fmt.Println(stringType)

		// (2) using nil to initialize a string slice
		stringSliceType := reflect.TypeOf([]string(nil))
		fmt.Println(stringSliceType)

		ssv := reflect.MakeSlice(stringSliceType, 0, 10)

		sv := reflect.New(stringType).Elem()
		sv.SetString("my first variable created with reflection")

		ssv = reflect.Append(ssv, sv)
		fmt.Println(ssv)
	}
}

// GENERAL NOTE:
// even though it is possible to detect an interface with a nil value,
// strive to write your code, so that it performs correctly,
// even when the value associated with an interface is nil.
//
// reserve this code for a situation where you have no other options.
func hasNoValue(i any) bool {
	iv := reflect.ValueOf(i)

	// IsValid returns true if [reflect.Value] holds anything other than a nil interface
	//
	// NOTE:
	// (1) a variable of an interface type is nil only if the associated type and value of the variable of an interface type are nil
	// (2) if a nil variable of a concrete type is assigned to a variable of an interface type, the variable of the interface type is not nil
	//
	// NOTE: you need to call IsValid first because calling any other method on [reflect.Value] will panic if IsValid is false (=nil interface)
	if !iv.IsValid() {
		return true
	}

	switch iv.Kind() {
	case reflect.Pointer, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		// IsNil returns true if the value of the [reflect.Value] is nil
		//
		// but it can be called only if the [reflect.Kind] is something that CAN be nil
		// if you call it on something whose zero values is not nil, it panics
		return iv.IsNil()
	default:
		return false
	}
}

// changeInt does effectively the same as changeIntReflect
func changeInt(i *int) {
	*i = 20
}

// changeIntReflect does effectively the same as changeInt
func changeIntReflect(i *int) {
	iv := reflect.ValueOf(i)
	iv.Elem().SetInt(20)
}
