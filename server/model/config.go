package model

import (
	"encoding/json"
	"github.com/apertoire/mlog"
	"log"
	"os"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`

	AppDir string `json:"appDir"`

	MediaPath []string `json:"mediaPath"`
}

func (self *Config) Init() {
	self.Host = GetOrDefaultString(os.Getenv("HOST"), "")
	self.Port = GetOrDefaultString(os.Getenv("PORT"), "8080")

	// self.AppDir = "/Volumes/Users/kayak/Library/Application Support/net.apertoire.mediabase/"
	self.AppDir = "."
}

func (self *Config) Load() {
	file, _ := os.Open("./config.json")

	log.Println("file: ", file)

	decoder := json.NewDecoder(file)

	log.Println("decoder: ", decoder)

	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		mlog.Fatalf("Unable to load configuration: %s", err)
	}

	self.Host = config.Host
	self.Port = config.Port
	self.AppDir = config.AppDir
	self.MediaPath = config.MediaPath
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
