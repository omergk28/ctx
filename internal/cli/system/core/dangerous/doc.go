//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

// Package dangerous hosts the pure-logic helpers for the
// block-dangerous-commands hook: the [Match] result type
// and the [Detect] function that maps a shell command to
// a dangerous-pattern variant.
//
// The cmd-side package
// (internal/cli/system/cmd/block_dangerous_commands) owns
// Cobra wiring, stdin parsing, and JSON output. This core
// package owns the rule set and is the unit-test target
// for new patterns.
package dangerous
