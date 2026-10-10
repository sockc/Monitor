# Monitor Windows x64 Agent installer.
# Run from an elevated Windows PowerShell 5.1+ or PowerShell 7 session.
[CmdletBinding()]
param(
 [ValidateSet('Install','Upgrade','Status','Uninstall')][string]$Action='Install',
 [string]$Server='',
 [string]$NodeName='',
 [string]$Token='',
 [string]$Version='v0.9.16.1'
)
Set-StrictMode -Version Latest
$ErrorActionPreference='Stop'
$ProgressPreference='SilentlyContinue'
$service='MonitorAgent'
$programDir=Join-Path $env:ProgramFiles 'Monitor'
$program=Join-Path $programDir 'monitor.exe'
$dataDir=Join-Path $env:ProgramData 'Monitor'
$config=Join-Path $dataDir 'agent.env'

function Require-Admin {
 $identity=[Security.Principal.WindowsIdentity]::GetCurrent()
 $principal=New-Object Security.Principal.WindowsPrincipal($identity)
 if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
  throw '请右键以管理员身份运行 PowerShell 后重新执行。'
 }
}
function Existing-Service {return Get-Service -Name $service -ErrorAction SilentlyContinue}
function Resolve-WindowsSystemTool([string]$FileName) {
 # Resolve native executables by their absolute system path. PATH may omit
 # System32 in restricted PowerShell sessions even on a healthy Windows host.
 $tool=Join-Path ([Environment]::SystemDirectory) $FileName
 if (-not (Test-Path -LiteralPath $tool -PathType Leaf)) {
  throw "找不到 Windows 系统程序：$tool。请确认系统文件完整。"
 }
 return $tool
}
function Check-ConnectionParameters {
 if ([string]::IsNullOrWhiteSpace($Server)) {throw '缺少 Monitor Server HTTPS 地址。'}
 if ($Server -notmatch '^https://[^\s/]+/?$') {throw 'Server 地址必须是 HTTPS 域名或 IP，不能包含路径和账号信息。'}
 $uri=$null
 if (-not [uri]::TryCreate($Server,[UriKind]::Absolute,[ref]$uri) -or $uri.Scheme -ne 'https' -or $uri.UserInfo -or $uri.Query -or $uri.Fragment) {throw 'Server URL 无效。'}
 if ($NodeName -notmatch '^[A-Za-z0-9_-]{1,64}$') {throw '节点 ID 无效。'}
 if (-not $Token) {
  $secret=Read-Host '输入节点专属 Token（不会回显）' -AsSecureString
  $ptr=[Runtime.InteropServices.Marshal]::SecureStringToBSTR($secret)
  try { $script:Token=[Runtime.InteropServices.Marshal]::PtrToStringBSTR($ptr) }
  finally { [Runtime.InteropServices.Marshal]::ZeroFreeBSTR($ptr) }
 }
 if ($Token -notmatch '^[0-9a-fA-F]{64}$') {throw '节点 Token 无效（应为 64 位十六进制）。'}
}
function Write-PrivateConfig {
 $icacls=Resolve-WindowsSystemTool 'icacls.exe'
 New-Item -ItemType Directory -Path $dataDir -Force | Out-Null
 # Protect the directory before writing the secret, including recovery after
 # an earlier failed installation that left a credential file behind.
 & $icacls $dataDir '/inheritance:r' '/grant:r' '*S-1-5-18:(OI)(CI)F' '*S-1-5-32-544:(OI)(CI)F' | Out-Null
 if ($LASTEXITCODE -ne 0) {throw '无法设置 Monitor 配置目录权限，已停止安装。'}
 # No BOM; never log credentials.
 $body="MONITOR_SERVER=$($Server.TrimEnd('/'))`nMONITOR_NODE_NAME=$NodeName`nMONITOR_AGENT_TOKEN=$Token`n"
 try {
  [IO.File]::WriteAllText($config,$body,(New-Object Text.UTF8Encoding($false)))
  # Locale-independent SIDs: LocalSystem and built-in Administrators.
  & $icacls $config '/inheritance:r' '/grant:r' '*S-1-5-18:F' '*S-1-5-32-544:F' | Out-Null
  if ($LASTEXITCODE -ne 0) {throw '无法设置 Agent 凭据文件权限。'}
 } catch {
  # Do not leave readable credentials on disk when ACL hardening fails.
  Remove-Item -LiteralPath $config -Force -ErrorAction SilentlyContinue
  throw
 }
}
function Download-Release([string]$targetDir) {
 $fileName='monitor-windows-amd64.exe'
 $base="https://github.com/sockc/Monitor/releases/download/$Version"
 $binary=Join-Path $targetDir $fileName
 $checksums=Join-Path $targetDir 'SHA256SUMS'
 Write-Host "正在下载 Monitor $Version（Windows x64）..."
 Invoke-WebRequest -UseBasicParsing -TimeoutSec 40 -Uri "$base/$fileName" -OutFile $binary
 Invoke-WebRequest -UseBasicParsing -TimeoutSec 40 -Uri "$base/SHA256SUMS" -OutFile $checksums
 $manifest=[IO.File]::ReadAllText($checksums)
 $expected=[regex]::Match($manifest,'(?im)^([0-9a-f]{64})\s+\*?monitor-windows-amd64\.exe\s*$')
 if (-not $expected.Success) {throw 'SHA256SUMS 未包含 Windows 安装包校验值。'}
 $actual=(Get-FileHash -Algorithm SHA256 -Path $binary).Hash
 if ($actual -ine $expected.Groups[1].Value) {throw 'Windows 程序 SHA256 校验失败，安装已停止。'}
 Write-Host 'SHA256 校验通过。'
 return $binary
}
function Show-Service {
 $s=Existing-Service
 if ($null -eq $s) {Write-Host 'MonitorAgent 未安装。';return}
 Write-Host ("MonitorAgent 服务状态：" + $s.Status)
 Write-Host "程序：$program"
 Write-Host "配置：$config"
}
Require-Admin
if ([Environment]::Is64BitOperatingSystem -ne $true -or $env:PROCESSOR_ARCHITECTURE -eq 'ARM64') {
 throw 'V0.9.16 仅发布 Windows x64 (AMD64) 安装包。'
}
if ($Action -eq 'Status') {Show-Service;exit 0}
if ($Action -eq 'Uninstall') {
 $old=Existing-Service
 if ($null -eq $old) {Write-Host 'MonitorAgent 尚未安装。';exit 0}
 if ($old.Status -ne 'Stopped') {Stop-Service -Name $service -Force -ErrorAction Stop}
 & (Resolve-WindowsSystemTool 'sc.exe') delete $service | Out-Host
 if ($LASTEXITCODE -ne 0) {throw '删除服务失败。'}
 Write-Host '服务已卸载。程序与 Agent 凭据保留，可用于后续恢复；需彻底删除请手动清理 ProgramData\Monitor。'
 exit 0
}
$old=Existing-Service
if ($Action -eq 'Install' -and $null -ne $old) {throw '已安装 MonitorAgent；请使用 -Action Upgrade，避免覆盖原节点身份。'}
if ($Action -eq 'Upgrade' -and $null -eq $old) {throw 'MonitorAgent 尚未安装，无法升级。'}
if ($Action -eq 'Install') {Check-ConnectionParameters}
# Upgrade retains the existing credential file and node identity.
if ($Action -eq 'Upgrade' -and -not (Test-Path $config)) {throw '找不到原来的 Agent 配置文件，为防止节点身份变化已停止升级。'}
$tmp=Join-Path $env:TEMP ("monitor-win-"+[Guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $tmp -Force | Out-Null
$backup=Join-Path $tmp 'previous-monitor.exe'
$replaced=$false
try {
 $binary=Download-Release $tmp
 New-Item -ItemType Directory -Path $programDir -Force | Out-Null
 if ($Action -eq 'Upgrade') {
  if ($old.Status -ne 'Stopped') {Stop-Service -Name $service -Force -ErrorAction Stop}
  if (Test-Path $program) {Copy-Item $program $backup -Force}
 }
 Copy-Item $binary $program -Force
 $replaced=$true
 if ($Action -eq 'Install') {
  Write-PrivateConfig
  $binPath='"'+$program+'" -mode agent -config "'+$config+'"'
  New-Service -Name $service -BinaryPathName $binPath -DisplayName 'Monitor Agent' -StartupType Automatic -Description 'Reports Windows server metrics to Monitor' | Out-Null
 }
 Start-Service -Name $service -ErrorAction Stop
 (Get-Service -Name $service).WaitForStatus('Running',[TimeSpan]::FromSeconds(15))
 Write-Host "Monitor Windows Agent 已安装／升级：$Version"
 Show-Service
} catch {
 if ($Action -eq 'Upgrade' -and $replaced -and (Test-Path $backup)) {
  Write-Warning '升级失败，正在回滚上一版本。'
  try {
   Stop-Service -Name $service -Force -ErrorAction SilentlyContinue
   Copy-Item $backup $program -Force
   Start-Service -Name $service -ErrorAction Stop
  } catch {Write-Warning '自动回滚启动失败，请手工检查 MonitorAgent 服务。'}
 }
 throw
} finally {
 Remove-Item -LiteralPath $tmp -Force -Recurse -ErrorAction SilentlyContinue
}
