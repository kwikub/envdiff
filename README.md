# envdiff

> Diff and reconcile `.env` files across environments with secret masking.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git && cd envdiff && go build -o envdiff .
```

---

## Usage

Compare two `.env` files and mask sensitive values:

```bash
envdiff .env.development .env.production
```

**Example output:**

```
~ DB_HOST        dev-db.local  →  prod-db.example.com
~ API_KEY        [masked]      →  [masked]
+ NEW_FEATURE_FLAG             →  true
- DEPRECATED_VAR  old_value    →  (missing)
```

### Flags

| Flag | Description |
|------|-------------|
| `--mask` | Comma-separated list of key patterns to mask (default: `*KEY*,*SECRET*,*TOKEN*,*PASSWORD*`) |
| `--output` | Output format: `text`, `json`, or `dotenv` (default: `text`) |
| `--reconcile` | Write a reconciled `.env` file merging both inputs |

### Reconcile example

```bash
envdiff .env.local .env.production --reconcile --output dotenv > .env.merged
```

---

## Why envdiff?

Managing `.env` files across staging, production, and local environments is error-prone. `envdiff` makes it easy to spot missing keys, changed values, and accidental secret exposure before deployments.

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

[MIT](LICENSE) © 2024 yourusername