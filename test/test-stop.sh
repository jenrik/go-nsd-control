#!/bin/sh

echo "NSDCT1 stop" | openssl s_client -cert ./config/domain.crt -key ./config/domain.key -cert_chain ./config/rootCA.crt -CAfile ./config/rootCA.crt -connect 127.0.0.1:8952 -servername localhost
