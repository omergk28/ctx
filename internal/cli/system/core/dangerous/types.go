//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dangerous

// Match maps a regex hit to the variant + fallback DescKey it
// resolves to. First-match-wins ordering in [Detect] is
// intentional: more specific patterns are listed first so a single
// command never produces an ambiguous variant.
type Match struct {
	Variant      string
	DescKey      string
	MatchedInput string
}
