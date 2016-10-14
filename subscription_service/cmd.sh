#!/bin/bash
set -e

if [ "$ENV" = 'DEV' ]; then
  echo "Krot, Subscription Service (Dev)"
else
  echo "Krot, Subscription Service (Prod)"
fi

exec go run "server.go"
