package mode

import (
	//"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//Chown file
// mode
// 	F （完全访问权限）,
//	M (修改访问权限）
//	RX（读取和执行访问权限）
//	R （只读）
// 	W （只写访问权限）
func Chown(fs, owner, group, mode string) {
	var modes []string
	for _, v := range strings.Split(mode, "|") {
		modes = append(modes, string(v))
	}
	var commandStr []string

	f, _ := os.Stat(fs)
	if runtime.GOOS == "windows" {
		commandStr = []string{"/C", "icacls", fs}
		if f.IsDir() {
			commandStr = append(commandStr, "/T")
		}
		if modes[0] != "-" {
			commandStr = append(commandStr, "/grant", owner+":"+modes[0])
		} else {
			panic("Must Config Mode User")
		}
		if modes[1] != "-" {
			commandStr = append(commandStr, group+":"+modes[1])
		}
		if modes[2] != "-" {
			commandStr = append(commandStr, "Everyone:"+modes[2])
		}
	}
	//fmt.Println(commandStr)

	// grant security and permission
	_, err := exec.Command("cmd", commandStr...).Output()
	if err != nil {
		panic(err)
	} //else {
	//	fmt.Println(string(d))
	//}
}

//reset file security and permission
func ResetMode(fs string) {
	var commandStr []string
	var resetStr []string
	f, _ := os.Stat(fs)
	if runtime.GOOS == "windows" {
		commandStr = []string{"/C", "icacls", fs}
		resetStr = append(commandStr, "/RESET")
		if f.IsDir() {
			resetStr = append(commandStr, "/T", "/RESET")
		}
		_, err := exec.Command("cmd", resetStr...).Output()
		if err != nil {
			panic(err)
		}
	}
}
