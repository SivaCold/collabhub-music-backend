# Test script for ZIP API endpoints
param(
    [string]$ZipFilePath = "C:\Users\Coldash (Dev)\Downloads\Overtime.zip",
    [string]$ServerUrl = "http://localhost:8081"
)

Write-Host "=== Testing CollabHub Music ZIP API ===" -ForegroundColor Green

# Check if ZIP file exists
if (-not (Test-Path $ZipFilePath)) {
    Write-Host "[ERROR] ZIP file not found: $ZipFilePath" -ForegroundColor Red
    Write-Host "Please provide a valid ZIP file path using -ZipFilePath parameter"
    exit 1
}

try {
    # Test 1: Upload ZIP file
    Write-Host "`n[TEST 1] Uploading ZIP file..." -ForegroundColor Yellow
    Write-Host "   File: $ZipFilePath"
    
    # Create multipart form data manually
    $boundary = [System.Guid]::NewGuid().ToString()
    $LF = "`r`n"
    
    $fileBytes = [System.IO.File]::ReadAllBytes($ZipFilePath)
    $fileName = [System.IO.Path]::GetFileName($ZipFilePath)
    
    $bodyLines = @(
        "--$boundary",
        "Content-Disposition: form-data; name=`"file`"; filename=`"$fileName`"",
        "Content-Type: application/zip",
        "",
        [System.Text.Encoding]::GetEncoding("ISO-8859-1").GetString($fileBytes),
        "--$boundary--"
    ) -join $LF
    
    $uploadResponse = Invoke-RestMethod -Uri "$ServerUrl/api/v1/files/zip/upload" `
        -Method Post `
        -Body $bodyLines `
        -ContentType "multipart/form-data; boundary=$boundary"

    if ($uploadResponse.status -eq "success") {
        Write-Host "[SUCCESS] Upload successful!" -ForegroundColor Green
        $fileId = $uploadResponse.data.file_id
        Write-Host "   File ID: $fileId"
        Write-Host "   Audio files: $($uploadResponse.data.audio_files)"
        Write-Host "   Total files: $($uploadResponse.data.total_files)"
        Write-Host "   Total size: $([math]::Round($uploadResponse.data.total_size / 1MB, 2)) MB"
    } else {
        Write-Host "[ERROR] Upload failed: $($uploadResponse.error)" -ForegroundColor Red
        exit 1
    }

    # Test 2: Get ZIP info
    Write-Host "`n[TEST 2] Getting ZIP information..." -ForegroundColor Yellow
    
    $infoResponse = Invoke-RestMethod -Uri "$ServerUrl/api/v1/files/zip/$fileId/info" `
        -Method Get

    if ($infoResponse.status -eq "success") {
        Write-Host "[SUCCESS] Info retrieved successfully!" -ForegroundColor Green
        Write-Host "   Supported files: $($infoResponse.data.supported_files.Count)"
        Write-Host "   Unsupported files: $($infoResponse.data.unsupported_files.Count)"
        Write-Host "   Total size: $([math]::Round($infoResponse.data.total_size / 1MB, 2)) MB"
        
        if ($infoResponse.data.supported_files.Count -gt 0) {
            Write-Host "   First few audio files:"
            $infoResponse.data.supported_files | Select-Object -First 3 | ForEach-Object {
                Write-Host "     - $_"
            }
        }
    }

    # Test 3: Extract ZIP
    Write-Host "`n[TEST 3] Extracting ZIP file..." -ForegroundColor Yellow
    
    $extractResponse = Invoke-RestMethod -Uri "$ServerUrl/api/v1/files/zip/$fileId/extract" `
        -Method Post `
        -ContentType "application/json"

    if ($extractResponse.status -eq "success") {
        Write-Host "[SUCCESS] Extraction successful!" -ForegroundColor Green
        $projectId = $extractResponse.data.project_id
        Write-Host "   Project ID: $projectId"
        Write-Host "   Extracted files: $($extractResponse.data.total_files)"
        Write-Host "   Audio files: $($extractResponse.data.audio_files.Count)"
        Write-Host "   Extracted to: $($extractResponse.data.extracted_path)"
    }

    # Test 4: List extracted files
    Write-Host "`n[TEST 4] Listing extracted files..." -ForegroundColor Yellow
    
    $filesResponse = Invoke-RestMethod -Uri "$ServerUrl/api/v1/files/projects/$projectId/files" `
        -Method Get

    if ($filesResponse.status -eq "success") {
        Write-Host "[SUCCESS] Files listed successfully!" -ForegroundColor Green
        Write-Host "   Total files: $($filesResponse.data.total_files)"
        Write-Host "   Audio files: $($filesResponse.data.audio_files)"
        
        # Show first few files
        if ($filesResponse.data.files.Count -gt 0) {
            Write-Host "   Sample files:"
            $filesResponse.data.files | Select-Object -First 5 | ForEach-Object {
                $audioIcon = if ($_.is_audio_file) { "[AUDIO]" } else { "[FILE]" }
                $sizeKB = [math]::Round($_.size / 1024, 1)
                $fileType = if ($_.is_directory) { "DIR" } else { "FILE" }
                Write-Host "     $audioIcon [$fileType] $($_.name) ($sizeKB KB)"
            }
        }
    }

    # Test 5: List only audio files
    Write-Host "`n[TEST 5] Listing only audio files..." -ForegroundColor Yellow
    
    $audioFilesResponse = Invoke-RestMethod -Uri "$ServerUrl/api/v1/files/projects/$projectId/files?audio_only=true" `
        -Method Get

    if ($audioFilesResponse.status -eq "success") {
        Write-Host "[SUCCESS] Audio files listed successfully!" -ForegroundColor Green
        Write-Host "   Audio files count: $($audioFilesResponse.data.audio_files)"
        
        if ($audioFilesResponse.data.files.Count -gt 0) {
            Write-Host "   Audio files:"
            $audioFilesResponse.data.files | ForEach-Object {
                $sizeKB = [math]::Round($_.size / 1024, 1)
                Write-Host "     [AUDIO] $($_.name) ($sizeKB KB) - $($_.content_type)"
            }
        }
    }

    # Test 6: Create project from ZIP
    Write-Host "`n[TEST 6] Creating project from ZIP..." -ForegroundColor Yellow
    
    $projectData = @{
        name = "PowerShell Test Project"
        description = "Project created via PowerShell API test"
        genre = "Electronic"
        bpm = 128
        key = "Am"
    } | ConvertTo-Json

    $projectResponse = Invoke-RestMethod -Uri "$ServerUrl/api/v1/files/zip/$fileId/project" `
        -Method Post `
        -Body $projectData `
        -ContentType "application/json"

    if ($projectResponse.status -eq "success") {
        Write-Host "[SUCCESS] Project created successfully!" -ForegroundColor Green
        Write-Host "   Project name: $($projectResponse.data.name)"
        Write-Host "   Project description: $($projectResponse.data.description)"
        Write-Host "   Audio files: $($projectResponse.data.audio_files)"
        Write-Host "   Extracted files: $($projectResponse.data.extracted_files)"
    }

    # Test 7: Health check
    Write-Host "`n[TEST 7] Health check..." -ForegroundColor Yellow
    
    $healthResponse = Invoke-RestMethod -Uri "$ServerUrl/api/v1/health" `
        -Method Get

    if ($healthResponse.status -eq "ok") {
        Write-Host "[SUCCESS] Server is healthy!" -ForegroundColor Green
        Write-Host "   Message: $($healthResponse.message)"
        Write-Host "   Version: $($healthResponse.version)"
    }

    Write-Host "`n[COMPLETED] All tests completed successfully!" -ForegroundColor Green
    Write-Host "Summary:" -ForegroundColor Cyan
    Write-Host "- ZIP file uploaded and validated" -ForegroundColor Cyan
    Write-Host "- Files extracted to project: $projectId" -ForegroundColor Cyan
    Write-Host "- Project created successfully" -ForegroundColor Cyan

} catch {
    Write-Host "`n[ERROR] Test failed with error:" -ForegroundColor Red
    Write-Host $_.Exception.Message -ForegroundColor Red
    
    if ($_.Exception.Response) {
        try {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $responseBody = $reader.ReadToEnd()
            Write-Host "Server Response: $responseBody" -ForegroundColor Red
            $reader.Close()
        } catch {
            Write-Host "Could not read error response" -ForegroundColor Red
        }
    }
    
    # Show stack trace for debugging
    Write-Host "Stack Trace:" -ForegroundColor Red
    Write-Host $_.ScriptStackTrace -ForegroundColor Red
}