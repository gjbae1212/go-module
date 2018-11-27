package ssh_cert_authority

import (
	"crypto/rand"
	"testing"

	"time"

	"encoding/pem"

	"github.com/mikesmitty/edkey"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

var (
		testCaPem = `-----BEGIN OPENSSH PRIVATE KEY-----
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
1wnFlQBrkKkAAAADYWFh
-----END OPENSSH PRIVATE KEY-----`
)

func TestNewCertificationED25519(t *testing.T) {
	assert := assert.New(t)

	pri, pub, cert, err := NewCertificationED25519([]byte(testCaPem), &CertMeta{
		KeyId:      "allan",
		Principals: []string{"ubuntu"},
		Before:     time.Now().Add(time.Hour),
		After:      time.Now(),
	})
	assert.NoError(err)

	block, _ := pem.Decode([]byte(pri))
	assert.Equal(block.Type, "OPENSSH PRIVATE KEY")
	pkeySigner, err := ssh.ParsePrivateKey([]byte(pri))
	assert.NoError(err)
	assert.Equal(string(ssh.MarshalAuthorizedKey(pkeySigner.PublicKey())), pub)
	_ = cert
}

func TestNewCertificationED25519WithPKey(t *testing.T) {
	assert := assert.New(t)

	_, pkey, err := ed25519.GenerateKey(rand.Reader)
	assert.NoError(err)
	pkeyBlock := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(pkey),
	}
	private := pem.EncodeToMemory(pkeyBlock)

	cert, err := NewCertificationED25519WithPKey([]byte(testCaPem), private, &CertMeta{
		KeyId:      "allan",
		Principals: []string{"ubuntu"},
		Before:     time.Now().Add(time.Hour),
		After:      time.Now(),
	})
	assert.NoError(err)
	_ = cert

	meta := &CertMeta{
		KeyId:      "allan",
		Principals: []string{"ubuntu"},
		Before:     time.Now().Add(2 * time.Minute),
		After:      time.Now(),
	}

	pri, _, cert1, err := NewCertificationED25519([]byte(testCaPem), meta)
	assert.NoError(err)

	cert2, err := NewCertificationED25519WithPKey([]byte(testCaPem), []byte(pri), meta)
	assert.NoError(err)
	assert.NotEqual(cert1, cert2)
}
