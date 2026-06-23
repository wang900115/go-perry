Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot

Push-Location $repoRoot
try {
    git config core.hooksPath .githooks
    Write-Host 'Configured core.hooksPath to .githooks'
    Write-Host 'Git will now auto-fill commit messages and sync info.md on commit.'
}
finally {
    Pop-Location
}