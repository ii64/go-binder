package binder

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

func LoadConfigTOML(path string) LoadConfigFunc {
	TagName = "toml"
	return func(mc *MappedConfiguration) (err error) {
		if _, err = toml.DecodeFile(path, mc); err != nil {
			err = errors.Wrap(err, "load_cfg_toml")
			return
		}
		return
	}
}

func SaveConfigTOML(path string) SaveConfigFunc {
	return func(mc *MappedConfiguration) (err error) {
		var f *os.File
		if f, err = os.Create(path); err != nil {
			err = errors.Wrap(err, "save_cfg_toml")
			return
		}
		defer f.Close()
		enc := toml.NewEncoder(f)
		enc.Indent = "  "
		if err = enc.Encode(mc); err != nil {
			err = errors.Wrap(err, "save_cfg_toml")
			return
		}
		return
	}
}
