package binder

import (
	"flag"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mem struct {
	saved MappedConfiguration
}

func (m *mem) memLoad(mc *MappedConfiguration) error {
	*mc = m.saved
	return nil
}

func (m *mem) memSave(mc *MappedConfiguration) error {
	m.saved = *mc
	return nil
}

func initTest() *mem {
	m := &mem{}
	LoadConfig = m.memLoad
	SaveConfig = m.memSave
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	return m
}

func TestBinderUseEnv(t *testing.T) {
	var err error
	m := initTest()
	m.memSave(&MappedConfiguration{
		// "conf": Mapp,
	})

	ir := myTestStruct{
		TestFieldString: new(string),
		TestFieldSlice:  &[]string{},
		TestFieldArray:  &[5]string{},
		MyFlag:          &myTestFlagValue{},
	}
	err = BindArgsConf(&ir, "conf")
	require.NoError(t, err, "bind args conf")

	for k, v := range map[string]string{
		"conf.test_field_string": "hello world1",
		"conf.test_field_slice":  "hello world2",
		// "conf.test_field_array":  "c",
		"conf.my_flag.xds_target": "xds://xx.cluster.local:10001",
		"conf.my_flag.bind":       "0:1000",
	} {
		os.Setenv(k, v)
	}

	err = Init()
	require.NoError(t, err, "bind init")

	In()

	require.NotNil(t, ir.TestFieldString)
	assert.Equal(t, "hello world1", *ir.TestFieldString)
	require.NotNil(t, ir.TestFieldSlice)
	assert.Equal(t, []string{"hello world2"}, *ir.TestFieldSlice)

	require.NotNil(t, ir.MyFlag)
	assert.Equal(t, "0:1000", ir.MyFlag.Bind)
	require.NotNil(t, ir.MyFlag.XDSTarget)
	assert.Equal(t, "xds://xx.cluster.local:10001", *ir.MyFlag.XDSTarget)

	spew.Dump(ir)
}
