//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package block_dangerous_commands implements the
// **`ctx system block-dangerous-commands`** hidden hook,
// which blocks shell commands matching a curated set of
// destructive or irreversible patterns.
//
// # What It Does
//
// The hook reads a JSON envelope from stdin and checks the
// command string against the dangerous-pattern set:
//
//   - **sudo**: privilege escalation
//   - **rm -rf /**: root deletion
//   - **rm -rf ~**: home deletion (and any subpath of ~)
//   - **chmod 777**: overly permissive permissions
//   - **git push --force / -f**: irreversible history rewrite
//     (--force-with-lease is allowed)
//   - **git reset --hard**: discards local changes
//
// When a match is found, a JSON block response is emitted to
// prevent execution and a relay notification is sent
// explaining why the invocation was rejected.
//
// # Input
//
// A JSON hook envelope on stdin with a ToolInput.Command
// field containing the shell command string.
//
// # Output
//
// On match: a JSON [entity.BlockResponse] with decision
// "block" and a human-readable reason. On no match: no
// output (silent allow).
//
// # Consumers
//
// The same Go binary serves multiple editor integrations:
//
//   - **Claude Code** PreToolUse hook (Bash matcher)
//   - **OpenCode** plugin tool.execute.before (parses stdout)
//   - **GitHub Copilot CLI** preToolUse script (delegates here)
//
// All three pipe the same JSON envelope shape; the plugin
// shim around each editor adapts the editor-specific block
// signal (throw, exit code, decision JSON).
package block_dangerous_commands
