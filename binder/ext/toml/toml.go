package toml

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ii64/go-binder/binder"
	"github.com/pkg/errors"
)

func LoadConfig(path string) binder.LoadConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		f, err := os.Open(path)
		if err != nil {
			err = errors.Wrap(err, "load_cfg_toml")
			return err
		}
		defer f.Close()
		return LoadConfigBuffer(f)(mc)
	}
}

func LoadConfigBuffer(r io.Reader) binder.LoadConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		b, err := io.ReadAll(r)
		if err != nil {
			err = errors.Wrap(err, "load_cfg_toml")
			return err
		}
		if _, err := toml.Decode(string(b), mc); err != nil {
			err = errors.Wrap(err, "load_cfg_toml")
			return err
		}
		return nil
	}
}

//

func SaveConfig(path string, indent string) binder.SaveConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		f, err := os.Create(path)
		if err != nil {
			err = errors.Wrapf(err, "save_cfg_toml")
			return err
		}
		defer f.Close()
		return SaveConfigBuffer(f, indent)(mc)
	}
}

func SaveConfigBuffer(w io.Writer, indent string) binder.SaveConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		enc := toml.NewEncoder(w)
		enc.Indent = indent
		if err := enc.Encode(mc); err != nil {
			err = errors.Wrap(err, "save_cfg_toml")
			return err
		}
		return nil
	}
}
