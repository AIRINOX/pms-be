# Authentication System Test Script
$baseUrl = "http://localhost:3000"
$token = ""

# Note: User registration is now handled by admin only
# This script tests login with existing user credentials

# Test 1: Login with existing user
Write-Host ""
Write-Host "1. Testing User Login..." -ForegroundColor Yellow
$loginData = @{
    username = "testuser"
    password = "password123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$baseUrl/auth/login" -Method POST -Body $loginData -ContentType "application/json"
    Write-Host "✓ Login successful!" -ForegroundColor Green
    Write-Host "Response: $($loginResponse | ConvertTo-Json)" -ForegroundColor Cyan
    
    # Extract token for subsequent requests
    $token = $loginResponse.token
    $headers = @{
        "Authorization" = "Bearer $token"
        "Content-Type" = "application/json"
    }
    
    # Test 2: Access protected endpoint (Get user info)
    Write-Host ""
    Write-Host "2. Testing Protected Endpoint (/auth/me)..." -ForegroundColor Yellow
    try {
        $meResponse = Invoke-RestMethod -Uri "$baseUrl/auth/me" -Method GET -Headers $headers
        Write-Host "✓ Protected endpoint access successful!" -ForegroundColor Green
        Write-Host "User info: $($meResponse | ConvertTo-Json)" -ForegroundColor Cyan
    } catch {
        Write-Host "✗ Protected endpoint access failed: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    # Test 3: Logout
    Write-Host ""
    Write-Host "3. Testing User Logout..." -ForegroundColor Yellow
    try {
        $logoutResponse = Invoke-RestMethod -Uri "$baseUrl/auth/logout" -Method POST -Headers $headers
        Write-Host "✓ Logout successful!" -ForegroundColor Green
        Write-Host "Response: $($logoutResponse | ConvertTo-Json)" -ForegroundColor Cyan
    } catch {
        Write-Host "✗ Logout failed: $($_.Exception.Message)" -ForegroundColor Red
    }
    
} catch {
    Write-Host "✗ Login failed: $($_.Exception.Message)" -ForegroundColor Red
}

# Test 4: Try to access protected endpoint without token (should fail)
Write-Host ""
Write-Host "4. Testing Access Without Token (should fail)..." -ForegroundColor Yellow
try {
    $unauthorizedResponse = Invoke-RestMethod -Uri "$baseUrl/auth/me" -Method GET
    Write-Host "✗ Unauthorized access should have failed!" -ForegroundColor Red
} catch {
    Write-Host "✓ Unauthorized access properly blocked!" -ForegroundColor Green
    Write-Host "Expected error: $($_.Exception.Message)" -ForegroundColor Cyan
}

Write-Host ""
Write-Host "Authentication system testing completed!" -ForegroundColor Green