package main

// imports cgo:

/*
	extern int mini_calc(char *op, int a, int b);
*/
import "C"

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unsafe"
)

type OrderInfo struct {
	OrderCode   rune     // 4 bytes, plus 4 for padding
	Amount      int      // 8 bytes, no padding
	OrderNumber uint16   // 2 bytes, plus 6 for padding
	Items       []string // 24 bytes, no padding
	IsReady     bool     // 1 byte, plus 7 for padding
}

type SmallOrderInfo struct {
	IsReady     bool     // 1 byte + 1 byte of padding
	OrderNumber uint16   // 2 bytes
	OrderCode   rune     // 4 bytes
	Amount      int      // 8 bytes
	Items       []string // 24 bytes
}

func main() {

	// ### - exercise 2: use [unsafe.Sizeof] and [unsafe.Offsetof] to print out the size and offsets for the [OrderInfo] struct
	{
		fmt.Println("OrderInfo",
			unsafe.Sizeof(OrderInfo{}),
			unsafe.Offsetof(OrderInfo{}.OrderCode),
			unsafe.Offsetof(OrderInfo{}.Amount),
			unsafe.Offsetof(OrderInfo{}.OrderNumber),
			unsafe.Offsetof(OrderInfo{}.Items),
			unsafe.Offsetof(OrderInfo{}.IsReady),
		)
		fmt.Println("size of OrderCode:", unsafe.Sizeof(OrderInfo{}.OrderCode))
		fmt.Println("size of Amount:", unsafe.Sizeof(OrderInfo{}.Amount))
		fmt.Println("size of OrderNumber:", unsafe.Sizeof(OrderInfo{}.OrderNumber))
		fmt.Println("size of Items:", unsafe.Sizeof(OrderInfo{}.Items))
		fmt.Println("size of IsReady:", unsafe.Sizeof(OrderInfo{}.IsReady))
		fmt.Println("---")

		fmt.Println("SmallOrderInfo",
			unsafe.Sizeof(SmallOrderInfo{}),
			unsafe.Offsetof(SmallOrderInfo{}.IsReady),
			unsafe.Offsetof(SmallOrderInfo{}.OrderNumber),
			unsafe.Offsetof(SmallOrderInfo{}.OrderCode),
			unsafe.Offsetof(SmallOrderInfo{}.Amount),
			unsafe.Offsetof(SmallOrderInfo{}.Items),
		)
		fmt.Println("size of IsReady:", unsafe.Sizeof(SmallOrderInfo{}.IsReady))
		fmt.Println("size of OrderNumber:", unsafe.Sizeof(SmallOrderInfo{}.OrderNumber))
		fmt.Println("size of OrderCode:", unsafe.Sizeof(SmallOrderInfo{}.OrderCode))
		fmt.Println("size of Amount:", unsafe.Sizeof(SmallOrderInfo{}.Amount))
		fmt.Println("size of Items:", unsafe.Sizeof(SmallOrderInfo{}.Items))

		fmt.Println("---")
		fmt.Println("")
	}

	// ### - exercise 3: use cgo to call code from mini_calc.c in your Go program
	{
		operation := "+"
		result := C.mini_calc((*C.char)(C.CString(operation)), 3, 5)
		fmt.Println("result:", result)
	}
}

// ### - exercise 1: use reflection to create a string-length validator for struct fields
//
// DONE error if ...
// (1) one or more of the fields is a string
// (2) && has a tag minStrLen
// (3) && the length of the string field is less than the value specified in the struct tag
//
// DONE non-string fields and string fields that do not have the minStrLen struct tag are ignored
//
// DONE use errors.Join() to report all invalid fields
//
// DONE be sure to validate a struct was passed in
//
// DONE return nil if all struct fields are of the proper length

func ValidateStringLength(s any) error {
	errs := make([]error, 0)

	sv := reflect.ValueOf(s)

	if sv.Kind() != reflect.Struct {
		return fmt.Errorf("parameter must be a struct")
	}

	for i := 0; i < sv.NumField(); i++ {
		fieldVal := sv.Field(i)
		if fieldVal.Kind() != reflect.String {
			continue
		}

		st := reflect.TypeOf(s)
		fieldType := st.Field(i)

		tag, ok := fieldType.Tag.Lookup("minStrLen")
		if !ok {
			continue
		}

		minStrLen, err := strconv.Atoi(tag)
		if err != nil {
			errs = append(errs, fmt.Errorf("value of tag 'minStrLen' must be of type int: %w", err))
		}

		if f := fieldVal.String(); len(f) < minStrLen {
			errs = append(errs, fmt.Errorf("field '%s' falls below minStrLen of at least %d runes", fieldType.Name, minStrLen))
		}
	}

	return errors.Join(errs...)
}
