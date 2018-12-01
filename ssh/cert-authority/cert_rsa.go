package ssh_cert_authority

import (
	"crypto/rand"
	"crypto/rsa"

	"encoding/pem"

	"crypto/x509"

	"encoding/base64"
	"strings"

	mssh "github.com/gjbae1212/go-module/ssh"
	"golang.org/x/crypto/ssh"
)

func NewCertificationRSA(caPem []byte, meta *CertMeta, bitSize int) (private string, public string, cert string, err error) {
	if len(caPem) == 0 || meta == nil || bitSize <= 0 {
		err = mssh.EmptyError.New("NewCertificationRSA")
		return
	}
	pkey, suberr := rsa.GenerateKey(rand.Reader, bitSize)
	if suberr != nil {
		err = suberr
		return
	}
	if suberr := pkey.Validate(); suberr != nil {
		err = suberr
		return
	}

	// convert pri rsa bytes to pri pem
	pkeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(pkey),
	}
	private = string(pem.EncodeToMemory(pkeyBlock))

	sshPub, suberr := ssh.NewPublicKey(&pkey.PublicKey)
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

	cert, err = createCertificationRSA(sshcert, caPem)
	return
}

func NewCertificationRSAWithPub(caPem []byte, pubAuthKey []byte, meta *CertMeta) (cert string, err error) {
	if len(caPem) == 0 || len(pubAuthKey) == 0 || meta == nil {
		err = mssh.EmptyError.New("NewCertificationRSAWithPub")
		return
	}

	public, _, _, _, suberr := ssh.ParseAuthorizedKey(pubAuthKey)
	if suberr != nil {
		err = suberr
		return
	}

	if public.Type() != "ssh-rsa" {
		err = mssh.InvalidParamsError.New("NewCertificationRSAWithPub")
		return
	}

	// generate certification
	sshcert := &ssh.Certificate{
		Key:             public, // public
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

	cert, err = createCertificationRSA(sshcert, caPem)
	return
}

func createCertificationRSA(cert *ssh.Certificate, caPem []byte) (string, error) {
	if cert == nil || len(caPem) == 0 {
		return "", mssh.EmptyError.New("createCertificationRSA")
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
	return strings.Join([]string{ssh.CertAlgoRSAv01, baseBody, cert.KeyId}, " "), nil
}
