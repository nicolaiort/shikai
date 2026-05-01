[CmdletBinding()]
param(
    [string]$InstallDir
)

$ErrorActionPreference = 'Stop'

function Get-ArchitectureName {
    $arch = [System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture
    if ($arch -eq [System.Runtime.InteropServices.Architecture]::Arm64) {
        return 'arm64'
    }

    if ($arch -eq [System.Runtime.InteropServices.Architecture]::X64) {
        return 'amd64'
    }

    throw "Unsupported architecture: $arch"
}

function Get-DefaultInstallDir {
    if ($IsWindows) {
        return Join-Path $HOME 'bin'
    }

    if ($IsMacOS) {
        return Join-Path $HOME 'bin'
    }

    if ($IsLinux) {
        return Join-Path $HOME 'bin'
    }

    throw 'This installer only supports Linux, macOS, and Windows.'
}

$architecture = Get-ArchitectureName
$installDir = if ($InstallDir) { $InstallDir } else { Get-DefaultInstallDir }
$binaryName = if ($IsWindows) { 'shikai.exe' } else { 'shikai' }
$assetName = if ($IsWindows) {
    "shikai-windows-$architecture.exe"
} elseif ($IsMacOS) {
    "shikai-darwin-$architecture"
} elseif ($IsLinux) {
    "shikai-linux-$architecture"
} else {
    throw 'This installer only supports Linux, macOS, and Windows.'
}

$headers = @{
    'Accept' = 'application/vnd.github+json'
    'User-Agent' = 'shikai-installer'
}

$release = Invoke-RestMethod -Uri 'https://api.github.com/repos/nicolaiort/shikai/releases/latest' -Headers $headers
$asset = $release.assets | Where-Object { $_.name -eq $assetName } | Select-Object -First 1

if (-not $asset) {
    throw "Could not find asset '$assetName' in the latest release."
}

New-Item -ItemType Directory -Force -Path $installDir | Out-Null

$destination = Join-Path $installDir $binaryName
Invoke-WebRequest -Uri $asset.browser_download_url -OutFile $destination -Headers $headers

if ($IsMacOS -or $IsLinux) {
    & chmod +x $destination
}

Write-Host "Installed $binaryName to $destination"
