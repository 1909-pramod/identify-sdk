import base64
from typing import Any

import requests

from .types import Identity, TokenResponse, ValidateResponse


class IdentifyError(Exception):
    def __init__(self, status_code: int, message: str) -> None:
        super().__init__(f"HTTP {status_code}: {message}")
        self.status_code = status_code


class IdentifyClient:
    def __init__(self, base_url: str) -> None:
        self._base_url = base_url.rstrip("/")

    def _post(self, path: str, auth_header: str) -> dict[str, Any]:
        resp = requests.post(
            f"{self._base_url}{path}",
            headers={"Authorization": auth_header},
        )
        body: dict[str, Any] = resp.json()
        if not resp.ok:
            raise IdentifyError(resp.status_code, body.get("error", "request failed"))
        return body

    def _get(self, path: str, token: str) -> dict[str, Any]:
        resp = requests.get(
            f"{self._base_url}{path}",
            headers={"Authorization": f"Bearer {token}"},
        )
        body: dict[str, Any] = resp.json()
        if not resp.ok:
            raise IdentifyError(resp.status_code, body.get("error", "request failed"))
        return body

    @staticmethod
    def _parse_identity(data: dict[str, Any]) -> Identity:
        return Identity(
            id=data["id"],
            primary_identifier=data["primary_identifier"],
            identifier_group_name=data.get("identifier_group_name"),
            identifier_group_type=data.get("identifier_group_type"),
            user_type=data.get("user_type"),
            expiry_at=data.get("expiry_at"),
            roles=data.get("roles"),
            is_super_admin=data.get("is_super_admin", False),
        )

    @staticmethod
    def _parse_token_response(data: dict[str, Any]) -> TokenResponse:
        return TokenResponse(
            token=data["token"],
            token_type=data["token_type"],
            expires_in=data["expires_in"],
        )

    def login_with_password(self, username: str, password: str) -> TokenResponse:
        """Exchange username/password for a JWT."""
        encoded = base64.b64encode(f"{username}:{password}".encode()).decode()
        return self._parse_token_response(self._post("/auth", f"Basic {encoded}"))

    def login_with_google(self, google_id_token: str) -> TokenResponse:
        """Exchange a Google ID token for an internal JWT."""
        return self._parse_token_response(
            self._post("/auth", f"Bearer {google_id_token}")
        )

    def refresh_token(self, token: str) -> TokenResponse:
        """Exchange an existing internal JWT for a fresh one."""
        return self._parse_token_response(self._post("/auth", f"Bearer {token}"))

    def validate_token(self, token: str) -> ValidateResponse:
        """Validate a JWT and return the resolved identity."""
        data = self._post("/token", f"Bearer {token}")
        return ValidateResponse(
            valid=data["valid"],
            identity=self._parse_identity(data["identity"]),
            expires_at=data.get("expires_at"),
        )

    def get_me(self, token: str) -> Identity:
        """Return the identity for the authenticated token."""
        return self._parse_identity(self._get("/me", token))
