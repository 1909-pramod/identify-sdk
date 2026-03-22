# identify-sdk

[![CI](https://github.com/1909-pramod/identify-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/1909-pramod/identify-sdk/actions/workflows/ci.yml)
[![npm](https://img.shields.io/npm/v/@1909-pramod/identify-sdk)](https://www.npmjs.com/package/@1909-pramod/identify-sdk)
[![PyPI](https://img.shields.io/pypi/v/identify-sdk)](https://pypi.org/project/identify-sdk/)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/1909-pramod/identify/blob/main/LICENSE)

Client SDKs for [**identify**](https://github.com/1909-pramod/identify) — a self-hosted, modular authentication service.

Instead of calling the identify HTTP API by hand, use these clients to log in, validate tokens, and fetch identities with a single method call.

## Packages

| Language | Package | Source | Docs |
|----------|---------|--------|------|
| TypeScript / JavaScript | `@1909-pramod/identify-sdk` | [`js/`](js/) | [js/README.md](js/README.md) |
| Python | `identify-sdk` | [`python/`](python/) | [python/README.md](python/README.md) |
| Go | `github.com/1909-pramod/identify-sdk/go` | [`go/`](go/) | [go/README.md](go/README.md) |

---

## Running identify locally

You need a running identify server. The quickest way:

```bash
git clone https://github.com/1909-pramod/identify
cd identify
ENABLED_PROVIDERS=basic,jwt JWT_SECRET=dev-secret go run main.go
```

Or with Docker / Podman:

```bash
docker compose up   # or: podman compose up
```

See the [identify README](https://github.com/1909-pramod/identify#readme) for the full list of providers (Google OIDC, GitHub, Azure AD, Okta, SAML) and configuration options.

---

## Contributing

Bug reports and pull requests welcome. New language SDKs (Ruby, Java, etc.) are a great contribution — follow the same method names and error type pattern used in the existing clients.

## License

MIT — see [LICENSE](LICENSE).
