package config

import "os"

func createDir(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.Mkdir(path, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	return err
}
