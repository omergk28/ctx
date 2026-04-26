//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package hook

// Hook variant constants: template selectors passed to Load and
// NewTemplateRef to choose the appropriate message for each trigger type.
const (
	// VariantChmod777 selects the chmod 777 block message.
	VariantChmod777 = "chmod-777"
	// VariantDotSlash selects the relative path (./ctx) block message.
	VariantDotSlash = "dot-slash"
	// VariantFormatVolume selects the PowerShell Format-Volume block
	// message.
	VariantFormatVolume = "format-volume"
	// VariantGitPushForce selects the git push --force block message.
	VariantGitPushForce = "git-push-force"
	// VariantGitResetHard selects the git reset --hard block message.
	VariantGitResetHard = "git-reset-hard"
	// VariantRemoveItemHome selects the PowerShell Remove-Item -Recurse
	// -Force $env:USERPROFILE block message.
	VariantRemoveItemHome = "remove-item-home"
	// VariantRemoveItemRoot selects the PowerShell Remove-Item -Recurse
	// -Force C:\ block message.
	VariantRemoveItemRoot = "remove-item-root"
	// VariantRmRfHome selects the rm -rf ~ block message.
	VariantRmRfHome = "rm-rf-home"
	// VariantRmRfRoot selects the rm -rf / block message.
	VariantRmRfRoot = "rm-rf-root"
	// VariantSudo selects the sudo escalation block message.
	VariantSudo = "sudo"
	// VariantGoRun selects the go run block message.
	VariantGoRun = "go-run"
	// VariantAbsolutePath selects the absolute path block message.
	VariantAbsolutePath = "absolute-path"
	// VariantBoth selects the template for both ceremonies missing.
	VariantBoth = "both"
	// VariantRemember selects the template for missing /ctx-remember.
	VariantRemember = "remember"
	// VariantWrapup selects the template for missing /ctx-wrap-up.
	VariantWrapup = "wrapup"
	// VariantUnimported selects the unimported journal entries variant.
	VariantUnimported = "unimported"
	// VariantUnenriched selects the unenriched journal entries variant.
	VariantUnenriched = "unenriched"
	// VariantWarning selects the generic warning variant.
	VariantWarning = "warning"
	// VariantAlert selects the alert variant.
	VariantAlert = "alert"
	// VariantBilling selects the billing threshold variant.
	VariantBilling = "billing"
	// VariantCheckpoint selects the checkpoint variant.
	VariantCheckpoint = "checkpoint"
	// VariantGate selects the gate variant.
	VariantGate = "gate"
	// VariantKeyRotation selects the key rotation variant.
	VariantKeyRotation = "key-rotation"
	// VariantMismatch selects the version mismatch variant.
	VariantMismatch = "mismatch"
	// VariantNudge selects the generic nudge variant.
	VariantNudge = "nudge"
	// VariantOversize selects the oversize threshold variant.
	VariantOversize = "oversize"
	// VariantPulse selects the heartbeat pulse variant.
	VariantPulse = "pulse"
	// VariantReminders selects the reminders variant.
	VariantReminders = "reminders"
	// VariantStale selects the staleness variant.
	VariantStale = "stale"
	// VariantWindow selects the context window variant.
	VariantWindow = "window"
)
