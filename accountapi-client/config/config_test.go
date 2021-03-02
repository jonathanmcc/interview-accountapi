package config

import (
	"accountapi-client/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

type config_test_1 struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
}

type config_test_2 struct {
	Foo struct {
		Bar string `yaml:"bar"`
		Baz string `yaml:"baz"`
	} `yaml:"foo"`
}

func TestReadConfig(t *testing.T) {
	c := config_test_1{}
	err := config.ReadConfig("examples/server.yaml", &c)
	want := "http://localhost"
	assert.NoError(t, err)
	assert.Equal(t, want, c.Server.Host, "got %s want %s", c.Server.Host, want)

	err = config.ReadConfig("nonexistent.yaml", &c)
	if assert.Error(t, err, "should return error when file does not exist") {
		assert.Equal(t, "open nonexistent.yaml: no such file or directory", err.Error())
	}
}
