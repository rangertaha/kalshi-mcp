// SPDX-License-Identifier: MIT

// Command kalshi runs the Kalshi Model Context Protocol server (`kalshi mcp`)
// and checks connectivity (`kalshi test`).
//
// Configuration is read from the environment (see package config). The `mcp`
// command communicates over stdio, the transport expected by MCP clients such
// as Claude Desktop/Code and Cursor.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v3"

	"github.com/rangertaha/kalshi-mcp/internal"
	"github.com/rangertaha/kalshi-mcp/internal/app"
	"github.com/rangertaha/kalshi-mcp/internal/config"
	"github.com/rangertaha/kalshi-mcp/internal/kalshi"
)

func main() {
	cmd := &cli.Command{
		Name:    "kalshi",
		Usage:   "Kalshi prediction markets as an MCP server",
		Version: internal.Version(),
		// A bare `kalshi` (no subcommand) runs the MCP server.
		Action: runMCP,
		Commands: []*cli.Command{
			mcpCommand(),
			testCommand(),
		},
		// Print errors ourselves so the MCP stdio stream is never touched.
		ExitErrHandler: func(context.Context, *cli.Command, error) {},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "kalshi: %v\n", err)
		os.Exit(1)
	}
}

// mcpCommand runs the MCP server over stdio.
func mcpCommand() *cli.Command {
	return &cli.Command{
		Name:   "mcp",
		Usage:  "Run the MCP server over stdio (for Claude Desktop/Code, Cursor, ...)",
		Action: runMCP,
	}
}

// runMCP assembles and serves the MCP server over stdio.
func runMCP(ctx context.Context, _ *cli.Command) error {
	if err := config.LoadEnvFile(config.EnvFile); err != nil {
		log.Printf("kalshi: reading %s: %v", config.EnvFile, err)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration error:\n%w", err)
	}

	ver := internal.Version()
	srv, cleanup, err := app.Assemble(cfg, ver)
	if err != nil {
		return err
	}
	defer cleanup()

	log.Printf("kalshi-mcp %s starting: %d tools, %d prompts across toolsets %v (read-only=%v)",
		ver, srv.ToolCount(), srv.PromptCount(), srv.Toolsets(), cfg.ReadOnly)

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	return srv.Run(ctx, &mcp.StdioTransport{})
}

// testCommand verifies connectivity against the Kalshi API.
func testCommand() *cli.Command {
	return &cli.Command{
		Name:  "test",
		Usage: "Test connectivity against the Kalshi API",
		Action: func(ctx context.Context, _ *cli.Command) error {
			if err := config.LoadEnvFile(config.EnvFile); err != nil {
				log.Printf("kalshi: reading %s: %v", config.EnvFile, err)
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("configuration error:\n%w", err)
			}

			clients, err := kalshi.NewClients(cfg.BaseURL, cfg.KeyID, cfg.PrivateKey)
			if err != nil {
				return err
			}

			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			n, err := kalshi.Check(ctx, clients)
			if err != nil {
				return fmt.Errorf("connecting to %s: %w", cfg.BaseURL, err)
			}

			fmt.Printf("OK  connected to %s (%d market(s) returned)\n", cfg.BaseURL, n)
			fmt.Printf("    authenticated=%v read-only=%v\n", cfg.KeyID != "", cfg.ReadOnly)
			return nil
		},
	}
}
