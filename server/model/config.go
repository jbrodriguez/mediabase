package model

import (
	"encoding/json"
	"github.com/apertoire/mlog"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`

	DataDir string `json:"-"`

	MediaFolders []string `json:"mediaFolders"`
	MediaRegexs  []string `json:"-"`

	Version string `json:"version"`
}

func (self *Config) Init(version string) {
	self.Version = version

	runtime := runtime.GOOS

	switch runtime {
	case "darwin":
		self.DataDir = filepath.Join(os.Getenv("HOME"), "Library/Application Support/net.apertoire.mediabase/")
	case "linux":
		self.DataDir = filepath.Join(os.Getenv("HOME"), ".mediabase/")
	case "windows":
		self.DataDir = filepath.Join(os.Getenv("APPDATA"), "mediabase\\")
	}

	path := filepath.Join(self.DataDir, "log")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil {
			log.Printf("FATAL: Unable to create folder %s: %s", path, err)
			os.Exit(255)
		}
	}

	// os.Setenv("GIN_MODE", "release")
	mlog.Start(mlog.LevelInfo, filepath.Join(self.DataDir, "log", "mediabase.log"))
	mlog.Info("mediabase v%s starting up on %s ...", self.Version, runtime)

	self.setupFolders()

	self.LoadConfig()
	self.LoadRegex()
}

func (self *Config) setupFolders() {
	path := filepath.Join(self.DataDir, "db")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", path, err)
		}
	}

	path = filepath.Join(self.DataDir, "log")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.Mkdir(path, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", path, err)
		}
	}

	path = filepath.Join(self.DataDir, "web", "img", "b")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", path, err)
		}
	}

	path = filepath.Join(self.DataDir, "web", "img", "bt")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", path, err)
		}
	}

	path = filepath.Join(self.DataDir, "web", "img", "p")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", path, err)
		}
	}

	path = filepath.Join(self.DataDir, "web", "img", "pt")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", path, err)
		}
	}

	path = filepath.Join(self.DataDir, "web", "img", "t")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", path, err)
		}
	}
}

func (self *Config) LoadConfig() {
	path := filepath.Join(self.DataDir, "config.json")
	file, err := os.Open(path)
	if err != nil {
		mlog.Warning("Config file %s doesn't exist. Creating one ...", path)

		self.Host = ""
		self.Port = "3267"
		self.MediaFolders = make([]string, 0)

		self.Save()

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
	self.MediaFolders = config.MediaFolders
}

func (self *Config) LoadRegex() {
	path := filepath.Join(self.DataDir, "regex.txt")
	content, err := ioutil.ReadFile(path)
	if err != nil {
		mlog.Warning("Regex file %s doesn't exist. Creating one ...", path)

		self.MediaRegexs = []string{
			`(?i)(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`,
			`(?i)(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
			`(?i)(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`,
			`(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`,
			`(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
			`(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`,
			`(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
		}

		var file *os.File
		if file, err = os.Create(path); err != nil {
			mlog.Warning("Unable to write to %s: %s", path, err)
			return
		}
		defer file.Close()

		for _, item := range self.MediaRegexs {
			_, err := file.WriteString(item + "\n")
			if err != nil {
				mlog.Warning("Unable to write to %s: %s", path, err)
				break
			}
		}

		return
	}

	self.MediaRegexs = strings.Split(string(content), "\n")
}

func (self *Config) Save() {
	b, err := json.MarshalIndent(self, "", "   ")
	if err != nil {
		mlog.Info("couldn't marshal: %s", err)
		return
	}

	path := filepath.Join(self.DataDir, "config.json")
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		mlog.Info("WriteFileJson ERROR: %+v", err)
	}

	mlog.Info("saved as: %s", string(b))
}

// func GetOrDefaultString(ask string, def string) string {
// 	if ask != "" {
// 		return ask
// 	}
// 	return def
// }

// func GetOrDefaultInt(ask int, def int) int {
// 	if ask != 0 {
// 		return ask
// 	}
// 	return def
// }
