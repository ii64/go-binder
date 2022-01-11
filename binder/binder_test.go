package binder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	hello string `environ:"HELLO"`
	Token string `environ:"TOKEN"`
	Count int    `json xml bson yaml toml arg:"count,omitempty" env:"COUNT" usage:"this is the usage"`
}

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
	return m
}

func TestBinderUseEnv(t *testing.T) {
	m := initTest()
	_ = m
	m.memSave(&MappedConfiguration{
		"base": MappedConfiguration{
			"hello": "hey",
			"Token": "",
			"Count": 221,
		},
	})
	var err error

	var td *testStruct = &testStruct{
		hello: "123",
		Token: "3242",
		Count: 999,
	}
	err = BindArgsConf(td, "base")
	assert.NoError(t, err, "bindArgsConf")

	err = Init()
	assert.NoError(t, err, "init")

	_ = td
}
