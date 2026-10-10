package main

const loginHTML=`<!doctype html><html lang="zh"><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Monitor 登录</title><style>body{background:#0b1220;color:#e5ecf7;font:16px system-ui;display:grid;place-items:center;min-height:90vh}form{display:grid;gap:16px;background:#172339;padding:32px;border-radius:16px;width:min(340px,80vw)}input,button{padding:13px;border-radius:9px;border:1px solid #52617a;background:#0b1220;color:white}button{background:#2d77ca;cursor:pointer}</style><form method="post"><h2>Monitor 管理登录</h2><input type="password" name="token" placeholder="管理员令牌" required autocomplete="current-password"><button>登录</button></form></html>`
const dashboardHTML=`<!doctype html><html lang="zh"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1"><title>Monitor</title><link rel="stylesheet" href="/static/style.css"></head><body><header><h2>◉ Monitor</h2><nav class="top-nav"><button type="button" class="nav-btn selected" data-view="overview">概览</button><button type="button" class="nav-btn" data-view="statistics">统计</button><button type="button" class="nav-btn" data-view="alerts">告警</button><button type="button" class="nav-btn" data-view="settings">设置</button></nav><div class="header-tools"><label class="theme-switch" for="theme-select"><span aria-hidden="true">✦</span><select id="theme-select" aria-label="页面主题"><option value="glass">梦幻玻璃</option><option value="dark">经典深色</option><option value="light">简约浅色</option></select></label><form method="post" action="/logout"><button class="logout-btn">退出</button></form></div></header><main><section class="view" id="view-overview">
 <div class="page-heading dashboard-heading"><div><span class="eyebrow">MONITOR · DASHBOARD</span><h1>服务器概览</h1><p id="overview-hint" class="muted">节点实时状态与资源使用</p></div><div class="dashboard-updated"><span class="refresh-dot"></span><span id="updated">加载中</span></div></div>
 <div class="stats summary-strip">
  <div><span class="summary-icon">▣</span><div><small>总节点</small><strong id="total">—</strong></div></div>
  <div><span class="summary-icon">●</span><div><small>在线节点</small><strong id="online">—</strong></div></div>
  <div><span class="summary-icon">◌</span><div><small>离线节点</small><strong id="offline">—</strong></div></div>
  <div><span class="summary-icon">◇</span><div><small>当前告警</small><strong id="alert-count">—</strong></div></div>
 </div>
 <div class="dashboard-toolbar"><div class="toolbar-label"><h3>服务器列表</h3><span id="node-visible-count" class="muted">—</span></div>
  <div class="dashboard-filters">
   <input id="node-search" type="search" aria-label="搜索服务器" placeholder="搜索名称、地区、分组">
   <select id="node-group" aria-label="筛选服务器分组"><option value="">全部分组</option></select>
   <select id="node-order" aria-label="排序方式"><option value="name">名称排序</option><option value="offline">离线优先</option><option value="cpu">CPU 由高到低</option></select>
   <div class="layout-switch" role="group" aria-label="服务器显示方式"><button type="button" class="layout-button" data-layout="cards">标准</button><button type="button" class="layout-button active" data-layout="compact">紧凑</button><button type="button" class="layout-button" data-layout="list">列表</button></div>
   <button type="button" class="primary-btn" id="add-node">＋ 添加节点</button>
  </div>
 </div>
 <section id="nodes" aria-label="服务器列表"></section>
</section><section class="view" id="view-detail" hidden>
 <div class="detail-heading detail-top"><button type="button" id="detail-back" class="subtle-btn">← 返回列表</button><div class="detail-identity"><span class="eyebrow">节点详情</span><h1 id="detail-title">服务器详情</h1><p id="detail-subtitle" class="muted"></p></div><div id="detail-status" class="detail-live muted"></div></div>
 <div class="section-tabs" id="detail-tabs" role="tablist" aria-label="服务器详情分类"><button type="button" class="detail-tab selected" data-detail-tab="performance" role="tab" aria-selected="true" aria-controls="detail-performance">性能</button><button type="button" class="detail-tab" data-detail-tab="network" role="tab" aria-selected="false" aria-controls="detail-network">网络</button><button type="button" class="detail-tab" data-detail-tab="system" role="tab" aria-selected="false" aria-controls="detail-system">系统</button></div>
 <div class="detail-pane" id="detail-performance" role="tabpanel">
  <div class="detail-metrics" id="detail-metrics"></div>
  <section class="history-panel ui-panel"><div class="panel-title"><div><h3>历史资源趋势</h3><p class="muted">CPU、内存及磁盘的使用率变化</p></div><select id="detail-hours" aria-label="趋势范围"><option value="24">24 小时</option><option value="168">7 天</option><option value="720">30 天</option></select></div><canvas id="detail-chart" width="900" height="240" aria-label="CPU、内存和磁盘趋势图"></canvas><p id="detail-chart-status" class="muted"></p></section>
 </div>
 <div class="detail-pane" id="detail-network" role="tabpanel" hidden>
  <section class="history-panel ui-panel"><div class="panel-title"><div><h3>网络概览</h3><p class="muted">实时速率、流量及统计口径</p></div></div><div id="detail-network-info" class="detail-data-grid"></div></section>
  <section class="history-panel ui-panel"><div class="panel-title"><h3>历史网络流量</h3></div><p id="detail-traffic" class="muted"></p><p class="field-hint">历史流量为采样估算，可能与服务商计费数据不同。</p></section>
 </div>
 <div class="detail-pane" id="detail-system" role="tabpanel" hidden>
  <section class="history-panel ui-panel"><div class="panel-title"><div><h3>系统与硬件</h3><p class="muted">操作系统、主机、Agent 与资源配置</p></div></div><div id="detail-info" class="detail-data-grid"></div></section>
 </div>
</section>
<section class="view" id="view-alerts" hidden>
 <div class="section-heading"><div><span class="eyebrow">系统事件</span><h1>告警记录</h1><p class="muted">查看触发中的问题与已经恢复的事件</p></div></div>
 <section class="history-panel ui-panel"><div class="panel-title"><div><h3>近期事件</h3><p class="muted">按最新时间排序</p></div><button type="button" id="refresh-alerts" class="subtle-btn">刷新</button></div><div id="alert-events" class="event-list"></div></section>
</section>
<section class="view" id="view-settings" hidden>
 <div class="section-heading"><div><span class="eyebrow">Monitor / 管理中心</span><h1>系统设置</h1><p class="muted">按任务分类管理服务器、流量提醒与安全选项</p></div></div>
 <div class="settings-shell">
  <nav class="settings-menu" role="tablist" aria-label="设置分类">
   <button type="button" class="settings-tab selected" data-settings-tab="nodes" role="tab" aria-selected="true" aria-controls="settings-nodes"><span class="menu-icon">▣</span><span class="menu-copy"><strong>节点管理</strong><small>名称、分组、添加和删除</small></span></button>
   <button type="button" class="settings-tab" data-settings-tab="limits" role="tab" aria-selected="false" aria-controls="settings-limits"><span class="menu-icon">◷</span><span class="menu-copy"><strong>流量与到期</strong><small>额度、时区和提醒</small></span></button>
   <button type="button" class="settings-tab" data-settings-tab="alerts" role="tab" aria-selected="false" aria-controls="settings-alerts"><span class="menu-icon">♢</span><span class="menu-copy"><strong>告警规则</strong><small>阈值、离线和通知</small></span></button>
   <button type="button" class="settings-tab" data-settings-tab="appearance" role="tab" aria-selected="false" aria-controls="settings-appearance"><span class="menu-icon">☼</span><span class="menu-copy"><strong>外观主题</strong><small>玻璃、深色、浅色</small></span></button>
   <button type="button" class="settings-tab" data-settings-tab="account" role="tab" aria-selected="false" aria-controls="settings-account"><span class="menu-icon">◈</span><span class="menu-copy"><strong>账户安全</strong><small>管理员密码</small></span></button>
  </nav>
  <div class="settings-content">
   <div class="settings-pane" id="settings-nodes" role="tabpanel">
    <div class="settings-pane-heading"><div><h2>节点管理</h2><p class="muted">维护名称、分组和备注，或生成新的 Agent 安装命令</p></div><button type="button" id="show-add-node" class="primary-btn">＋ 添加节点</button></div>
    <section class="settings-card"><div class="settings-card-heading"><h3>基本资料</h3><p class="muted">仅修改监控面板中的显示信息，不会重命名系统主机</p></div>
     <form id="metadata-form" class="settings-form">
      <label class="form-field form-wide"><span>选择服务器</span><select id="metadata-name" aria-label="选择服务器"></select></label>
      <label class="form-field"><span>显示名称</span><input id="metadata-display" maxlength="60" placeholder="例如 香港主服务器"></label>
      <label class="form-field"><span>所属分组</span><input id="metadata-group" maxlength="40" placeholder="例如 香港"></label><label class="form-field"><span>服务器位置（手动覆盖）</span><input id="metadata-location" maxlength="80" placeholder="留空使用 IP 自动定位" autocomplete="off"><small id="metadata-auto-location">等待 Agent 获取公网 IP 位置</small></label>
      <label class="form-field form-wide"><span>备注信息</span><textarea id="metadata-notes" maxlength="500" rows="3" placeholder="可选：用途、配置说明或运营商"></textarea></label>
      <div class="form-actions"><button type="submit" class="primary-btn">保存基本资料</button><span id="metadata-status" class="form-feedback" role="status"></span></div>
     </form>
    </section>
    <section class="settings-card settings-danger"><div class="settings-card-heading"><h3>节点维护</h3><p class="muted">重命名、撤销令牌或删除节点，请确认操作对象</p></div>
     <div class="settings-form"><label class="form-field"><span>操作节点</span><select id="manage-name" aria-label="管理节点"></select></label><label class="form-field"><span>新节点标识（仅用于重命名）</span><input id="rename-target" placeholder="英文、数字、下划线、横线"></label>
     <div class="form-actions"><button type="button" id="rename-node" class="subtle-btn">重命名</button><button type="button" id="revoke-node" class="subtle-btn">撤销令牌</button><button type="button" id="delete-node" class="danger-btn">删除节点</button></div></div>
     <p class="field-hint">删除会清理历史采样且拒绝旧令牌；重命名后需要重新配置 Agent。</p><p id="manage-status" class="form-feedback" role="status"></p>
    </section>
    <section class="settings-card" id="add-panel" hidden><div class="settings-card-heading"><h3>添加服务器 / 轮换令牌</h3><button type="button" id="hide-add-node" class="subtle-btn">收起</button></div><p class="field-hint">节点标识仅支持字母、数字、横线和下划线。对现有节点重新生成令牌，会立即使旧 Agent 认证失效。</p>
     <form id="create-node" class="settings-form"><label class="form-field"><span>节点标识</span><input id="node-name" required maxlength="64" pattern="[A-Za-z0-9_-]+" placeholder="例如 hk-01"></label><label class="form-field"><span>Monitor 服务器地址</span><input id="server-url" placeholder="https://monitor.example.com" type="url" required></label><div class="form-actions"><button type="submit" class="primary-btn">生成安装命令</button><span id="create-status" class="form-feedback" role="status"></span></div></form>
     <label class="form-field"><span>专属安装命令</span><textarea id="install-command" rows="3" readonly spellcheck="false" aria-label="安装命令"></textarea></label><button type="button" id="copy-command" class="subtle-btn">复制命令</button><p class="field-hint">命令包含一次性显示的认证令牌，不要公开分享；执行后建议清理终端历史记录。</p>
    </section>
   </div>
   <div class="settings-pane" id="settings-limits" role="tabpanel" hidden>
    <div class="settings-pane-heading"><div><h2>流量与到期提醒</h2><p class="muted">按服务器分别设置统计时区、月额度及续费提醒</p></div></div>
    <section class="settings-card"><div class="settings-card-heading"><h3>节点流量规则</h3><p class="muted">流量汇总按采样计数，无法补齐掉线期间的数据</p></div>
     <form id="node-limits-form" class="settings-form">
      <label class="form-field form-wide"><span>选择服务器</span><select id="limit-name" aria-label="服务器节点"></select></label>
      <label class="form-field"><span>统计时区</span><input id="limit-timezone" list="monitor-timezones" maxlength="64" placeholder="Asia/Shanghai" required><small>决定今日与本月的日期边界</small></label>
      <datalist id="monitor-timezones"><option value="UTC"><option value="Asia/Shanghai"><option value="Asia/Hong_Kong"><option value="Asia/Tokyo"><option value="Europe/Amsterdam"><option value="Europe/London"><option value="America/Los_Angeles"><option value="America/New_York"></datalist>
      <label class="form-field"><span>月流量额度（GB）</span><input id="limit-quota" type="number" min="0" max="1000000" step="0.1" value="0"><small>填 0 表示不启用额度告警</small></label>
      <label class="form-field"><span>到期日期</span><input id="limit-expires" type="date"><small>留空则不提醒到期</small></label>
      <div class="form-actions"><button type="submit" class="primary-btn">保存流量与到期设置</button><span id="limit-status" class="form-feedback" role="status"></span></div>
     </form>
    </section>
    <div class="settings-note"><strong>提醒触发条件</strong><p>本月流量达到 80%、90%、100%；到期前 30、15、7 天。统计流量并非云服务商的计费数据。</p></div>
   </div>
   <div class="settings-pane" id="settings-alerts" role="tabpanel" hidden>
    <div class="settings-pane-heading"><div><h2>告警规则</h2><p class="muted">调整触发阈值与 Webhook 通知地址</p></div></div>
    <section class="settings-card"><div class="settings-card-heading"><h3>资源与离线阈值</h3><p class="muted">资源异常达到持续时间后触发，避免瞬时波动误报</p></div>
     <form id="alerts-form" class="settings-form">
      <label class="form-field"><span>离线判定（秒）</span><input id="a-offline" type="number" min="30" max="3600"></label>
      <label class="form-field"><span>持续时间（秒）</span><input id="a-duration" type="number" min="30" max="3600"></label>
      <label class="form-field"><span>CPU 告警阈值（%）</span><input id="a-cpu" type="number" min="1" max="100"></label>
      <label class="form-field"><span>内存告警阈值（%）</span><input id="a-memory" type="number" min="1" max="100"></label>
      <label class="form-field"><span>磁盘告警阈值（%）</span><input id="a-disk" type="number" min="1" max="100"></label>
      <label class="form-field form-wide"><span>Webhook 通知地址</span><input id="a-webhook" type="url" placeholder="https://example.com/webhook"><small>可选，仅允许 HTTPS 地址</small></label>
      <div class="form-actions"><button type="submit" class="primary-btn">保存告警规则</button><span id="alerts-status" class="form-feedback" role="status"></span></div>
     </form>
    </section>
   </div>
   <div class="settings-pane" id="settings-appearance" role="tabpanel" hidden>
    <div class="settings-pane-heading"><div><h2>外观主题</h2><p class="muted">主题只影响当前浏览器，不修改服务器运行配置</p></div></div>
    <section class="settings-card"><div class="settings-card-heading"><h3>选择面板外观</h3><p class="muted">默认使用参考图中的柔和背景和磨砂玻璃卡片风格</p></div>
      <div class="theme-options" role="group" aria-label="选择监控面板主题">
       <button class="theme-choice" type="button" data-theme-choice="glass"><span class="theme-thumb theme-thumb-glass"></span><strong>梦幻玻璃</strong><small>默认 · 柔和背景与透明卡片</small></button>
       <button class="theme-choice" type="button" data-theme-choice="dark"><span class="theme-thumb theme-thumb-dark"></span><strong>经典深色</strong><small>高对比度的深色监控面板</small></button>
       <button class="theme-choice" type="button" data-theme-choice="light"><span class="theme-thumb theme-thumb-light"></span><strong>简约浅色</strong><small>清爽亮色，适合白天使用</small></button>
      </div>
      <p id="theme-status" class="field-hint" role="status">主题选择将保存在当前浏览器中。</p>
    </section>
   </div>
   <div class="settings-pane" id="settings-account" role="tabpanel" hidden>
    <div class="settings-pane-heading"><div><h2>账户安全</h2><p class="muted">管理管理员凭据与账户访问</p></div></div>
    <section class="settings-card"><div class="settings-card-heading"><h3>管理员密码</h3><p class="muted">定期更换管理员密码，避免在不受信任的环境中保存令牌</p></div><a class="primary-link" href="/password">修改管理员密码 →</a></section>
   </div>
  </div>
 </div>
</section>
<section class="view" id="view-statistics" hidden>
 <div class="section-heading"><div><span class="eyebrow">趋势与用量</span><h1>历史统计</h1><p class="muted">切换节点和时间范围，分析资源与网络使用情况</p></div></div>
 <section class="history-panel ui-panel"><div class="panel-title"><div><h3>资源使用趋势</h3><p class="muted">CPU、内存、磁盘历史使用率</p></div></div>
  <div class="history-controls stat-filters"><select id="history-name" aria-label="服务器"></select><select id="history-hours" aria-label="范围"><option value="24">24 小时</option><option value="168">7 天</option><option value="720">30 天</option></select><select id="history-metric" aria-label="指标"><option value="cpu">CPU</option><option value="memory">内存</option><option value="disk">磁盘</option></select></div>
  <canvas id="history-chart" width="900" height="260" aria-label="历史趋势图"></canvas><p id="history-status" class="muted"></p>
 </section>
 <section class="history-panel ui-panel"><div class="panel-title"><div><h3>历史网络流量</h3><p class="muted">所选时间范围内的接收与发送用量</p></div></div><p id="traffic-summary" class="muted">选择节点和时间范围查看接收、发送量。</p><div id="traffic-bars"></div></section>
</section><footer class="muted">Monitor · 每 5 秒更新</footer></main><script src="/static/app.js" defer></script></body></html>`
