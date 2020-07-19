package main

import (
	"io/ioutil"

	"github.com/satori/uuid"
	"gopkg.in/yaml.v3"
)

// ConfigureGaspaNode as is
type ConfigureGaspaNode struct {
	Gaspa *ConfigureBridgeNode `yaml:"gaspa"`
}

// ConfigureBridgeNode as is
type ConfigureBridgeNode struct {
	Bridge     *ConfigurePrototype      `yaml:"bridge"`
	Connection *ConfigureConnectionNode `yaml:"connection"`
}

// ConfigureConnectionNode as is
type ConfigureConnectionNode struct {
	KeepAlive int `yaml:"keep-alive"`
}

// ConfigurePrototype as is
type ConfigurePrototype struct {
	Name   string `yaml:"name"`
	UUID   string `yaml:"uuid"`
	Listen string `yaml:"listen"`
}

type configure struct {
	name      string
	uuid      uuid.UUID
	listen    string
	keepAlive int
}

var config *configure

func loadConfigure(filename string) error {
	origin, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	confPrototype := new(ConfigureGaspaNode)
	err = yaml.Unmarshal(origin, confPrototype)
	if err != nil {
		return err
	}
	uuid, err := uuid.FromString(confPrototype.Gaspa.Bridge.UUID)
	if err != nil {
		return err
	}
	config = &configure{
		name:      confPrototype.Gaspa.Bridge.Name,
		uuid:      uuid,
		listen:    confPrototype.Gaspa.Bridge.Listen,
		keepAlive: confPrototype.Gaspa.Connection.KeepAlive,
	}
	return nil
}
