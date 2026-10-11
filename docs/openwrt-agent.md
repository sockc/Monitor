# Monitor OpenWrt ARM64 Agent (V0.9.17)

适用：OpenWrt aarch64，使用 procd。独立 Agent，不安装数据库、Web 或 Docker。

## 首次安装

先在 Monitor 网页「设置 → 节点管理 → 添加」创建节点，取得服务器 HTTPS 地址、节点 ID 和 Agent 专属令牌。

路由器执行：

```sh
uname -m
df -h /overlay
wget -O /tmp/install-openwrt.sh https://github.com/sockc/Monitor/releases/download/v0.9.18/install-openwrt.sh
sh /tmp/install-openwrt.sh https://你的-monitor-域名 节点ID
```

首次安装会在终端中提示输入令牌（不回显）。不要将令牌直接写进命令参数或者公开日志。现有配置位于 `/etc/monitor-agent.json`，权限应为 0600。

安装程序必须下载校验 OpenWrt Agent 和管理工具，至少保留约 10 MiB 可写空间。不要用旧安装脚本覆盖正在运行的 Agent。

## 日常维护

```sh
monitor-agentctl status
monitor-agentctl upgrade
monitor-agentctl uninstall
logread -e monitor
```

升级前程序下载至 `/tmp`（内存文件系统），核对发布 SHA256；旧文件临时保存在 `/tmp` 以备启动失败回滚。卸载默认保留配置和令牌，以便恢复；如需删除配置，请确认后自行处理。

## 指标口径

- 使用 netifd/ubus 识别逻辑 `wan` 和 `wan6`，PPPoE 优先读取 `l3_device`。
- WAN 流量统计选中的主 WAN 接口，WAN6 用于识别，不重复叠加网卡计数。
- ubus 不可用时回退读取默认路由，过滤常见 TUN/WireGuard 接口；无法准确识别则不假定为物理 WAN。
- CPU 温度来自 sysfs 的 CPU/SoC 温度传感器；无法取得则不显示。
- 存储占用以 `/overlay` 的可写空间为准。
- 断网时自动延长上报重试间隔，恢复后自动恢复正常采样。
- 设备监控数据经 HTTPS 上传原有 Monitor 服务端；Agent 不维护本地监控数据库。

## 注意

V0.9.17 仍需在不同 OpenWrt 固件、PPPoE 和代理环境下实机验证。若路由器使用自定义网络接口、内存紧张或 `pidof` 行为有差异，请先使用 `status` 和 `logread` 验证服务是否正常。
