Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$repoRoot = Split-Path -Parent $PSScriptRoot
$infoPath = Join-Path $repoRoot 'README.md'

if (-not (Test-Path $infoPath)) {
    Write-Error 'README.md not found at repository root.'
}

Push-Location $repoRoot
try {
    $stagedPaths = git diff --cached --name-only --diff-filter=ACMR
    if (-not $stagedPaths) {
        exit 0
    }

    $topLevelDirs = @(
        $stagedPaths |
            ForEach-Object { $_.Trim() } |
            Where-Object {
                $_ -and
                $_ -match '[\\/]' -and
                $_ -notmatch '^\.githooks[\\/]' -and
                $_ -notmatch '^scripts[\\/]' -and
                $_ -notmatch '^\.git[\\/]'
            } |
            ForEach-Object {
                ($_ -split '[\\/]')[0]
            } |
            Sort-Object -Unique
    )

    if (-not $topLevelDirs -or $topLevelDirs.Count -eq 0) {
        exit 0
    }

    $content = Get-Content -Path $infoPath -Raw -Encoding UTF8
    $existing = [System.Collections.Generic.HashSet[string]]::new([System.StringComparer]::OrdinalIgnoreCase)

    foreach ($match in [regex]::Matches($content, '(?m)^\|\s*([^|]+?)\s*\|')) {
        $name = $match.Groups[1].Value.Trim()
        if ($name -and $name -ne '資料夾' -and $name -ne '---') {
            [void]$existing.Add($name)
        }
    }

    $missing = @(
        $topLevelDirs |
            Where-Object {
                -not $existing.Contains($_) -and (Test-Path (Join-Path $repoRoot $_))
            }
    )

    if (-not $missing -or $missing.Count -eq 0) {
        exit 0
    }

    $newRows = ($missing | Sort-Object | ForEach-Object {
        "| $_ | TODO: 補上此資料夾的專案主題說明。 |"
    }) -join "`r`n"

    $marker = '## 補充'
    if ($content.Contains($marker)) {
        $updated = $content.Replace($marker, "$newRows`r`n`r`n$marker")
    }
    else {
        $updated = $content.TrimEnd() + "`r`n`r`n$newRows`r`n"
    }

    [System.IO.File]::WriteAllText($infoPath, $updated, [System.Text.UTF8Encoding]::new($false))
    git add -- README.md | Out-Null

    Write-Host ('Updated README.md with: ' + ($missing -join ', '))
}
finally {
    Pop-Location
}