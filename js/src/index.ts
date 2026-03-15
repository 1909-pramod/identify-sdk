export type {
  Identity,
  TokenResponse,
  ValidateResponse,
  IdentifyClientOptions,
} from "./types.js";

import type {
  Identity,
  TokenResponse,
  ValidateResponse,
  IdentifyClientOptions,
} from "./types.js";

export class IdentifyError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.name = "IdentifyError";
    this.status = status;
  }
}

export class IdentifyClient {
  private readonly baseUrl: string;

  constructor(options: IdentifyClientOptions) {
    this.baseUrl = options.baseUrl.replace(/\/$/, "");
  }

  private async post<T>(path: string, authHeader: string): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method: "POST",
      headers: { Authorization: authHeader },
    });
    const body = await res.json();
    if (!res.ok) {
      throw new IdentifyError(res.status, body.error ?? "request failed");
    }
    return body as T;
  }

  private async get<T>(path: string, token: string): Promise<T> {
    const res = await fetch(`${this.baseUrl}${path}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    const body = await res.json();
    if (!res.ok) {
      throw new IdentifyError(res.status, body.error ?? "request failed");
    }
    return body as T;
  }

  /** Exchange username + password for a JWT. */
  async loginWithPassword(
    username: string,
    password: string
  ): Promise<TokenResponse> {
    const encoded = btoa(`${username}:${password}`);
    return this.post<TokenResponse>("/auth", `Basic ${encoded}`);
  }

  /** Exchange a Google ID token for an internal JWT. */
  async loginWithGoogle(googleIdToken: string): Promise<TokenResponse> {
    return this.post<TokenResponse>("/auth", `Bearer ${googleIdToken}`);
  }

  /** Exchange an existing internal JWT for a fresh one. */
  async refreshToken(token: string): Promise<TokenResponse> {
    return this.post<TokenResponse>("/auth", `Bearer ${token}`);
  }

  /** Validate a JWT and return the resolved identity. */
  async validateToken(token: string): Promise<ValidateResponse> {
    return this.post<ValidateResponse>("/token", `Bearer ${token}`);
  }

  /** Return the identity for the authenticated token. */
  async getMe(token: string): Promise<Identity> {
    return this.get<Identity>("/me", token);
  }
}
