#!/usr/bin/env bash

echo "ENCRYPTION_KEY=\"$(openssl rand -base64 32)\"" >> app.env
