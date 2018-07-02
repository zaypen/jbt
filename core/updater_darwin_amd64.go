package core

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"howett.net/plist"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

var fileNames = map[string]string{
	AppCode:   "AppCode.app",
	CLion:     "CLion.app",
	DataGrip:  "DataGrip.app",
	GoLand:    "GoLand.app",
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
	GoLand:    "GoLand",
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
	location := filepath.Join(string(os.PathSeparator), "Applications", fileNames[code])
	logrus.Debugf("Detecting installation of %s...", productName)
	if !exists(location) {
		return false, ""
	}
	logrus.Debugf("Reading version of %s...", productName)
	if version, err := readVersion(filepath.Join(location, "Contents", "Info.plist")); err == nil {
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

func list(codes []string) map[string]Installation {
	installations := make(map[string]Installation)
	for _, code := range codes {
		installed, version := detectInstallation(code)
		installations[code] = Installation{installed, version}
	}
	return installations
}

func doCopy(source string, destination string) (err error) {
	if totalSize, err := sizeOf(source); err == nil {
		writer := newProgressWriter(totalSize)
		defer writer.Clear()
		mask := syscall.Umask(0)
		defer syscall.Umask(mask)
		return copyDir(source, destination, writer)
	}
	return
}

func install(code string, file string) (err error) {
	productName := ProductNames[code]
	fmt.Println(fmt.Sprintf("Mounting %s...", filepath.Base(file)))
	if err = exec.Command("hdiutil", "attach", file).Run(); err == nil {
		volume := filepath.Join(string(filepath.Separator), "Volumes", volumeNames[code])
		defer exec.Command("hdiutil", "detach", volume).Run()
		source := filepath.Join(volume, fileNames[code])
		destination := filepath.Join(string(os.PathSeparator), "Applications", fileNames[code])
		backup := filepath.Join(string(os.PathSeparator), "Applications", fileNames[code]+".old")
		fmt.Println(fmt.Sprintf("Backuping %s...", productName))
		if err = os.Rename(destination, backup); err == nil {
			defer func() {
				if err == nil {
					fmt.Println(fmt.Sprintf("Cleaning backup of %s...", productName))
					if err = os.RemoveAll(backup); err == nil {
						fmt.Println(fmt.Sprintf("Cleaning temporary files..."))
						err = os.Remove(file)
					}
				} else {
					fmt.Println(fmt.Sprintf("Restoring backup of %s...", productName))
					err = os.Rename(backup, destination)
				}
			}()
			fmt.Println(fmt.Sprintf("Installing %s...", productName))
			return doCopy(source, destination)
		}
	}
	return
}
