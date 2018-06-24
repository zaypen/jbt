package core

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"howett.net/plist"
	"os"
	"os/exec"
	"path"
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

type bundleInfo struct {
	Version string `plist:"CFBundleShortVersionString"`
}

func detectInstallation(code string) (bool, string) {
	productName := ProductNames[code]
	fileName := fileNames[code]
	location := path.Join(string(os.PathSeparator), "Applications", fileName)
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

func darwinInstall(_ string, file string) error {
	cmd := exec.Command("open", file)
	return cmd.Run()
}
