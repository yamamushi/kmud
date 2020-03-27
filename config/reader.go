package config

import (
	"errors"
	"flag"
	"github.com/BurntSushi/toml"
	"os"
	"sync"
)

// Config struct
type Config struct {
	Server serverConfig   `toml:"server"`
	DB     databaseConfig `toml:"database"`
	Crypt  cryptConfig    `toml:"crypt"`
}

var configquerylocker sync.Mutex

// read function
func read(path string) (config *Config, err error) {
	conf := &Config{}
	if _, err := toml.DecodeFile(path, conf); err != nil {
		return nil, err
	}

	return conf, nil
}

// GetConfig function
func GetConfig(name string) (conf *Config, err error) {
	configquerylocker.Lock()
	defer configquerylocker.Unlock()

	var confPath string
	flag.StringVar(&confPath, "c", name, "Path to Config File")
	flag.Parse()

	_, err = os.Stat(confPath)
	if err != nil {
		return nil, errors.New("Config file is missing: " + confPath)
	}

	conf, err = read(confPath)
	if err != nil {
		return nil, errors.New("Error reading config file: " + confPath + " - " + err.Error())
	}

	return conf, nil
}

// databaseConfig struct
type databaseConfig struct {
	DBFile    string `toml:"filename"`
	MongoHost string `toml:"mongohost"`
	MongoDB   string `toml:"mongodb"`
	MongoUser string `toml:"mongouser"`
	MongoPass string `toml:"mongopass"`
}

// serverConfig struct
type serverConfig struct {
	Port         string `toml:"port"`
	Interface    string `toml:"interface"`
	Debug        bool   `toml:"debug"`
	LoggingLevel string `toml:"verbosity"`
}

type cryptConfig struct {
	AccountManagerSecret string `toml:"account_manager_secret"`
}
