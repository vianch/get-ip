#!/usr/bin/env bash
# Tag a new release and print the formula update snippet.
#
# Usage: scripts/release.sh v0.1.0
#
# What it does:
#   1. Verifies the working tree is clean and on main.
#   2. Checks the tag does not already exist.
#   3. Creates an annotated git tag.
#   4. Pushes the tag to origin (which must point at github.com/vianch/get-ip).
#   5. Downloads the GitHub-generated source tarball, computes its SHA256,
#      and prints the exact lines to paste into Formula/get-ip.rb.
#
# After running this, edit Formula/get-ip.rb with the printed url + sha256,
# commit the formula change, and push.

set -euo pipefail

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
  echo "usage: $0 <version>   e.g. $0 v0.1.0" >&2
  exit 2
fi

if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "error: version must look like vMAJOR.MINOR.PATCH (got: $VERSION)" >&2
  exit 2
fi

BRANCH="$(git rev-parse --abbrev-ref HEAD)"
if [[ "$BRANCH" != "main" ]]; then
  echo "error: must be on main, currently on $BRANCH" >&2
  exit 1
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  echo "error: working tree has uncommitted changes" >&2
  exit 1
fi

if git rev-parse "$VERSION" >/dev/null 2>&1; then
  echo "error: tag $VERSION already exists" >&2
  exit 1
fi

ORIGIN_URL="$(git remote get-url origin 2>/dev/null || true)"
if [[ "$ORIGIN_URL" != *"vianch/get-ip"* ]]; then
  echo "error: origin must be a vianch/get-ip remote (got: $ORIGIN_URL)" >&2
  exit 1
fi

echo "==> tagging $VERSION"
git tag -a "$VERSION" -m "$VERSION"
git push origin "$VERSION"

TARBALL_URL="https://github.com/vianch/get-ip/archive/refs/tags/${VERSION}.tar.gz"
echo "==> waiting for GitHub to publish the tarball at $TARBALL_URL"
sleep 5

TMPFILE="$(mktemp -t get-ip-release.XXXXXX)"
trap 'rm -f "$TMPFILE"' EXIT

for i in 1 2 3 4 5; do
  if curl -fsSL "$TARBALL_URL" -o "$TMPFILE"; then
    break
  fi
  echo "    tarball not ready yet, retrying ($i/5)..."
  sleep 5
done

if ! [[ -s "$TMPFILE" ]]; then
  echo "error: failed to download $TARBALL_URL" >&2
  exit 1
fi

SHA="$(shasum -a 256 "$TMPFILE" | awk '{print $1}')"

cat <<EOF

==> Update Formula/get-ip.rb with these values:

  url "$TARBALL_URL"
  sha256 "$SHA"

Then:
  git add Formula/get-ip.rb
  git commit -m "chore: bump formula to $VERSION"
  git push

Users can then install with:
  brew tap vianch/get-ip https://github.com/vianch/get-ip
  brew install vianch/get-ip/get-ip

EOF
