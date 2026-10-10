#!/usr/bin/env bash
set -euo pipefail

# Snap's FFmpeg cannot see /tmp created by VHS. Keep VHS's temporary frames
# under this repository instead, where the snap sandbox has access.
repo_root=$(cd "$(dirname "$0")/.." && pwd)
tmp_dir="$repo_root/showcase/.tmp"

mkdir -p "$tmp_dir" "$repo_root/showcase/assets"
for tape in "$repo_root"/showcase/*.tape "$repo_root"/showcase/silly/*.tape; do
  TMPDIR="$tmp_dir" vhs "$tape"
done
