
$latestRelease = Invoke-WebRequest "https://api.github.com/repos/GetStream/stream-cli/releases/latest"
$json = $latestRelease.Content | ConvertFrom-Json

$url = $json.assets | ? { $_.name -match "Windows_x86" } | select -expand browser_download_url
if ($env:PROCESSOR_ARCHITECTURE -ne $null -and $env:PROCESSOR_ARCHITECTURE.Contains("arm")) {
    $url = $json.assets | ? { $_.name -match "Windows_arm" } | select -expand browser_download_url
}

# Typically C:\Users\User\AppData\Roaming
$appDataFolder = [Environment]::GetFolderPath([Environment+SpecialFolder]::ApplicationData)
$cliFolder = Join-Path $appDataFolder "stream-cli"
$zipFile = Join-Path $cliFolder "stream-cli.zip"
$exeFile = Join-Path $cliFolder "stream-cli.exe"
if (-not (Test-Path $cliFolder)) {
    New-Item -Path $cliFolder -ItemType Directory | Out-Null
}
if (Test-Path $zipFile) {
    Remove-Item -Force $zipFile
}
if (Test-Path $exeFile) {
    Remove-Item -Force $exeFile
}

# Download the zip file
Write-Host " > Downloading stream-cli zip file from GitHub..."
Invoke-WebRequest -Uri $url -OutFile $zipFile

# Unzip it and remove the zip file
Write-Host " > Unzipping stream-cli zip file in $cliFolder"
Expand-Archive -Path $zipFile -DestinationPath $cliFolder
Remove-Item $zipFile
Write-Host " > The exe is available at $exeFile"

# Add to $PATH
# It requires elevated access, so we'll rather ask the user to do it

Write-Host " > Done! Please add $cliFolder to your PATH environment variable."