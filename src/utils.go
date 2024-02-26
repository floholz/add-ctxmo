package src

import (
	"fmt"
	"github.com/flytam/filenamify"
	"golang.org/x/sys/windows"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func SaveToDisk(file []byte, fpath string) error {
	fpath, err := Pathify(fpath)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
	if err != nil {
		return err
	}
	err = os.WriteFile(fpath, file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func Pathify(fpath string) (string, error) {
	fpath, err := filepath.Abs(fpath)
	if err != nil {
		return "", err
	}
	dir, file := filepath.Split(fpath)
	file, err = filenamify.Filenamify(file, filenamify.Options{Replacement: "_"})
	if err != nil {
		return "", err
	}
	file = strings.ReplaceAll(file, " ", "_")
	fpath = filepath.Join(dir, file)
	fpath = filepath.Clean(fpath)
	return fpath, nil
}

func IsAdmin() bool {
	return windows.GetCurrentProcessToken().IsElevated()
}

func RelaunchWithElevatedPerms() {
	verb := "runas"
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()
	args := strings.Join(os.Args[1:], " ")

	verbPtr, _ := syscall.UTF16PtrFromString(verb)
	exePtr, _ := syscall.UTF16PtrFromString(exe)
	cwdPtr, _ := syscall.UTF16PtrFromString(cwd)
	argPtr, _ := syscall.UTF16PtrFromString(args)

	var showCmd int32 = 1 //SW_NORMAL

	err := windows.ShellExecute(0, verbPtr, exePtr, argPtr, cwdPtr, showCmd)
	if err != nil {
		fmt.Println(err)
	}

	// close this instance
	os.Exit(0)
}
