package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/zaypen/jbt/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func subVersion(array []string, index int) string {
	if index < len(array) {
		return array[index]
	}
	return "0"
}

func compareVersion(a string, b string) int {
	arrayA, arrayB := strings.Split(a, "."), strings.Split(b, ".")
	lenA, lenB := len(arrayA), len(arrayB)
	n := util.If(lenA > lenB, lenA, lenB).(int)
	for i := 0; i < n; i++ {
		subA := util.Atoi(subVersion(arrayA, i))
		subB := util.Atoi(subVersion(arrayB, i))
		if subA != subB {
			return subA - subB
		}
	}
	return 0
}

func fetchReleases(code string) ([]Release, error) {
	if resp, err := http.Get(fmt.Sprintf(fmtReleaseEntry, code)); err == nil {
		defer resp.Body.Close()
		response := make(map[string][]Release)
		err = json.NewDecoder(resp.Body).Decode(&response)
		return response[code], err
	} else {
		return nil, err
	}
}

func latestOf(releases []Release) Release {
	sort.Slice(releases, func(i, j int) bool {
		return compareVersion(releases[i].Version, releases[j].Version) > 0
	})
	return releases[0]
}

func checkUpdate(code string, version string) (*Release, bool) {
	productName := ProductNames[code]
	logrus.Debugf("Fetching releases of %s...", productName)
	if releases, err := fetchReleases(code); err == nil && len(releases) > 0 {
		if release := latestOf(releases); compareVersion(version, release.Version) < 0 {
			return &release, true
		}
	} else {
		logrus.Warnf("Error while fetching releases of %s, %v", productName, err)
	}
	return nil, false
}

func doDownload(url string, file string) (*string, error) {
	if out, err := os.Create(file); err == nil {
		defer out.Close()
		if resp, err := http.Get(url); err == nil {
			defer resp.Body.Close()
			hash := sha256.New()
			progress := newProgressWriter(uint64(resp.ContentLength))
			writer := io.MultiWriter(out, hash, progress)
			if _, err = io.Copy(writer, resp.Body); err == nil {
				defer progress.Clear()
				checksum := hex.EncodeToString(hash.Sum(nil))
				return &checksum, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func readChecksum(checksumFile string) (*string, error) {
	if checksumBytes, err := ioutil.ReadFile(checksumFile); err == nil {
		if checksumStrings := strings.Split(string(checksumBytes), " "); len(checksumStrings) > 0 {
			return &checksumStrings[0], nil
		} else {
			return nil, errors.New(fmt.Sprintf("empty checksum %s", checksumFile))
		}
	} else {
		return nil, err
	}
}

func download(platform string, code string, release Release) (*string, error) {
	productName := ProductNames[code]
	if downloads, ok := release.Downloads[platform]; ok {
		fmt.Printf("Updating %s to %s...\n", productName, release.Version)
		productLink, checksumLink := downloads.Link, downloads.ChecksumLink
		productTemp := filepath.Join(os.TempDir(), filepath.Base(productLink))
		logrus.Debugf("Downloading release of %s to %s...", productName, productTemp)
		if productChecksum, err := doDownload(productLink, productTemp); err == nil {
			checksumTemp := filepath.Join(os.TempDir(), filepath.Base(checksumLink))
			logrus.Debugf("Downloading checksum of %s...", productName)
			if _, err := doDownload(checksumLink, checksumTemp); err == nil {
				if checksum, err := readChecksum(checksumTemp); err == nil {
					if *productChecksum == *checksum {
						fmt.Printf("Finished downloading of %s: %s\n", productName, productTemp)
						return &productTemp, nil
					} else {
						return nil, errors.New(fmt.Sprintf("checksum not match: %s != %s", *productChecksum, *checksum))
					}
				} else {
					return nil, err
				}
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("no download for platform %s", platform))
	}
}

func copyFile(source string, destination string, writer *ProgressWriter) (err error) {
	if fi, err := os.Lstat(source); err == nil {
		if fi.Mode().IsRegular() {
			if sf, err := os.Open(source); err == nil {
				defer sf.Close()
				if df, err := os.Create(destination); err == nil {
					defer df.Close()
					writer := io.MultiWriter(df, writer)
					if _, err = io.Copy(writer, sf); err == nil {
						if si, err := os.Stat(source); err == nil {
							return os.Chmod(destination, si.Mode())
						}
					}
				}
			}
		} else {
			if target, err := os.Readlink(source); err == nil {
				return os.Symlink(target, destination)
			}
		}
	}
	return
}

func copyDir(source string, destination string, writer *ProgressWriter) (err error) {
	if fi, err := os.Lstat(source); err == nil {
		if err = os.MkdirAll(destination, fi.Mode()); err == nil {
			if fis, err := ioutil.ReadDir(source); err == nil {
				for _, fi := range fis {
					s := filepath.Join(source, fi.Name())
					d := filepath.Join(destination, fi.Name())
					if fi.IsDir() {
						if err = copyDir(s, d, writer); err != nil {
							return err
						}
					} else {
						if err = copyFile(s, d, writer); err != nil {
							return err
						}
					}
				}
			}
		}
	}
	return
}

func sizeOf(location string) (size uint64, err error) {
	err = filepath.Walk(location, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() && fi.Mode().IsRegular() {
			size += uint64(fi.Size())
		}
		return nil
	})
	return
}

func (e *executor) List(codes []string) map[string]Installation {
	return e.list(codes)
}

func (e *executor) Check(installations map[string]Installation) map[string]Release {
	updates := make(map[string]Release)
	for code, installation := range installations {
		if installation.Installed {
			if release, ok := checkUpdate(code, installation.Version); ok {
				updates[code] = *release
			}
		}
	}
	return updates
}

func (e *executor) Update(releases map[string]Release) {
	for code, release := range releases {
		if file, err := download(e.platform, code, release); err == nil {
			if err := e.install(code, *file); err != nil {
				fmt.Printf("Error: %v\n", err.Error())
				break
			}
		} else {
			fmt.Printf("Error: %v\n", err.Error())
			break
		}
	}
}

func New(os string) (Executor, error) {
	switch os {
	case "darwin":
		return &executor{"mac", list, install}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Not supported OS: %s", os))
	}
}
