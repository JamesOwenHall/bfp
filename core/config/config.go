package config

import (
	"fmt"
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
func ReadConfig(filename string) (*Configuration, error) {
	parsed, err := parseJsonFile(filename)
	if err != nil {
		return nil, err
	}
	if len(parsed.Directions) == 0 {
		return nil, fmt.Errorf("no directions defined.")
	}
	if parsed.ListenType != "unix" && parsed.ListenType != "tcp" {
		return nil, fmt.Errorf(`unknown listen type "%s"`, parsed.ListenType)
	}

	result := new(Configuration)
	result.ListenAddress = parsed.ListenAddress
	result.ListenType = parsed.ListenType
	result.Directions = make([]hitcounter.Direction, 0, len(parsed.Directions))

	for _, jsonDir := range parsed.Directions {
		// Validate the input (only positive numbers allowed)
		if jsonDir.MaxHits <= 0 || jsonDir.WindowSize <= 0 {
			return nil, fmt.Errorf(`direction named "%s" contains non-positive parameters.`, jsonDir.Name)
		}

		// Set defaults
		if int32(jsonDir.CleanUpTime) == 0 {
			jsonDir.CleanUpTime = 5
		}

		// Create the direction according to its type
		var dir hitcounter.Direction
		switch jsonDir.Typ {
		case "string":
			dir = hitcounter.NewStringDirection(
				jsonDir.Name,
				int32(jsonDir.WindowSize),
				int32(jsonDir.MaxHits),
				int32(jsonDir.CleanUpTime),
			)
		case "int32":
			dir = hitcounter.NewInt32Direction(
				jsonDir.Name,
				int32(jsonDir.WindowSize),
				int32(jsonDir.MaxHits),
				int32(jsonDir.CleanUpTime),
			)
		default:
			return nil, fmt.Errorf("invalid direction type %s.", jsonDir.Typ)
		}

		// Add it to the list
		result.Directions = append(result.Directions, dir)
	}

	return result, nil
}
