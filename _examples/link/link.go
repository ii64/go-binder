package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/ii64/go-binder/binder"
)

type MyConfig struct {
	Token string
	Count int

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
	Loaded = &MyConfig{}
)

func main() {
	st, in, out := binder.Link(Loaded)
	_, _ = in, out
	spew.Dump(st)
}
