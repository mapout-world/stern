package config

import (
	"os"
	stdfilepath "path/filepath"

	"github.com/go-micro/plugins/v4/config/encoder/yaml"
	"go-micro.dev/v4/config"
	"go-micro.dev/v4/config/reader"
	"go-micro.dev/v4/config/reader/json"
	"go-micro.dev/v4/config/source/file"
)

var (
	cfg *DynamicConfig
)

type DynamicConfig struct {
	config.Config
	file string
}

func MustLoadDynamic(filepath string) {
	opts := []config.Option{
		config.WithReader(json.NewReader(reader.WithEncoder(yaml.NewEncoder()))),
	}

	c, err := config.NewConfig(opts...)
	if err != nil {
		panic(err)
	}

	if err := c.Load(file.NewSource(file.WithPath(filepath))); err != nil {
		panic(err)
	}

	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file, err := stdfilepath.Rel(pwd, filepath)
	if err != nil {
		panic(err)
	}

	cfg = &DynamicConfig{
		Config: c,
		file:   file,
	}
}

func Stop() error {
	return cfg.Close()
}

type Value interface {
	reader.Value
}

func Get(path ...string) Value {
	return cfg.Get(path...)
}

func All() map[string]any {
	return cfg.Map()
}
