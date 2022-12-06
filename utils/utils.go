package utils

import (
	"os"
	"strings"
)

func GetFileLastName(filename string) string {
	str := strings.Split(filename, ".")
	return str[1]
}
func IsExistPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
