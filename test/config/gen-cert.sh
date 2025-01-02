#!/usr/bin/env bash
set -ex

openssl req -newkey rsa:2048 -noenc -keyout domain.key -out domain.csr
openssl req -x509 -sha256 -days 1825 -newkey rsa:2048 -keyout rootCA.key -out rootCA.crt -noenc
openssl x509 -req -CA rootCA.crt -CAkey rootCA.key -in domain.csr -out domain.crt -days 365 -CAcreateserial -extfile domain.ext
cat domain.crt rootCA.crt > domain-bundle.crt
