param(
    [Parameter(Mandatory = $true)]
    [string]$MessageFile,
    [string]$MessageSource,
    [string]$CommitSha
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

if ($MessageSource -in @('message', 'commit', 'merge', 'squash')) {
    exit 0
}

if (-not (Test-Path $MessageFile)) {
    Write-Error 'Commit message file not found.'
}

$repoRoot = Split-Path -Parent $PSScriptRoot

function Get-TopLevelDirectories {
    Push-Location $repoRoot
    try {
        $paths = git diff --cached --name-only --diff-filter=ACMR
        if (-not $paths) {
            return @()
        }

        return @(
            $paths |
                ForEach-Object { $_.Trim() } |
                Where-Object {
                    $_ -and
                    $_ -match '[\\/]' -and
                    $_ -notmatch '^\.githooks[\\/]' -and
                    $_ -notmatch '^scripts[\\/]'
                } |
                ForEach-Object { ($_ -split '[\\/]')[0] } |
                Sort-Object -Unique
        )
    }
    finally {
        Pop-Location
    }
}

function Get-DefaultMessage {
    $dirs = @(Get-TopLevelDirectories)
    if (-not $dirs -or $dirs.Count -eq 0) {
        return 'chore: update repository notes'
    }

    if ($dirs.Count -eq 1) {
        $dir = $dirs[0]

        Push-Location $repoRoot
        try {
            $statuses = git diff --cached --name-status -- $dir
        }
        finally {
            Pop-Location
        }

        $added = @($statuses | Where-Object { $_ -match '^A\s' }).Count
        $modified = @($statuses | Where-Object { $_ -match '^(M|R|C)\s' }).Count

        if ($added -gt 0 -and $modified -eq 0) {
            return "feat: add $dir topic"
        }

        return "chore: update $dir examples"
    }

    $preview = ($dirs | Select-Object -First 3) -join ', '
    if ($dirs.Count -gt 3) {
        return "chore: update $preview and more"
    }

    return "chore: update $preview"
}

$lines = Get-Content -Path $MessageFile -Encoding UTF8
$hasUserMessage = $false
foreach ($line in $lines) {
    if ($line -and -not $line.TrimStart().StartsWith('#')) {
        $hasUserMessage = $true
        break
    }
}

if ($hasUserMessage) {
    exit 0
}

$defaultMessage = Get-DefaultMessage
[System.IO.File]::WriteAllText($MessageFile, $defaultMessage + "`r`n", [System.Text.UTF8Encoding]::new($false))