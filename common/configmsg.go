package common

import (
	"log"
	"os"
)

const (
	MaxDownActiveNum = 10000
)

var CommonSourceAddr = ""

func init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println("get dir err", err)
		return
	}
	CommonSourceAddr = path + "\\" + "common-data"
}
