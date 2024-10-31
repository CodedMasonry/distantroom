package main

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/log"
)

var CONFIG_PATH = xdg.ConfigHome + "/distantroom"
var args struct {
	Debug bool `arg:"-d, --debug" help:"sets log level to debug"`
}

var Error = "[-]"
var Warn = "[!]"
var Info = "[+]"
var Debug = "[*]"
var Trace = "[$]"

func main() {
	arg.MustParse(&args)
	// Inits Logging
	initMain()

	_, err := parseProfile(CONFIG_PATH + "/test.toml")
	if err != nil {
		panic(err)
	}
}

func initMain() {
	logger := log.New(os.Stderr)
	if args.Debug {
		logger.SetReportCaller(true)
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
	log.SetDefault(logger)

	if err := os.MkdirAll(CONFIG_PATH, 0740); err != nil {
		log.Fatal("Failed to create Configuration Directory", "error", err)
	}
}
