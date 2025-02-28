package utils

import (
	"os"
	"slices"
)

func GetBoolFlag(name string, shortcut string) bool {
	flags := getFlags(name, shortcut)
	for _, arg := range os.Args {
		if slices.Contains(flags, arg) {
			return true
		}
	}
	return false
}

func GetStringFlag(name string, shortcut string) string {
	flags := getFlags(name, shortcut)
	for i, arg := range os.Args {
		if slices.Contains(flags, arg) && len(os.Args) > i+1 {
			return os.Args[i+1]
		}
	}
	return ""
}

func getFlags(name string, shortcut string) []string {
	if name == "" {
		return nil
	}
	if shortcut == "" {
		return []string{"--" + name}
	}
	return []string{"--" + name, "-" + shortcut}
}
