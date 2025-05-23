package main

import (
	"bytes"
	"crypto/md5"
	"embed"
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
	user       = flag.String("u", "", "user of which images shall be taken")
	inputPath  = flag.String("i", "", "directory to read the images")
	outputPath = flag.String("o", "", "directory to store the images")
	minSize    = flag.Int("s", 1000000, "minimum size of image file")
)

//go:embed go.mod
var resources embed.FS

func init() {
	common.Init("", "", "", "", "Windows background image getter", "", "", "", &resources, nil, nil, run, time.Minute*5)

	if !common.IsWindows() {
		common.Panic(fmt.Errorf("Runs only on Windows"))
	}

	common.Events.AddListener(common.EventFlags{}, func(event common.Event) {
		if *inputPath == "" {
			*inputPath = filepath.Join("c:"+string(os.PathSeparator), "users", *user, "AppData", "Local", "Packages", "Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy", "LocalState", "Assets")
		}
	})
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
	err := common.WalkFiles(filepath.Join(*inputPath, "*"), false, false, func(path string, f os.FileInfo) error {
		if f.IsDir() || int(f.Size()) < *minSize {
			return nil
		}

		return processImage(path)
	})
	if common.Error(err) {
		return err
	}

	return nil
}

func main() {
	common.Run([]string{"u", "o"})
}
