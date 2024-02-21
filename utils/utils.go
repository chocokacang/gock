package utils

import (
	"log"
	"os"
	"path"
	"reflect"
	"runtime"
	"strings"
)

func GetEnv(key string, d string) string {
	val := os.Getenv(key)
	if val == "" {
		return d
	}
	return val
}

func GetEnvBool(key string, d bool) bool {
	val := os.Getenv(key)
	if val == "" {
		return d
	}
	if strings.ToLower(val) == "true" {
		return true
	}

	return false
}

func GetFunctionName(f any) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}

func LastChar(str string) uint8 {
	if str == "" {
		log.Panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

func JointPath(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if LastChar(relativePath) == '/' && LastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}
