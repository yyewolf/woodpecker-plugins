#!/bin/bash

# Ask for a name
read -p "Enter the name for the certificate set (default: github-app): " CERT_NAME
CERT_NAME=${CERT_NAME:-github-app}

mkdir -p .certs/$CERT_NAME

# Create CA
openssl genrsa -out .certs/$CERT_NAME/ca.key 2048
openssl req -x509 -new -nodes -key .certs/$CERT_NAME/ca.key -subj "/CN=GitHub App" -days 365 -out .certs/$CERT_NAME/ca.crt
# Server cert
openssl genrsa -out .certs/$CERT_NAME/server.key 2048
openssl req -new -key .certs/$CERT_NAME/server.key -subj "/CN=Service" -out .certs/$CERT_NAME/server.csr
openssl x509 -req -in .certs/$CERT_NAME/server.csr -CA .certs/$CERT_NAME/ca.crt -CAkey .certs/$CERT_NAME/ca.key -CAcreateserial -out .certs/$CERT_NAME/server.crt -days 365
# Client cert
openssl genrsa -out .certs/$CERT_NAME/client.key 2048
openssl req -new -key .certs/$CERT_NAME/client.key -subj "/CN=Client" -out .certs/$CERT_NAME/client.csr
openssl x509 -req -in .certs/$CERT_NAME/client.csr -CA .certs/$CERT_NAME/ca.crt -CAkey .certs/$CERT_NAME/ca.key -CAcreateserial -out .certs/$CERT_NAME/client.crt -days 365

rm .certs/$CERT_NAME/*.csr .certs/$CERT_NAME/*.srl
echo "Certificates generated in .certs/$CERT_NAME/ directory"
