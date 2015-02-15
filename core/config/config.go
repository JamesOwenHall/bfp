package config

import (
	"github.com/JamesOwenHall/BruteForceProtection/core/hitcounter"
)

// Configuration is a struct that represents the contents of a configuration
// file.
type Configuration struct {
	Directions       []hitcounter.Direction
	ListenAddress    string
	ListenType       string
	DashboardAddress string
}

// ReadConfig parses a configuration file and returns an instance of
// Configuration.
func ReadConfig(filename string) (*Configuration, []error) {
	parsed, err := parseJsonFile(filename)
	if err != nil {
		return nil, []error{err}
	}
	if errs := parsed.Validate(); len(errs) != 0 {
		return nil, errs
	}

	result := new(Configuration)
	result.ListenAddress = parsed.ListenAddress
	result.ListenType = parsed.ListenType
	result.DashboardAddress = parsed.DashboardAddress
	result.Directions = make([]hitcounter.Direction, 0, len(parsed.Directions))

	for _, jsonDir := range parsed.Directions {
		// Create the direction according to its type
		dir := hitcounter.Direction{
			Name:        jsonDir.Name,
			CleanUpTime: jsonDir.CleanUpTime,
			MaxHits:     jsonDir.MaxHits,
			WindowSize:  jsonDir.WindowSize,
			History:     hitcounter.DefaultHistory(),
		}

		switch jsonDir.Typ {
		case "string":
			dir.Store = hitcounter.NewStringMap(int64(jsonDir.MaxTracked))
		case "int32":
			dir.Store = hitcounter.NewInt32Map(int64(jsonDir.MaxTracked))
		}

		// Add it to the list
		result.Directions = append(result.Directions, dir)
	}

	return result, nil
}
