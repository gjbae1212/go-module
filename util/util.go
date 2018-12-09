package util

import (
	"log"
	"path"
	"path/filepath"
	"runtime"
)

var (
	modulePath string
)

func init() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("No caller information")
	}
	modulePath = filepath.Join(path.Dir(filename), "../")
}

func GetModulePath() string {
	return modulePath
}
