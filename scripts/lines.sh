#!/usr/bin/env bash

set -euo pipefail

GOSRC=$(find ./ -name "*.go" -print0 |xargs -0 grep -v '^$' | wc -l)

TSSRC=$(find ./ui -name "*.ts" -not -path './ui/node_modules/*' -print0 |xargs -0 grep -v '^$' | wc -l)
TSXSRC=$(find ./ui -name "*.tsx" -not -path './ui/node_modules/*' -print0 |xargs -0 grep -v '^$' | wc -l)
UISRC=$(expr $TSSRC + $TSXSRC)

echo "Go:    ${GOSRC}"
echo "UI:    ${UISRC}"

# shellcheck disable=SC2003
echo "total: $(expr "${GOSRC}" + "${UISRC}")"
