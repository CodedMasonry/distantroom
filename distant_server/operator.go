package main

import (
	"crypto/x509/pkix"
	"fmt"
	"math/rand/v2"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/pelletier/go-toml/v2"
)

// CLI format
type NewOperatorCmd struct {
	output string `arg:"-o, --output" help:"Directory to output file; defaults to current directory"`
	host   string `arg:"-h, --host" help:"The host that the client should connect to; defaults to 0.0.0.0"`
	port   uint16 `arg:"-p, --port" help:"The port the client should connect to; defaults to 8080"`
	name   string `arg:"-n, --name" help:"Name of the file; defaults to date + random number"`
}

// The TOML template for `operator.toml` files
// Format complying with snake_case
type OperatorTemplate struct {
	Host              string `toml:"host" comment:"Host that the server is listening on"`
	Port              uint16 `toml:"port" comment:"Port the server is listening on"`
	Certificate       string `toml:"certificate" comment:"The Client's Certificate"`
	PrivateKey        string `toml:"private_key" comment:"The Client's Private Key"`
	ServerCertificate string `toml:"server_certificate" comment:"The Servers's Certificate"`
}

func NewOperator(cmd *NewOperatorCmd) {
	subject := &pkix.Name{
		CommonName:   "R11",
		Organization: []string{"Let's Encrypt"},
		Country:      []string{"US"},
	}

	// Handling default values
	if cmd.output == "" {
		cmd.output = "."
	}
	if cmd.host == "" {
		cmd.host = "0.0.0.0"
	}
	if cmd.port == 0 {
		cmd.port = 8080
	}
	if cmd.name == "" {
		cmd.name = time.Now().Format(time.DateOnly) + "-" + fmt.Sprint(rand.UintN(4096))
	}

	// Generate a certificate for client
	if err := makeOperatorCert(GLOBAL_STATE.caCert, GLOBAL_STATE.caKey, subject, []string{cmd.host}, cmd.name); err != nil {
		log.Error("Failed to create Operator Certificate", "error", err)
	}

	// Load all the PEM files
	certPEM, err := os.ReadFile(CONFIG_PATH + "/operators/" + cmd.name + ".cert")
	if err != nil {
		log.Fatal("Failed to load file", "error", err)
	}
	keyPEM, err := os.ReadFile(CONFIG_PATH + "/operators/" + cmd.name + ".key")
	if err != nil {
		log.Fatal("Failed to load file", "error", err)
	}
	caPEM, err := os.ReadFile(CONFIG_PATH + "/ca.cert")
	if err != nil {
		log.Fatal("Failed to load file", "error", err)
	}

	// Prepare the template
	template := OperatorTemplate{
		Host:              cmd.host,
		Port:              cmd.port,
		Certificate:       string(certPEM),
		PrivateKey:        string(keyPEM),
		ServerCertificate: string(caPEM),
	}

	// Convert to TOML
	bytes, err := toml.Marshal(template)
	if err != nil {
		log.Fatal("Failed to encode client config", "error", err)
	}

	// Write to `output` directory
	path := cmd.output + "/" + cmd.name + ".toml"
	if err := os.WriteFile(path, bytes, 0640); err != nil {
		log.Fatal("Failed to write config to file", "error", err)
	}
	log.Info("Configuration saved", "path", path)
}
