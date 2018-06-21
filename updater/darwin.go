package updater

import (
	"errors"
	"fmt"
	"howett.net/plist"
	"os"
	"path"
	"sort"
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

type darwinUpdater struct {
	baseUpdater
}

func detectInstallation(location string) (bool, string) {
	installed := exists(location)
	if !installed {
		return false, ""
	}
	version, err := readVersion(path.Join(location, "Contents", "Info.plist"))
	if err != nil {
		fmt.Printf("WARNING: %s\n", err.Error())
		return true, "Unknown"
	}
	return true, *version
}

func readVersion(name string) (*string, error) {
	type bundleInfo struct {
		Version string `plist:"CFBundleShortVersionString"`
	}
	var info bundleInfo
	f, err := os.Open(name)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed opening file %s", name))
	}
	defer f.Close()
	decoder := plist.NewDecoder(f)
	err = decoder.Decode(&info)
	if err != nil {
		return nil, errors.New("failed decoding info")
	}
	return &info.Version, nil
}

func (updater *darwinUpdater) List() map[string]Installation {
	installations := make(map[string]Installation)
	for key, fileName := range fileNames {
		location := path.Join(string(os.PathSeparator), "Applications", fileName)
		installed, version := detectInstallation(location)
		installations[key] = Installation{installed, version}
	}
	return installations
}

func (updater *darwinUpdater) Check(installations map[string]Installation) map[string]Release {
	updates := make(map[string]Release, len(installations))
	for key, installation := range installations {
		if installation.Installed {
			releases, err := fetchReleases(key)
			if err != nil {
				fmt.Printf("WARNING: %s\n", err.Error())
			}
			if len(releases) > 0 {
				sort.Slice(releases, func(i, j int) bool {
					return compareVersion(releases[i].Version, releases[j].Version) > 0
				})
				release := releases[0]
				if compareVersion(installation.Version, release.Version) < 0 {
					updates[key] = release
				}
			}
		}
	}
	return updates
}
