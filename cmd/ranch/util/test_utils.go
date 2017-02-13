package util

import (
	"os"
)

const TEST_RANCHY = ".ranch.yaml"
const TEST_FILE_1 = ".ranch.test1.yaml"
const TEST_FILE_2 = ".ranch.test2.yaml"

func fileExists(path string) (bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	} else {
		// this isn't actually always false, there might have been an error
		// reading the file system, permission erorr, who knows...
		return false
	}
}

func removeFileIfExists(path string) {
	if fileExists(path) {
		os.Remove(path)
	}
}

func removeTestFiles() {
	// make sure files do not exist
	removeFileIfExists(TEST_RANCHY)
	removeFileIfExists(TEST_FILE_1)
	removeFileIfExists(TEST_FILE_2)
}
