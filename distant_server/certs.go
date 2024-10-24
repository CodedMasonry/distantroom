package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func makeCA(subject *pkix.Name) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	caCert := &x509.Certificate{
		Subject:               *subject,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 365),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	// Generate Cert
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Error("Failed to generate encryption key", "error", err)
		return nil, nil, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, caCert, caCert, caKey.PublicKey, caKey)
	if err != nil {
		log.Error("Failed to generate certificate", "error", err)
		return nil, nil, err
	}

	// Logging handled in function
	if err := saveCert(caBytes, caKey, "ca"); err != nil {
		return nil, nil, err
	}

	return caCert, caKey, nil
}

func makeCertWithDNS(caCert *x509.Certificate, caKey *ecdsa.PrivateKey, subject *pkix.Name, dnsNames []string, fileName string) error {
	cert := &x509.Certificate{
		Subject:     *subject,
		DNSNames:    dnsNames,
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(time.Hour * 24 * 365),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// Generate Cert
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Error("Failed to generate encryption key", "error", err)
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, cert.PublicKey, caKey)
	if err != nil {
		log.Error("Failed to generate certificate", "error", err)
		return err
	}

	// Logging handled in function
	if err := saveCert(certBytes, caKey, fileName); err != nil {
		return err
	}

	return nil
}

func makeCertWithIPAddress(caCert *x509.Certificate, caKey *ecdsa.PrivateKey, subject *pkix.Name, ipAddresses []net.IP, fileName string) error {
	cert := &x509.Certificate{
		Subject:     *subject,
		IPAddresses: ipAddresses,
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(time.Hour * 24 * 365),
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// Generate Cert
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Error("Failed to generate encryption key", "error", err)
		return err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, cert.PublicKey, caKey)
	if err != nil {
		log.Error("Failed to generate certificate", "error", err)
		return err
	}

	// Logging handled in function
	if err := saveCert(certBytes, caKey, fileName); err != nil {
		return err
	}

	return nil
}

func saveCert(certBytes []byte, key *ecdsa.PrivateKey, name string) error {
	// Encode Public
	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err := os.WriteFile(configPath+"/"+name+".cert", caPEM.Bytes(), 0640); err != nil {
		log.Error("Failed to save certificae", "error", err)
		return err
	}

	// Encode Private
	caKeyBytes, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		log.Error("Failed to parse private key", "error", err)
		return err
	}

	PrivateKeyPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "ECC PRIVATE KEY",
		Bytes: caKeyBytes,
	})
	if err := os.WriteFile(configPath+"/"+name+".key", PrivateKeyPEM.Bytes(), 0640); err != nil {
		log.Error("Failed to save certificae", "error", err)
		return err
	}

	return nil
}
