#!/bin/sh

ENDPOINT=${API_BASE}/next/$1

>>/tmp/liquidsoap.log echo curl ${ENDPOINT}

sleep 1

curl $ENDPOINT
