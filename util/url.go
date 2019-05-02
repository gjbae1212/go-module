package util

import (
	"fmt"

	"github.com/goware/urlx"
)

func ParseURL(s string) (schema, host, port, path, query, fragment string, err error) {
	if s == "" {
		err = fmt.Errorf("[err] ParseURI empty uri")
	}

	url, suberr := urlx.Parse(s)
	if suberr != nil {
		err = suberr
		return
	}

	schema = url.Scheme

	host, port, err = urlx.SplitHostPort(url)
	if err != nil {
		return
	}
	if schema == "http" && port == "" {
		port = "80"
	} else if schema == "https" && port == "" {
		port = "443"
	}

	path = url.Path
	query = url.RawQuery
	fragment = url.Fragment
	return
}
