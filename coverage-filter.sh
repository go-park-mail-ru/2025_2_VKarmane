#!/usr/bin/env bash
set -euo pipefail

PKGS=$(go list ./... \
  | grep -vE '(^|/)(mocks|cmd|docs|proto)(/|$)' \
  | grep -vE 'internal/app/[^/]+$' \
  | grep -vE 'internal/app/.+/handlers' \
  | grep -vE 'internal/app/image/repository' \
  | grep -vE 'internal/metrics$' \
  | grep -vE 'internal/models$' \
  | grep -vE 'scripts$')

GOFLAGS= go test -covermode=atomic -coverprofile=coverage.out $PKGS

go tool cover -func=coverage.out | grep total: