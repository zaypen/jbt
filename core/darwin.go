package core

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"howett.net/plist"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"os/user"
	"io/ioutil"
	"github.com/zaypen/jbt/util"
)

var fileNames = map[string]string{
	AppCode:   "AppCode.app",
	CLion:     "CLion.app",
	DataGrip:  "DataGrip.app",
	Idea:      "IntelliJ IDEA.app",
	IdeaCE:    "IntelliJ IDEA CE.app",
	PyCharm:   "PyCharm.app",
	PyCharmCE: "PyCharm CE.app",
	PhpStorm:  "PhpStorm.app",
	RubyMine:  "RubyMine.app",
	WebStorm:  "WebStorm.app",
}

var volumeNames = map[string]string{
	AppCode:   "AppCode",
	CLion:     "CLion",
	DataGrip:  "DataGrip",
	Idea:      "IntelliJ IDEA",
	IdeaCE:    "IntelliJ IDEA CE",
	PyCharm:   "PyCharm",
	PyCharmCE: "PyCharm CE",
	PhpStorm:  "PhpStorm",
	RubyMine:  "RubyMine",
	WebStorm:  "WebStorm",
}

type bundleInfo struct {
	Version string `plist:"CFBundleShortVersionString"`
}

func detectInstallation(code string) (bool, string) {
	productName := ProductNames[code]
	location := path.Join(string(os.PathSeparator), "Applications", fileNames[code])
	logrus.Debugf("Detecting installation of %s...", productName)
	if !exists(location) {
		return false, ""
	}
	logrus.Debugf("Reading version of %s...", productName)
	if version, err := readVersion(path.Join(location, "Contents", "Info.plist")); err == nil {
		return true, *version
	} else {
		logrus.Warnf("Error while reading version of %s, %v", productName, err.Error())
		return true, "Unknown"
	}
}

func readVersion(name string) (*string, error) {
	var info bundleInfo
	f, err := os.Open(name)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open %s", name))
	}
	defer f.Close()
	if err = plist.NewDecoder(f).Decode(&info); err == nil {
		return &info.Version, nil
	} else {
		return nil, errors.New(fmt.Sprintf("failed to decode %s", name))
	}
}

func darwinList(codes []string) map[string]Installation {
	installations := make(map[string]Installation)
	for _, code := range codes {
		installed, version := detectInstallation(code)
		installations[code] = Installation{installed, version}
	}
	return installations
}

func copyWithProgress(source string, destination string) error {
	if totalSize, err := dirSize(source); err == nil {
		writer := newProgressWriter(totalSize)
		return copyDir(source, destination, writer)
	} else {
		return err
	}
}

func copyDir(source string, destination string, writer *ProgressWriter) error {
	if fi, err := os.Stat(source); err == nil {
		if err = os.MkdirAll(destination, fi.Mode()); err == nil {
			if fis, err := ioutil.ReadDir(source); err == nil {
				for _, fi := range fis {
					s := filepath.Join(source, fi.Name())
					d := filepath.Join(destination, fi.Name())
					if err := util.Iff(fi.IsDir(), func() interface{} {
						return copyDir(s, d, writer)
					}, func() interface{} {
						return copyFile(s, d, writer)
					}); err != nil {
						return err.(error)
					}
				}
				return nil
			} else {
				return err
			}
		} else {
			return err
		}
	} else {
		return err
	}
}
func copyFile(source string, destination string, writer *ProgressWriter) error {
	if sf, err := os.Open(source); err == nil {
		defer
	}
}

func dirSize(location string) (size uint64, err error) {
	err = filepath.Walk(location, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			size += uint64(fi.Size())
		}
		return nil
	})
	return
}

func darwinInstall(code string, file string) (err error) {
	productName := ProductNames[code]
	fmt.Println(fmt.Sprintf("Mounting %s...", filepath.Base(file)))
	if err = exec.Command("hdiutil", "attach", file).Run(); err == nil {
		volume := filepath.Join(string(filepath.Separator), "Volumes", volumeNames[code])
		source := filepath.Join(volume, fileNames[code])
		destination := path.Join(string(os.PathSeparator), "Applications", fileNames[code])
		backup := path.Join(user.Current(), ".jbt_backup", fileNames[code])
		fmt.Println(fmt.Sprintf("Backuping %s...", productName))
		if err = os.Rename(destination, backup); err == nil {
			defer func() {
				if err != nil {
					err = os.Rename(backup, destination)
				}
			}()
		}
		if err := exec.Command("hdiutil", "detach", volume).Run(); err == nil {
			return nil
		}
	}
	return
}
