package yaml

import (
	"io"
	"os"

	"github.com/ii64/go-binder/binder"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

func LoadConfig(path string) binder.LoadConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		f, err := os.Open(path)
		if err != nil {
			err = errors.Wrap(err, "load_cfg_yaml")
			return err
		}
		defer f.Close()
		return LoadConfigBuffer(f)(mc)
	}
}

func LoadConfigBuffer(r io.Reader) binder.LoadConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		dec := yaml.NewDecoder(r)
		if err := dec.Decode(&mc); err != nil {
			err = errors.Wrap(err, "load_cfg_yaml")
			return err
		}
		return nil
	}
}

//

func SaveConfig(path string, indent int) binder.SaveConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		f, err := os.Create(path)
		if err != nil {
			err = errors.Wrap(err, "save_cfg_yaml")
			return err
		}
		defer f.Close()
		return SaveConfigBuffer(f, indent)(mc)
	}
}

func SaveConfigBuffer(w io.Writer, indent int) binder.SaveConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		enc := yaml.NewEncoder(w)
		enc.SetIndent(indent)
		if err := enc.Encode(mc); err != nil {
			err = errors.Wrap(err, "save_cfg_yaml")
			return err
		}
		return nil
	}
}
