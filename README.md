# randomkit

> Agentic-first random data generation service. Generate random strings, numbers, passwords, tokens, colors, dates, choices, coin flips, dice rolls, and UUIDs. Plain text API, agent-driven, single Go binary.

## Quick Start

```bash
# Build
make build

# Run
./randomkit
# → randomkit starting on :8787

# Use
curl http://localhost:8787/random/string?length=32&charset=hex
curl http://localhost:8787/password?length=24
curl http://localhost:8787/dice?count=3&sides=20
curl http://localhost:8787/lorem/paragraphs?count=2
curl http://localhost:8787/color?format=hsl
curl http://localhost:8787/uuid
```

## Principles

- **The agent IS the interface** — No UI, no SDK. The API is the product.
- **Plain text by default** — One labeled, grepable line per record. JSON on demand.
- **Instructive errors** — Every 4xx includes a hint for self-correction.
- **Self-documenting** — `GET /help` returns a one-page operating manual.
- **Simple auth** — OTP via email → long-lived bearer token.
- **Single static binary** — Go, CGO_ENABLED=0, zero external dependencies.
- **Zero config defaults** — Runs out of the box.
- **MCP connector** — Speaks Model Context Protocol at `/mcp`.

## API Reference

### Random String
```
GET /random/string?length=32&charset=alnum
→ value=aB3xK9mP2nQr... length=32 charset=alnum
```
Charsets: `lower`, `upper`, `digits`, `symbols`, `hex`, `alpha`, `alnum`, or custom.

### Random Integer
```
GET /random/int?min=1&max=100
→ value=42 min=1 max=100
```

### Random Float
```
GET /random/float?min=0&max=1
→ value=0.582934 min=0 max=1
```

### Password
```
GET /password?length=20&lower=true&upper=true&digits=true&symbols=true
→ password=K7$mP2nQr@9x... length=20
```

### Random Hex Token
```
GET /random/hex?bytes=32
→ token=a1b2c3d4e5f6... bytes=32
```

### Random Bytes
```
GET /random/bytes?bytes=16
→ bytes=a1b2c3d4e5f6... length=16
```

### Lorem Ipsum
```
GET /lorem/words?count=50
GET /lorem/sentences?count=3
GET /lorem/paragraphs?count=2
```

### Random Choice
```
GET /random/choice?items=red,green,blue
→ choice=green count=3
```

### Random Sample (without replacement)
```
GET /random/sample?items=a,b,c,d&count=2
→ item=c
  item=a
```

### Shuffle
```
GET /random/shuffle?items=a,b,c,d
→ item=c
  item=a
  item=d
  item=b
```

### Random Boolean
```
GET /random/bool
→ value=true
```

### Random Color
```
GET /color?format=hex    → color=#a1b2c3 format=hex
GET /color?format=rgb    → color=rgb(163, 178, 195) format=rgb
GET /color?format=hsl    → color=hsl(210, 75%, 50%) format=hsl
```

### Coin Flip
```
GET /coin
→ result=heads
```

### Dice Roll
```
GET /dice?count=2&sides=6
→ total=8 rolls=3+5 count=2 sides=6
```

### Random Date
```
GET /random/date?start=2020-01-01&end=2025-12-31
→ date=2023-06-15T00:00:00Z unix=1686787200
```

### UUID v4
```
GET /uuid
→ uuid=550e8400-e29b-41d4-a716-446655440000
```

### JSON Responses
Add `Accept: application/json` header or `?format=json` query param.

### Help
```
GET /help        — Full API reference
GET /health      — Health check
GET /.well-known/agent.md  — Agent-readable manual
```

## Configuration

| Flag | Env | Default | Description |
|------|-----|---------|-------------|
| `-addr` | `RANDOMKIT_ADDR` | `:8787` | Listen address |
| `-secret` | `RANDOMKIT_SECRET` | (auto) | Token signing secret |

## Build

```bash
make build    # CGO_ENABLED=0 go build -trimpath
make test     # go test -race ./...
make vet      # go vet ./...
```

## License

MIT
