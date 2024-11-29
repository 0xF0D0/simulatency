package config

import (
	"os"

	"github.com/0xf0d0/simulatency-controller/client/yaml"
)

type SimulatencyConfig struct {
	Regions []SimulatencyRegion
}

type SimulatencyRegion struct {
	Name      string `yaml:"name"`
	Longitude int    `yaml:"longitude"`
}


func GetMesaGrpcGatewayConfig(path string) (*SimulatencyConfig, error) {
	parser := yaml.NewYamlParser()

	bz, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &SimulatencyConfig{}
	err = parser.UnMarshal(bz, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}