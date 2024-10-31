package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

// The TOML template for `operator.toml` files
// Format complying with snake_case
type OperatorTemplate struct {
	Host              string `toml:"host" comment:"Host that the server is listening on"`
	Port              uint16 `toml:"port" comment:"Port the server is listening on"`
	Certificate       string `toml:"certificate" comment:"The Client's Certificate"`
	PrivateKey        string `toml:"private_key" comment:"The Client's Private Key"`
	ServerCertificate string `toml:"server_certificate" comment:"The Servers's Certificate"`
}

func parseProfile(file string) (*OperatorTemplate, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var profile OperatorTemplate
	if err := toml.Unmarshal(bytes, profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
