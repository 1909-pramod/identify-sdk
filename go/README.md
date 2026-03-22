# identify-sdk (Go)

Go client for [identify](https://github.com/1909-pramod/identify-sdk) — a self-hosted modular authentication service.

## Installation

```bash
go get github.com/1909-pramod/identify-sdk/go@latest
```

## Usage

```go
import "github.com/1909-pramod/identify-sdk/go/identify"

client := identify.New("http://localhost:8080")

// Exchange username + password for a JWT
resp, err := client.LoginWithPassword(ctx, "user@example.com", "secret")
if err != nil {
    log.Fatal(err)
}
token := resp.Token

// Use a Google ID token
resp, err = client.LoginWithGoogle(ctx, googleIDToken)

// Validate a JWT and get the resolved identity
result, err := client.ValidateToken(ctx, token)
fmt.Println(result.Valid, result.Identity.PrimaryIdentifier)

// Get identity for an authenticated token
identity, err := client.GetMe(ctx, token)

// Refresh an existing JWT
resp, err = client.RefreshToken(ctx, token)
```

## Error handling

```go
resp, err := client.LoginWithPassword(ctx, "user@example.com", "wrong")
if err != nil {
    var idErr *identify.Error
    if errors.As(err, &idErr) {
        fmt.Println(idErr.StatusCode, idErr.Message) // 401, "..."
    }
}
```

## API reference

| Method | identify endpoint | Description |
|--------|-------------------|-------------|
| `LoginWithPassword(ctx, username, password)` | `POST /auth` | Basic credentials → JWT |
| `LoginWithGoogle(ctx, googleIDToken)` | `POST /auth` | Google ID token → JWT |
| `RefreshToken(ctx, token)` | `POST /auth` | Existing JWT → fresh JWT |
| `ValidateToken(ctx, token)` | `POST /token` | Validate JWT, return identity |
| `GetMe(ctx, token)` | `GET /me` | Return identity for a token |

Full documentation: [github.com/1909-pramod/identify-sdk](https://github.com/1909-pramod/identify-sdk)
