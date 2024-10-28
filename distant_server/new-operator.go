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
	Output string `arg:"--output" help:"Directory to output file" default:"."`
	Name   string `arg:"-n, --name" help:"Name of the file"`
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
		CommonName:   "Cloudflare",
		Country:      []string{"US"},
	}

	// Handling default values
	if cmd.Name == "" {
		cmd.Name = time.Now().Format(time.DateOnly) + "-" + fmt.Sprint(rand.UintN(4096))
	}

	// Generate a certificate for client
	certPEM, keyPEM, err := makeOperatorCert(GLOBAL_STATE.caCert, GLOBAL_STATE.caKey, subject, []string{args.Host})
	if err != nil {
		log.Error("Failed to create Operator Certificate", "error", err)
	}

	// Load all the PEM files
	caPEM, err := os.ReadFile(CONFIG_PATH + "/ca.cert")
	if err != nil {
		log.Fatal("Failed to load file", "error", err)
	}

	// Prepare the template
	template := OperatorTemplate{
		Host:              args.Host,
		Port:              args.Port,
		Certificate:       string(certPEM.Bytes()),
		PrivateKey:        string(keyPEM.Bytes()),
		ServerCertificate: string(caPEM),
	}

	// Convert to TOML
	bytes, err := toml.Marshal(template)
	if err != nil {
		log.Fatal("Failed to encode client config", "error", err)
	}

	// Write to `output` directory
	path := cmd.Output + "/" + cmd.Name + ".toml"
	if err := os.WriteFile(path, bytes, 0640); err != nil {
		log.Fatal("Failed to write config to file", "error", err)
	}
	log.Info("Operator Saved", "path", path)
}
