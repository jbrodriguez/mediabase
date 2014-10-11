package model

import (
	"encoding/json"
	"github.com/apertoire/mlog"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`

	AppDir string `json:"appDir"`

	MediaFolders []string `json:"mediaFolders"`
	MediaRegexs  []string `json:"-"`
}

func (self *Config) Init() {
	self.Host = GetOrDefaultString(os.Getenv("HOST"), "")
	self.Port = GetOrDefaultString(os.Getenv("PORT"), "8080")

	// self.AppDir = "/Volumes/Users/kayak/Library/Application Support/net.apertoire.mediabase/"
	self.AppDir = "."
}

func (self *Config) Load() {
	file, err := os.Open("./config.json")
	if err != nil {
		mlog.Fatalf("unable to open config.json: %s", err)
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		mlog.Fatalf("Unable to load configuration: %s", err)
	}

	self.Host = config.Host
	self.Port = config.Port
	self.AppDir = config.AppDir
	self.MediaFolders = config.MediaFolders

	content, err := ioutil.ReadFile("./regex.txt")
	if err != nil {
		mlog.Fatalf("unable to open regex.txt: %s", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	self.MediaRegexs = lines
}

func (self *Config) Save() {
	b, err := json.MarshalIndent(self, "", "   ")
	if err != nil {
		mlog.Info("couldn't marshal: %s", err)
		return
	}

	err = ioutil.WriteFile("./config.json.tmp", b, 0644)
	if err != nil {
		mlog.Info("WriteFileJson ERROR: %+v", err)
	}

	mlog.Info("saved as: %s", string(b))
}

func GetOrDefaultString(ask string, def string) string {
	if ask != "" {
		return ask
	}
	return def
}

func GetOrDefaultInt(ask int, def int) int {
	if ask != 0 {
		return ask
	}
	return def
}
