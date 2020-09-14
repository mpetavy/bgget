package main

import (
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	directory *string
)

func init() {
	common.Init(true, "1.0.0", "2019", "gnome background get", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, run, time.Duration(60)*time.Second)

	directory = flag.String("d", "", "target path")
}

func run() error {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.background", "picture-uri")

	ba, err := cmd.Output()
	if common.Error(err) {
		return err
	}

	srcFile := string(ba)
	srcFile = srcFile[1 : len(srcFile)-2]
	srcFile = srcFile[7:]

	destFile := common.CleanPath(filepath.Join(*directory, filepath.Base(srcFile)))

	b, err := common.FileExists(destFile)
	if common.Error(err) {
		return err
	}

	if b {
		return nil
	}

	err = common.FileCopy(srcFile, destFile)
	if common.Error(err) {
		return err
	}

	fmt.Printf("src: %v\n", srcFile)
	fmt.Printf("dest: %v\n", destFile)

	return nil
}

func main() {
	defer common.Done()

	common.Run(nil)
}
