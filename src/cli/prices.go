package main

import (
	"fmt"
	"github.com/Farengier/gotools/logging"
	"github.com/Farengier/gotools/routine"
	"github.com/Farengier/prices/src/pkg/server"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
)

type opts struct {
	Config string `short:"c" long:"config" description:"General config file" default:"config.yml"`
}

type cfg struct {
	Log struct {
		Level string `yaml:"level"`
		File  string `yaml:"file"`
	} `yaml:"log"`
	ServerConfig `yaml:"server"`
}

type ServerConfig struct {
	Listen string `yaml:"listen"`
}

func (sc ServerConfig) Addr() string {
	return sc.Listen
}

func main() {
	defer errHandler()
	opts := initOpts()
	conf := readConf(opts.Config)
	log := configLogger(conf)
	srv := server.Server{Log: log, Cfg: conf.ServerConfig}
	srv.Start()

	routine.WaitForExit()
}

func initOpts() opts {
	o := opts{}
	_, err := flags.NewParser(&o, flags.HelpFlag|flags.PassDoubleDash).Parse()
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

func configLogger(c cfg) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	lvl, err := zerolog.ParseLevel(c.Log.Level)
	if err != nil {
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	var writer io.Writer
	writer = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05-07:00"}
	if c.Log.File != "" {
		if f, err := os.OpenFile(c.Log.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err == nil {
			writer = zerolog.SyncWriter(f)
		}
	}

	l := zerolog.New(writer).Level(lvl).With().Timestamp().Logger()
	logging.SetLoggerZeroLog(l)
	return l
}

func errHandler() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}
