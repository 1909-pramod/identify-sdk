export interface Identity {
  id: number;
  primary_identifier: string;
  identifier_group_name?: string;
  identifier_group_type?: string;
  user_type?: string;
  expiry_at?: string;
  roles?: string[];
  is_super_admin?: boolean;
}

export interface TokenResponse {
  token: string;
  token_type: string;
  expires_in: number;
}

export interface ValidateResponse {
  valid: boolean;
  identity: Identity;
  expires_at?: string;
}

export interface IdentifyClientOptions {
  baseUrl: string;
}
