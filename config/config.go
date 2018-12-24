package config

import (
	"github.com/bcext/gcash/chaincfg"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

const (
	ProjectLastDir = "/huhongjia/bitcoin-abc-tech"
)

func GetChainParam() *chaincfg.Params {
	conf := GetConf()
	if conf.TestNet {
		return &chaincfg.TestNet3Params
	}

	return &chaincfg.MainNetParams
}

var conf *configuration

type configuration struct {
	TestNet  bool `mapstructure:"testnet"`
	Electron struct {
		Host string `mapstructure:"host"`
		Port string `mapstructure:"port"`
	}
}

func GetConf() *configuration {
	if conf != nil {
		return conf
	}

	config := &configuration{}
	viper.SetEnvPrefix("whc")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")
	viper.SetDefault("conf", "./conf.yml")

	// get config file path from environment
	confFile := viper.GetString("conf")

	var realPath string

	// conf.go unit testing
	if viper.GetString("test_conf") != "" {
		realPath = viper.GetString("test_conf")
	} else {
		path, err := filepath.Abs("./")
		if err != nil {
			panic(err)
		}

		correctPath := path + ProjectLastDir
		realPath = filepath.Join(correctPath, confFile)
	}

	// parse config
	file, err := os.Open(realPath)
	if err != nil {
		panic("Open config file error: " + err.Error())
	}
	defer file.Close()

	err = viper.ReadConfig(file)
	if err != nil {
		panic("Read config file error: " + err.Error())
	}

	err = viper.Unmarshal(config)
	if err != nil {
		panic("Parse config file error: " + err.Error())
	}

	conf = config
	return config
}
