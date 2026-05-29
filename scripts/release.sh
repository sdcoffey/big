#!/bin/bash

set -euf -o pipefail

echo -n "Version: "
read -r newversion

tag_name="v$newversion"

if git rev-parse -q --verify "refs/tags/$tag_name" >/dev/null; then
  echo "Tag: $tag_name already exists"
  exit 1
else
  echo "Releasing $newversion"
fi

echo "Update CHANGELOG.md and press enter"
read -r

git add CHANGELOG.md

added_count=$(git status --porcelain | grep "CHANGELOG.md" | wc -l | tr -d '[:space:]' || true)
if [[ $added_count -gt 0 ]]; then
  git commit -m"Release version $newversion"
fi

git tag "$tag_name" -m "$tag_name"

git push origin main
git push --tags
