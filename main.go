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
	inputPath  *string
	outputPath *string
	minSize    *int
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	common.Panic(err)

	common.Init("1.0.0", "", "", "2022", "Windows background image getter", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, time.Hour)

	inputPath = flag.String("i", filepath.Join(userHomeDir, "AppData", "Local", "Packages", "Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy", "LocalState", "Assets"), "directory to store the images")
	outputPath = flag.String("o", filepath.Join(userHomeDir, "bgget"), "directory to store the images")
	minSize = flag.Int("s", 1000000, "minimum size of image file")
}

func processImage(path string) error {
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
	filename := filepath.Join(*outputPath, hashStr+".jpg")

	if common.FileExists(filename) {
		return nil
	}

	err = os.MkdirAll(*outputPath, common.DefaultDirMode)
	if common.Error(err) {
		return err
	}

	err = os.WriteFile(filename, ba, common.DefaultFileMode)
	if common.Error(err) {
		return err
	}

	return nil
}

func run() error {
	fw, err := common.NewFilewalker(filepath.Join(*inputPath, "*"), false, false, func(path string, f os.FileInfo) error {
		if f.IsDir() || int(f.Size()) < *minSize {
			return nil
		}

		return processImage(path)
	})
	if common.Error(err) {
		return err
	}

	err = fw.Run()
	if common.Error(err) {
		return err
	}

	return nil
}

func main() {
	defer common.Done()

	common.Run(nil)
}
