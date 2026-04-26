//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package regex

import "regexp"

// cmdStart is the leading anchor for "this is the start of a real
// command" — either the beginning of input or immediately after a
// shell separator (`;`, `&&`, `||`, `|`). Plain whitespace is NOT
// in the alternation: it would let `echo sudo …` and similar print
// statements trip the guard. Trailing `\s*` swallows any spaces
// between the separator and the dangerous verb.
const cmdStart = `(?:^|;|&&|\|\||\|)\s*`

// DangerousSudo matches sudo invocations at the start of a real
// command. `pseudo-foo`, `echo sudo`, and `"sudo"` (quoted strings)
// fall through to a benign allow.
var DangerousSudo = regexp.MustCompile(cmdStart + `sudo\s`)

// DangerousRmRfRoot matches rm -rf / (root deletion).
// `/` must be followed by whitespace or end-of-string so that
// `rm -rf /var/log` (a normal operation) does not match.
var DangerousRmRfRoot = regexp.MustCompile(
	cmdStart + `rm\s+-[rR][fF]?\s+/(\s|$)` +
		`|` + cmdStart + `rm\s+-[fF][rR]\s+/(\s|$)`)

// DangerousRmRfHome matches rm -rf ~ (home deletion or any subpath).
// Mirrors the copilot-cli substring guard `*"rm -rf ~"*`: blocks
// both `rm -rf ~` and `rm -rf ~/Downloads`.
var DangerousRmRfHome = regexp.MustCompile(
	cmdStart + `rm\s+-[rR][fF]?\s+~` +
		`|` + cmdStart + `rm\s+-[fF][rR]\s+~`)

// DangerousChmod777 matches chmod 777 (overly permissive
// permissions). Allows any flag arguments between `chmod` and the
// mode (e.g. `chmod -R 777`).
var DangerousChmod777 = regexp.MustCompile(
	cmdStart + `chmod\s+(-\S+\s+)*777\b`)

// DangerousGitPushForce matches git push --force or git push -f
// (including combined short flags like `-fu`). The `--force` arm
// requires trailing whitespace or end-of-string to exclude
// `--force-with-lease`, the safer alternative we allow. The
// short-flag arm matches any single-dash flag bundle containing
// `f` (e.g. `-f`, `-fu`, `-uf`).
var DangerousGitPushForce = regexp.MustCompile(
	cmdStart + `git\s+push\b[^\n]*?\s` +
		`(--force(\s|$)|-[a-zA-Z]*f[a-zA-Z]*(\s|$))`)

// DangerousGitResetHard matches git reset --hard.
var DangerousGitResetHard = regexp.MustCompile(
	cmdStart + `git\s+reset\s+--hard\b`)

// removeItemAlt is the alternation set for "Remove-Item with
// -Recurse, -Force, and the target path, in any order". RE2
// doesn't support lookahead, so we enumerate the six permutations
// of three tokens. Each arm uses [^\n]*? so the tokens can be
// separated by other arguments without escaping the line.
//
// Parameters:
//   - target: regex-escaped path literal to require alongside
//     -Recurse and -Force
//
// Returns:
//   - string: regex source with six permutations joined by '|',
//     suitable for [regexp.MustCompile]
func removeItemAlt(target string) string {
	t := target
	parts := [][3]string{
		{"-Recurse", "-Force", t},
		{"-Recurse", t, "-Force"},
		{"-Force", "-Recurse", t},
		{"-Force", t, "-Recurse"},
		{t, "-Recurse", "-Force"},
		{t, "-Force", "-Recurse"},
	}
	const head = `Remove-Item\b[^\n]*?`
	const sep = `[^\n]*?`
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += "|"
		}
		out += cmdStart + head + p[0] + sep + p[1] + sep + p[2]
	}
	return out
}

// DangerousRemoveItemRoot matches PowerShell Remove-Item against
// the system root in any flag order: `-Recurse -Force C:\`,
// `-Force -Recurse C:\`, `-Recurse C:\ -Force`, etc.
var DangerousRemoveItemRoot = regexp.MustCompile(
	removeItemAlt(`C:\\`))

// DangerousRemoveItemHome matches PowerShell Remove-Item against
// `$env:USERPROFILE` in any flag order.
var DangerousRemoveItemHome = regexp.MustCompile(
	removeItemAlt(`\$env:USERPROFILE`))

// DangerousFormatVolume matches PowerShell Format-Volume (drive
// reformat). Word-anchored so substrings in unrelated identifiers
// don't trip it.
var DangerousFormatVolume = regexp.MustCompile(`\bFormat-Volume\b`)
