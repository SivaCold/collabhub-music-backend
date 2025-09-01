# Clean script to fix import paths
Write-Host "Starting import path fixes..." -ForegroundColor Green

# Get all Go files recursively
$goFiles = Get-ChildItem -Path . -Filter "*.go" -Recurse

$oldImportPath = "github.com/collabhub/music-backend/"
$newImportPath = "collabhub-music-backend/"

Write-Host "Found $($goFiles.Count) Go files to process" -ForegroundColor Yellow

foreach ($file in $goFiles) {
    Write-Host "Processing: $($file.FullName)" -ForegroundColor Cyan
    
    $content = Get-Content -Path $file.FullName -Raw
    $originalContent = $content
    
    # Replace the import paths
    $content = $content -replace [regex]::Escape($oldImportPath), $newImportPath
    
    # Only write if content changed
    if ($content -ne $originalContent) {
        Set-Content -Path $file.FullName -Value $content -NoNewline
        Write-Host "  -> Updated imports in $($file.Name)" -ForegroundColor Green
    }
}

Write-Host "Import fixing complete!" -ForegroundColor Green
