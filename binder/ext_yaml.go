package binder

import (
	"os"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

func LoadConfigYAML(path string) LoadConfigFunc {
	TagName = "yaml"
	return func(mc *MappedConfiguration) (err error) {
		var f *os.File
		if f, err = os.Open(path); err != nil {
			err = errors.Wrap(err, "load_cfg_yaml")
			return
		}
		defer f.Close()
		dec := yaml.NewDecoder(f)
		if err = dec.Decode(mc); err != nil {
			err = errors.Wrap(err, "load_cfg_yaml")
			return
		}
		return
	}
}

func SaveConfigYAML(path string) SaveConfigFunc {
	return func(mc *MappedConfiguration) (err error) {
		var f *os.File
		if f, err = os.Create(path); err != nil {
			err = errors.Wrap(err, "save_cfg_yaml")
			return
		}
		defer f.Close()
		enc := yaml.NewEncoder(f)
		enc.SetIndent(2)
		if err = enc.Encode(mc); err != nil {
			err = errors.Wrap(err, "save_cfg_yaml")
			return
		}
		return
	}
}
