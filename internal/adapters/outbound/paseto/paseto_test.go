package paseto_test

import (
	"encoding/base64"
	"testing"
	"time"
	"uuid"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/paseto"
	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"github.com/go-openapi/testify/require"
)

func decodeB64Key(t *testing.T, key string) []byte {
	decodedKey, err := base64.StdEncoding.DecodeString(key)
	require.NoError(t, err)
	return decodedKey
}
func newTestClaims(t *testing.T, email string, userID int64, duration time.Duration) ports.Claims {
	claim, err := ports.NewClaims(email, userID, duration)
	require.NoError(t, err)
	return claim
}
func TestPasetoAuthenticator_CreateToken(t *testing.T) {
	keyB64 := "0wKIJsaVqOKCE9rWYnTZD8ASPn8Ks9YCXpQfwDjQ/b4="
	decodedKey := decodeB64Key(t, keyB64)
	correctClaim := newTestClaims(t, "test@test.com", 1, 5*time.Minute)
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		symmetricKey []byte
		// Named input parameters for target function.
		claims  ports.Claims
		wantErr bool
	}{
		{
			name:         "should create valid token",
			symmetricKey: decodedKey,
			claims:       correctClaim,
			wantErr:      false,
		}, {
			name:         "should not create a token if UUID is invalid",
			symmetricKey: decodedKey,
			claims:       ports.Claims{ID: uuid.Nil(), UserID: 1, Email: "false@email.com", IssuedAt: time.Now(), ExpiresAt: time.Now().Add(5 * time.Minute)},
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := paseto.NewPasetoAuthenticator(tt.symmetricKey)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, gotErr := p.CreateToken(tt.claims)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateToken() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateToken() succeeded unexpectedly")
			}
			verifiedClaims, err := p.VerifyToken(got)
			require.NoError(t, err)
			require.Equal(t, tt.claims.ID, verifiedClaims.ID)
			require.Equal(t, tt.claims.Email, verifiedClaims.Email)
			require.Equal(t, tt.claims.UserID, verifiedClaims.UserID)
			require.WithinDuration(t, tt.claims.IssuedAt, verifiedClaims.IssuedAt, time.Second)
			require.WithinDuration(t, tt.claims.ExpiresAt, verifiedClaims.ExpiresAt, time.Second)
		})
	}
}
