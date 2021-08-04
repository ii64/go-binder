<img src="go-binder.png" alt="go-binder" width="35%"/>

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/ii64/go-binder.svg)](https://pkg.go.dev/github.com/ii64/go-binder)
[![DeepSource](https://deepsource.io/gh/ii64/go-binder.svg/?label=active+issues&show_trend=true&token=vJ6MknNoKFjEe4Fr1cTdsP0x)](https://deepsource.io/gh/ii64/go-binder/?ref=repository-badge)

✨Binding configuration and command flag made easy!✨

You can use multiple keys tag to simplify the look like this [(supported feature\*\*)](https://github.com/golang/go/issues/40281):

```go
// single tag key
type MyConfig struct {
    Token string `json:token_json" xml:"token_xml" arg:"token" env:"token" environ:"TOKEN"`
}
// multiple tag keys
type MyConfig struct {
    Token string `json xml bson yaml toml arg env:"token" environ:"TOKEN"`
}
```

Below is default mapping implementation of `binder.RegisterCmdArgs = defaultRegisterCmdArgsFlagStd` that use standard golang `flag` package to perform command flag.

The `<parent>` is a placeholder for parent key, `binder.BindArgs(Loaded, "my")` this case `<parent>` will be replaced with `my`, if there's field with type `struct` in the component, it'll be replaced to `my.<struct name | arg value>.<field name | arg>`

| Tag               | Go Code                                                 | Description                                                                       |
| ----------------- | ------------------------------------------------------- | --------------------------------------------------------------------------------- |
| `arg:"token"`     | `flag.StringVar(val, "<parent>.token", *val, argUsage)` | Used for binding flag with contextual key `<parent>`                              |
| `argx:"token"`    | `flag.StringVar(val, "token", *val, argUsage)`          | Used for binding flag                                                             |
| `bind:"log"`      | _No equivalent_                                         | Used for binder to differ `struct` parent sub context `<parent>.log.<field name>` |
| `env:"token"`     | `os env("<parent>.token")`                              | Used for binding to environment variable with contextual key `<parent>`           |
| `environ:"token"` | `os env("token")`                                       | Used for binding to environment                                                   |
| `usage:"<DESC>"`  | Used as `argUsage`                                      | Description for flag                                                              |
| `ignore:"true"`   | _No equivalent_                                         | Ignore struct field                                                               |
| `bind:"abc"`      | _No equivalent_                                         | Used for mapstructure (`bind` is default value of `binder.TagName`)               |

Other:

- `arg` and `argx` (dedicated) basically has same function.
- `env` and `environ` (dedicated) basically has same function.
- Currently you can't have _dedicated key_ for configuration because the way it parsed is from _Unmarshaller_ that results in `map[string]interface{}`, but this definitely possible to implement.

More thing you can learn from the example below.

## Example

```go
package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/ii64/go-binder/binder"
    "github.com/pkg/errors"
)

type MyConfig struct {
    Token string `json xml bson yaml toml arg:"token,omitempty" env:"TOKEN" environ:"TOKEN"`
    Count int    `json xml bson yaml toml arg:"count,omitempty" env:"COUNT" usage:"this is the usage"`

    Ktes *int
    Sub  **struct {
        Hello    *string
        SubOfSub struct {
            InSub **bool
        }
        PtrOfSub *struct {
            YourName **string `json xml bson yaml toml bind:"your_name,omitempty" env:"COUNT" usage:"this is the usage"`
        }
    }
    Log struct {
        SubLog       int
        Directory    string
        Filename     string `json xml bson yaml toml arg:"filename" env:"FILENAME"`
        DedicatedArg string `json xml bson yaml toml argx:"dedicatedArg" env:"DEDICATED_ARG"`
    } `json xml bson yaml toml arg env bind:"log"`
}
var (
    configFile = os    env("CONFIG_FILE")
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
    binder.LoadConfig = binder.LoadConfigJSON(configFile)
    binder.SaveConfig = binder.SaveConfigJSON(configFile)
    binder.SaveOnClose = true
    // register component to binder
    registerToBinder()
    // perform binding
    if err = binder.Init(); err != nil {
        if errors.Is(err, os.ErrNotExist) {
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
    fmt.Printf("%+#v\n", Loaded)
}
```

Output help:

```text
$ main -h
Usage of main:
  -dedicatedArg string

  -my.Sub.Hello string

  -my.count int
        this is the usage (default 121)
  -my.log.Directory string

  -my.log.filename string

  -my.token string
         (default "some default value")
```

Output JSON:

```json
{
  "my": {
    "token": "some default value",
    "count": 121,
    "Ktes": 0,
    "Sub": {
      "Hello": "",
      "SubOfSub": {
        "InSub": false
      },
      "PtrOfSub": {
        "your_name": ""
      }
    },
    "log": {
      "SubLog": 0,
      "Directory": "",
      "filename": "",
      "dedicatedArg": ""
    }
  }
}
```

Output TOML:

```toml
[my]
  token = "some default value"
  count = 121
  Ktes = 0
  [my.Sub]
    Hello = ""
    [my.Sub.SubOfSub]
      InSub = false
    [my.Sub.PtrOfSub]
      your_name = ""
  [my.log]
    SubLog = 0
    Directory = ""
    filename = ""
    dedicatedArg = ""

```

Output YAML:

```yaml
my:
  token: some default value
  count: 121
  ktes: 0
  sub:
    hello: ""
    subofsub:
      insub: false
    ptrofsub:
      your_name: ""
  log:
    sublog: 0
    directory: ""
    filename: ""
    dedicatedArg: ""
```

### Note

Contributions are welcome

```text
**) Reverted feature as from 1.16 but found it useful
```

## License

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
