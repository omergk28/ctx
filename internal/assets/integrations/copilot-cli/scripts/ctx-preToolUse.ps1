# ctx preToolUse hook for GitHub Copilot CLI (Windows / PowerShell).
# Reshapes the Copilot envelope into the ctx hook envelope and
# delegates the dangerous-command decision to
# `ctx system block-dangerous-commands` (single source of truth
# shared with Claude Code and OpenCode integrations).
$ErrorActionPreference = 'SilentlyContinue'

$RawInput = $input | Out-String
if (-not $RawInput) { exit 0 }

try {
    $Data = $RawInput | ConvertFrom-Json
} catch {
    exit 0
}

$ToolName = if ($Data.tool_name) { $Data.tool_name } elseif ($Data.tool) { $Data.tool } else { '' }
if ($ToolName -ne 'shell' -and $ToolName -ne 'powershell' -and $ToolName -ne 'bash') {
    exit 0
}

$Command = ''
if ($Data.input -and $Data.input.command) {
    $Command = [string]$Data.input.command
}
if (-not $Command) { exit 0 }

$Envelope = @{
    session_id = 'copilot-cli'
    tool_input = @{ command = $Command }
} | ConvertTo-Json -Compress

# Run the Go hook. Missing binary → fail open.
try {
    $Decision = $Envelope | ctx system block-dangerous-commands 2>$null
} catch {
    exit 0
}
if (-not $Decision) { exit 0 }

try {
    $Parsed = $Decision | ConvertFrom-Json
} catch {
    exit 0
}

if ($Parsed.decision -eq 'block') {
    Write-Error $Decision
    exit 1
}
