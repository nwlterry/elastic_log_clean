# Windows 11 runner for elastic_log_clean.
# Use after copying an Elastic Support Diagnostics zip or server log off the host.
#
#   powershell -ExecutionPolicy Bypass -File .\elastic_log_clean.ps1 -InputPath .\local-diagnostics.zip
#   .\elastic_log_clean.ps1 -InputPath C:\logs\elasticsearch -VerboseLog
#   .\elastic_log_clean.ps1 -InputPath .\cluster.log -Config .\elastic_ip_name_map.yaml

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true, Position = 0)]
    [string]$InputPath,

    [string]$Output,
    [string]$Config,
    [string]$Report = "report.yaml",
    [int]$Workers = 0,
    [switch]$VerboseLog,
    [switch]$Stdin
)

$ErrorActionPreference = "Stop"
$here = Split-Path -Parent $MyInvocation.MyCommand.Path
$search = @(
    (Join-Path $here "elastic_log_clean.exe"),
    (Join-Path $here "elastic_log_clean"),
    (Join-Path (Split-Path $here -Parent) "bin\elastic_log_clean.exe"),
    (Join-Path (Split-Path $here -Parent) "elastic_log_clean.exe")
)
$exe = $search | Where-Object { Test-Path $_ } | Select-Object -First 1
if (-not $exe) {
    throw "elastic_log_clean.exe not found next to this script. Unpack the windows_amd64 zip first."
}

if (-not $Stdin) {
    if (-not (Test-Path -LiteralPath $InputPath)) {
        throw "Input not found: $InputPath"
    }
    $InputPath = (Resolve-Path -LiteralPath $InputPath).Path
}

$argsList = @()
if ($Config) { $argsList += @("-c", $Config) }
if ($InputPath -and -not $Stdin) { $argsList += @("-i", $InputPath) }
if ($Output) { $argsList += @("-o", $Output) }
if ($Report) { $argsList += @("-r", $Report) }
if ($Workers -gt 0) { $argsList += @("-w", "$Workers") }
if ($VerboseLog) { $argsList += "-v" }
if ($Stdin) { $argsList += "-stdin" }

Write-Host "[elastic_log_clean] $exe $($argsList -join ' ')"
& $exe @argsList
exit $LASTEXITCODE
