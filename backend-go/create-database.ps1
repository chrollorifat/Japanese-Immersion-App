# PostgreSQL Database Creation Script
# This script helps create the database for your Japanese Learning App

Write-Host "🐘 PostgreSQL Database Setup for Japanese Learning App" -ForegroundColor Green
Write-Host ""

# PostgreSQL paths
$PG_PATH = "C:\Program Files\PostgreSQL\17\bin"
$PSQL_EXE = "$PG_PATH\psql.exe"
$CREATEDB_EXE = "$PG_PATH\createdb.exe"

# Check if PostgreSQL is installed
if (-not (Test-Path $PSQL_EXE)) {
    Write-Host "❌ PostgreSQL not found at expected location!" -ForegroundColor Red
    Write-Host "Please check your PostgreSQL installation." -ForegroundColor Red
    exit 1
}

Write-Host "✅ PostgreSQL found!" -ForegroundColor Green
Write-Host ""

# Get password securely
Write-Host "Please enter your PostgreSQL password (the one you set during installation):"
$password = Read-Host -AsSecureString
$plainPassword = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))

Write-Host ""
Write-Host "🔨 Creating database 'japanese_learning_app'..." -ForegroundColor Yellow

# Set environment variable for password
$env:PGPASSWORD = $plainPassword

try {
    # Create the database
    & $CREATEDB_EXE -U postgres -h localhost japanese_learning_app
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ Database 'japanese_learning_app' created successfully!" -ForegroundColor Green
        Write-Host ""
        Write-Host "🎯 Next steps:" -ForegroundColor Cyan
        Write-Host "1. Edit your .env file with the correct password" -ForegroundColor White
        Write-Host "2. Run: go run main.go" -ForegroundColor White
        Write-Host ""
    } else {
        Write-Host "❌ Failed to create database. Please check your password." -ForegroundColor Red
    }
} catch {
    Write-Host "❌ Error: $_" -ForegroundColor Red
} finally {
    # Clear password from environment
    $env:PGPASSWORD = $null
}

Write-Host "Press any key to continue..."
$null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
