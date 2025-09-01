# Script to fix all imports from github.com/collabhub/music-backend to collabhub-music-backend

Write-Host "🔧 Fixing imports in all Go files..." -ForegroundColor Yellow

$files = Get-ChildItem -Path . -Include "*.go" -Recurse

foreach ($file in $files) {
    $content = Get-Content $file.FullName -Raw
    $originalContent = $content
    
    # Replace all github.com/collabhub/music-backend imports
    $content = $content -replace 'github\.com/collabhub/music-backend/', 'collabhub-music-backend/'
    
    if ($content -ne $originalContent) {
        Set-Content -Path $file.FullName -Value $content -NoNewline
        Write-Host "✅ Fixed imports in: $($file.FullName)" -ForegroundColor Green
    }
}

Write-Host "Import fixing complete!" -ForegroundColor Cyan
