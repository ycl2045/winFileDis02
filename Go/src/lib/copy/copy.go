package copy

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

//CopyFile Copies file source to destination dest.
func CopyFile(source string, dest string) (err error) {
	sf, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	if err == nil {
		si, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, si.Mode())
			if err != nil {
				panic(err)
			}
		}

	}

	return
}

// CopyDir copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return &CustomError{"Source is not a directory"}
	}

	// ensure dest dir does not already exist

	//_, err = os.Open(dest)
	//if !os.IsNotExist(err) {
	//	return &CustomError{dest + "Destination already exists"}
	//}

	// create dest dir

	err = os.MkdirAll(dest, fi.Mode())
	if err != nil {
		return err
	}

	entries, err := ioutil.ReadDir(source)

	for _, entry := range entries {

		sfp := path.Join(source, entry.Name())
		dfp := path.Join(dest, entry.Name())
		if entry.IsDir() {
			err = CopyDir(sfp, dfp)
			if err != nil {
				log.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sfp, dfp)
			if err == nil {
				log.Println(err)
			}
		}

	}
	return
}

//CustomError A struct for returning custom error messages
type CustomError struct {
	What string
}

// Returns the error message defined in What as a string
func (e *CustomError) Error() string {
	return e.What
}

// Copy Deal copy like os

func Copy(source string, dest string) (err error) {
	// get properties of source dir
	argv := []string{}

	fi, err := os.Stat(source)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		argv = append(argv, "echo", "D|")
	} else {
		argv = append(argv, "echo", "F|")
	}

	dest_win := strings.Replace(dest, "/", "\\\\", -1)
	src_win := strings.Replace(source, "/", "\\\\", -1)

	argv = append(argv, "xcopy")
	argv = append(argv, src_win, dest_win)
	argv = append(argv, "/E", "/C", "/H", "/R", "/K", "/O", "/Y", "/Q")
	argvs := strings.Join(argv, " ")

	c, err := exec.Command("cmd", "/C", argvs).Output()
	if err != nil {
		return err
	}
	fmt.Println(string(c))
	return
}
