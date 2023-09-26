package main

import (
	"fmt"
	"log"
	"os"

	"git.s8k.top/SeraphJACK/beanbot/config"
	"git.s8k.top/SeraphJACK/beanbot/repo"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

var exampleConfig = pflag.BoolP("example", "e", false, "Print example config and exit")
var confPath = pflag.StringP("conf", "c", "config.yml", "Path to the configuration file")

func main() {
	pflag.Parse()

	if *exampleConfig {
		b, _ := yaml.Marshal(config.Cfg)
		fmt.Print(string(b))
		os.Exit(0)
	}

	if err := config.Load(*confPath); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	if err := repo.Init(); err != nil {
		log.Fatalf("Failed to init beancount repository: %v", err)
		return
	}
}
