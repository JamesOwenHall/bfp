package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// jsonConfiguration is a struct that mirrors the data as it should be found in
// the configuration file.
type jsonConfiguration struct {
	Directions    []jsonDirection `json:"directions"`
	ListenAddress string          `json:"listen address"`
}

// jsonDirection is a struct that mirrors the direction objects as they should
// be found in the configuration file.
type jsonDirection struct {
	Name       string  `json:"name"`
	Typ        string  `json:"type"`
	WindowSize float64 `json:"window size"`
	MaxHits    float64 `json:"max hits"`
}

// parseJsonFile will read the contents of a file and return its structure as a
// jsonConfiguration.
func parseJsonFile(filename string) (*jsonConfiguration, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("can't find configuration file %s.", filename)
	}

	parsed := new(jsonConfiguration)
	err = json.Unmarshal(data, parsed)
	if err != nil {
		return nil, fmt.Errorf("configuration file is not valid.")
	}

	return parsed, nil
}
