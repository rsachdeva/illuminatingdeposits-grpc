package testcredentials

import (
	"log"
	"path/filepath"
	"runtime"
)

func Path(rel string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(currentFile)
	pa := filepath.Join(basepath, rel)
	log.Printf("pa for test is %v", pa)
	return pa
}
