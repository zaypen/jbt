package core

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

type Executor interface {
	List([]string) map[string]Installation
	Check(map[string]Installation) map[string]Release
	Update(map[string]Release)
}

type List func([]string) map[string]Installation

type Install func(string, string) error

type executor struct {
	platform string
	list     List
	install  Install
}
