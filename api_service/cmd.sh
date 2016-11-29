#!/bin/bash
set -e

if [ "$ENV" = 'DEV' ]; then
  echo "Krot, API Service (Dev)"
else
  echo "Krot, API Service (Prod)"
fi

exec go run "server.go"
