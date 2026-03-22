# identify-sdk (JavaScript / TypeScript)

TypeScript client for [identify](https://github.com/1909-pramod/identify-sdk) — a self-hosted modular authentication service.

## Installation

```bash
npm install @1909-pramod/identify-sdk
# or
pnpm add @1909-pramod/identify-sdk
```

> ESM-only. Requires Node 18+ or any modern browser runtime.

## Usage

```ts
import { IdentifyClient, IdentifyError } from "@1909-pramod/identify-sdk";

const client = new IdentifyClient({ baseUrl: "http://localhost:8080" });

// Exchange username + password for a JWT
const { token } = await client.loginWithPassword("user@example.com", "secret");

// Use a Google ID token
const { token } = await client.loginWithGoogle(googleIdToken);

// Validate a JWT and get the resolved identity
const { valid, identity, expires_at } = await client.validateToken(token);

// Get identity for an authenticated token
const identity = await client.getMe(token);

// Refresh an existing JWT
const { token: newToken } = await client.refreshToken(token);
```

## Error handling

```ts
try {
  const { token } = await client.loginWithPassword("user@example.com", "wrong");
} catch (err) {
  if (err instanceof IdentifyError) {
    console.error(err.status, err.message); // 401, "..."
  }
}
```

## API reference

| Method | identify endpoint | Description |
|--------|-------------------|-------------|
| `loginWithPassword(username, password)` | `POST /auth` | Basic credentials → JWT |
| `loginWithGoogle(googleIdToken)` | `POST /auth` | Google ID token → JWT |
| `refreshToken(token)` | `POST /auth` | Existing JWT → fresh JWT |
| `validateToken(token)` | `POST /token` | Validate JWT, return identity |
| `getMe(token)` | `GET /me` | Return identity for a token |

## Building from source

```bash
npm install
npm run build      # outputs to dist/
npm run typecheck  # type-check without emitting
```

Full documentation: [github.com/1909-pramod/identify-sdk](https://github.com/1909-pramod/identify-sdk)
