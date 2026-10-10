#!/usr/bin/env bash
set -euo pipefail

# gitlings-ko 빌드 스크립트
# - mac: arm64 + amd64 유니버설 바이너리 (lipo)
# - Windows: amd64 exe

BINARY=gitlings
PKG=./cmd/gitlings

echo "=== gitlings-ko 빌드 시작 ==="

# 1. mac arm64
echo "[1/4] mac arm64 빌드..."
GOOS=darwin GOARCH=arm64 go build -o "${BINARY}-darwin-arm64" "${PKG}"

# 2. mac amd64
echo "[2/4] mac amd64 빌드..."
GOOS=darwin GOARCH=amd64 go build -o "${BINARY}-darwin-amd64" "${PKG}"

# 3. mac 유니버설 바이너리
echo "[3/4] mac 유니버설 바이너리 합치기 (lipo)..."
lipo -create -output "${BINARY}" "${BINARY}-darwin-arm64" "${BINARY}-darwin-amd64"
chmod +x "${BINARY}"
rm "${BINARY}-darwin-arm64" "${BINARY}-darwin-amd64"

# 4. Windows amd64
echo "[4/4] Windows amd64 빌드..."
GOOS=windows GOARCH=amd64 go build -o "${BINARY}.exe" "${PKG}"

echo ""
echo "=== 빌드 완료 ==="
ls -lh "${BINARY}" "${BINARY}.exe"
echo ""
echo "바이너리 확인:"
file "${BINARY}"
