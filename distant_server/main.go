package main

import (
	"crypto/x509"
	"flag"
	"os"

	"github.com/adrg/xdg"
)

var port = flag.Int("port", 8080, "The port to listen on")
var configPath = xdg.ConfigHome + "distant_server"

func main() {
	flag.Parse()

	caCertFile, err := os.ReadFile(configPath + "/ca.cert")
	if err != nil {
		panic("todo!")
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)
}
