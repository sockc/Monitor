# Monitor

A lightweight, **Docker-free**, self-hosted Linux server monitor. V0.1 is an initial working MVP.

**Stack:** Go standard library, systemd, embedded dashboard. No Docker, Node.js, Redis, or database service required on hosts.

## 统一安装与管理（V0.9.3 起，推荐）

Server 和 Agent 现在使用同一个中文管理器 **monitorctl**，适用 Linux AMD64、ARM64，均无需 Docker。

首次在目标服务器上运行（Server 或 Agent 使用同一条命令）：

```bash
curl -fsSL https://raw.githubusercontent.com/sockc/Monitor/main/scripts/monitorctl -o /tmp/monitorctl \
  && sudo install -m 0755 /tmp/monitorctl /usr/local/bin/monitorctl \
  && sudo monitorctl
```

进入交互菜单后选择安装 Server 或 Agent。以后在已安装管理器的服务器上直接执行 `sudo monitorctl`。

也支持直接执行命令：

```bash
sudo monitorctl install server
sudo monitorctl install agent
sudo monitorctl upgrade server
sudo monitorctl upgrade agent
monitorctl status
sudo monitorctl logs server
sudo monitorctl logs agent
sudo monitorctl backup
sudo monitorctl uninstall server
sudo monitorctl uninstall agent
sudo monitorctl self-update
```

Server 初次安装会提示监听端口（默认 **127.0.0.1:8090**），并在
`/etc/monitor/server.env` 中一次性生成独立的 Agent 与管理员令牌。
端口只监听本机，需配置 Nginx/Caddy 等 HTTPS 反向代理。
管理员首访使用管理端 `/setup` 页面设置登录账号及密码。
原有服务不能通过 `install` 覆盖，请使用 `upgrade`；升级 Server 自动备份数据库，
通过 SHA256 校验发布二进制并检查健康状态，失败时自动回滚程序。
原有端口、令牌和数据库不会被安装管理器改写。

Agent 首次安装会提示：**Monitor Server HTTPS 地址**、节点标识与节点专属令牌。
先在 Web 面板的「设置 → 节点管理 → 添加节点」生成该令牌。
输入时令牌不回显，安装后保存在 `/etc/monitor/agent.env`（0600）。
升级 Agent 会保留以上配置。**卸载默认仅移除 systemd 服务，保留配置、数据库、备份及程序**，
便于重新安装或恢复，卸载需要人工输入 REMOVE 确认。

新版本发布后先运行 `sudo monitorctl self-update` 更新管理器以获取新的默认版本；
也可以通过 `sudo env MONITOR_VERSION=v0.9.3 monitorctl upgrade server` 指定版本。
若同一机器同时运行 Server 和 Agent，升级时优先备份 Server 并同步重启 Agent。
升级前建议检查 [Releases](https://github.com/sockc/Monitor/releases)。

## V0.9.7：修复国旗、按需显示套餐与 IPv4/IPv6 自动检测

- 卡片上的国旗现在由 Server 本地提供 **SVG 图标**（GB、US、NL、HK 等 30 个常见地区；其他国家代码使用中性代码图形），解决部分 Windows 系统缺少国旗 Emoji 导致只出现 `GB`、`US` 的问题。素材来自 [flag-icons](https://github.com/lipis/flag-icons)，MIT 授权声明见 `third_party/flag-icons-LICENSE`
- 未设置的价格、到期日期、带宽端口、流量额度、国家代码和 IPv4／IPv6 信息，在卡片上**直接隐藏对应区块**；填写后才显示。累计上传／下载改为同一行
- Agent 独立使用 `api4.ipify.org`（强制 TCP/IPv4）与 `api6.ipify.org`（强制 TCP/IPv6）检测公网**出站** IP，启动时开始后台检测，每 12 小时刷新一次，不阻塞原 5 秒监控上报。任一地址检测失败视为“未知”，不会武断显示 IPv4/IPv6 不支持
- 服务器详情会显示分别检测的公网 IPv4、IPv6 地址。节点管理原来的三态设置仍保留：手动“支持”优先显示标签，手动“不支持”优先隐藏标签，“自动检测／未知时隐藏”则只在检测成功时显示
- API 保持兼容旧 Agent，节点原数据及流量不会重置；新 Agent 无需重新安装或重新填写令牌
- 自动 IPv4／IPv6 检测默认开启，若需要完全关闭第三方 IP 查询，可在 Agent 环境配置 `MONITOR_IP_DETECTION_ENABLED=0` 和 `MONITOR_GEO_ENABLED=0`，再重启服务。关闭后不会继续发起新查询
- 检测到的地址为 **VPS 的公网出口 IP**，不保证对应入站连通性，代理或特殊网络路径可能造成差异。查询所用的第三方 HTTPS 服务会看到出口 IP

升级：先在 Server 执行 `sudo monitorctl self-update && sudo monitorctl upgrade server`，再在每台 Agent 执行 `sudo monitorctl self-update && sudo monitorctl upgrade agent`。

## V0.9.6：参考图同款服务器信息卡片

- 首页玻璃卡片重排为：**在线状态／国旗／名称** → **价格和剩余天数** → **CPU、内存、磁盘、实时上传、实时下载五列** → **累计上传／累计下载** → **带宽口／月流量／IPv4／IPv6 标签** → **月流量使用进度**；保持宽屏每行 4 张卡
- 节点设置「基本资料」新增服务商、国家代码（两位字母）、价格、币种、月／季／年／一次性计费周期、带宽端口（Mbps）、IPv4 和 IPv6 支持状态。端口是套餐宣称值，**不代表实测速度**；IPv4/IPv6 支持是手动标注的三态（未确认／支持／不支持），不臆测
- 国旗优先使用手动国家代码；若留空，则使用新版 Agent 经 IP 定位查询并上报的自动国家代码（可能受代理影响）
- 套餐表单还可设置本期开始日期、到期日期和每月流量额度；后两项直接使用原有 `node_limits`，设置界面与原「流量与到期」页面保持同步
- 到期剩余天数来自真实到期日期；**到期进度条必须填写本期开始日期**，否则只显示剩余天数。套餐流量进度基于 Server 已采集的分时流量，并非运营商官方计费流量
- 累计上传和累计下载为 Agent 网卡自本次启动以来的累计字节计数，系统重启后可能归零；底部月流量为 Server 保存的分时累计，二者口径不同
- SQLite 自动新增 `node_profile` 并为 `node_geo` 补充国家代码列；不删除或覆盖既有节点资料、令牌、历史和设置
- 升级：Server 执行 `sudo monitorctl self-update && sudo monitorctl upgrade server`；每台 Agent 执行 `sudo monitorctl self-update && sudo monitorctl upgrade agent`。旧 Agent 能继续监控，但无法自动上报国家代码

## V0.9.5：玻璃主题与公网 IP 自动定位

- **默认主题**：参考提供的截图采用粉紫柔和背景与半透明磨砂玻璃卡片；顶栏「主题」及「设置 → 外观主题」均可切换梦幻玻璃、经典深色、简约浅色
- 浏览器使用 `localStorage` 记住选择，不同步到 Server；若浏览器禁止存储，默认仍为玻璃主题
- **自动位置**：新版 Agent 通过 `https://ipwho.is/?lang=zh-CN` 获取自身请求的公网出口 IP 及近似国家/城市，最多每 12 小时更新一次；失败后 1 小时再尝试。不会将反代 Nginx/Cloudflare 的地址当作 Agent 所在位置
- 位置存于 Server 的 SQLite `node_geo` 表；原有节点位置 `node_metadata.location` 视为**手动覆盖**，优先级高于自动识别。清空手动值并保存后可恢复自动显示
- 服务器位置基于公网出口 IP，不保证真实机房地址：代理、NAT、隧道、运营商 IP 库可能引起偏差
- Agent 请求第三方 geolocation API 时，第三方可获知 Agent 的公网出口 IP。需要关闭此功能时在 `/etc/monitor/agent.env` 加入 `MONITOR_GEO_ENABLED=0`，再运行 `sudo systemctl restart monitor-agent`。关闭后已缓存的历史自动位置不会被新数据更新；管理员仍可手动设置位置
- 升级流程：在 Server 运行 `sudo monitorctl self-update && sudo monitorctl upgrade server`，然后在**每台 Agent** 上运行 `sudo monitorctl self-update && sudo monitorctl upgrade agent`。旧 Agent 保持在线但不会上报自动位置
- 截图展示了参考界面，并非无界面原始壁纸；当前默认主题还原其色彩与玻璃质感，不直接把带文字的截图作为背景

## V0.9.4：四列服务器驾驶舱与位置管理

- 桌面宽屏（1250px 起）每行 4 张服务器卡片，中屏 2 张、手机 1 张；默认保持紧凑卡片，也可切换标准或列表
- 首页总览改为扁平统计条，搜索、分组、排序、布局切换和添加节点整合到一行工具栏
- 卡片显示服务器名称、**所在国家/城市（自定义位置）**、系统、CPU 核心/内存/磁盘、资源进度、实时速率、今日/本月流量
- 网络计数累计值继续在详情页查看；临近到期、额度使用过高、资源高占用才显示提醒
- 登录后到「设置 → 节点管理 → 基本资料 → 服务器位置」填写位置，例「中国香港 · 沙田」「美国 · 洛杉矶」。不使用 IP 自动定位，避免代理和云服务出口造成误判
- Server 首次启动自动为已有 SQLite 节点资料增加 `location` 字段，保留已有名称、分组、备注和历史采样；旧 Agent 继续兼容
- 通过 `sudo monitorctl self-update && sudo monitorctl upgrade server` 升级 Server，Agent 无需升级

## Features

- Linux CPU, memory, root filesystem, total network traffic, uptime
- Multiple named agents; outbound HTTPS reporting every 5 seconds
- Private browser dashboard with admin token login and 5-second auto-refresh
- Server-side agent bearer-token authentication
- Atomic JSON snapshot persistence (every 30 seconds); **history charts and SQLite are planned for future releases**
- GitHub Actions builds for linux/amd64, linux/arm64; server and agent use the same binary with different mode flags

## Run server

Download the matching binary from a GitHub Actions build / release (when available), or build on a development machine:

```sh
go build -trimpath -ldflags="-s -w" -o monitor ./cmd/monitor
MONITOR_AGENT_TOKEN="$(openssl rand -hex 32)"  # save this token for the agents
MONITOR_ADMIN_TOKEN="$(openssl rand -hex 32)"  # save separately
sudo mkdir -p /var/lib/monitor
sudo env MONITOR_AGENT_TOKEN="$MONITOR_AGENT_TOKEN" MONITOR_ADMIN_TOKEN="$MONITOR_ADMIN_TOKEN" ./monitor -mode server -listen 127.0.0.1:8090
```

Place an HTTPS reverse proxy (Nginx/Caddy) in front of `127.0.0.1:8090`. Do not expose the HTTP port to the public internet. Open the proxy URL in your browser and log in using the admin token.

## Run agent

```sh
MONITOR_AGENT_TOKEN="<the server's agent token>" ./monitor -mode agent -server https://monitor.example.com -name server-1
```

**Security:** Use a dedicated long random ingestion token and a different admin token. The MVP uses one shared ingestion token across agents: rotating it affects all agents. Cookies contain the admin token and must be protected by HTTPS. Use a private deployment first; per-node credentials, session hardening, encrypted history, and monitoring alerts are future work.

## systemd

See `deploy/systemd` for example units and `scripts/install.sh` for installation guidance. Review local secrets before enabling services.

## Roadmap

- V0.1: Agent reporting, live dashboard, simple persistence, release automation
- V0.2: SQLite, history retention and charts, per-agent token provisioning
- V0.3: HTTP/TCP checks and alerting
- V0.4: safe signed upgrades, backups, detailed access control

## License

MIT.


## V0.2: history and node-scoped tokens

Server now stores samples in SQLite at `/var/lib/monitor/monitor.db` (WAL mode, rolling 30-day history). Use the dashboard's history range and metric selectors for charts. Migration preserves the old snapshot if present, while the new SQLite database becomes the main persistent metric history.

To create or rotate a **per-node** token, log into the admin dashboard first, then send:

```sh
curl -b cookies.txt -X POST https://monitor.example.com/api/v1/tokens -H 'Content-Type: application/json' -d '{"name":"server-1"}'
```

The endpoint requires the authenticated `monitor_session` cookie; tokens are only returned once. The legacy shared `MONITOR_AGENT_TOKEN` remains a fallback **only for nodes without an issued token**. After every node has its own token, disable the shared fallback in a future breaking change. Each token is scoped to its exact agent name; give the agent `-name server-1` when using its token. Ensure transport uses HTTPS.

## V0.2 installation from tagged GitHub Release

After [the V0.2 release](https://github.com/sockc/Monitor/releases/tag/v0.2.0) is published, on Linux AMD64/ARM64:

```sh
curl -fsSL https://raw.githubusercontent.com/sockc/Monitor/main/scripts/install-release.sh -o install-monitor.sh
sudo bash install-monitor.sh server
```

For an agent, set `MONITOR_SERVER` and `MONITOR_AGENT_TOKEN` in the environment before installation, or configure `/etc/monitor/agent.env` before restarting the service. The server installs to localhost-only port 8090; configure HTTPS reverse proxy separately.

Server credentials: `sudo cat /etc/monitor/server.env`. Treat them as secrets. Use `journalctl -u monitor-server -f` / `journalctl -u monitor-agent -f` for diagnostics.


## V0.4

Node management provides rename, revoke, delete. Rename changes the identity and immediately revokes its old token; register and reconfigure the Agent using the newly named identity. Deletion removes historical samples and leaves a revoked token tombstone to block legacy shared-token fallback. Server state backup uses SQLite online VACUUM INTO; see scripts/backup.sh. Updates are available through scripts/upgrade.sh and verify published SHA256 checksums. Verify your deployment before enabling automatic scheduled updates.

## V0.5: administrator password login

Existing installations retain `MONITOR_ADMIN_TOKEN` only as a **one-time setup verification secret**. On first visit after upgrade, go to `/setup` to create a username (3–64 ASCII letters/digits/-/_) and a password (12–72 bytes). Enter the original admin token from `/etc/monitor/server.env`. A new installation creates its token in this file. Once configured, `/login` accepts only username and password; token login is disabled.

Passwords use bcrypt. Sessions use random 256-bit secrets; their SHA-256 hashes are stored in SQLite with 7-day expiry, HttpOnly and SameSite Strict cookies. After a password change, all existing sessions are invalidated. Five failed attempts per remote IP lead to a 15-minute temporary lockout (in-memory counter, resets after restart). Use HTTPS to protect passwords and tokens; the server remains bound to localhost by default.

To upgrade an existing VPS, **after the V0.5.0 Release becomes available**, download `scripts/backup.sh` and `scripts/upgrade.sh` to the *same* folder, and invoke `sudo bash upgrade.sh server`. The backup helper requires `python3`. The updater verifies the SHA256 checksum and checks systemd activation and server `/healthz`, restoring the old binary if these fail.

Recovery from a forgotten password is not yet automated. Keep the backup and console access; do not delete the existing admin credential before completing setup.
