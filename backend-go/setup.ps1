Write-Host "🚀 Japanese Learning App - Go Backend Setup" -ForegroundColor Green
Write-Host ""

$envFile = ".env"

Write-Host "📝 Configuring database connection..." -ForegroundColor Yellow
Write-Host "Please enter your PostgreSQL password (the one you set during installation):"
$password = Read-Host -MaskInput

# Update .env file
$envContent = @"
# PostgreSQL Database Configuration
DATABASE_URL=postgres://postgres:$password@localhost:5432/japanese_learning_app?sslmode=disable
PORT=8000
DEBUG=true
JWT_SECRET=your-super-secret-jwt-key-for-development-only
"@

$envContent | Out-File -FilePath $envFile -Encoding UTF8

Write-Host "✅ Configuration saved to .env file!" -ForegroundColor Green
Write-Host ""
Write-Host "🎯 Testing database connection..." -ForegroundColor Yellow

# Test the Go application
try {
    $process = Start-Process -FilePath "go" -ArgumentList "run", "main.go" -NoNewWindow -PassThru -RedirectStandardOutput "output.log" -RedirectStandardError "error.log"
    Start-Sleep -Seconds 3
    
    if (!$process.HasExited) {
        Write-Host "✅ Server started successfully!" -ForegroundColor Green
        Write-Host "🌐 Your API is running at: http://localhost:8000" -ForegroundColor Cyan
        Write-Host "🔍 Health check: http://localhost:8000/health" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "Press any key to stop the server..."
        $null = $Host.UI.RawUI.ReadKey("NoEcho,IncludeKeyDown")
        Stop-Process -Id $process.Id -Force
    } else {
        Write-Host "❌ Server failed to start. Check error.log for details." -ForegroundColor Red
        if (Test-Path "error.log") {
            Write-Host "Error details:" -ForegroundColor Red
            Get-Content "error.log" | Write-Host -ForegroundColor Red
        }
    }
} catch {
    Write-Host "❌ Error starting server: $_" -ForegroundColor Red
}

# Cleanup log files
if (Test-Path "output.log") { Remove-Item "output.log" }
if (Test-Path "error.log") { Remove-Item "error.log" }
