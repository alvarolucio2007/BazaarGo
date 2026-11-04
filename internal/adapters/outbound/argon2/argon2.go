package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	ports "github.com/alvarolucio2007/bazaargo/internal/ports/outbound"
	"golang.org/x/crypto/argon2"
)

const (
	SaltLength = 16
	format     = "$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s"
)

var InvalidHashFormatError = errors.New("argon2: invalid hash format")

var _ ports.Hasher = (*Argon2Hasher)(nil)

type Argon2Hasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func NewArgon2Hasher(memory, iterations uint32, parallelism uint8) *Argon2Hasher {
	return &Argon2Hasher{
		memory:      memory,      //64*1024
		iterations:  iterations,  // 3
		parallelism: parallelism, // 2
		keyLength:   32,
	}
}

// Hash implements [ports.Hasher].
// It hashes a password with the argon2id algorithm and returns a standard formatted string.
//
// Returns an error if the system's secure random number generator fails.
func (a *Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hashedBytes := argon2.IDKey([]byte(password), salt, a.iterations, a.memory, a.parallelism, a.keyLength)

	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hashedBytes)

	fullHash := fmt.Sprintf(format, a.memory, a.iterations, a.parallelism, saltBase64, hashBase64)
	return fullHash, nil
}

// Check implements [ports.Hasher].
// It compares a plaintext password with a hashed password, returning true if they match.
//
// Returns [InvalidHashFormatError] if the hash format is malformed, or a base64 decoding
// error if the salt or hash cannot be decoded.
func (a *Argon2Hasher) Check(password string, hash string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		return false, InvalidHashFormatError
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8

	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return false, InvalidHashFormatError
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	originalHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		iterations,
		memory,
		parallelism,
		uint32(len(originalHash)),
	)

	if subtle.ConstantTimeCompare(computedHash, originalHash) == 1 {
		return true, nil
	}
	return false, nil
}
