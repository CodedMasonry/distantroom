package main

import (
	"crypto/x509"
	"os"
	"time"

	"github.com/adrg/xdg"
	"github.com/alexflint/go-arg"
	"github.com/charmbracelet/log"
)

type NewOperatorCmd struct {
	output string `arg:"-o, --output" help:"Directory to output file" default:"."`
}

var CONFIG_PATH = xdg.ConfigHome + "/distantroom"
var args struct {
	NewOperator *NewOperatorCmd `arg:"subcommand:new-operator"`
	Debug       bool            `arg:"-d, --debug" help:"sets log level to debug"`
}

func main() {
	arg.MustParse(&args)

	logger := log.New(os.Stderr)
	logger.SetTimeFormat(time.Stamp)
	logger.SetReportTimestamp(true)
	if args.Debug {
		logger.SetReportCaller(true)
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.DebugLevel)
	}
	log.SetDefault(logger)

	if err := os.MkdirAll(CONFIG_PATH, 0740); err != nil {
		log.Fatal("Failed to create Configuration Directory", "error", err)
	}

	initState()

	// Handle sub commands
	switch {
	case args.NewOperator != nil:
		log.Info("Generating New Operator")
	default:
		log.Info("Starting Server")
		runServer()
	}
}

func runServer() {
	caCertFile, err := os.ReadFile(CONFIG_PATH + "/ca.cert")
	if err != nil {
		// Retry with new CA
		NewCA()
		caCertFile, err = os.ReadFile(CONFIG_PATH + "/ca.cert")
		if err != nil {
			log.Fatal("Failed to create CA Cert", "error", err)
		}
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)
}
