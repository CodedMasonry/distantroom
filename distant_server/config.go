package main

import (
	"crypto/ecdsa"
	"crypto/x509"

	"github.com/charmbracelet/log"
)

type State struct {
	caCert *x509.Certificate
	caKey  *ecdsa.PrivateKey
}

type Config struct {
}

var GLOBAL_STATE State

func initState() {
	caCert, caKey, err := LoadCA()
	if err != nil {
		log.Fatal("Failed to create CA", "error", err)
	}
	log.Debug("Certificate Authority Loaded")
	log.Info("Config loaded")

	GLOBAL_STATE = State{
		caCert: caCert,
		caKey:  caKey,
	}
}

/*
func LoadConfig(path string) (*Config, error) {
	if _, err := os.Stat(path); err != nil {
		return createConfig(path)
	}

	var cfg Config
	_, err := toml.DecodeFile(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func createConfig(path string) (*Config, error) {
	log.Info("Generating Config")
	cfg := &Config{}

	bytes, err := toml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	if err = os.WriteFile(path, bytes, 0640); err != nil {
		return nil, err
	}

	return cfg, nil
}
*/
