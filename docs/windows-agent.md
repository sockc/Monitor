# Windows Agent (V0.9.16)

Windows x64 Agent reports to the **existing Linux Monitor Server**. Windows Server mode is not a supported deployment.

## Requirements

- Windows 10/11 x64 or Windows Server 2019/2022/2025 x64
- Administrator PowerShell 5.1+; outbound HTTPS (443) access to the Monitor Server and GitHub Releases
- Monitor Server upgraded to V0.9.16 to show the Windows install option

## First installation

1. In Monitor Web -> Settings -> Nodes -> Add node, choose **Windows (x64)**.
2. Generate and copy the Windows-specific command.
3. Open PowerShell **as Administrator** on the Windows machine, paste and run.
4. Wait for the installation to report the **MonitorAgent Running** service, then refresh the Monitor dashboard.

The generated command includes the new node's token. **Never share screenshots or terminal histories containing this command.** The installer writes credentials to `C:\ProgramData\Monitor\agent.env` and restricts the file to SYSTEM and local Administrators. The long-running Windows service command line contains **only the path to that file**, not the token.

## Upgrade / status / uninstall

Retrieve the official PowerShell script from `https://raw.githubusercontent.com/sockc/Monitor/main/scripts/install-windows.ps1` and save it as `install-windows.ps1`. In elevated PowerShell:

```powershell
.\install-windows.ps1 -Action Upgrade
.\install-windows.ps1 -Action Status
.\install-windows.ps1 -Action Uninstall
```

Upgrades preserve the existing credentials, identity, and historical data. Uninstall stops/removes the service but retains installed files and credentials by default. Deleting a node in the Monitor Web UI is separate and **can delete monitoring history**.

## Agent metrics in V0.9.16

- CPU utilization: Win32 GetSystemTimes delta
- Physical memory: GlobalMemoryStatusEx
- System drive (normally C:) usage: GetDiskFreeSpaceExW
- Network cumulative Ethernet bytes: Windows `netstat -e`; this is *not per-interface* and may omit virtual adapters
- System uptime: GetTickCount64
- Hostname, CPU model, architecture, and outgoing IPv4/IPv6 detection

Windows load averages and total disk I/O are not reported by this initial agent. Separate volumes (D:/E:) and per-adapter counters are planned for a later release.

## Troubleshooting

```powershell
Get-Service MonitorAgent
Get-CimInstance Win32_Service -Filter "Name='MonitorAgent'" | Select Name,State,PathName
Test-NetConnection your-monitor-domain.example -Port 443
```

If the service fails to start, check **Event Viewer -> Windows Logs -> System** for Service Control Manager errors. The service's credential file and downloaded binary are separate. Windows Agent uses the same HTTPS ingest API as Linux and does not require any inbound firewall rule.

The CI workflow validates PowerShell syntax and cross-compiles the Windows EXE, but the first Windows release still warrants a real-machine smoke test on both desktop and Windows Server.
