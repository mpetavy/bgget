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
	user      *string
	directory *string
)

func init() {
	common.Init(true, "1.0.0", "", "2019", "gnome background get", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, run, time.Duration(60)*time.Second)

	user = flag.String("u", "", "run as user")
	directory = flag.String("d", "", "target path")
}

func run() error {
	cmd := exec.Command("runuser", "-l", *user, "-c", "gsettings get org.gnome.desktop.background picture-uri")

	ba, err := cmd.Output()
	if common.Error(err) {
		return err
	}

	srcFile := string(ba)
	srcFile = srcFile[1 : len(srcFile)-2]
	srcFile = srcFile[7:]

	common.Info("Found: %v", srcFile)

	destFile := common.CleanPath(filepath.Join(*directory, filepath.Base(srcFile)))

	b, err := common.FileExists(destFile)
	if common.Error(err) {
		return err
	}

	common.Info("Exists: %v", b)

	if b {
		return nil
	}

	err = common.FileCopy(srcFile, destFile)
	if common.Error(err) {
		return err
	}

	common.Info("Saved: %v", destFile)

	cmd = exec.Command("chown", fmt.Sprintf("%v:%v", *user, *user), destFile)

	err = cmd.Run()
	if common.Error(err) {
		return err
	}

	return nil
}

func main() {
	defer common.Done()

	common.Run(nil)
}
