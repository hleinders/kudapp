#!/usr/bin/bash

openssl req -new -x509 -sha256 -key cert.key -out cert.pem -days 3650 -config input.req
