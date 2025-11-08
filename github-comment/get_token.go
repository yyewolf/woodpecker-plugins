package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetToken retrieves the token from an mTLS-protected server
func GetToken(serverURL, clientCert, clientKey, caCert string) (string, error) {
	// Load client cert + key
	clientKeyPair, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
	if err != nil {
		return "", fmt.Errorf("failed to load client certificate/key: %w", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM([]byte(caCert)) {
		return "", fmt.Errorf("failed to append CA certificate")
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientKeyPair},
		RootCAs:            caPool,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	client := &http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	}

	var resp *http.Response
	for i := range 5 {
		resp, err = client.Get(serverURL)
		if err != nil {
			time.Sleep(time.Duration(200*i) * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		return "", fmt.Errorf("request failed after retries: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("server returned %s: %s", resp.Status, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	return string(body), nil
}
