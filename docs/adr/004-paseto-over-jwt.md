# 4. Use PASETO v4 instead of JWT for Authentication

## Status

Accepted

## Context

JSON Web Tokens (JWT) are the industry standard, but they suffer from design flaws, most notably the algorithm confusion attack (where an attacker changes the header algorithm from RS256 to HS256). Even with libraries, misconfiguration is a common security risk.

## Decision

We will use **PASETO (Platform-Agnostic Security Tokens) v4**, specifically the `v4.local` (symmetric encryption) variant, using the `aidantwoods/go-paseto` library.

## Consequences

- **Positive:** Security by default. PASETO eliminates algorithm negotiation, forcing the use of secure cryptographic primitives (XChaCha20 + BLAKE2b).
- **Negative:** Smaller ecosystem compared to JWT. Fewer ready-made tools and libraries for debugging.
