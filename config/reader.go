package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

// Config struct
type Config struct {
	Server serverConfig   `toml:"server"`
	DB     databaseConfig `toml:"database"`
}

// ReadConfig function
func ReadConfig(path string) (config Config, err error) {

	var conf Config
	if _, err := toml.DecodeFile(path, &conf); err != nil {
		log.Println(err)
		return conf, err
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
	Port      string `toml:"port"`
	Interface string `toml:"interface"`
	Debug     bool   `toml:"debug"`
}
