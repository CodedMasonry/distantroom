package main

import (
	"crypto/x509/pkix"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func statusHandler(c *fiber.Ctx) error {
	return c.SendStatus(200)
}

func runServer(host string, port uint16) {
	// Handler Config
	app := fiber.New()
	app.Get("/status", statusHandler)

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Print(c.Locals("allowed"))  // true
		log.Print(c.Params("id"))       // 123
		log.Print(c.Query("v"))         // 1.0
		log.Print(c.Cookies("session")) // ""

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			mt  int
			msg []byte
			err error
		)
		for {
			if mt, msg, err = c.ReadMessage(); err != nil {
				log.Printf("read: %v", err)
				break
			}
			log.Printf("recv: %s", msg)

			if err = c.WriteMessage(mt, msg); err != nil {
				log.Printf("write: %v", err)
				break
			}
		}

	}))

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
