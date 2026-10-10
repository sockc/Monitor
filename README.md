# Monitor


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

Agent 手动运行 `monitorctl install agent` 时仍可按生成的安装命令填写 Server 地址、系统自动生成的内部 ID 与专属令牌。推荐直接在 Web 面板「设置 → 节点管理 → 添加节点」一键生成完整安装命令，无需手动输入内部 ID。
输入时令牌不回显，安装后保存在 `/etc/monitor/agent.env`（0600）。
升级 Agent 会保留以上配置。**卸载默认仅移除 systemd 服务，保留配置、数据库、备份及程序**，
便于重新安装或恢复，卸载需要人工输入 REMOVE 确认。

新版本发布后先运行 `sudo monitorctl self-update` 更新管理器以获取新的默认版本；
也可以通过 `sudo env MONITOR_VERSION=v0.9.3 monitorctl upgrade server` 指定版本。
若同一机器同时运行 Server 和 Agent，升级时优先备份 Server 并同步重启 Agent。
升级前建议检查 [Releases](https://github.com/sockc/Monitor/releases)。

