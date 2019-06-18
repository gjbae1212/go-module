package ip

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

var (
	// https://www.ipify.org
	ipifyApiV4 = "https://api.ipify.org"
	ipifyApiV6 = "https://api6.ipify.org"
)

func GetPublicIPV4() (string, error) {
	resp, err := http.Get(ipifyApiV4)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(string(data))
	if ip == nil {
		return "", fmt.Errorf("[err] not found ipv4")
	}
	return ip.String(), nil
}

func GetPublicIPV6() (string, error) {
	resp, err := http.Get(ipifyApiV6)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := net.ParseIP(string(data))
	if ip == nil {
		return "", fmt.Errorf("[err] not found ipv6")
	}
	return ip.String(), nil
}
