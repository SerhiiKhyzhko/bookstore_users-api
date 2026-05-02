package cryptoutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBcrypt(t *testing.T) {
	t.Run("HashAndCompareSuccess", func(t *testing.T) {
		hash, err := GetBcrypt("qwerty")
		assert.NoError(t, err)
		assert.NotEmpty(t, hash)

		err = ComparePassword(hash, "qwerty")
		assert.NoError(t, err)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		hash, err := GetBcrypt("qwerty")
		assert.NoError(t, err)

		err = ComparePassword(hash, "wrong")
		assert.Error(t, err)
		assert.Equal(t, "invalid credentials", err.Error())
	})

	t.Run("EmptyPassword", func(t *testing.T) {
		hash, err := GetBcrypt("")
		assert.NoError(t, err)

		err = ComparePassword(hash, "")
		assert.NoError(t, err)
	})
}