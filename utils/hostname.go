package utils

import (
	"os"
)

func GetHostName() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return name, nil
}
