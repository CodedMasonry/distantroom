package main

import (
	"os"
	"time"

	"github.com/adrg/xdg"
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/log"
)

var CONFIG_PATH = xdg.ConfigHome + "/distantroom"
var args struct {
	NewOperator *NewOperatorCmd `arg:"subcommand:new-operator"`
	Port        uint16          `arg:"-p, --port" help:"Port to listen on" default:"3000"`
	Host        string          `arg:"-h, --host" help:"The host that the client should connect to" default:"0.0.0.0"`
	Debug       bool            `arg:"-d, --debug" help:"sets log level to debug"`
}

func main() {
	arg.MustParse(&args)
	// Inits Logging
	initMain()
	// Inits Config
	initState()

	// Handle sub commands
	switch {
	case args.NewOperator != nil:
		log.Info("Generating New Operator")
		NewOperator(args.NewOperator)
	default:
		runServer(args.Host, args.Port)
	}
}

func initMain() {
	logger := log.New(os.Stderr)
	logger.SetTimeFormat(time.Stamp)
	logger.SetReportTimestamp(true)
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
