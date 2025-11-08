#!/bin/bash

mkdir -p .certs

# Create CA
openssl genrsa -out .certs/ca.key 2048
openssl req -x509 -new -nodes -key .certs/ca.key -subj "/CN=GitHub App" -days 365 -out .certs/ca.crt

# Server cert
openssl genrsa -out .certs/server.key 2048
openssl req -new -key .certs/server.key -subj "/CN=Service" -out .certs/server.csr
openssl x509 -req -in .certs/server.csr -CA .certs/ca.crt -CAkey .certs/ca.key -CAcreateserial -out .certs/server.crt -days 365

# Client cert
openssl genrsa -out .certs/client.key 2048
openssl req -new -key .certs/client.key -subj "/CN=Client" -out .certs/client.csr
openssl x509 -req -in .certs/client.csr -CA .certs/ca.crt -CAkey .certs/ca.key -CAcreateserial -out .certs/client.crt -days 365

rm .certs/*.csr .certs/*.srl
echo "Certificates generated in .certs/ directory"