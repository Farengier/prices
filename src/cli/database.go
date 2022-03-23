package main

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type opts struct {
	Config string `short:"c" long:"config" description:"General config file" default:"config.yml"`
}

type cfg struct {
	Log struct {
		Level string `yaml:"level"`
		File string `yaml:"file"`
	} `yaml:"log"`
}

func main() {
	defer errHandler()
	opts := initOpts()
	conf := readConf(opts.Config)
	configLogger(conf)

	fmt.Println("Config is", conf)
}

func initOpts() opts {
	o := opts{}
	_, err := flags.NewParser(&o, flags.HelpFlag | flags.PassDoubleDash).Parse()
	if err != nil {
		panic(err)
	}
	return o
}

func readConf(fn string) cfg {
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		panic("Failed reading config file: " + err.Error())
	}

	c := cfg{}
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		panic("Failed reading config file: " + err.Error())
	}
	return c
}

func configLogger(c cfg) {
	lvl, err := log.ParseLevel(c.Log.Level)
	if err != nil {
		return
	}
	log.SetLevel(lvl)

	if c.Log.File != "" {
		f, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(f)
		}
	}
}

func errHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}