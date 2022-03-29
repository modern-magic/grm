package internal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Actions interface {
	ShowList()
}

func getCurrentRegisitry() string {
	ReadFile(Npmrc)
	return ""
}

func ShowList() {
	showList()
}

func showList() {
	getCurrentRegisitry()
	getPresetRegistries()
}

type Registries struct {
	Npm       RegistryInner `json:"npm"`
	Yarn      RegistryInner `json:"yarn"`
	Tencent   RegistryInner `json:"tencent"`
	Cnpm      RegistryInner `json:"cnpm"`
	Taobao    RegistryInner `json:"taobao"`
	NpmMirror RegistryInner `json:"npmMirror"`
}
type RegistryInner struct {
	Home     string `json:"home"`
	Registry string `json:"registry"`
}

func getPresetRegistries() {
	presets, _ := ioutil.ReadFile("../registries.json")

	registries := &Registries{}
	json.Unmarshal(presets, registries)
	fmt.Printf("%+v\n", registries)

}
