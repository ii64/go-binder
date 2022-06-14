package binder

import (
	"flag"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestFlagBinder(t *testing.T) {
	var expectedValue []string

	type irStruct struct {
		TestFieldString *string    `arg:"myfield1" usage:"the usage"`
		TestFieldSlice  *[]string  `arg:"myfield2" usage:"the usage"`
		TestFieldArray  *[5]string `arg:"myfield3" usage:"the usage"`
	}

	ir := irStruct{
		TestFieldString: new(string),
		TestFieldSlice:  &[]string{},
		TestFieldArray:  &[5]string{},
	}

	val := reflect.ValueOf(ir)
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		tf := typ.Field(i)
		tv := val.Field(i)

		ns := "conf"
		argName := tf.Tag.Get("arg")
		if argName == "" {
			argName = tf.Name
		}
		assert.NotEqual(t, "", argName)

		expectedValue = append(expectedValue, ns+"."+argName)

		defaultRegisterCmdArgsFlagStd(ns, tf, tv)
	}

	var actualValue = map[string][]*flag.Flag{}
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		v, _ := actualValue[f.Name]
		actualValue[f.Name] = append(v, f)
	})

	for _, argKey := range expectedValue {
		fs, exist := actualValue[argKey]
		assert.Equal(t, true, exist, argKey)
		for _, f := range fs {
			assert.Equal(t, "the usage", f.Usage)
		}
	}

	flag.CommandLine.Parse([]string{
		"-conf.myfield1", "hello world1",
		"-conf.myfield1", "hello world2",
		"-conf.myfield2", "hello world3",
		"-conf.myfield2", "hello world4",
		"-conf.myfield3", "hello world5",
		"-conf.myfield3", "hello world6",
	})

	spew.Dump(ir)
}
