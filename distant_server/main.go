package main

import (
	"crypto/x509"
	"flag"
	"os"
	"time"

	"github.com/adrg/xdg"
	"github.com/charmbracelet/log"
)

var port = flag.Int("port", 8080, "The port to listen on")
var CONFIG_PATH = xdg.ConfigHome + "/distantroom"

func main() {
	flag.Parse()

	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	logger.SetTimeFormat(time.Kitchen)
	logger.SetReportTimestamp(true)
	logger.SetReportCaller(true)
	log.SetDefault(logger)

	initState()

	caCertFile, err := os.ReadFile(CONFIG_PATH + "/ca.cert")
	if err != nil {
		panic("todo!")
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)
}
