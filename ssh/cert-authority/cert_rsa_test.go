package ssh_cert_authority

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"crypto/rsa"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/ssh"
)

func TestNewCertificationRSA(t *testing.T) {
	assert := assert.New(t)

	pri, pub, cert, err := NewCertificationRSA([]byte(testCaPem), &CertMeta{
		KeyId:      "allan",
		Principals: []string{"ubuntu"},
		Before:     time.Now().Add(time.Hour),
		After:      time.Now(),
	}, 4096)
	assert.NoError(err)

	block, _ := pem.Decode([]byte(pri))
	assert.Equal(block.Type, "RSA PRIVATE KEY")
	pkeySigner, err := ssh.ParsePrivateKey([]byte(pri))
	assert.NoError(err)
	assert.Equal(string(ssh.MarshalAuthorizedKey(pkeySigner.PublicKey())), pub)
	_ = cert
}

func TestNewCertificationRSAWithPub(t *testing.T) {
	assert := assert.New(t)

	pkey, err := rsa.GenerateKey(rand.Reader, 4096)
	assert.NoError(err)

	pkeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pkey),
	}
	private := string(pem.EncodeToMemory(pkeyBlock))
	sshPub, err := ssh.NewPublicKey(&pkey.PublicKey)
	assert.NoError(err)
	public := string(ssh.MarshalAuthorizedKey(sshPub))

	cert, err := NewCertificationRSAWithPub([]byte(testCaPem), []byte(public), &CertMeta{
		KeyId:      "allan",
		Principals: []string{"ubuntu"},
		Before:     time.Now().Add(time.Hour),
		After:      time.Now(),
	})

	assert.NoError(err)
	_ = cert
	_ = private
}
