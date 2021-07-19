package binder

import (
	"encoding/json"
	"os"

	"github.com/pkg/errors"
)

func LoadConfigJSON(path string) LoadConfigFunc {
	TagMapDefault = "json"
	return func(mc *MappedConfiguration) (err error) {
		var f *os.File
		if f, err = os.Open(path); err != nil {
			err = errors.Wrap(err, "load_cfg_json")
			return
		}
		defer f.Close()
		dec := json.NewDecoder(f)
		if err = dec.Decode(mc); err != nil {
			err = errors.Wrap(err, "load_cfg_json")
			return
		}
		return
	}
}

func SaveConfigJSON(path string) SaveConfigFunc {
	return func(mc *MappedConfiguration) (err error) {
		var f *os.File
		if f, err = os.Create(path); err != nil {
			err = errors.Wrap(err, "save_cfg_json")
			return
		}
		defer f.Close()
		enc := json.NewEncoder(f)
		enc.SetIndent("", "\t")
		if err = enc.Encode(mc); err != nil {
			err = errors.Wrap(err, "save_cfg_json")
			return
		}
		return
	}
}
