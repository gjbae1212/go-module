package ssh_cert_authority

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/pem"
	"strings"
	"time"

	mssh "github.com/gjbae1212/go-module/ssh"
	"github.com/mikesmitty/edkey"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
)

type CertMeta struct {
	KeyId      string
	Principals []string
	Before     time.Time
	After      time.Time
}

// Only Support ED25519 Algorithm
func NewCertificationED25519(caPem []byte, meta *CertMeta) (private string, public string, cert string, err error) {
	if len(caPem) == 0 || meta == nil {
		err = mssh.EmptyError.New("NewCertification")
		return
	}
	// generate ed25519
	_, pkey, suberr := ed25519.GenerateKey(rand.Reader)
	if suberr != nil {
		err = suberr
		return
	}
	// convert pri ed25519 bytes to pri pem
	pkeyBlock := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: edkey.MarshalED25519PrivateKey(pkey),
	}
	private = string(pem.EncodeToMemory(pkeyBlock))

	sshPub, suberr := ssh.NewPublicKey(pkey.Public())
	if suberr != nil {
		err = suberr
		return
	}
	public = string(ssh.MarshalAuthorizedKey(sshPub))
	// generate certification
	sshcert := &ssh.Certificate{
		Key:             sshPub, // pubkey
		Serial:          0,
		CertType:        ssh.UserCert,
		KeyId:           meta.KeyId,                 // comment
		ValidPrincipals: meta.Principals,            // user
		ValidBefore:     uint64(meta.Before.Unix()), // end
		ValidAfter:      uint64(meta.After.Unix()),  // start
	}
	sshcert.Extensions = make(map[string]string)
	sshcert.Extensions["permit-X11-forwarding"] = ""
	sshcert.Extensions["permit-agent-forwarding"] = ""
	sshcert.Extensions["permit-port-forwarding"] = ""
	sshcert.Extensions["permit-pty"] = ""
	sshcert.Extensions["permit-user-rc"] = ""

	cert, err = createCertification(sshcert, caPem)
	return
}

func NewCertificationED25519WithPKey(caPem []byte, pkeyPem []byte, meta *CertMeta) (cert string, err error) {
	if len(caPem) == 0 || len(pkeyPem) == 0 || meta == nil {
		err = mssh.EmptyError.New("NewCertificationWithPKEY")
		return
	}
	block, _ := pem.Decode(pkeyPem)
	if block.Type != "OPENSSH PRIVATE KEY" {
		err = mssh.InvalidParamsError.New("NewCertificationED25519WithPKey")
		return
	}
	pkey, suberr := ssh.ParsePrivateKey(pkeyPem)
	if suberr != nil {
		err = suberr
		return
	}
	// generate certification
	sshcert := &ssh.Certificate{
		Key:             pkey.PublicKey(), // pubkey
		Serial:          0,
		CertType:        ssh.UserCert,
		KeyId:           meta.KeyId,                 // comment
		ValidPrincipals: meta.Principals,            // user
		ValidBefore:     uint64(meta.Before.Unix()), // end
		ValidAfter:      uint64(meta.After.Unix()),  // start
	}
	sshcert.Extensions = make(map[string]string)
	sshcert.Extensions["permit-X11-forwarding"] = ""
	sshcert.Extensions["permit-agent-forwarding"] = ""
	sshcert.Extensions["permit-port-forwarding"] = ""
	sshcert.Extensions["permit-pty"] = ""
	sshcert.Extensions["permit-user-rc"] = ""

	cert, err = createCertification(sshcert, caPem)
	return
}

func createCertification(cert *ssh.Certificate, caPem []byte) (string, error) {
	if cert == nil || len(caPem) == 0 {
		return "", mssh.EmptyError.New("createCertification")
	}
	signer, err := ssh.ParsePrivateKey(caPem)
	if err != nil {
		return "", err
	}
	// sign certification
	if err = cert.SignCert(rand.Reader, signer); err != nil {
		return "", err
	}
	body := cert.Marshal()
	baseBody := base64.StdEncoding.EncodeToString(body)
	return strings.Join([]string{ssh.CertAlgoED25519v01, baseBody, cert.KeyId}, " "), nil
}
