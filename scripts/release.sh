#!/usr/bin/env bash

set -Eeu

CHART_YAML=./deploy/chart/k8s-log-proxy/Chart.yaml

OLD_VERSION=$(sed -nr 's/version: (.+)/\1/p' "$CHART_YAML" | tr -d '[:space:]')  # need to remove trailing whitespace
NEW_VERSION=${1:-}

if [ -z "$NEW_VERSION" ]; then
    echo "Usage: $(basename "$0") <NEW_VERSION>"
    exit 1
fi

PATTERN=$(printf 's/appVersion: "v%s"/appVersion: "v%s"/g' "$OLD_VERSION" "$NEW_VERSION")
sed -ri "$PATTERN" "$CHART_YAML"
PATTERN=$(printf 's/version: %s/version: %s/g' "$OLD_VERSION" "$NEW_VERSION")
sed -ri "$PATTERN" "$CHART_YAML"

git add "$CHART_YAML"
git commit -m"chore: Release v$NEW_VERSION"
git tag "v$NEW_VERSION"
git push --tags origin :