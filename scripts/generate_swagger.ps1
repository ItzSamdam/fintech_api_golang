# scripts/generate_swagger.ps1
Import-Module powershell-yaml

# Base Swagger object
$swagger = @{
  openapi = "3.0.3"
  info = @{
    title = "Fintech API"
    description = "Fintech API Documentation"
    version = "1.0.0"
  }
  servers = @(@{ url = "http://localhost:5000/api/v1" })
  paths = @{}
}

function Add-Route($paths, $method, $path, $tag) {
    if (-not $paths.ContainsKey($path)) {
        $paths[$path] = @{}
    }
    $paths[$path][$method] = @{
        tags = @($tag)
        summary = "$method $path"
        responses = @{
            "200" = @{ description = "Successful response" }
            "400" = @{ description = "Bad request" }
            "401" = @{ description = "Unauthorized" }
            "500" = @{ description = "Internal server error" }
        }
    }
}

# === Route Lists ===
$authRoutes = @(
  "post:/auth/register/phone","post:/auth/verify/otp","post:/auth/register/bvn","post:/auth/verify/face",
  "post:/auth/login","post:/auth/reset-password","post:/auth/refresh","get:/auth/me","put:/auth/me",
  "post:/auth/change-password","post:/auth/logout"
)

$walletRoutes = @("post:/wallets/create","get:/wallets/balance","get:/wallets/transactions","get:/wallets/limits",
  "post:/wallets/lock","post:/wallets/unlock","get:/wallets/statement")

$transferRoutes = @("post:/transfers/send","post:/transfers/send-to-wallet","get:/transfers/status/{reference}",
  "post:/transfers/retry","get:/transfers/history","get:/banks/list","post:/transfers/name-enquiry")

$airtimeRoutes = @("get:/airtime/networks","get:/airtime/denominations","post:/airtime/purchase","get:/airtime/history")

$dataRoutes = @("get:/data/networks","get:/data/plans/{network}","post:/data/purchase","get:/data/history")

$electricityRoutes = @("get:/electricity/providers","post:/electricity/validate-meter","post:/electricity/pay-prepaid",
  "post:/electricity/pay-postpaid","get:/electricity/token/{transaction_id}","get:/electricity/history")

$bettingRoutes = @("get:/betting/providers","post:/betting/validate-account","post:/betting/fund","get:/betting/history")

$savingsRoutes = @("post:/savings/goals/create","post:/savings/goals/contribute","get:/savings/goals","get:/savings/goals/{id}",
  "put:/savings/goals/{id}","delete:/savings/goals/{id}","post:/savings/roundup/activate","post:/savings/roundup/deactivate",
  "get:/savings/roundup/status")

$transactionRoutes = @("get:/transactions","get:/transactions/{id}","get:/bills/history")

$complianceRoutes = @("post:/compliance/report/suspicious","get:/compliance/limits/check","post:/security/sim-swap/check",
  "post:/security/device/trust","post:/security/2fa/enable","post:/security/2fa/verify","get:/security/sessions",
  "delete:/security/sessions/{id}")

$notificationRoutes = @("get:/notifications/in-app","put:/notifications/{id}/read")

$supportRoutes = @("post:/support/tickets/create","get:/support/tickets","get:/support/tickets/{id}",
  "post:/support/tickets/{id}/reply","put:/support/tickets/{id}/status")

$adminUserRoutes = @("get:/admin/users","get:/admin/users/{id}","post:/admin/users/{id}/tier/upgrade",
  "post:/admin/users/{id}/suspend","post:/admin/users/{id}/unsuspend","delete:/admin/users/{id}",
  "put:/admin/users/{id}/limits","get:/admin/users/search")

$adminTransactionRoutes = @("get:/admin/transactions","get:/admin/transactions/{id}","post:/admin/transactions/reverse",
  "post:/admin/transactions/void","get:/admin/transactions/summary")

$adminWalletRoutes = @("get:/admin/wallets","get:/admin/wallets/{id}","post:/admin/wallets/credit","post:/admin/wallets/debit",
  "post:/admin/wallets/freeze","post:/admin/wallets/unfreeze","get:/admin/wallets/balances/summary")

$adminKYCRoutes = @("get:/admin/kyc/pending","get:/admin/kyc/{id}","post:/admin/kyc/{id}/approve","post:/admin/kyc/{id}/reject")

$adminProviderRoutes = @("get:/admin/providers","put:/admin/providers/{id}/toggle","put:/admin/providers/{id}/priority",
  "get:/admin/providers/{id}/health","get:/admin/providers/logs")

$adminFeeRoutes = @("get:/admin/fees","put:/admin/fees/{bill_type}","get:/admin/margins","put:/admin/margins/{provider_id}")

$adminReportRoutes = @("get:/admin/reports/daily","get:/admin/reports/monthly","get:/admin/reports/revenue/by-bill-type",
  "get:/admin/reports/top-users","get:/admin/reports/fraud-attempts","get:/admin/reports/provider-performance",
  "post:/admin/reports/export")

$adminSystemRoutes = @("get:/admin/settings","put:/admin/settings","get:/admin/settings/health","get:/admin/settings/audit-logs",
  "get:/admin/settings/audit-logs/{id}","post:/admin/settings/backup/database","get:/admin/settings/metrics")

$adminRoleRoutes = @("get:/admin/roles","post:/admin/roles","put:/admin/roles/{id}","delete:/admin/roles/{id}",
  "get:/admin/staff","post:/admin/staff/invite","put:/admin/staff/{id}/role","delete:/admin/staff/{id}",
  "get:/admin/staff/{id}/audit")

# === Grouping by Tags ===
$routeGroups = @{
    Authentication = $authRoutes
    Wallet         = $walletRoutes
    Transfer       = $transferRoutes
    Airtime        = $airtimeRoutes
    Data           = $dataRoutes
    Electricity    = $electricityRoutes
    Betting        = $bettingRoutes
    Savings        = $savingsRoutes
    Transactions   = $transactionRoutes
    Compliance     = $complianceRoutes
    Notifications  = $notificationRoutes
    Support        = $supportRoutes
    "Admin Users"  = $adminUserRoutes
    "Admin Transactions" = $adminTransactionRoutes
    "Admin Wallets" = $adminWalletRoutes
    "Admin KYC"    = $adminKYCRoutes
    "Admin Providers" = $adminProviderRoutes
    "Admin Fees"   = $adminFeeRoutes
    "Admin Reports" = $adminReportRoutes
    "Admin Settings" = $adminSystemRoutes
    "Admin Roles"  = $adminRoleRoutes
}

# Add all routes grouped by tag
foreach ($group in $routeGroups.Keys) {
    foreach ($route in $routeGroups[$group]) {
        $parts = $route -split ":", 2
        $method = $parts[0]
        $path = $parts[1]
        Add-Route $swagger.paths $method $path $group
    }
}

# Add health check and root routes
Add-Route $swagger.paths "get" "/health" "Health"
Add-Route $swagger.paths "get" "/" "Root"

# Export to YAML
$swagger | ConvertTo-Yaml | Out-File -FilePath "api/docs/swagger.yaml" -Encoding UTF8

Write-Host "Swagger file generated successfully with $($swagger.paths.Keys.Count) paths!" -ForegroundColor Green
