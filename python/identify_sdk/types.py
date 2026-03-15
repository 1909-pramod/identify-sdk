from dataclasses import dataclass, field
from typing import Optional


@dataclass
class Identity:
    id: int
    primary_identifier: str
    identifier_group_name: Optional[str] = None
    identifier_group_type: Optional[str] = None
    user_type: Optional[str] = None
    expiry_at: Optional[str] = None
    roles: Optional[list[str]] = None
    is_super_admin: bool = False


@dataclass
class TokenResponse:
    token: str
    token_type: str
    expires_in: int


@dataclass
class ValidateResponse:
    valid: bool
    identity: Identity
    expires_at: Optional[str] = None
