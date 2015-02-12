package config

import (
	"github.com/JamesOwenHall/BruteForceProtection/core/hitcounter"
)

// Configuration is a struct that represents the contents of a configuration
// file.
type Configuration struct {
	Directions    []hitcounter.Direction
	ListenAddress string
	ListenType    string
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
	result.Directions = make([]hitcounter.Direction, 0, len(parsed.Directions))

	for _, jsonDir := range parsed.Directions {
		// Create the direction according to its type
		var dir hitcounter.Direction
		switch jsonDir.Typ {
		case "string":
			dir = hitcounter.NewStringDirection(
				jsonDir.Name,
				jsonDir.WindowSize,
				jsonDir.MaxHits,
				jsonDir.CleanUpTime,
			)
		case "int32":
			dir = hitcounter.NewInt32Direction(
				jsonDir.Name,
				jsonDir.WindowSize,
				jsonDir.MaxHits,
				jsonDir.CleanUpTime,
			)
		}

		// Add it to the list
		result.Directions = append(result.Directions, dir)
	}

	return result, nil
}
