//   /    ctx:                         https://ctx.ist
// ,'`./    do you remember?
// `.,'\
//   \    Copyright 2026-present Context contributors.
//                 SPDX-License-Identifier: Apache-2.0

package block_dangerous_commands

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"

	"github.com/ActiveMemory/ctx/internal/config/hook"
	"github.com/ActiveMemory/ctx/internal/entity"
)

// runCmd writes a hook envelope containing command to a temp stdin
// file, invokes Run, and returns the captured output. Empty output
// signals a silent allow; non-empty output should parse as a
// BlockResponse.
func runCmd(t *testing.T, command string) string {
	t.Helper()

	envelope, err := json.Marshal(entity.HookInput{
		SessionID: "test-session",
		ToolInput: entity.ToolInput{Command: command},
	})
	if err != nil {
		t.Fatalf("marshal envelope: %v", err)
	}

	stdinPath := filepath.Join(t.TempDir(), "stdin")
	if writeErr := os.WriteFile(stdinPath, envelope, 0o600); writeErr != nil {
		t.Fatalf("write stdin: %v", writeErr)
	}
	stdin, openErr := os.Open(stdinPath)
	if openErr != nil {
		t.Fatalf("open stdin: %v", openErr)
	}
	t.Cleanup(func() { _ = stdin.Close() })

	c := &cobra.Command{}
	var out bytes.Buffer
	c.SetOut(&out)
	c.SetErr(&out)
	if runErr := Run(c, stdin); runErr != nil {
		t.Fatalf("Run() err = %v, want nil", runErr)
	}
	return out.String()
}

// requireBlock asserts the captured output contains a BlockResponse
// with decision "block" and a non-empty reason. Returns the parsed
// reason for further checks.
func requireBlock(t *testing.T, out string) string {
	t.Helper()
	if out == "" {
		t.Fatal("expected block response, got silent allow")
	}
	// Output may include surrounding text or JSON-only; locate the JSON.
	start := strings.Index(out, "{")
	end := strings.LastIndex(out, "}")
	if start < 0 || end < 0 || end <= start {
		t.Fatalf("no JSON object in output: %q", out)
	}
	var resp entity.BlockResponse
	if err := json.Unmarshal([]byte(out[start:end+1]), &resp); err != nil {
		t.Fatalf("unmarshal block response: %v (raw: %q)", err, out)
	}
	if resp.Decision != hook.DecisionBlock {
		t.Errorf("decision = %q, want %q", resp.Decision, hook.DecisionBlock)
	}
	if resp.Reason == "" {
		t.Error("reason should not be empty")
	}
	return resp.Reason
}

// TestRun_EmptyCommand: zero-value envelope is a silent allow.
func TestRun_EmptyCommand(t *testing.T) {
	out := runCmd(t, "")
	if out != "" {
		t.Errorf("empty command should be silent, got %q", out)
	}
}

// TestRun_BenignCommands: common commands must NOT trip the hook.
// These represent the false-positive risk the regex set has to avoid.
func TestRun_BenignCommands(t *testing.T) {
	benign := []string{
		"ls -la",
		"git status",
		"git push origin main",
		"git push --force-with-lease origin main",
		"git push origin feat/branch --force-with-lease",
		"git reset HEAD~1",
		"git reset --soft HEAD~1",
		"rm -rf /var/log/old",
		"rm -rf node_modules",
		"rm -rf ./build",
		"chmod 755 script.sh",
		"chmod 644 file.txt",
		"echo 'sudo make install'",
		"grep pseudo file.txt",
		"go test ./...",
		"make build",
		// Print/log statements naming dangerous commands must
		// NOT trip the guard — they're print arguments, not
		// commands being run.
		"echo rm -rf /",
		"echo sudo make install",
		`echo "sudo rm -rf /"`,
		`echo "Remove-Item -Recurse -Force C:\\"`,
		`git log --grep="rm -rf /"`,
		"$(echo sudo)",
	}
	for _, cmd := range benign {
		t.Run(cmd, func(t *testing.T) {
			if out := runCmd(t, cmd); out != "" {
				t.Errorf("benign command %q produced block: %s", cmd, out)
			}
		})
	}
}

// TestRun_Sudo: sudo must be blocked at start and after separators.
func TestRun_Sudo(t *testing.T) {
	for _, cmd := range []string{
		"sudo make install",
		"sudo rm /tmp/foo",
		"cd /tmp && sudo make install",
		"echo hi; sudo cat /etc/shadow",
	} {
		t.Run(cmd, func(t *testing.T) {
			reason := requireBlock(t, runCmd(t, cmd))
			if !strings.Contains(strings.ToLower(reason), "sudo") &&
				!strings.Contains(strings.ToLower(reason), "privilege") {
				t.Errorf("reason %q should mention sudo/privilege", reason)
			}
		})
	}
}

// TestRun_RmRfRoot: rm -rf / must be blocked but rm -rf /var/log
// must be allowed (anchored on trailing whitespace/EOL).
func TestRun_RmRfRoot(t *testing.T) {
	for _, cmd := range []string{
		"rm -rf /",
		"rm -rf / ",
		"rm -Rf /",
		"rm -fr /",
	} {
		t.Run(cmd, func(t *testing.T) {
			requireBlock(t, runCmd(t, cmd))
		})
	}
}

// TestRun_RmRfHome: rm -rf ~ and any subpath of ~ must be blocked.
func TestRun_RmRfHome(t *testing.T) {
	for _, cmd := range []string{
		"rm -rf ~",
		"rm -rf ~/Downloads",
		"rm -Rf ~/.config",
	} {
		t.Run(cmd, func(t *testing.T) {
			requireBlock(t, runCmd(t, cmd))
		})
	}
}

// TestRun_Chmod777: chmod 777 must be blocked.
func TestRun_Chmod777(t *testing.T) {
	for _, cmd := range []string{
		"chmod 777 file",
		"chmod -R 777 dir",
		"chmod 777 /tmp/x",
	} {
		t.Run(cmd, func(t *testing.T) {
			requireBlock(t, runCmd(t, cmd))
		})
	}
}

// TestRun_GitPushForce: --force / -f (including combined short
// flag bundles like -fu) must be blocked; --force-with-lease must
// be allowed.
func TestRun_GitPushForce(t *testing.T) {
	blocked := []string{
		"git push --force",
		"git push -f",
		"git push origin main --force",
		"git push origin main -f",
		"git push -fu origin main",
		"git push -uf origin main",
	}
	for _, cmd := range blocked {
		t.Run("block/"+cmd, func(t *testing.T) {
			requireBlock(t, runCmd(t, cmd))
		})
	}

	allowed := []string{
		"git push --force-with-lease",
		"git push origin main --force-with-lease",
	}
	for _, cmd := range allowed {
		t.Run("allow/"+cmd, func(t *testing.T) {
			if out := runCmd(t, cmd); out != "" {
				t.Errorf("--force-with-lease must be allowed, got: %s", out)
			}
		})
	}
}

// TestRun_GitResetHard: git reset --hard must be blocked.
func TestRun_GitResetHard(t *testing.T) {
	for _, cmd := range []string{
		"git reset --hard",
		"git reset --hard HEAD~3",
		"git reset --hard origin/main",
	} {
		t.Run(cmd, func(t *testing.T) {
			requireBlock(t, runCmd(t, cmd))
		})
	}
}

// TestRun_BlockResponseShape: a blocked command yields valid JSON
// with the constitution suffix appended to the reason.
func TestRun_BlockResponseShape(t *testing.T) {
	out := runCmd(t, "sudo make install")
	reason := requireBlock(t, out)
	if !strings.Contains(reason, "CONSTITUTION") {
		t.Errorf("reason should append constitution suffix, got %q", reason)
	}
}

// TestRun_PowerShellRemoveItem: PowerShell Remove-Item targeting
// system root or user home must be blocked across all flag
// orderings (canonical, reversed, and target-between-flags).
func TestRun_PowerShellRemoveItem(t *testing.T) {
	for _, cmd := range []string{
		`Remove-Item -Recurse -Force C:\`,
		`Remove-Item -Recurse -Force C:\Users`,
		`Remove-Item -Force -Recurse C:\`,
		`Remove-Item -Recurse C:\ -Force`,
		`Remove-Item -Force C:\ -Recurse`,
		`Remove-Item -Recurse -Force $env:USERPROFILE`,
		`Remove-Item -Recurse -Force $env:USERPROFILE\Documents`,
		`Remove-Item -Force -Recurse $env:USERPROFILE`,
		`Remove-Item -Recurse $env:USERPROFILE -Force`,
	} {
		t.Run(cmd, func(t *testing.T) {
			requireBlock(t, runCmd(t, cmd))
		})
	}
}

// TestRun_PowerShellFormatVolume: Format-Volume must be blocked.
func TestRun_PowerShellFormatVolume(t *testing.T) {
	for _, cmd := range []string{
		`Format-Volume -DriveLetter D`,
		`Get-Volume D | Format-Volume`,
	} {
		t.Run(cmd, func(t *testing.T) {
			requireBlock(t, runCmd(t, cmd))
		})
	}
}
