# identify-sdk

Client SDKs for [**identify**](https://github.com/1909-pramod/identify) — a self-hosted, modular authentication service.

Instead of calling the identify HTTP API by hand, use these clients to log in, validate tokens, and fetch identities with a single method call.

## Packages

| Language | Package | Source |
|----------|---------|--------|
| TypeScript / JavaScript | `@identify/sdk` | [`js/`](js/) |
| Python | `identify-sdk` | [`python/`](python/) |

---

## JavaScript / TypeScript

### Installation

```bash
npm install @identify/sdk
# or
pnpm add @identify/sdk
```

> The package is ESM-only and requires Node 18+ or any modern browser runtime.

### Usage

```ts
import { IdentifyClient, IdentifyError } from "@identify/sdk";

const client = new IdentifyClient({ baseUrl: "http://localhost:8080" });

// Exchange username + password for a JWT
const { token } = await client.loginWithPassword("user@example.com", "secret");

// Use a Google ID token (when the google provider is enabled on the server)
const { token } = await client.loginWithGoogle(googleIdToken);

// Validate a JWT and get the resolved identity
const { valid, identity, expires_at } = await client.validateToken(token);

// Fetch the identity for an already-authenticated token
const identity = await client.getMe(token);

// Refresh an existing JWT
const { token: newToken } = await client.refreshToken(token);
```

### Error handling

All methods throw `IdentifyError` on non-2xx responses:

```ts
try {
  const { token } = await client.loginWithPassword("user@example.com", "wrong");
} catch (err) {
  if (err instanceof IdentifyError) {
    console.error(err.status, err.message); // 401, "..."
  }
}
```

### Building from source

```bash
cd js
npm install
npm run build      # outputs to js/dist/
npm run typecheck  # type-check without emitting
```

---

## Python

### Installation

```bash
pip install identify-sdk
```

> Requires Python 3.10+ and [`requests`](https://pypi.org/project/requests/).

### Usage

```python
from identify_sdk import IdentifyClient, IdentifyError

client = IdentifyClient(base_url="http://localhost:8080")

# Exchange username + password for a JWT
resp = client.login_with_password("user@example.com", "secret")
token = resp.token

# Use a Google ID token (when the google provider is enabled on the server)
resp = client.login_with_google(google_id_token)

# Validate a JWT and get the resolved identity
result = client.validate_token(token)
print(result.valid, result.identity.primary_identifier)

# Fetch the identity for an already-authenticated token
identity = client.get_me(token)

# Refresh an existing JWT
resp = client.refresh_token(token)
```

### Error handling

All methods raise `IdentifyError` on non-2xx responses:

```python
from identify_sdk import IdentifyClient, IdentifyError

try:
    resp = client.login_with_password("user@example.com", "wrong")
except IdentifyError as e:
    print(e.status_code, str(e))  # 401, "HTTP 401: ..."
```

### Development setup

```bash
cd python
pip install -e ".[dev]"
mypy identify_sdk   # strict type checking
```

---

## API reference

Both clients cover the same surface of the identify HTTP API:

| Method | identify endpoint | Description |
|--------|-------------------|-------------|
| `loginWithPassword` / `login_with_password` | `POST /auth` | Basic credentials → JWT |
| `loginWithGoogle` / `login_with_google` | `POST /auth` | Google ID token → JWT |
| `refreshToken` / `refresh_token` | `POST /auth` | Existing JWT → fresh JWT |
| `validateToken` / `validate_token` | `POST /token` | Validate JWT, return identity |
| `getMe` / `get_me` | `GET /me` | Return identity for a token |

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

Bug reports and pull requests welcome. New language SDKs (Go, Ruby, Java, etc.) are a great contribution — follow the same method names and error type pattern used in the existing clients.

## License

MIT — see [LICENSE](../identify/LICENSE) in the main repository.
