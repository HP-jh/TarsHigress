# TarsHigress

MCP integration for Higress (Alibaba Cloud AI gateway).

## Tools

- `higress_routes` — List AI gateway routes
- `higress_health` — Check Higress gateway health

## Quick Start

```bash
go build -o tarshigress
TARS_MCP_ENDPOINT=https://api.example.com TARS_MCP_API_KEY=xxx ./tarshigress
```

## Environment Variables

- `TARS_MCP_ENDPOINT` — upstream API endpoint
- `TARS_MCP_API_KEY` — upstream API key
- `TARS_MCP_MODE` — `stdio` (default) or `http`

## License

MIT
