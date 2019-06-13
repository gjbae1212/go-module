package encryption

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptAES(t *testing.T) {
	assert := assert.New(t)

	tests := map[string]struct {
		key  []byte
		text []byte
		err  bool
	}{
		"error1":       {key: []byte("invalid key"), text: []byte{}, err: true},
		"empty-string": {key: []byte("1111111111111111"), text: []byte{}, err: false},
		"good":         {key: []byte("1111111111111111"), text: []byte("hi"), err: false},
	}
	for _, t := range tests {
		_, err := EncryptAES(t.key, t.text)
		if t.err {
			assert.NotNil(err)
		} else {
			assert.Nil(err)
		}
	}

}

func TestDecryptAES(t *testing.T) {
	assert := assert.New(t)

	validkey := "allanallanallan-"
	text := "hello world"

	encrypt, err := EncryptAES([]byte(validkey), []byte(text))
	assert.NoError(err)

	decrypt, err := DecryptAES([]byte(validkey), encrypt)
	assert.NoError(err)
	assert.Equal(text, string(decrypt))
}
