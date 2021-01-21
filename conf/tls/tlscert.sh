#!/bin/bash
# This script is executed by Dockerfile.openssl

if [[ -z $DEPOSITS_GRPC_SERVICE_ADDRESS ]]; then
  echo DEPOSITS_GRPC_SERVICE_ADDRESS not provided so please export DEPOSITS_GRPC_SERVICE_ADDRESS
  exit 1
else
  echo EPOSITS_GRPC_SERVICE_ADDRESS=$DEPOSITS_GRPC_SERVICE_ADDRESS is provided
fi

HOST=$DEPOSITS_GRPC_SERVICE_ADDRESS

# out ca.key
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096

#openssl req -passin pass:1111 -new -x509 -days 3650 -key ca.key -out ca.crt -subj "/CN=${CN}"
# penssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout ${KEY_FILE} -out ${CERT_FILE} -subj "/CN=${HOST}/O=${HOST}"
openssl req -passin pass:1111 -new -x509 -days 3650 -key ca.key -out cacrtto.pem -subj "/CN=${HOST}/O=${HOST}"

# out server.pem format from server.key -- server.pem used by server, hence shared
openssl genrsa -passout pass:1111 -des3 -out server.key 4096
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out serverkeyto.pem

# out server.csr from server.key
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${HOST}/O=${HOST}"

# out server.crt from server.csr and ca.key and ca.crt -- server.crt used by server, hence shared
#openssl x509 -req -passin pass:1111 -days 3650 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt
openssl x509 -req -passin pass:1111 -days 3650 -in server.csr -CA cacrtto.pem -CAkey ca.key -set_serial 01 -out servercrtto.pem
