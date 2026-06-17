#!/bin/bash
set -e

BIN_DIR="/usr/local/bin"
CONFIG_DIR="/etc/sean"
DATA_DIR="/var/lib/sean"

echo "[sean] Installing binary..."
cp sean "$BIN_DIR/sean"
chmod +x "$BIN_DIR/sean"

echo "[sean] Installing configs..."
mkdir -p "$CONFIG_DIR/manifests"
cp configs/installed.yaml "$CONFIG_DIR/installed.yaml"
cp configs/manifests/*.yaml "$CONFIG_DIR/manifests/"

echo "[sean] Creating data directory..."
mkdir -p "$DATA_DIR/tools"

echo "[sean] Done. Run: sean --help"