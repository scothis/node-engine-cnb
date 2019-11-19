package node

import (
	"os"

	"gopkg.in/yaml.v2"
)

type BuildpackYMLVersionParser struct{}

func NewBuildpackYMLVersionParser() BuildpackYMLVersionParser {
	return BuildpackYMLVersionParser{}
}

func (p BuildpackYMLVersionParser) Parse(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	defer file.Close()

	var config struct {
		NodeJS struct {
			Version string `yaml:"version"`
		} `yaml:"nodejs"`
	}

	err = yaml.NewDecoder(file).Decode(&config)
	if err != nil {
		return "", err
	}

	return config.NodeJS.Version, nil
}
