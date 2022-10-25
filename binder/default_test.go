package binder

import (
	"encoding/json"
	"flag"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type myTestMapFlagValue map[string]string

var _ flag.Value = (myTestMapFlagValue)(nil)

func (m myTestMapFlagValue) Set(v string) (err error) {
	return json.Unmarshal([]byte(v), &m)
}

func (m myTestMapFlagValue) String() string {
	return "{}"
}

type myTestFlagValue struct {
	XDSTarget *string `json:"xds,omitempty" env:"xds_target"`
	Bind      string  `json:"bind,omitempty" env:"bind"`
}

var _ flag.Value = (*myTestFlagValue)(nil)

func (m *myTestFlagValue) Set(v string) (err error) {
	return json.Unmarshal([]byte(v), m)
}

func (m *myTestFlagValue) String() string {
	return "{}"
}

type myTestStruct struct {
	TestFieldString *string            `arg:"myfield1" env:"test_field_string" usage:"the usage"`
	TestFieldSlice  *[]string          `arg:"myfield2" env:"test_field_slice" usage:"the usage"`
	TestFieldArray  *[5]string         `arg:"myfield3" env:"test_field_array" usage:"the usage"`
	MyFlag          *myTestFlagValue   `arg:"my_flag" bind:"my_flag" usage:"the usage"`
	MyMapFlag       myTestMapFlagValue `arg:"my_map" bind:"my_map" env:"my_map" usage:"the usage"`
}

func TestFlagBinder(t *testing.T) {
	var expectedValue []string
	m := initTest()
	m.memSave(&MappedConfiguration{
		// "conf": Mapp,
	})

	ir := myTestStruct{
		TestFieldString: new(string),
		TestFieldSlice:  &[]string{},
		TestFieldArray:  &[5]string{},
		MyFlag:          &myTestFlagValue{},
		MyMapFlag:       myTestMapFlagValue{},
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
		"-conf.my_flag", `{"xds": "xds://xx.cluster.local:10001", "bind": "0:1000"}`,
		"-conf.my_map", `{"xds_c0": "xds://c0", "xds_c1": "xds://c1"}`,
	})

	require.NotNil(t, ir.TestFieldString)
	assert.Equal(t, "hello world2", *ir.TestFieldString)
	require.NotNil(t, ir.TestFieldSlice)
	assert.Equal(t, []string{"hello world3", "hello world4"}, *ir.TestFieldSlice)
	require.NotNil(t, ir.TestFieldArray)
	assert.Equal(t, [5]string{"hello world5", "hello world6"}, *ir.TestFieldArray)

	require.NotNil(t, ir.MyFlag)
	assert.Equal(t, "0:1000", ir.MyFlag.Bind)
	require.NotNil(t, ir.MyFlag.XDSTarget)
	assert.Equal(t, "xds://xx.cluster.local:10001", *ir.MyFlag.XDSTarget)

	require.NotNil(t, ir.MyMapFlag)
	require.Contains(t, ir.MyMapFlag, "xds_c0")
	require.Equal(t, "xds://c0", ir.MyMapFlag["xds_c0"])
	require.Contains(t, ir.MyMapFlag, "xds_c1")
	require.Equal(t, "xds://c1", ir.MyMapFlag["xds_c1"])

	spew.Dump(ir)
}
