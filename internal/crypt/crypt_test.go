package crypt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncDec(t *testing.T) {
	salt := "34e0320b94132dfe03960845786790e8"
	token := "d5fc830f9efcea6a8381a3ecd7360f13ca178e31b25096638a326cef9ac138"

	crypter := NewCrypter([]byte(salt))
	encToken, err := crypter.Encrypt(token)
	require.NoError(t, err)

	fmt.Println("enc", encToken)

	decToken, err := crypter.Decrypt(encToken)
	require.NoError(t, err)

	assert.Equal(t, token, decToken)

}
