# Check if the OS is 64-bit
if ([Environment]::Is64BitOperatingSystem -eq $false) {
    Write-Host "This application requires a 64-bit version of Windows. Installation aborted."
    exit
}

# Set variables
$githubUser = "dwiandhikaap"
$githubRepo = "rawdog-md"
$baseUrl = "https://api.github.com/repos/$githubUser/$githubRepo/releases/latest"
$downloadFilePattern = "rawd-*-windows-amd64.zip"

# Get the latest release information
$response = Invoke-RestMethod -Uri $baseUrl -Headers @{ "User-Agent" = "PowerShell" }

# Find the download URL for the appropriate file
$downloadUrl = $response.assets | Where-Object { $_.name -like $downloadFilePattern } | Select-Object -First 1 -ExpandProperty browser_download_url

if (-not $downloadUrl) {
    Write-Host "No suitable release found."
    exit
}

# Define download location and extraction path
$downloadPath = "$env:TEMP\$($downloadUrl.Split('/')[-1])"
$extractPath = "$env:APPDATA\rawdog-md"

# Download the zip file
Invoke-WebRequest -Uri $downloadUrl -OutFile $downloadPath

# Create the extraction directory if it doesn't exist
if (-not (Test-Path -Path $extractPath)) {
    New-Item -ItemType Directory -Path $extractPath
}

# Extract the zip file
Expand-Archive -Path $downloadPath -DestinationPath $extractPath -Force

# Optionally, add the extracted folder to the PATH environment variable
$existingPath = [Environment]::GetEnvironmentVariable("Path", [EnvironmentVariableTarget]::User)
if ($existingPath -notlike "*$extractPath*") {
    [Environment]::SetEnvironmentVariable("Path", "$existingPath;$extractPath", [EnvironmentVariableTarget]::User)
    Write-Host "Added $extractPath to PATH."
} else {
    Write-Host "$extractPath is already in PATH."
}

# Clean up the downloaded zip file
Remove-Item -Path $downloadPath -Force

Write-Host "Installation complete! You can now use 'rawd' command in your terminal."