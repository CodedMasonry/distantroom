package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var MaxSerialNumber = big.NewInt(0).SetBytes(bytes.Repeat([]byte{255}, 20))

// Generates a template Certificate & a Keypair.
//
// If hosts are provided, they are added to the certificate
// else if nil, none are added.
func templateCertificate(subject *pkix.Name, hosts *[]string) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	serial, err := generateSerialNumber()
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      *subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Hour * 24 * 365),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	if hosts != nil {
		for _, h := range *hosts {
			if ip := net.ParseIP(h); ip != nil {
				template.IPAddresses = append(template.IPAddresses, ip)
			} else {
				template.DNSNames = append(template.DNSNames, h)
			}
		}
	}
	log.Debug("Generating Certificate")

	// Generate Cert
	certKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Error("Failed to generate encryption key", "error", err)
		return nil, nil, err
	}

	return template, certKey, nil
}

func makeCA(subject *pkix.Name) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	serial, err := generateSerialNumber()
	if err != nil {
		return nil, nil, err
	}

	template := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               *subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
	log.Debug("Generating Certificate")

	// Generate Cert
	certKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Error("Failed to generate encryption key", "error", err)
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &certKey.PublicKey, certKey)
	if err != nil {
		log.Error("Failed to generate certificate", "error", err)
		return nil, nil, err
	}

	// Logging handled in function
	if err := saveCert(certBytes, certKey, "ca"); err != nil {
		return nil, nil, err
	}

	return template, certKey, nil
}

func makeServerCert(caCert *x509.Certificate, caKey *ecdsa.PrivateKey, subject *pkix.Name, hosts []string) error {
	template, certKey, err := templateCertificate(subject, &hosts)
	if err != nil {
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, &certKey.PublicKey, caKey)
	if err != nil {
		log.Error("Failed to generate certificate", "error", err)
		return err
	}

	// Logging handled in function
	if err := saveCert(certBytes, certKey, "server"); err != nil {
		return err
	}

	return err
}

func makeOperatorCert(caCert *x509.Certificate, caKey *ecdsa.PrivateKey, subject *pkix.Name, hosts []string) (*bytes.Buffer, *bytes.Buffer, error) {
	template, certKey, err := templateCertificate(subject, &hosts)
	if err != nil {
		return nil, nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, template, caCert, &certKey.PublicKey, caKey)
	if err != nil {
		log.Error("Failed to generate certificate", "error", err)
		return nil, nil, err
	}

	// Encode Public
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	// Encode Private
	keyBytes, err := x509.MarshalPKCS8PrivateKey(certKey)
	if err != nil {
		log.Error("Failed to parse private key", "error", err)
		return nil, nil, err
	}

	privateKeyPEM := new(bytes.Buffer)
	pem.Encode(privateKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})

	return certPEM, privateKeyPEM, err
}

func saveCert(certBytes []byte, key *ecdsa.PrivateKey, name string) error {
	// Encode Public
	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err := os.WriteFile(CONFIG_PATH+"/"+name+".cert", certPEM.Bytes(), 0640); err != nil {
		log.Error("Failed to save certificae", "error", err)
		return err
	}

	// Encode Private
	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		log.Error("Failed to parse private key", "error", err)
		return err
	}

	privateKeyPEM := new(bytes.Buffer)
	pem.Encode(privateKeyPEM, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	})
	if err := os.WriteFile(CONFIG_PATH+"/"+name+".key", privateKeyPEM.Bytes(), 0640); err != nil {
		log.Error("Failed to save certificae", "error", err)
		return err
	}
	log.Debug("Successfully saved certificate", "path", CONFIG_PATH+"/"+name+".cert")

	return nil
}

func generateSerialNumber() (*big.Int, error) {
	return rand.Int(rand.Reader, MaxSerialNumber)
}

func NewCA() (*x509.Certificate, *ecdsa.PrivateKey, error) {
	log.Debug("Generating Certificate Authority")
	subject := pkix.Name{
		CommonName:   "WE1",
		Organization: []string{"Google Trust Services"},
		Country:      []string{"US"},
	}

	caCert, caKey, err := makeCA(&subject)
	if err != nil {
		return nil, nil, err
	}
	log.Info("Created Root Certificate", "path", CONFIG_PATH+"/ca.cert")
	return caCert, caKey, nil
}

func LoadCA() (*x509.Certificate, *ecdsa.PrivateKey, error) {
	log.Debug("Loading Certificate Authority", "path", CONFIG_PATH+"/ca.cert")
	certFile, err := os.ReadFile(CONFIG_PATH + "/ca.cert")
	if err != nil {
		return NewCA()
	}
	keyFile, err := os.ReadFile(CONFIG_PATH + "/ca.key")
	if err != nil {
		return NewCA()
	}

	certBlock, _ := pem.Decode(certFile)
	keyBlock, _ := pem.Decode(keyFile)
	if certBlock == nil || keyBlock == nil {
		log.Fatal("Failed to parse CA files; Not valid PEM files")
	}

	caCert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	caKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	switch caKey.(type) {
	case *ecdsa.PrivateKey:
		return caCert, caKey.(*ecdsa.PrivateKey), nil
	default:
		return nil, nil, fmt.Errorf("Invalid Key Type; only ECDSA keys are supporteds")
	}
}
