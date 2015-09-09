package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"lib/copy"
	"lib/mode"
	"lib/tar"
	"lib/util"
	"log"
	"os"
	"path"
	"strings"
)

const (
	cmd     = `C:\Program Files (x86)\WinRAR\WinRAR.exe`
	version = "2.0"
)

// Global user,group,allauth
type Global struct {
	USER    string
	GROUP   string
	ALLAUTH string
}

// Special define special mode
type Special map[string]string

//Detail define source --> destination
type Detail map[string]string

type apiR struct {
	Global  Global
	Special Special
	Detail  Detail
}

func main() {
	// check input parameter,-p package.rar
	if len(os.Args) < 2 {
		panic("Usage: distribute.exe -p xxxxx.rar[zip]")
	}
	dstPkg := flag.String("p", "C:/distribute004v1.0.zip", "for distribute file package")
	//dstPkg := os.Args[0]
	flag.Parse()

	dstPkgU := util.ReplaceWindowsPathSeparator(*dstPkg)
	//check file if not exist
	err := util.CheckFile(dstPkgU)

	if err != nil {
		log.Fatal(err)
	}

	// get config filename from package name
	confName := path.Base(dstPkgU)
	t1 := strings.Split(confName, ".")
	headerName := strings.Join(t1[:len(t1)-1], ".")
	cfgFile := headerName + ".json"
	fullCfgFile := path.Join(path.Dir(dstPkgU), cfgFile)

	// Change path to os temporary path
	os.Chdir(os.TempDir())

	// Create sub temporary directory
	tmpName, err := ioutil.TempDir(os.TempDir(), "Dst")
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(tmpName)
	//if crush clear temporary dir
	defer func() {
		os.RemoveAll(tmpName)
	}()

	fmt.Println("Begin uncompress file")

	// Uncompress pkg file to sub temporary path
	if path.Ext(dstPkgU) == ".rar" {
		if err := util.CheckFile(cmd); err != nil {
			fmt.Printf("Uncompress file error: %v", err)
			os.Exit(1)
		}
		if err := tar.UnRar(cmd, dstPkgU, tmpName); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := tar.UnZip(dstPkgU, tmpName); err != nil {
			log.Fatal(err)
		}
	}

	// Change dir to sub temporary dir
	os.Chdir(tmpName)

	// Check Json
	// Read json config file
	body, _ := util.ReadTxt(fullCfgFile)
	var r apiR

	if err := json.Unmarshal([]byte(body), &r); err != nil {
		log.Fatal(err)
	}

	// Global parameter
	ownerVar := r.Global.USER
	groupVar := r.Global.GROUP
	modeVar := r.Global.ALLAUTH

	//check user and group in window os
	fmt.Println("Check owner")
	if err := util.CheckUG(ownerVar); err != nil {
		log.Fatalf("The user [%v] is not found\n", ownerVar)
	}

	//check mode,only contain fwr-
	fmt.Println("Check mode")
	if !util.CheckM(modeVar) {
		log.Fatalf("AllAUTH %v FORMAT ERROR,ONLY f,w,r", modeVar)
	}

	//check source file(s)
	fmt.Println("Check src")
	var errPath []string
	for key, _ := range r.Special {
		if util.CheckFile(key) != nil {
			errPath = append(errPath, key)
		}
	}

	for key, _ := range r.Detail {
		if util.CheckFile(key) != nil {
			errPath = append(errPath, key)
		}
	}

	util.RemoveDuplicate(&errPath)

	if len(errPath) > 0 {
		log.Fatal("Src File(s) :" + strings.Join(errPath, ",") + " is not exist!!!")
	}

	// dependence config file
	// Chmod all file;
	fmt.Println("Begin Change Globle Mode")
	for src := range r.Detail {
		//fmt.Println(key)
		mode.ResetMode(src)
		mode.Chown(src, ownerVar, groupVar, modeVar)
	}

	fmt.Println("Begin Change Special Mode")
	for key, value := range r.Special {
		mode.Chown(key, ownerVar, groupVar, value)
	}

	fmt.Println("Begin copy files")
	// manager file
	//fmt.Println(tmpName)
	for key, value := range r.Detail {
		fi, _ := os.Stat(key)
		if fi.IsDir() {
			//fmt.Printf("Src:%v,Dst:%v\n", key, value)
			err := copy.Copy(key, path.Join(value, key))
			if err != nil {
				fmt.Printf("Copy Dir Error ---- %s\n", err)
			}
		} else {
			//fmt.Printf("Src:%v,Dst:%v\n", key, value)
			err := copy.Copy(key, path.Join(value, key))
			if err != nil {
				fmt.Printf("Copy File Error ---- %s\n", err)
			}
		}

	}

}
