package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base32"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	srcPath *string
	dstPath *string
	timeout *int
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	common.Panic(err)

	common.Init(true, "1.0.0", "", "", "2022", "Windows background image getter", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, time.Hour)

	srcPath = flag.String("src", filepath.Join(userHomeDir, "AppData", "Local", "Packages", "Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy", "LocalState", "Assets"), "directory to store the images")
	dstPath = flag.String("dst", filepath.Join(userHomeDir, "bgget"), "directory to store the images")
	timeout = flag.Int("timeout", 3600000, "timeout to look for new images")

	common.Events.NewFuncReceiver(common.EventFlagsParsed{}, func(event common.Event) {
		common.App().RunTime = common.MillisecondToDuration(*timeout)

		common.Panic(os.MkdirAll(*dstPath, common.DefaultDirMode))
	})
}

func processImage(path string) error {
	fi, err := os.Stat(path)
	if common.Error(err) {
		return err
	}

	if fi.IsDir() {
		return nil
	}

	ba, err := os.ReadFile(path)
	if common.Error(err) {
		return err
	}

	hash := md5.New()
	_, err = io.Copy(hash, bytes.NewReader(ba))
	if common.Error(err) {
		return err
	}

	hashStr := base32.StdEncoding.EncodeToString(hash.Sum(nil))
	filename := filepath.Join(*dstPath, hashStr+".jpg")

	if common.FileExists(filename) {
		return nil
	}

	err = os.WriteFile(filename, ba, common.DefaultFileMode)
	if common.Error(err) {
		return err
	}

	return nil
}

func run() error {
	fw := common.NewFilewalker(filepath.Join(*srcPath, "*"), false, false, func(path string) error {
		return processImage(path)
	})

	err := fw.Run()
	if common.Error(err) {
		return err
	}

	return nil
}

func main() {
	defer common.Done()

	common.Run(nil)
}
