package main

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func printConnState(state *tls.ConnectionState) {
	log.Print(">>>>>>>>>>>>>>>> State <<<<<<<<<<<<<<<<")

	log.Printf("Version: %x", state.Version)
	log.Printf("HandshakeComplete: %t", state.HandshakeComplete)
	log.Printf("DidResume: %t", state.DidResume)
	log.Printf("CipherSuite: %x", state.CipherSuite)
	log.Printf("NegotiatedProtocol: %s", state.NegotiatedProtocol)

	log.Print("Certificate chain:")
	for i, cert := range state.PeerCertificates {
		subject := cert.Subject
		issuer := cert.Issuer
		log.Printf(" %d s:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", i, subject.Country, subject.Province, subject.Locality, subject.Organization, subject.OrganizationalUnit, subject.CommonName)
		log.Printf("   i:/C=%v/ST=%v/L=%v/O=%v/OU=%v/CN=%s", issuer.Country, issuer.Province, issuer.Locality, issuer.Organization, issuer.OrganizationalUnit, issuer.CommonName)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil {
		printConnState(r.TLS)
	}
	io.WriteString(w, "Status 200")
}

func runServer(port uint16) {
	// Handler Config
	handler := http.NewServeMux()
	handler.HandleFunc("/status", statusHandler)

	// Root CA Config
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
	// TLS Config
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
		MinVersion: tls.VersionTLS13,
	}
	// Logging Config
	logger := log.NewWithOptions(os.Stderr, log.Options{Prefix: "HTTPS", TimeFormat: time.Stamp})
	logger.SetReportTimestamp(true)
	stdlog := logger.StandardLog(log.StandardLogOptions{
		ForceLevel: log.ErrorLevel,
	})
	// Server Config
	server := http.Server{
		Addr:      fmt.Sprintf(":%d", port),
		Handler:   handler,
		TLSConfig: tlsConfig,
		ErrorLog:  stdlog,
	}
	// Server Certificate
	if _, err := os.Stat(CONFIG_PATH + "/server.cert"); err != nil {
		subject := &pkix.Name{
			CommonName: args.Host,
		}
		if err := makeServerCert(subject); err != nil {
			log.Fatal("Failed to create server certificate", "error", err)
		}
		log.Info("Created Server Certificate", "path", CONFIG_PATH+"/server.cert")
	}
	// Start the Server
	logger.Infof("Listening on :%d\n", port)
	if err := server.ListenAndServeTLS(CONFIG_PATH+"/server.cert", CONFIG_PATH+"/server.key"); err != nil {
		logger.Fatal("Server Stopped", "error", err)
	}
}
