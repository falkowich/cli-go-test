#!/bin/bash

# Usage: ./scripts/bump_version.sh [major|minor|patch]

if [ -z "$1" ]; then
  echo "Usage: $0 [major|minor|patch]"
  exit 1
fi

CURRENT_VERSION=$(git describe --tags --abbrev=0)
IFS='.' read -r -a VERSION_PARTS <<< "$CURRENT_VERSION"

case "$1" in
  major)
    VERSION_PARTS[0]=$((VERSION_PARTS[0] + 1))
    VERSION_PARTS[1]=0
    VERSION_PARTS[2]=0
    ;;
  minor)
    VERSION_PARTS[1]=$((VERSION_PARTS[1] + 1))
    VERSION_PARTS[2]=0
    ;;
  patch)
    VERSION_PARTS[2]=$((VERSION_PARTS[2] + 1))
    ;;
  *)
    echo "Unknown version type: $1"
    exit 1
    ;;
esac

NEW_VERSION="${VERSION_PARTS[0]}.${VERSION_PARTS[1]}.${VERSION_PARTS[2]}"
git tag -a "$NEW_VERSION" -m "Release $NEW_VERSION"
git push origin "$NEW_VERSION"

echo "Bumped version to $NEW_VERSION"

