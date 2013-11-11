package helper

import (
	"os"
)

type Config struct {
	Host string
	Port string
}

func (self *Config) Init() {
	self.Host = GetOrDefaultString(os.Getenv("HOST"), "blackbeard.apertoire.org")
	self.Port = GetOrDefaultString(os.Getenv("PORT"), "8080")
}
