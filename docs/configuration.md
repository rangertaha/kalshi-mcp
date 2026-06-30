# Configuration

All configuration is read from the environment.

| Variable             | Required | Description                                              |
| -------------------- | :------: | -------------------------------------------------------- |
| `KALSHI_API_KEY_ID`  |    no    | Access key ID (enables authenticated tools).             |
| `KALSHI_PRIVATE_KEY` |    no    | PEM-encoded RSA private key for request signing.         |
| `KALSHI_BASE_URL`    |    no    | API base URL (default elections trade-api/v2).           |
| `KALSHI_TOOLSETS`    |    no    | Comma-separated toolset names to enable, or `all`.       |
| `KALSHI_READONLY`    |    no    | `true` to expose only read-only tools.                   |

Public market-data tools work without credentials.
