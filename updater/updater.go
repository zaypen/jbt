package updater

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/zaypen/jbt/util"
)

const fmtReleaseEntry = "https://data.services.jetbrains.com/products/releases?platform=&code=%s"

const (
	AppCode   = "AC"
	CLion     = "CL"
	DataGrip  = "DG"
	Idea      = "IIU"
	IdeaCE    = "IIC"
	PhpStorm  = "PS"
	PyCharm   = "PCP"
	PyCharmCE = "PCC"
	RubyMine  = "RM"
	WebStorm  = "WS"
)

var ProductCodes = []string{AppCode, CLion, DataGrip, Idea, IdeaCE, PhpStorm, PyCharm, PyCharmCE, RubyMine, WebStorm}

var ProductNames = map[string]string{
	AppCode:   "AppCode",
	CLion:     "CLion",
	DataGrip:  "DataGrip",
	Idea:      "IntelliJ IDEA",
	IdeaCE:    "IntelliJ IDEA CE",
	PhpStorm:  "PhpStorm",
	PyCharm:   "PyCharm",
	PyCharmCE: "PyCharm CE",
	RubyMine:  "RubyMine",
	WebStorm:  "WebStorm",
}

type Installation struct {
	Installed bool
	Version   string
}

type Download struct {
	Link         string `json:"link"`
	Size         int    `json:"size"`
	ChecksumLink string `json:"checksumLink"`
}

type Release struct {
	Version   string              `json:"version"`
	Downloads map[string]Download `json:"downloads"`
}

type Updater interface {
	List() map[string]Installation
	Check(map[string]Installation) map[string]Release
}

type baseUpdater struct {
	Platform string
}

func exists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func atoi(a string) int {
	i, err := strconv.Atoi(a)
	if err != nil {
		return 0
	}
	return i
}

func compareVersion(a string, b string) int {
	arrayA, arrayB := strings.Split(a, "."), strings.Split(b, ".")
	lenA, lenB := len(arrayA), len(arrayB)
	n := util.If(lenA < lenB, lenA, lenB).(int)
	for i := 0; i < n; i++ {
		subA, subB := atoi(arrayA[i]), atoi(arrayB[i])
		if subA != subB {
			return subA - subB
		}
	}
	return 0
}

func fetchReleases(code string) ([]Release, error) {
	resp, err := http.Get(fmt.Sprintf(fmtReleaseEntry, code))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	response := make(map[string][]Release)
	err = json.NewDecoder(resp.Body).Decode(&response)
	return response[code], err
}

func NewUpdater(os string) (Updater, error) {
	switch os {
	case "darwin":
		return &darwinUpdater{baseUpdater{Platform: "mac"}}, nil
	default:
		return nil, errors.New(fmt.Sprintf("Not support %s", os))
	}
}
