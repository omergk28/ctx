#!/usr/bin/env bash
# ctx preToolUse hook for GitHub Copilot CLI
# Reshapes the Copilot envelope into the ctx hook envelope and
# delegates the dangerous-command decision to
# `ctx system block-dangerous-commands` (single source of truth
# shared with Claude Code and OpenCode integrations).
set -euo pipefail

INPUT=$(cat)

# Without jq we can't reshape envelopes safely. Fail open.
if ! command -v jq >/dev/null 2>&1; then
  exit 0
fi

TOOL=$(echo "$INPUT" | jq -r '.tool_name // .tool // empty' 2>/dev/null || true)
case "$TOOL" in
  shell|bash) ;;
  *) exit 0 ;;
esac

COMMAND=$(echo "$INPUT" | jq -r '.input.command // empty' 2>/dev/null || true)
[ -n "$COMMAND" ] || exit 0

# Reshape into the Claude-style envelope the Go hook expects.
ENVELOPE=$(jq -nc --arg cmd "$COMMAND" \
  '{session_id: "copilot-cli", tool_input: {command: $cmd}}')

# Run the Go hook. Missing binary or non-zero exit → fail open;
# re-installing ctx restores the safety net.
DECISION=$(echo "$ENVELOPE" | ctx system block-dangerous-commands 2>/dev/null || true)
[ -n "$DECISION" ] || exit 0

# If the hook decided to block, surface its JSON to copilot-cli
# (decision/reason format) and exit 1.
if echo "$DECISION" | jq -e '.decision == "block"' >/dev/null 2>&1; then
  echo "$DECISION" >&2
  exit 1
fi
