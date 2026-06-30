# kalshi-mcp

[![CI](https://github.com/rangertaha/kalshi-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/rangertaha/kalshi-mcp/actions/workflows/ci.yml)
[![Status: under construction](https://img.shields.io/badge/status-under%20construction-orange)](#-under-construction)

<div align="center">

## 🚧 &nbsp; UNDER CONSTRUCTION &nbsp; 🚧

**This server is an early scaffold — a work in progress.**

It runs over stdio with **one read-only toolset** wired end-to-end.<br>
More toolsets are on the way (see the **TODO** list below).<br>
APIs, configuration, and tool names may still change.

</div>

---

A [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server, written
in Go, exposing the **Kalshi** prediction-markets API as tools an LLM client
(Claude Desktop/Code, Cursor, and others) can call.

## Features

- **Typed tools with schemas**: every tool has an auto-generated JSON Schema for
  its input and output, inferred from Go structs.
- **Read-only switch**: `KALSHI_READONLY=true` hides every mutating tool.
- **Toolset filtering**: enable only the areas you need with `KALSHI_TOOLSETS`.
- **Public by default**: market-data tools need no credentials. Supplying an API
  key signs requests with Kalshi's RSA-PSS scheme for authenticated endpoints.

## Install

```sh
go install github.com/rangertaha/kalshi-mcp/cmd/kalshi@latest
```

Or build from source:

```sh
git clone https://github.com/rangertaha/kalshi-mcp
cd kalshi-mcp
make build        # produces ./bin/kalshi
```

## CLI

```sh
kalshi mcp      # run the MCP server over stdio (default when no subcommand)
kalshi test     # verify connectivity
```

## Configuration

| Variable             | Required | Description                                                   |
| -------------------- | :------: | ------------------------------------------------------------- |
| `KALSHI_API_KEY_ID`  |    no    | Access key ID (enables authenticated tools).                  |
| `KALSHI_PRIVATE_KEY` |    no    | PEM-encoded RSA private key for request signing.              |
| `KALSHI_BASE_URL`    |    no    | API base URL (default `https://api.elections.kalshi.com/trade-api/v2`). |
| `KALSHI_TOOLSETS`    |    no    | Comma-separated toolset names to enable, or `all`.            |
| `KALSHI_READONLY`    |    no    | `true` to expose only read-only tools.                        |

## Toolsets

| Toolset   | Covers                                                          |
| --------- | -------------------------------------------------------------- |
| `markets` | list markets (`markets_list`) and get one (`markets_get`) — public market data |

### TODO toolsets

- `events` — list/get events and their nested markets.
- `portfolio` — balance and positions (needs auth).
- `orders` — list/create/cancel orders (write; needs auth).

## License

MIT — see [LICENSE](LICENSE).
