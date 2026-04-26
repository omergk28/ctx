//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package block_dangerous_commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/assets/read/desc"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/dangerous"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/message"
	"github.com/ActiveMemory/ctx/internal/cli/system/core/nudge"
	coreSession "github.com/ActiveMemory/ctx/internal/cli/system/core/session"
	"github.com/ActiveMemory/ctx/internal/config/embed/text"
	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/config/token"
	"github.com/ActiveMemory/ctx/internal/entity"
	"github.com/ActiveMemory/ctx/internal/notify"
	writeSetup "github.com/ActiveMemory/ctx/internal/write/setup"
)

// Run executes the block-dangerous-commands hook logic.
//
// Reads a hook input from stdin, checks the command against the
// dangerous-pattern set, and emits a block response if matched.
// Uses first-match-wins ordering so a single command never
// produces an ambiguous variant.
//
// Parameters:
//   - cmd: Cobra command for output
//   - stdin: standard input for hook JSON
//
// Returns:
//   - error: Always nil (hook errors are non-fatal)
func Run(cmd *cobra.Command, stdin *os.File) error {
	input := coreSession.ReadInput(stdin)
	command := input.ToolInput.Command
	if command == "" {
		return nil
	}

	hit := dangerous.Detect(command)
	if hit.Variant == "" {
		return nil
	}

	fallback := desc.Text(hit.DescKey)
	reason := message.Load(
		hook.BlockDangerousCommands, hit.Variant, nil, fallback,
	)
	if reason == "" {
		return nil
	}

	resp := entity.BlockResponse{
		Decision: hook.DecisionBlock,
		Reason: reason + token.NewlineLF + token.NewlineLF +
			desc.Text(text.DescKeyBlockConstitutionSuffix),
	}
	data, _ := json.Marshal(resp)
	writeSetup.BlockResponse(cmd, string(data))

	blockRef := notify.NewTemplateRef(
		hook.BlockDangerousCommands, hit.Variant, nil,
	)
	return nudge.Relay(fmt.Sprintf(
		desc.Text(text.DescKeyRelayPrefixFormat),
		hook.BlockDangerousCommands,
		desc.Text(text.DescKeyBlockDangerousRelayMessage),
	), input.SessionID, blockRef)
}
