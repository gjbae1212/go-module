package ssh_cert_authority

import "time"

type CertMeta struct {
	KeyId      string
	Principals []string
	Before     time.Time
	After      time.Time
}
