package env

import (
	"fmt"
	"os"
)

func GetString(name, defaultValue string) string {
	got := os.Getenv(name)
	if got == "" {
		return defaultValue
	}
	return got
}

func MustGetString(name string) string {
	got := os.Getenv(name)
	if got == "" {
		panic("environment variable " + name + " is required")
	}
	return got
}

func GetBool(name string, defaultValue bool) bool {
	got := os.Getenv(name)
	if got == "" {
		return defaultValue
	}
	if got == "true" || got == "1" {
		return true
	}
	return false
}

func MustGetBool(name string) bool {
	got := os.Getenv(name)
	if got == "" {
		panic("environment variable " + name + " is required")
	}
	if got == "true" || got == "1" {
		return true
	}
	return false
}

func GetInt(name string, defaultValue int) int {
	got := os.Getenv(name)
	if got == "" {
		return defaultValue
	}
	var intValue int
	_, err := fmt.Sscanf(got, "%d", &intValue)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func MustGetInt(name string) int {
	got := os.Getenv(name)
	if got == "" {
		panic("environment variable " + name + " is required")
	}
	var intValue int
	_, err := fmt.Sscanf(got, "%d", &intValue)
	if err != nil {
		panic("environment variable " + name + " must be an integer")
	}
	return intValue
}
