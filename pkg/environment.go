package lemin

import (
	"os"
	"strconv"
)

type Environment struct {
	Port             uint16
	EnableVisualizer bool
	NoPathFinding    bool
	NoSanityCheck    bool
}

func ParseBooleanEnv(key string, def bool) bool {
	valueString, exists := os.LookupEnv(key)
	if exists && (valueString != "no" && valueString != "true") {
		return true
	}
	return def
}

func ParseEnvironment() Environment {

	portString := os.Getenv("PORT")
	portU64, err := strconv.ParseUint(portString, 10, 16)
	if err != nil {
		portU64 = 8080
	}

	return Environment{
		Port:             uint16(portU64),
		EnableVisualizer: ParseBooleanEnv("VISUALIZER", false),
		NoPathFinding:    ParseBooleanEnv("NO_PATH_FINDING", false),
		NoSanityCheck:    ParseBooleanEnv("NO_SANITY_CHECK", false),
	}
}
