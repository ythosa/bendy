#!/usr/bin/env bash

# Check code with golangci-lint
echo "==> Checking that the code satisfies the linter..."
lint_files=$("${GOPATH}"/bin/golangci-lint run)
if [[ -n ${lint_files} ]]; then
  echo "there are some linter errors:"
  echo "${lint_files}"
  exit 1
fi

exit 0
