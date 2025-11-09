#!/bin/bash

# Check that woodpecker CLI is installed
if ! command -v woodpecker &> /dev/null; then
  echo "woodpecker CLI could not be found. Please install it from https://github.com/woodpecker-ci/woodpecker-cli"
  exit 1
fi

# Ask where to find the certificate files
read -p "Enter the path to the .certs directory (default: .certs): " CERTS_DIR
CERTS_DIR=${CERTS_DIR:-.certs}

# Check if the directory exists
if [ ! -d "$CERTS_DIR" ]; then
  echo "Directory $CERTS_DIR does not exist."
  exit 1
fi

# Check if required files exist
for file in ca.crt client.crt client.key server.crt server.key; do
  if [ ! -f "$CERTS_DIR/$file" ]; then
    echo "File $CERTS_DIR/$file does not exist."
    exit 1
  fi
done

# Ask for the Woodpecker Instance URL if needed
if [ -z "$WOODPECKER_SERVER" ]; then
read -p "Enter the Woodpecker Instance URL (e.g., https://ci.example.com): " WOODPECKER_SERVER
fi

# Ask for the token if not set
if [ -z "$WOODPECKER_TOKEN" ]; then
  read -p "Enter your Woodpecker token: " WOODPECKER_TOKEN
fi

# Ask which repository to install the certificate to
read -p "Enter the GitHub repository (owner/repo) to install the certificate to: " REPO

# Ask for the env variable name with defaults
DEFAULT_CA_ENV_VAR="MTLS_CA_CERT"
DEFAULT_SERVER_CERT_ENV_VAR="MTLS_SERVER_CERT"
DEFAULT_SERVER_KEY_ENV_VAR="MTLS_SERVER_KEY"
DEFAULT_CLIENT_CERT_ENV_VAR="MTLS_CLIENT_CERT"
DEFAULT_CLIENT_KEY_ENV_VAR="MTLS_CLIENT_KEY"

read -p "Enter the environment variable name to store the CA Cert (default: MTLS_CA_CERT): " CA_ENV_VAR
CA_ENV_VAR=${CA_ENV_VAR:-MTLS_CA_CERT}
read -p "Enter the environment variable name to store the Server Cert (default: MTLS_SERVER_CERT): " SERVER_CERT_ENV_VAR
SERVER_CERT_ENV_VAR=${SERVER_CERT_ENV_VAR:-MTLS_SERVER_CERT}
read -p "Enter the environment variable name to store the Server Key (default: MTLS_SERVER_KEY): " SERVER_KEY_ENV_VAR
SERVER_KEY_ENV_VAR=${SERVER_KEY_ENV_VAR:-MTLS_SERVER_KEY}
read -p "Enter the environment variable name to store the Client Cert (default: MTLS_CLIENT_CERT): " CLIENT_CERT_ENV_VAR
CLIENT_CERT_ENV_VAR=${CLIENT_CERT_ENV_VAR:-MTLS_CLIENT_CERT}
read -p "Enter the environment variable name to store the Client Key (default: MTLS_CLIENT_KEY): " CLIENT_KEY_ENV_VAR
CLIENT_KEY_ENV_VAR=${CLIENT_KEY_ENV_VAR:-MTLS_CLIENT_KEY} 

# Install the certificates as repository secrets
woodpecker repo secret add --server "$WOODPECKER_SERVER" --token "$WOODPECKER_TOKEN" --repo "$REPO" --name "$CA_ENV_VAR" --value "$(cat "$CERTS_DIR/ca.crt")"
woodpecker repo secret add --server "$WOODPECKER_SERVER" --token "$WOODPECKER_TOKEN" --repo "$REPO" --name "$SERVER_CERT_ENV_VAR" --value "$(cat "$CERTS_DIR/server.crt")" --image="ghcr.io/yyewolf/woodpecker-plugins/github-app-token-svc"
woodpecker repo secret add --server "$WOODPECKER_SERVER" --token "$WOODPECKER_TOKEN" --repo "$REPO" --name "$SERVER_KEY_ENV_VAR" --value "$(cat "$CERTS_DIR/server.key")" --image="ghcr.io/yyewolf/woodpecker-plugins/github-app-token-svc"
woodpecker repo secret add --server "$WOODPECKER_SERVER" --token "$WOODPECKER_TOKEN" --repo "$REPO" --name "$CLIENT_CERT_ENV_VAR" --value "$(cat "$CERTS_DIR/client.crt")"
woodpecker repo secret add --server "$WOODPECKER_SERVER" --token "$WOODPECKER_TOKEN" --repo "$REPO" --name "$CLIENT_KEY_ENV_VAR" --value "$(cat "$CERTS_DIR/client.key")"

echo "Certificates installed as secrets in repository $REPO."       