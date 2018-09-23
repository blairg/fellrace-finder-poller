package storage

import (
	"io/ioutil"
	"os"
	"strings"
)

// Store a html recipe file to the filesystem
func Store(filename, path string, html string) {
	fullPath := buildFilename(filename, path)
	lastIndexOfSlash := strings.LastIndex(fullPath, "/")
	pathToSaveTo := fullPath[0:lastIndexOfSlash]

	createPath(pathToSaveTo)

	if html != "" {
		d1 := []byte(html)

		if !Exists(fullPath) {
			err := ioutil.WriteFile(fullPath, d1, 0644)
			check(err)
		} else {
			//fmt.Println("already exists - " + fullPath)
		}
	} else {
		emptyFile := []byte("")
		err := ioutil.WriteFile(fullPath, emptyFile, 0644)

		if err != nil {
			panic(err)
		}
	}
}

// GetFromFileSystem returns the content of a saved file if its exists, otherwise blank string
func GetFromFileSystem(filename, path string) string {
	fullPath := buildFilename(filename, path)

	if Exists(fullPath) {
		fileData, err := ioutil.ReadFile(fullPath)

		if err != nil {
			panic(err)
		}

		return string(fileData)
	}

	return ""
}

// Exists checks if a file exists
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func createPath(directory string) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.MkdirAll(directory, os.ModePerm)
	}
}

func buildFilename(filename, path string) string {
	return path + strings.Replace(filename, " ", "", -1) + ".html"
}
