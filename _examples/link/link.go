package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/ii64/go-binder/binder"
)

type myint int
type mystring string

type MyConfig struct {
	Token string
	Count int

	MyCounter   myint
	MyString    mystring
	StringSlice []string

	Ktes *int
	Sub  *struct {
		Hello    *string
		SubOfSub struct {
			InSub bool
		}
		PtrOfSub *struct {
			YourName *string
		}
	}
	Log struct {
		SubLog       int
		Directory    string
		Filename     string
		DedicatedArg string
	}
}

var (
	TMap   = map[string]interface{}{}
	Loaded = &MyConfig{
		Token:     "",
		Count:     0,
		MyCounter: 0,
		Ktes:      new(int),
		Sub: &struct {
			Hello    *string
			SubOfSub struct{ InSub bool }
			PtrOfSub *struct{ YourName *string }
		}{},
		Log: struct {
			SubLog       int
			Directory    string
			Filename     string
			DedicatedArg string
		}{},
	}
)

func main() {
	st, in, out := binder.Link(Loaded)
	_, _ = in, out
	spew.Dump(st)
}
