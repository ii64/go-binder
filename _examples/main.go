package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/ii64/go-binder/binder"
	"github.com/pkg/errors"
)

type MyConfig struct {
	Token string `json xml bson yaml toml arg:"token,omitempty" env:"TOKEN" environ:"TOKEN"`
	Count int    `json xml bson yaml toml arg:"count,omitempty" env:"COUNT" usage:"this is the usage"`

	Ktes *int
	Sub  *struct {
		Hello *string
	}
	Log struct {
		Directory    string
		Filename     string `json xml bson yaml toml arg:"filename" env:"FILENAME"`
		DedicatedArg string `json xml bson yaml toml argx:"dedicatedArg" env:"DEDICATED_ARG"`
	} `json xml bson yaml toml arg env bind:"log"`
}

var (
	configFile = os.Getenv("CONFIG_FILE")
	Loaded     *MyConfig
)

func registerToBinder() {
	Loaded = &MyConfig{
		Token: "some default value",
		Count: 121,
	}
	binder.BindArgsConf(Loaded, "my")
}

func main() {
	var err error
	if configFile == "" {
		configFile = "config.json"
	}
	ext := filepath.Ext(configFile)
	fmt.Printf("Filename %q ext %q\n", configFile, ext)
	switch ext {
	case ".json":
		// json
		binder.LoadConfig = binder.LoadConfigJSON(configFile)
		binder.SaveConfig = binder.SaveConfigJSON(configFile)
	case ".yaml":
		// yaml
		binder.LoadConfig = binder.LoadConfigYAML(configFile)
		binder.SaveConfig = binder.SaveConfigYAML(configFile)
	case ".toml":
		// toml
		binder.LoadConfig = binder.LoadConfigTOML(configFile)
		binder.SaveConfig = binder.SaveConfigTOML(configFile)
	}
	binder.SaveOnClose = true
	// register component to binder
	registerToBinder()
	// perform binding
	if err = binder.Init(); err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, io.EOF) {
			if err = binder.Save(); err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	flag.Parse()
	// reflect back to component
	binder.In()
	defer binder.Close()

	// runtime
	spew.Dump(Loaded)
	// fmt.Printf("%+#v\n", Loaded)
}
