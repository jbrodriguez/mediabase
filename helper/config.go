package helper

import (
	"os"
)

type Config struct {
	Host string
	Port string

	AppDir string
}

func (self *Config) Init() {
	self.Host = GetOrDefaultString(os.Getenv("HOST"), "localhost")
	self.Port = GetOrDefaultString(os.Getenv("PORT"), "8080")

	self.AppDir = "~/Library/Application Support/net.apertoire.mediabase/"
}
