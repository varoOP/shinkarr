package config

import "net/http"

type ApiKeyTransport struct {
	Transport http.RoundTripper
	ApiKey    string
}

func (c *ApiKeyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if c.Transport == nil {
		c.Transport = http.DefaultTransport
	}

	req.Header.Add("X-Api-Key", c.ApiKey)
	return c.Transport.RoundTrip(req)
}
