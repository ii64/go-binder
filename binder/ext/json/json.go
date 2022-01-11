package json

import (
	"encoding/json"
	"os"

	"io"

	"github.com/ii64/go-binder/binder"
	"github.com/pkg/errors"
)

func LoadConfig(path string) binder.LoadConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		f, err := os.Open(path)
		if err != nil {
			err = errors.Wrap(err, "")
			return err
		}
		defer f.Close()
		return LoadConfigBuffer(f)(mc)
	}
}

func LoadConfigBuffer(r io.Reader) binder.LoadConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		dec := json.NewDecoder(r)
		if err := dec.Decode(mc); err != nil {
			err = errors.Wrap(err, "load_cfg_json")
			return err
		}
		return nil
	}
}

//

func SaveConfig(path, indent string) binder.SaveConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		f, err := os.Create(path)
		if err != nil {
			err = errors.Wrap(err, "save_cfg_json")
			return err
		}
		defer f.Close()
		return SaveConfigBuffer(f, indent)(mc)
	}
}

func SaveConfigBuffer(w io.Writer, indent string) binder.SaveConfigFunc {
	return func(mc *binder.MappedConfiguration) error {
		enc := json.NewEncoder(w)
		enc.SetIndent("", indent)
		if err := enc.Encode(mc); err != nil {
			err = errors.Wrap(err, "save_cfg_json")
			return err
		}
		return nil
	}
}
