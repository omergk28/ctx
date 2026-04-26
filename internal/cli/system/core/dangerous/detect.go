//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package dangerous

import (
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/regex"
)

// Detect returns the first matching variant for the given command,
// or an empty Match if no pattern fires.
//
// Parameters:
//   - command: the shell command string from the hook envelope
//
// Returns:
//   - Match: variant + descKey, both empty when no pattern matches
func Detect(command string) Match {
	switch {
	case regex.DangerousSudo.MatchString(command):
		return Match{
			hook.VariantSudo,
			text.DescKeyBlockDangerousSudo,
			command,
		}
	case regex.DangerousRmRfRoot.MatchString(command):
		return Match{
			hook.VariantRmRfRoot,
			text.DescKeyBlockDangerousRmRfRoot,
			command,
		}
	case regex.DangerousRmRfHome.MatchString(command):
		return Match{
			hook.VariantRmRfHome,
			text.DescKeyBlockDangerousRmRfHome,
			command,
		}
	case regex.DangerousChmod777.MatchString(command):
		return Match{
			hook.VariantChmod777,
			text.DescKeyBlockDangerousChmod777,
			command,
		}
	case regex.DangerousGitPushForce.MatchString(command):
		return Match{
			hook.VariantGitPushForce,
			text.DescKeyBlockDangerousGitPushForce,
			command,
		}
	case regex.DangerousGitResetHard.MatchString(command):
		return Match{
			hook.VariantGitResetHard,
			text.DescKeyBlockDangerousGitResetHard,
			command,
		}
	case regex.DangerousRemoveItemRoot.MatchString(command):
		return Match{
			hook.VariantRemoveItemRoot,
			text.DescKeyBlockDangerousRemoveItemRoot,
			command,
		}
	case regex.DangerousRemoveItemHome.MatchString(command):
		return Match{
			hook.VariantRemoveItemHome,
			text.DescKeyBlockDangerousRemoveItemHome,
			command,
		}
	case regex.DangerousFormatVolume.MatchString(command):
		return Match{
			hook.VariantFormatVolume,
			text.DescKeyBlockDangerousFormatVolume,
			command,
		}
	}
	return Match{}
}
