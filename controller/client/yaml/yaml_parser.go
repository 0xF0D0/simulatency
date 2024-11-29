package yaml

import (
	"gopkg.in/yaml.v3"
)

type YamlParser interface {
	UnMarshal(in []byte, args interface{}) error
}

type yamlParser struct{}

func NewYamlParser() YamlParser {
	return &yamlParser{}
}

func (p *yamlParser) UnMarshal(in []byte, arg interface{}) error {
	return yaml.Unmarshal(in, arg)
}
