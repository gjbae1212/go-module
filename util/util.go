package util

import (
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	modulePath  string
	httpMethods = []string{
		http.MethodGet, http.MethodHead,
		http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete,
		http.MethodConnect, http.MethodOptions,
		http.MethodTrace,
	}
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

func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func CheckHttpMethod(method string) bool {
	return StringInSlice(strings.ToUpper(method), httpMethods)
}
