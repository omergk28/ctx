// ctx OpenCode plugin — thin shim to ctx system subcommands.
// All real logic lives in the ctx Go binary; this plugin just
// wires OpenCode lifecycle hooks to ctx system calls.
//
// Tool names below match @opencode-ai/plugin v1.4.x. If the
// upstream renames a tool, the corresponding branch silently
// no-ops; verify against the OpenCode plugin docs when bumping.
import type { Plugin } from "@opencode-ai/plugin"

const SHELL_TOOLS = new Set(["shell", "bash"])
const EDIT_TOOLS = new Set(["edit", "write", "file_edit"])
// Match `git commit` but not `git commit-tree` / `git commit-graph`.
// The negative lookahead rejects `-` immediately after the boundary.
const GIT_COMMIT_RE = /\bgit\s+commit\b(?!-)/

function extractCommand(input: unknown): string {
  if (typeof input === "string") return input
  if (input && typeof input === "object") {
    const cmd = (input as { command?: unknown }).command
    if (typeof cmd === "string") return cmd
  }
  return ""
}

export default ((ctx) => ({
  "shell.env": () => ({
    CTX_DIR: ".context",
  }),
  event: {
    "session.created": async () => {
      await ctx.$`ctx system bootstrap 2>/dev/null || true`
      await ctx.$`ctx agent --budget 4000 2>/dev/null || true`
    },
    "session.idle": async () => {
      await ctx.$`ctx system check-persistence 2>/dev/null || true`
      await ctx.$`ctx system check-task-completion 2>/dev/null || true`
    },
  },
  "tool.execute.before": async ({ tool, sessionID }, output) => {
    if (!SHELL_TOOLS.has(tool)) return
    const cmd = extractCommand(output?.args)
    if (!cmd) return
    const envelope = JSON.stringify({
      session_id: sessionID,
      tool_input: { command: cmd },
    })
    let stdout = ""
    try {
      const result = await ctx.$`echo ${envelope} | ctx system block-dangerous-commands`.quiet()
      stdout = result.stdout?.toString() ?? ""
    } catch {
      // Subcommand missing or failed — fail open. Re-installing ctx
      // restores the safety net.
      return
    }
    const trimmed = stdout.trim()
    if (!trimmed) return
    let decision: { decision?: string; reason?: string } | null = null
    try {
      decision = JSON.parse(trimmed) as { decision?: string; reason?: string }
    } catch {
      // Unparseable output → fail open.
      return
    }
    if (decision?.decision === "block") {
      throw new Error(decision.reason ?? "ctx: blocked dangerous command")
    }
  },
  "tool.execute.after": async ({ tool, args }) => {
    if (SHELL_TOOLS.has(tool)) {
      const cmd = extractCommand(args)
      if (GIT_COMMIT_RE.test(cmd)) {
        await ctx.$`ctx system post-commit 2>/dev/null || true`
      }
    }
    if (EDIT_TOOLS.has(tool)) {
      await ctx.$`ctx system check-task-completion 2>/dev/null || true`
    }
  },
})) satisfies Plugin
