package ip

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicIPV4(t *testing.T) {
	assert := assert.New(t)
	ip, err := GetPublicIPV4()
	assert.NoError(err)
	log.Println(ip)
}

func TestGetPublicIPV6(t *testing.T) {
	assert := assert.New(t)
	ip, err := GetPublicIPV6()
	assert.NoError(err)
	log.Println(ip)
}
