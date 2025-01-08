package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

type Foo struct {
	a int
}

type Bar struct {
	s string
}

type BarB struct {
	s       string
	another string `minStrLen:"a"`
}

type Too struct {
	a string `minStrLen:"15"`
}

type Large struct {
	a string `minStrLen:"10"`
	b string `minStrLen:"10"`
	c string `minStrLen:"10"`
	d string `minStrLen:"10"`
	e string `minStrLen:"10"`
	f string `minStrLen:"10"`
	g string `minStrLen:"10"`
	h string `minStrLen:"10"`
}

func TestValidateStringLength_validBehavior(t *testing.T) {
	t.Parallel()

	// struct w/o string
	{
		f := Foo{1}
		err := ValidateStringLength(f)
		if err != nil {
			t.Fatal("wanted NO error but got:", err)
		}
	}

	// struct w/ string field w/o tag
	{
		b := Bar{s: "this is a non validated string"}
		err := ValidateStringLength(b)
		if err != nil {
			t.Fatal("wanted NO error, but got:", err)
		}
	}

	// struct w/ string field w/ valid tag && valid value
	{
		too := Too{a: "this is a non validated string"}
		err := ValidateStringLength(too)
		if err != nil {
			t.Fatal("wanted NO error, but got:", err)
		}
	}

	// struct w/ string field w/ valid tag && valid values
	{
		l2 := Large{
			a: "sufficient length",
			b: "sufficient length",
			c: "sufficient length",
			d: "sufficient length",
			e: "sufficient length",
			f: "sufficient length",
			g: "sufficient length",
			h: "sufficient length",
		}
		err := ValidateStringLength(l2)
		if err != nil {
			t.Fatal("wanted NO error, but got:", err)
		}
	}
}

func TestValidateStringLength_invalidBehavior(t *testing.T) {
	t.Parallel()

	// no struct
	{
		err := ValidateStringLength("1")
		if err == nil {
			t.Fatal("wanted an error, but got none")
		}
		if !cmp.Equal(err.Error(), "parameter must be a struct") {
			t.Error(cmp.Diff(err.Error(), "parameter must be a struct"))
		}
	}

	// struct w/ string field w/o tag && string field w/ invalid tag
	{
		bb := BarB{s: "this is a non validated string", another: "1234567899"}
		err := ValidateStringLength(bb)
		if err == nil {
			t.Fatal("wanted an error, but got none")
		}
		msg := "value of tag 'minStrLen' must be of type int: strconv.Atoi: parsing \"a\": invalid syntax"
		if !cmp.Equal(err.Error(), msg) {
			t.Error(cmp.Diff(err.Error(), msg))
		}
	}

	// struct w/ string field w/ valid tag && invalid value
	{
		l := Large{
			a: "too short",
			b: "too short",
			c: "too short",
			d: "too short",
			e: "too short",
			f: "too short",
			g: "too short",
			h: "too short",
		}
		err := ValidateStringLength(l)
		if err == nil {
			t.Fatal("wanted an error, but got none")
		}
		msgs := `field 'a' falls below minStrLen of at least 10 runes
field 'b' falls below minStrLen of at least 10 runes
field 'c' falls below minStrLen of at least 10 runes
field 'd' falls below minStrLen of at least 10 runes
field 'e' falls below minStrLen of at least 10 runes
field 'f' falls below minStrLen of at least 10 runes
field 'g' falls below minStrLen of at least 10 runes
field 'h' falls below minStrLen of at least 10 runes`
		if !cmp.Equal(err.Error(), msgs) {
			t.Error(cmp.Diff(err.Error(), msgs))
		}
	}
}
