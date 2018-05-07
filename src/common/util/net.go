package util

import (
	"net/url"
)

// GetDailAddress returns the address for net.Dail
func GetDailAddress(URL string) (string, error) {
	uri, err := url.Parse(URL)
	if err != nil {
		return "", err
	}
	var port = uri.Port()
	if uri.Port() == "" {
		port = "80"
	}
	return uri.Hostname() + ":" + port, err
}
