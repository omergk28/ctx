//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// DangerousSudo matches sudo invocations.
// Anchored at start-of-string or after a shell separator so that
// substrings like "pseudo-foo" or "echo 'sudo'" don't trip it.
var DangerousSudo = regexp.MustCompile(`(^|[\s;&|])sudo\s`)

// DangerousRmRfRoot matches rm -rf / (root deletion).
// `/` must be followed by whitespace or end-of-string so that
// `rm -rf /var/log` (a normal operation) does not match.
var DangerousRmRfRoot = regexp.MustCompile(
	`\brm\s+-[rR][fF]?\s+/(\s|$)|\brm\s+-[fF][rR]\s+/(\s|$)`)

// DangerousRmRfHome matches rm -rf ~ (home deletion or any subpath).
// Mirrors the copilot-cli substring guard `*"rm -rf ~"*`: blocks
// both `rm -rf ~` and `rm -rf ~/Downloads`.
var DangerousRmRfHome = regexp.MustCompile(
	`\brm\s+-[rR][fF]?\s+~|\brm\s+-[fF][rR]\s+~`)

// DangerousChmod777 matches chmod 777 (overly permissive permissions).
// Allows any flag arguments before the mode (e.g. `chmod -R 777`).
var DangerousChmod777 = regexp.MustCompile(`\bchmod\s+(-\S+\s+)*777\b`)

// DangerousGitPushForce matches git push --force or git push -f.
// The trailing `(\s|$)` anchor excludes `--force-with-lease`, which
// is the safer alternative we deliberately allow.
var DangerousGitPushForce = regexp.MustCompile(
	`\bgit\s+push\b[^\n]*?\s(--force|-f)(\s|$)`)

// DangerousGitResetHard matches git reset --hard.
var DangerousGitResetHard = regexp.MustCompile(`\bgit\s+reset\s+--hard\b`)

// DangerousRemoveItemRoot matches PowerShell `Remove-Item -Recurse
// -Force C:\` (Windows root deletion). The pattern matches the
// canonical flag order; reordered flags fall through to a benign
// allow.
var DangerousRemoveItemRoot = regexp.MustCompile(
	`Remove-Item\s+-Recurse\s+-Force\s+C:\\`)

// DangerousRemoveItemHome matches PowerShell `Remove-Item -Recurse
// -Force $env:USERPROFILE` (Windows home deletion).
var DangerousRemoveItemHome = regexp.MustCompile(
	`Remove-Item\s+-Recurse\s+-Force\s+\$env:USERPROFILE`)

// DangerousFormatVolume matches PowerShell Format-Volume (drive
// reformat). Word-anchored so substrings in unrelated identifiers
// don't trip it.
var DangerousFormatVolume = regexp.MustCompile(`\bFormat-Volume\b`)
