#!/usr/bin/env zsh

# Build app with go build
echo "==> Checking that the code is building..."
lint_files=$(go build -o cli ./cmd/cli/main.go)
if [[ -n ${lint_files} ]]; then
  echo "there are some building errors:"
  echo "${lint_files}"
  exit 1
fi

rm ./cli
exit 0
