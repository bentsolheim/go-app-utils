package utils

import (
	"os"
)

func GetEnvOrDefault(name string, defaultValue string) string {
	value, isSet := os.LookupEnv(name)
	if !isSet {
		value = defaultValue
	}
	return value
}
