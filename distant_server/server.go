package main

import (
	"crypto/x509/pkix"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
)

func statusHandler(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func runServer(host string, port uint16) {
	// Handler Config
	app := fiber.New()
	app.Get("/status", statusHandler)

	// Create Server Certificate
	if _, err := os.Stat(CONFIG_PATH + "/server.cert"); err != nil {
		subject := &pkix.Name{
			CommonName: args.Host,
		}
		if err := makeServerCert(GLOBAL_STATE.caCert, GLOBAL_STATE.caKey, subject, []string{host}); err != nil {
			log.Fatal("Failed to create server certificate", "error", err)
		}
		log.Info("Created Server Certificate", "path", CONFIG_PATH+"/server.cert")
	}

	// Run listener
	portStr := fmt.Sprintf(":%d", port)
	log.Fatal(app.ListenMutualTLS(portStr, CONFIG_PATH+"/server.cert", CONFIG_PATH+"/server.key", CONFIG_PATH+"/ca.cert"))
}

/*
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
*/
