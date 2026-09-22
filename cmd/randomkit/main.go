package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/relentlessworks/randomkit/internal/api"
	"github.com/relentlessworks/randomkit/internal/auth"
	"github.com/relentlessworks/randomkit/internal/config"
	"github.com/relentlessworks/randomkit/internal/mcp"
)

func main() {
	cfg := config.Load()

	// Initialize auth
	a := auth.New(cfg.Secret)

	// Initialize API handler
	h := api.New(a)

	// Initialize MCP handler
	mcpHandler := mcp.New()

	// Register MCP endpoint
	mux := h.Routes()

	// Wrap the mux to add MCP route
	wrappedMux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/mcp" {
			mcpHandler.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	// Start server
	fmt.Fprintf(os.Stderr, "randomkit starting on %s\n", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, wrappedMux))
}
