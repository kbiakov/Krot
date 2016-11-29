#!/bin/bash
set -e

if [ "$ENV" = 'DEV' ]; then
  echo "Krot, Notification Service (Dev)"
else
  echo "Krot, Notification Service (Prod)"
fi

exec go run "server.go"
