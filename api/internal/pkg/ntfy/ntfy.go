package ntfy

import (
	"encoding/base64"
	"errors"
	"fmt"
)

// ntfySettings represents custom options to create NtfyClient.
type ntfySettings struct {
	protocol string
	host     string
	port     string
	token    string
	username string
	password string
}

// Option represents an option for Client initializing.
type Option func(*ntfySettings)

// Client represents a client for ntfy server.
type Client struct {
	addr       string
	authHeader string
}

// NewClient returns a new instance Client.
func NewClient(options ...Option) (*Client, error) {
	settings := &ntfySettings{
		protocol: "https",
		host:     "ntfy.sh",
		port:     "80",
	}
	for _, option := range options {
		option(settings)
	}

	client := &Client{
		addr: fmt.Sprintf("%s://%s:%s/", settings.protocol, settings.host, settings.port),
	}

	switch {
	case settings.token != "":
		client.authHeader = fmt.Sprintf("Bearer %s", settings.token)

	case settings.username != "" && settings.password != "":
		credentials := []byte(settings.username + ":" + settings.password)
		encoded := base64.StdEncoding.EncodeToString(credentials)
		client.authHeader = fmt.Sprintf("Basic %s", encoded)

	default:
		return nil, errors.New("at least one auth method must be presented")
	}

	return client, nil
}

// WithHTTPS sets secure protocol for ntfy client.
func WithHTTPS() Option {
	return func(s *ntfySettings) {
		s.protocol = "https"
	}
}

// WithHTTP sets insecure protocol for ntfy client.
func WithHTTP() Option {
	return func(s *ntfySettings) {
		s.protocol = "http"
	}
}

// WithHost sets custom host for ntfy client.
func WithHost(host string) Option {
	return func(s *ntfySettings) {
		s.host = host
	}
}

// WithPort sets custom port for ntfy client.
func WithPort(port string) Option {
	return func(s *ntfySettings) {
		s.port = port
	}
}

// WithTokenAuth sets auth token for ntfy client.
func WithTokenAuth(token string) Option {
	return func(s *ntfySettings) {
		s.token = token
	}
}

// WithBasicAuth sets username and password for ntfy client.
func WithBasicAuth(username, password string) Option {
	return func(s *ntfySettings) {
		s.username = username
		s.password = password
	}
}
