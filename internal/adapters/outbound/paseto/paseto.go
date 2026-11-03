package paseto

import (
	"fmt"
	"uuid"

	"aidanwoods.dev/go-paseto"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
)

var _ ports.Authenticator = (*PasetoAuthenticator)(nil)

type PasetoAuthenticator struct {
	key paseto.V4SymmetricKey
}

// NewPasetoAuthenticator creates a new [PasetoAuthenticator].
//
// Returns an error if, for whatever reason, the symmetric key is invalid.
func NewPasetoAuthenticator(symmetricKey []byte) (*PasetoAuthenticator, error) {
	key, err := paseto.V4SymmetricKeyFromBytes(symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("paseto: invalid symmetric key: %w", err)
	}
	return &PasetoAuthenticator{key: key}, nil
}

// CreateToken implements [ports.Authenticator].
// CreateToken creates a secure PASETO V4 token for user authentication.
//
// Returns an error if the user_id cannot be set in the token payload.
func (p *PasetoAuthenticator) CreateToken(claims ports.Claims) (string, error) {
	token := paseto.NewToken()
	token.SetIssuedAt(claims.IssuedAt)
	token.SetExpiration(claims.ExpiresAt)
	token.SetString("id", claims.ID.String())
	token.SetString("email", claims.Email)
	err := token.Set("user_id", claims.UserID)
	if err != nil {
		return "", fmt.Errorf("paseto: CreateToken: couldn't set user_id: %w", err)
	}
	return token.V4Encrypt(p.key, nil), nil
}

// VerifyToken implements [ports.Authenticator].
// VerifyToken verifies a secure PASETO V4 token for user authentication.
//
// Returns [ErrInvalidToken] if the token is expired, malformed, or has an invalid signature.
// It may also return custom wrapping errors if the token payload is missing expected fields.
func (p *PasetoAuthenticator) VerifyToken(tokenString string) (ports.Claims, error) {
	parser := paseto.NewParser()
	parsedToken, err := parser.ParseV4Local(p.key, tokenString, nil)
	if err != nil {
		return ports.Claims{}, ErrInvalidToken
	}
	idStr, err := parsedToken.GetString("id")
	if err != nil {
		return ports.Claims{}, fmt.Errorf("paseto: failed to parse id: %w", err)
	}
	parsedID, err := uuid.Parse(idStr)
	if err != nil {
		return ports.Claims{}, fmt.Errorf("paseto: failed to parse uuid: %w", err)
	}
	email, err := parsedToken.GetString("email")
	if err != nil {
		return ports.Claims{}, fmt.Errorf("paseto: failed to parse email: %w", err)
	}
	issuedAt, err := parsedToken.GetIssuedAt()
	if err != nil {
		return ports.Claims{}, fmt.Errorf("paseto: failed to parse issued_at: %w", err)
	}
	expiresAt, err := parsedToken.GetExpiration()
	if err != nil {
		return ports.Claims{}, fmt.Errorf("paseto: failed to parse expires_at: %w", err)
	}
	var userID int64
	err = parsedToken.Get("user_id", &userID)
	if err != nil {
		return ports.Claims{}, fmt.Errorf("paseto: failed to parse user_id: %w", err)
	}
	claims := ports.Claims{
		ID:        parsedID,
		Email:     email,
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		UserID:    userID,
	}
	return claims, nil
}
