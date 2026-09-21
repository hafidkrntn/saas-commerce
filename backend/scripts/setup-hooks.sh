#!/bin/sh
# Alternative: manually set git hooks path.
# This runs automatically via `make run/build/test`, so you usually don't need this.
git config core.hooksPath scripts/hooks
chmod +x scripts/hooks/*
echo "✅ Git hooks configured!"
