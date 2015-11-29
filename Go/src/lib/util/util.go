package util

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"regexp"
)

//ReadTxt ...
func ReadTxt(logPath string) ([]byte, error) {
	var r string
	flag := false
	file, err := os.OpenFile(logPath, os.O_RDONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		str, err := reader.ReadSlice('\n')

		if m, _ := regexp.MatchString("^\\/\\*", string(str)); m {
			flag = false
			continue
		}
		if m, _ := regexp.MatchString("^\\*\\/", string(str)); m {
			flag = true
			continue
		}

		if flag {
			r = r + " " + string(str)
		}

		if err != nil {
			break
		}
	}
	return []byte(r), err
}

//CheckFile ..
func CheckFile(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// ReadFile deal  json file to map
func ReadFile(filename string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return nil, err
	}

	return bytes, nil
}

//CheckUG ...
func CheckUG(userVar string) error {
	// Check all parameter ,owner,group,source is exist
	if _, err := user.Lookup(userVar); err != nil {
		return err
	}

	return nil
}

// CheckM ...
func CheckM(modeVar string) bool {
	modeVarR := []byte(modeVar)
	checkType := []byte("fwr|-")
	for _, t := range modeVarR {
		s := []byte{t}
		if !bytes.Contains(checkType, s) {
			return false
		}
	}
	return true
}

// ReplaceWindowsPathSeparator xxx
func ReplaceWindowsPathSeparator(path string) string {
	bs := regexp.MustCompile(`\\`)
	unix := bs.ReplaceAllLiteralString(path, "/")

	driveRe := regexp.MustCompile(`\A[A-Z]/`)
	return driveRe.ReplaceAllLiteralString(unix, "")
}

// IsDirectoryExisted xxx
func IsDirectoryExisted(dir string) bool {
	f, err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	return fi.Mode().IsDir()
}

// RemoveDuplicate array Duplicate item
func RemoveDuplicate(slis *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, val := range *slis {
		if _, ok := found[val]; !ok {
			found[val] = true
			(*slis)[j] = (*slis)[i]
			j++
		}
	}
	*slis = (*slis)[:j]
}
