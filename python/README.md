# identify-sdk (Python)

Python client for [identify](https://github.com/1909-pramod/identify-sdk) — a self-hosted modular authentication service.

## Installation

```bash
pip install identify-sdk
```

## Usage

```python
from identify_sdk import IdentifyClient, IdentifyError

client = IdentifyClient(base_url="http://localhost:8080")

# Exchange username + password for a JWT
resp = client.login_with_password("user@example.com", "secret")
token = resp.token

# Validate a JWT
result = client.validate_token(token)
print(result.valid, result.identity.primary_identifier)

# Get identity for an authenticated token
identity = client.get_me(token)
```

## Error handling

```python
try:
    resp = client.login_with_password("user@example.com", "wrong")
except IdentifyError as e:
    print(e.status_code, str(e))  # 401, "HTTP 401: ..."
```

Full documentation: [github.com/1909-pramod/identify-sdk](https://github.com/1909-pramod/identify-sdk)
