package main

import (
	"errors"
	"flag"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ii64/go-binder/binder"
	"github.com/ii64/go-binder/binder/ext/yaml"

	_ "github.com/joho/godotenv/autoload"
)

type MyConfig struct {
	A string `environ:"TOKEN"`
}

var Loaded = &MyConfig{
	A: "default",
}

const configfile = "config.yaml"

func init() {
	binder.LoadConfig = yaml.LoadConfig(configfile)
	binder.SaveConfig = yaml.SaveConfig(configfile, 2)
	binder.SaveOnClose = true

	binder.BindArgsConf(Loaded, "main")
}

func main() {
	var err error
	if err = binder.Init(); err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, io.EOF) {
			binder.In()
			if err = binder.Save(); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	flag.Parse()

	binder.In()
	defer binder.Close()

	spew.Dump(Loaded)
}
