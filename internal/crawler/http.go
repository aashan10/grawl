package crawler

import (
	"crypto/tls"
	"net/http"
)

func GetClient(skipCertCheck bool) *http.Client {
	client := &http.Client{}

	if skipCertCheck {
		// Set custom transport to skip SSL certificate verification
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return client
}
