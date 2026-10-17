# TarsHigress

TSG (TarsSecureGuard) MCP adapter for **Higress** — Cloud-native API gateway with WAF.

This is one of the **15 MCP integration projects** under github.com/HP-jh/, designed to be:

- **One-click downloadable** from TSG's MCP Marketplace
- **One-click enableable** in the TSG gateway
- **Independently versioned** (this repo) with its own release cycle
- **Tightly integrated** with TSG via the standard MCP protocol

## Quick Start

### Standalone

```bash
# Download from GitHub Releases
curl -L -o tars-higress https://github.com/HP-jh/TarsHigress/releases/latest/download/tars-higress-$(uname -s | tr A-Z a-z)-amd64
chmod +x tars-higress
./tars-higress --upstream-endpoint <your-Higress-endpoint>
```

### Via TSG MCP Marketplace

1. Open TSG admin → **MCP Marketplace**
2. Find **TarsHigress** in the list
3. Click **Download** → wait for completion
4. Click **Enable** → enter your Higress endpoint
5. The MCP server is now available to all TSG clients

## MCP Tools Exposed

| Tool | Description |
|------|-------------|
| `higress_list_models` | List available Higress models |
| `higress_invoke` | Invoke Higress with prompt + params |
| `higress_health` | Health check of upstream connection |
| `higress_status` | Current TSG integration status |

## Configuration

Configurable via env vars or `~/.tars/higress/config.json`:

- `TARS_HIGRESS_ENDPOINT` — upstream endpoint URL
- `TARS_HIGRESS_API_KEY` — upstream API key (optional, for hosted services)
- `TARS_HIGRESS_TIMEOUT` — request timeout (default 30s)
- `TARS_HIGRESS_LOG_LEVEL` — log verbosity (default info)

## Version

- **v0.1.0** — Initial MVP (basic MCP adapter)
- See [Releases](https://github.com/HP-jh/TarsHigress/releases) for full changelog

## Cross-Repo Coordination

This repo is part of the **15 MCP integration projects** in github.com/HP-jh/. See [HP-jh/TarsSecureGuard](https://github.com/HP-jh/TarsSecureGuard) for the master gateway that orchestrates all 15.

## License

MIT — see [LICENSE](LICENSE).
