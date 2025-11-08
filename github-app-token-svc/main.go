package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"net/http"

	"github.com/bradleyfalzon/ghinstallation/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

var githubToken string

func main() {
	env, err := loadEnv()
	if err != nil {
		logrus.Fatalf("Error loading environment variables: %v", err)
	}

	keyString := env.Plugin.GithubPrivateKeyPEM

	logrus.Infof("Getting Github Token for repo: %s", env.CI.Repo)
	logrus.Infof("Github App ID: %d", env.Plugin.GithubAppID)
	logrus.Infof("Github Installation ID: %d", env.Plugin.GithubInstallationID)

	// Parse the private key to an *rsa.PrivateKey
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(keyString))
	if err != nil {
		logrus.Fatalf("Error parsing private key: %v", err)
	}

	atr := ghinstallation.NewAppsTransportFromPrivateKey(http.DefaultTransport, env.Plugin.GithubAppID, privateKey)
	itr := ghinstallation.NewFromAppsTransport(atr, env.Plugin.GithubInstallationID)

	token, err := itr.Token(context.TODO())
	if err != nil {
		logrus.Fatalf("Error getting installation token: %v", err)
	}

	githubToken = token
	logrus.Infof("Github Token retrieved successfully")

	serverCert := []byte(env.Plugin.MtlsServerCert)
	serverKey := []byte(env.Plugin.MtlsServerKey)
	caCert := []byte(env.Plugin.MtlsCACert)

	cert, err := tls.X509KeyPair(serverCert, serverKey)
	if err != nil {
		logrus.Fatalf("Error loading server certificate: %v", err)
	}

	caPool := x509.NewCertPool()
	if !caPool.AppendCertsFromPEM(caCert) {
		logrus.Fatalf("Error loading CA cert")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(githubToken))
	})

	srv := &http.Server{
		Addr:      ":8443",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	logrus.Infof("mTLS server running on :8443")
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		logrus.Fatalf("mTLS Server error: %v", err)
	}
}
