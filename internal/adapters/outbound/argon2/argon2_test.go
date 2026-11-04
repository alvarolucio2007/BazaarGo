package argon2_test

import (
	"testing"

	"github.com/alvarolucio2007/bazaargo/internal/adapters/outbound/argon2"
	"github.com/go-faker/faker/v4"
	"github.com/go-openapi/testify/require"
)

type FakePassword struct {
	Password string `faker:"word"`
}

func TestPassword(t *testing.T) {
	hasher := argon2.NewArgon2Hasher(64*1024, 2, 3)
	fakePassword := FakePassword{}
	err := faker.FakeData(&fakePassword)
	require.NoError(t, err)
	hashedPassword, err := hasher.Hash(fakePassword.Password)

	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword)

	isSame, err := hasher.Check(fakePassword.Password, hashedPassword)
	require.NoError(t, err)
	require.True(t, isSame)

	hashedPassword2, err := hasher.Hash(fakePassword.Password)
	require.NoError(t, err)
	require.NotEmpty(t, hashedPassword2)
	require.NotEqual(t, hashedPassword, hashedPassword2)

	wrongPassword := "THIS IS A WRONG PASSWORD"
	isSame, err = hasher.Check(wrongPassword, hashedPassword)
	require.NoError(t, err)
	require.False(t, isSame)
}
