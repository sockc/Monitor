package main

// Monitor stylesheet segment; concatenated in original cascade order.
const styleCSSDashboard=`/* V0.9.2: unified dense desktop dashboard and redesigned management UI */
:root{color-scheme:dark;--ui-bg:#0b1422;--ui-card:#15243a;--ui-card-2:#192a43;--ui-border:#29415d;--ui-text:#edf5ff;--ui-muted:#9db0c9;--ui-blue:#69adff;--ui-green:#46d6ac}
body{background:radial-gradient(ellipse 1100px 420px at 30% -220px,#1b3553 0%,transparent 80%),var(--ui-bg)!important;color:var(--ui-text)!important;line-height:1.48}
header{border-bottom:1px solid #233954;background:#101d30!important;padding:11px clamp(16px,3.5vw,58px)!important;gap:14px}
header h2{font-size:18px!important;letter-spacing:.2px;color:#f4f9ff}
header form button{background:#20334d;border:1px solid #304a69}
.top-nav{background:#0c1728!important;padding:4px!important;gap:4px!important;border:1px solid #22364e}
.top-nav button{font-size:13px!important;padding:8px 15px!important;font-weight:550}
.top-nav button.selected{background:#28496f!important;color:#f5faff;box-shadow:0 1px 4px #0004}
main{max-width:1420px!important;margin:17px auto!important;padding:0 clamp(14px,2.2vw,28px)!important}
.view{min-height:0!important}.view[hidden]{display:none!important}
.section-heading,.page-heading{display:flex;justify-content:space-between;align-items:flex-end;gap:12px;margin:4px 0 18px!important}
.section-heading h1,.page-heading h1,#view-statistics h1{font-size:20px!important;line-height:1.35;font-weight:700!important;margin:2px 0 4px!important;color:#f5f9ff}
.section-heading p,.page-heading p{font-size:12px;margin:3px 0 0;color:var(--ui-muted)}
.eyebrow{font-size:10px;font-weight:700;letter-spacing:1.3px;text-transform:uppercase;color:#74b5ff}
.muted{color:var(--ui-muted)}
button,select,input,textarea{font-family:inherit}
button{transition:background .15s,border-color .15s,box-shadow .15s}
button:focus-visible,select:focus-visible,input:focus-visible,textarea:focus-visible,a:focus-visible{outline:2px solid #77b9ff;outline-offset:2px}
button:hover{filter:brightness(1.1)}
.primary-btn{background:#3278c4;color:#fff;font-weight:650;border:1px solid #4a96eb;border-radius:8px;padding:9px 13px;white-space:nowrap}
.subtle-btn{background:#1c2d44;color:#dce9fc;border:1px solid #36516f;border-radius:8px;padding:8px 12px;white-space:nowrap}
.danger-btn{background:#4a2634;color:#ffb9c6;border:1px solid #854052;border-radius:8px;padding:8px 12px;white-space:nowrap}
.primary-link{display:inline-flex;align-items:center;justify-content:center;padding:10px 14px;background:#2e6faf;border:1px solid #498ddd;color:white;border-radius:8px;text-decoration:none;font-size:13px;font-weight:650}
.primary-link:hover{text-decoration:none;background:#347bc4}
.ui-panel,.history-panel{background:var(--ui-card)!important;border:1px solid var(--ui-border)!important;border-radius:13px!important;box-shadow:0 4px 18px #050b1615}
.panel-title{display:flex;justify-content:space-between;align-items:center;gap:12px;margin-bottom:12px}
.panel-title h3{font-size:15px!important;margin:0!important;color:#f0f6ff}.panel-title p{font-size:11px;margin:3px 0 0}
#view-statistics .history-panel,#view-alerts .history-panel,#view-detail .history-panel{padding:16px 18px!important;margin-top:12px!important}
#history-chart,#detail-chart{max-height:280px}
#history-status,#detail-chart-status,#traffic-summary{font-size:12px}
.stat-filters{display:flex;gap:7px!important;margin:6px 0 14px!important}
.stat-filters select{font-size:12px!important;padding:8px 10px!important;min-width:120px}
#traffic-bars{margin-top:15px}
#traffic-bars .item{display:flex;gap:12px;margin:9px 0!important;font-size:12px!important}
#traffic-bars .track{height:8px!important;background:#283e56}
#traffic-bars .value{background:linear-gradient(90deg,#4c90e3,#67c5fa)}
.stats{gap:10px!important}.stats>div{padding:12px 15px!important;min-height:75px!important;border-radius:12px!important}.stats small{font-size:11px!important}.stats strong{font-size:22px!important;line-height:1.12;margin-top:6px!important}
.toolbar{margin-top:16px!important;margin-bottom:10px!important}
.toolbar h3{font-size:15px!important;font-weight:650}
.history-controls{gap:8px!important;align-items:center}
#node-search,#node-group,#node-order,.history-controls select{font-size:12px;padding:8px 9px!important;border:1px solid #304c6a;border-radius:8px;background:#13243a}
#node-search{flex:1 1 210px;min-width:150px}
.layout-switch{margin-left:auto!important;border-radius:9px;padding:3px!important}
.layout-button{font-size:11px!important;padding:7px 10px!important;border-radius:7px}
.server-grid{gap:11px!important}
.server-grid[data-layout=compact]{grid-template-columns:repeat(auto-fill,minmax(288px,1fr))!important}
.server-grid[data-layout=cards]{grid-template-columns:repeat(auto-fill,minmax(325px,1fr))!important}
.server-grid .node{padding:13px 14px!important;background:var(--ui-card)!important;border:1px solid var(--ui-border)!important;border-left:3px solid #36c995!important;box-shadow:0 3px 13px #0207121b;transform:translateY(0)}
.server-grid .node.node-offline{border-left-color:#e47b85!important}
.server-grid .node:hover{background:#1b2f49!important;border-color:#4c7099!important;transform:translateY(-2px)}
.server-grid .node-head{margin-bottom:9px!important;gap:7px;min-width:0}
.server-grid .node-name-wrap{min-width:0;overflow:hidden}
.server-grid .node-title{font-size:14px!important;font-weight:690;overflow:hidden;text-overflow:ellipsis}
.server-grid .node-subtitle{color:#96b0cc;font-size:11px!important;margin-top:3px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
.server-grid .pill{font-size:11px;white-space:nowrap}
.server-grid .hardware-line{margin:7px 0 9px!important;gap:5px}
.server-grid .hardware-line span{font-size:10px!important;padding:3px 6px!important;border-radius:5px;background:#233954}
.server-grid .node .meters{gap:8px!important}
.server-grid .node .metric-row{font-size:11px!important;margin-bottom:4px!important}
.server-grid .node .bar{height:5px!important}
.server-grid .node .resource-info{font-size:10px!important;margin-top:4px}
.server-grid .node .meta{font-size:11px!important;padding-top:9px!important;margin-top:11px!important}
.server-grid .period-traffic{gap:6px;margin-top:10px}
.server-grid .period-traffic>div{padding:8px 9px;background:#0f2034;border-radius:8px;border-color:#263f5c}
.server-grid .period-traffic small{font-size:10px}.server-grid .period-traffic strong{font-size:14px;line-height:1.35;display:block;margin-top:3px}
.server-grid .period-traffic em{font-size:10px;line-height:1.4}
.server-grid .period-traffic .coverage{font-size:10px!important;margin-top:6px}
.server-grid[data-layout=compact] .node{padding:12px!important}
.server-grid[data-layout=list] .node{grid-template-columns:minmax(200px,.9fr) minmax(320px,1.3fr);padding:12px 14px!important}
.server-grid[data-layout=list] .period-traffic{grid-column:1/-1;margin-top:0}
.server-grid[data-layout=list] .hardware-line{grid-column:1/-1}
.detail-top{display:flex;align-items:center;gap:15px;margin:4px 0 10px!important}
.detail-identity{min-width:0;flex:1}
#view-detail .detail-top h1{font-size:19px!important;margin:2px 0!important}
#view-detail .detail-top p{margin:0;font-size:11px}
.detail-live{font-size:12px;padding:6px 10px;background:#142d39;border:1px solid #285452;border-radius:8px;color:#91e8ca;max-width:280px}
.section-tabs{display:flex;gap:5px;background:#0c1b2d;padding:5px;border:1px solid #273f59;border-radius:10px;width:max-content;max-width:100%;margin:10px 0}
.detail-tab{background:transparent;padding:8px 26px;border-radius:8px;color:#a5b8d1;font-size:12px;font-weight:650}
.detail-tab.selected{background:#2d5179;color:#fff}
.detail-pane[hidden],.settings-pane[hidden],#add-panel[hidden]{display:none!important}
#view-detail .detail-metrics{margin-top:10px!important;gap:10px!important}
#view-detail .detail-metrics>div{padding:13px 16px!important}
#view-detail .detail-metrics small{font-size:11px;margin-bottom:5px}
#view-detail .detail-metrics strong{font-size:23px!important}
.detail-data-grid{display:grid!important;grid-template-columns:repeat(3,minmax(0,1fr))!important;gap:8px!important}#view-detail #detail-info,#view-detail #detail-network-info{grid-template-columns:repeat(3,minmax(0,1fr))!important}
.detail-data-grid>div{padding:12px 13px!important;background:#0f2035!important;border:1px solid #203a55;border-radius:9px!important;font-size:11px;line-height:1.5;overflow-wrap:anywhere}
.detail-data-grid strong{font-size:13px!important;font-weight:650;color:#e6f0ff;display:inline-block;margin-top:4px}
.field-hint{font-size:11px;color:#8fa8c4;line-height:1.6;margin:9px 0}
#view-detail #detail-traffic{font-size:13px;margin:4px 0}
.event-list{display:flex;flex-direction:column;gap:8px;margin-top:4px}
.event-item{display:flex;gap:12px;align-items:center;padding:12px 14px;border:1px solid #40506a;border-radius:9px;background:#1b2d44}
.event-item.recovered{background:#12283a;border-color:#2a4355}
.event-dot{width:9px;height:9px;flex:0 0 9px;background:#fa8b99;box-shadow:0 0 0 3px #fa8b9920;border-radius:50%}
.event-item.recovered .event-dot{background:#50c59f;box-shadow:0 0 0 3px #50c59f20}
.event-content{flex:1;min-width:0;display:flex;flex-direction:column;gap:3px}
.event-content strong{font-size:13px;font-weight:650;color:#ecf5ff}
.event-content small{font-size:11px;color:#9db4cc}
.event-status{white-space:nowrap;font-size:11px;padding:4px 8px;border-radius:6px;color:#ffa9b7;background:#4a2634}
.event-item.recovered .event-status{color:#88e5c3;background:#184c44}
.event-empty{border:1px dashed #38536f;border-radius:9px;padding:27px;text-align:center;color:#9db1cb;font-size:12px}
.settings-shell{display:grid;grid-template-columns:220px minmax(0,1fr);gap:17px;align-items:start}
.settings-menu{position:sticky;top:20px;display:flex;flex-direction:column;gap:5px;padding:8px;border:1px solid var(--ui-border);border-radius:13px;background:#142339}
.settings-tab{display:flex;align-items:center;gap:10px;text-align:left;background:transparent;color:#a4b7cf;padding:12px 10px;border:1px solid transparent;border-radius:9px;min-width:0}
.settings-tab.selected{background:#233f60;border-color:#365d84;color:#eff6ff}
.menu-icon{display:grid;place-items:center;width:28px;height:28px;flex:0 0 28px;border-radius:8px;background:#273b55;color:#a8cfff;font-size:17px}
.settings-tab.selected .menu-icon{background:#366a9e;color:#fff}
.menu-copy{min-width:0;display:flex;flex-direction:column;gap:3px}.menu-copy strong{font-weight:650;font-size:12px}.menu-copy small{font-size:10px;color:#92a9c3}
.settings-content{min-width:0}
.settings-pane-heading{display:flex;justify-content:space-between;align-items:center;gap:12px;margin:0 0 12px}
.settings-pane-heading h2{font-size:18px;margin:0 0 5px;color:#f4f8ff}.settings-pane-heading p{font-size:11px;margin:0}
.settings-card{background:var(--ui-card);border:1px solid var(--ui-border);border-radius:13px;padding:18px 19px;margin:0 0 13px;box-shadow:0 4px 17px #03091521}
.settings-card-heading{display:flex;justify-content:space-between;gap:8px;align-items:flex-start;flex-wrap:wrap;margin-bottom:13px}
.settings-card-heading h3{font-size:14px;margin:0;color:#f4f8ff;font-weight:700}.settings-card-heading p{font-size:11px;margin:5px 0 0}
#view-settings .settings-form,#view-settings #alerts-form,#view-settings #metadata-form,#view-settings #node-limits-form,#view-settings #create-node{display:grid!important;grid-template-columns:repeat(2,minmax(0,1fr));gap:13px 15px!important;margin:0!important;align-items:start!important;max-width:none!important}
#view-settings .form-field{display:flex;flex-direction:column;gap:6px;font-size:12px;color:#b5c8df;font-weight:550;min-width:0}
#view-settings .form-field small{font-size:10px;color:#8fa6c1;font-weight:400}
#view-settings .form-wide,#view-settings .form-actions{grid-column:1/-1}
#view-settings .settings-form input,#view-settings .settings-form select,#view-settings .settings-form textarea,#view-settings #manage-name,#view-settings #rename-target,#view-settings #install-command{box-sizing:border-box;width:100%!important;min-width:0!important;max-width:none!important;padding:10px 11px!important;border:1px solid #365170!important;border-radius:8px!important;background:#0d1b2e!important;color:#eaf3ff!important;font-size:12px!important;min-height:39px}
#view-settings textarea{resize:vertical;line-height:1.5}
#view-settings .form-actions{display:flex;align-items:center;flex-wrap:wrap;gap:11px;margin-top:2px}
.form-feedback{font-size:11px;color:#8be2c3;min-height:16px}
.settings-danger{border-color:#513b4a}
.settings-danger .form-actions{margin:0!important}
#view-settings .settings-note{border:1px solid #284767;border-radius:11px;background:#112a40;padding:14px 16px;font-size:12px;color:#b4cce4;line-height:1.6}
#view-settings .settings-note strong{color:#d3e9ff}.settings-note p{margin:5px 0 0}
#view-settings #install-command{margin:8px 0 10px}
footer{padding:22px 0!important;font-size:11px!important;color:#718aa9!important}
@media(min-width:701px) and (max-width:940px){#view-detail #detail-info,#view-detail #detail-network-info{grid-template-columns:repeat(2,minmax(0,1fr))!important}}@media(max-width:940px){.settings-shell{grid-template-columns:1fr;gap:13px}.settings-menu{position:static;flex-direction:row;overflow-x:auto;scrollbar-width:thin;gap:6px}.settings-tab{flex:1 0 auto;min-width:145px;padding:10px 9px}.menu-icon{width:25px;height:25px;flex-basis:25px}.menu-copy small{display:none}.detail-data-grid{grid-template-columns:repeat(2,minmax(0,1fr))!important}}
@media(max-width:700px){#view-detail #detail-info,#view-detail #detail-network-info{grid-template-columns:1fr!important}header{padding:12px 14px!important}.top-nav{width:100%;order:3;justify-content:space-between}.top-nav button{flex:1;padding:8px 4px!important;font-size:12px!important}main{margin:12px auto!important;padding:0 12px!important}.section-heading h1,.page-heading h1{font-size:18px!important}.section-heading{margin-bottom:12px!important}.section-heading p{font-size:11px}.stats{grid-template-columns:repeat(2,minmax(0,1fr))!important;gap:7px!important}.stats>div{padding:10px 11px!important;min-height:66px!important}.stats strong{font-size:20px!important}.toolbar{flex-wrap:wrap;gap:8px}.history-controls{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:6px!important}.layout-switch{grid-column:1/-1;width:100%;margin:0!important}.layout-button{flex:1}.server-grid[data-layout=compact],.server-grid[data-layout=cards],.server-grid[data-layout=list]{grid-template-columns:1fr!important}.server-grid[data-layout=list] .node{display:block!important}.settings-card{padding:14px!important}.settings-pane-heading{align-items:flex-start;flex-wrap:wrap}.settings-pane-heading h2{font-size:16px}.settings-menu{padding:6px}.settings-tab{min-width:118px;padding:8px 7px;gap:6px}.menu-icon{display:none}.menu-copy strong{font-size:11px}#view-settings .settings-form,#view-settings #alerts-form,#view-settings #metadata-form,#view-settings #node-limits-form,#view-settings #create-node{grid-template-columns:minmax(0,1fr)!important;gap:11px!important}#view-settings .form-wide,#view-settings .form-actions{grid-column:1}.detail-top{flex-wrap:wrap;gap:8px}.detail-identity{flex:1 1 60%}.detail-live{max-width:none;width:100%}.detail-tab{flex:1;padding:9px 8px}.section-tabs{width:100%}#view-detail .detail-metrics>div{padding:10px 9px!important}#view-detail .detail-metrics strong{font-size:18px!important}.detail-data-grid{grid-template-columns:1fr!important}#view-statistics .history-panel,#view-alerts .history-panel,#view-detail .history-panel{padding:12px!important}.panel-title{align-items:flex-start}.stat-filters{display:flex;flex-wrap:wrap}.stat-filters select{flex:1;min-width:95px!important}.event-item{padding:10px 9px;gap:8px}}

/* Monitor 0.9.4 dashboard redesign: measured 4-column desktop grid */
#view-overview .dashboard-heading{align-items:center;margin:0 0 13px!important}
#view-overview .dashboard-heading h1{font-size:21px!important;letter-spacing:-.3px;line-height:1.25;margin:4px 0!important}
#view-overview .dashboard-heading .eyebrow{font-size:10px;letter-spacing:1.5px;color:#7d9fc8}
#view-overview #overview-hint{font-size:11px;color:#91a8c3;margin:3px 0}
#view-overview .dashboard-updated{display:flex;align-items:center;gap:7px;align-self:center;color:#91a8c3;font-size:11px;white-space:nowrap}
#view-overview .refresh-dot{width:7px;height:7px;border-radius:50%;background:#38cb9d;box-shadow:0 0 0 3px #38cb9d19;flex:0 0 7px}
#view-overview .summary-strip{display:grid!important;grid-template-columns:repeat(4,minmax(0,1fr))!important;gap:10px!important;margin-bottom:16px}
#view-overview .summary-strip>div{min-height:64px!important;display:flex!important;gap:10px;align-items:center;padding:10px 14px!important;background:#15243a!important;border:1px solid #2b4059!important;border-left:1px solid #2b4059!important;box-shadow:none!important}
#view-overview .summary-strip>div:nth-child(1){border-left:2px solid #63a8fa!important}
#view-overview .summary-strip>div:nth-child(2){border-left:2px solid #38d7a2!important}
#view-overview .summary-strip>div:nth-child(3){border-left:2px solid #f28c9d!important}
#view-overview .summary-strip>div:nth-child(4){border-left:2px solid #f4c45a!important}
#view-overview .summary-strip .summary-icon{display:grid;place-items:center;width:30px;height:30px;flex:0 0 30px;border-radius:9px;background:#223550;color:#87b8f4;font-size:16px;font-weight:700}
#view-overview .summary-strip>div:nth-child(2) .summary-icon{background:#183d38;color:#4cdaac}
#view-overview .summary-strip>div:nth-child(3) .summary-icon{background:#3e293a;color:#f38b9a}
#view-overview .summary-strip>div:nth-child(4) .summary-icon{background:#3c3426;color:#f8c76d}
#view-overview .summary-strip small{font-size:11px!important;color:#a1b6ce!important;margin:0!important}
#view-overview .summary-strip strong{font-size:20px!important;line-height:1.08!important;margin:4px 0 0!important}
#view-overview .dashboard-toolbar{display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:10px 18px;margin:0 0 13px;padding:0}
#view-overview .toolbar-label{display:flex;align-items:baseline;gap:10px;flex:0 0 auto}
#view-overview .toolbar-label h3{font-size:15px;line-height:1.4;color:#f3f8ff;margin:0}
#view-overview #node-visible-count{font-size:11px;color:#7f9cb9}
#view-overview .dashboard-filters{display:flex;justify-content:flex-end;align-items:center;flex:1;gap:7px;min-width:0}
#view-overview .dashboard-filters input,#view-overview .dashboard-filters select{height:35px!important;min-height:35px;font-size:11px!important;background:#14243a!important;border:1px solid #304863!important;padding:7px 9px!important;border-radius:8px!important;color:#d5e4f8}
#view-overview .dashboard-filters input{width:190px;max-width:215px;min-width:120px;flex:1 1 130px}
#view-overview .dashboard-filters select{width:auto;max-width:150px;flex:0 0 auto}
#view-overview .dashboard-filters .layout-switch{width:auto!important;display:flex;gap:2px!important;flex:0 0 auto;margin:0!important;border-color:#304863!important}
#view-overview .dashboard-filters .layout-button{height:30px;padding:5px 9px!important;font-size:11px!important}
#view-overview .dashboard-filters .primary-btn{height:35px;padding:7px 12px;font-size:11px;border-radius:8px;white-space:nowrap}
#view-overview .server-grid{gap:11px!important;align-items:stretch}
#view-overview .server-grid[data-layout=compact],#view-overview .server-grid[data-layout=cards]{grid-template-columns:repeat(3,minmax(0,1fr))!important}
#view-overview .server-grid[data-layout=list]{grid-template-columns:1fr!important}
#view-overview .server-grid[data-layout=compact] .dashboard-node,#view-overview .server-grid[data-layout=cards] .dashboard-node{display:flex!important;flex-direction:column!important;min-width:0;min-height:0;padding:14px 14px 13px!important;gap:0;border-radius:13px!important;border:1px solid #2c4059!important;background:linear-gradient(150deg,#172840,#142239)!important;box-shadow:0 4px 17px #050f1b22!important}
#view-overview .server-grid .dashboard-node{border-left:1px solid #2c4059!important}
#view-overview .server-grid .dashboard-node.node-offline{border-color:#744559!important}
#view-overview .server-grid .dashboard-node:hover{transform:translateY(-2px)!important;border-color:#5681a8!important;box-shadow:0 7px 23px #0003!important}
#view-overview .dashboard-node .node-head{display:flex;align-items:flex-start;justify-content:space-between;gap:8px;margin:0 0 10px!important}
#view-overview .dashboard-node .node-title{font-size:14px!important;font-weight:720;color:#f3f8ff;line-height:1.3;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:100%}
#view-overview .dashboard-node .node-location{display:flex;align-items:center;gap:4px;margin-top:5px;font-size:11px;color:#93c6ed;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:100%}
#view-overview .dashboard-node .location-symbol{font-size:15px;line-height:11px;color:#78add8}
#view-overview .dashboard-node .node-subtitle{font-size:10px!important;margin-top:5px;color:#829ab6;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
#view-overview .dashboard-node .pill{font-size:10px!important;font-weight:650;flex:0 0 auto;border-radius:7px;padding:4px 6px;background:#173d37}
#view-overview .dashboard-node .pill.bad{background:#472835}
#view-overview .dashboard-node .hardware-line{display:flex;gap:4px!important;margin:0 0 12px!important;flex-wrap:nowrap}
#view-overview .dashboard-node .hardware-line span{font-size:10px!important;display:block;min-width:0;padding:4px 6px!important;border-radius:6px!important;background:#233955!important;color:#c9dcf5;white-space:nowrap;text-overflow:ellipsis;overflow:hidden}
#view-overview .dashboard-node .meters{display:flex!important;flex-direction:column;gap:11px!important;grid-template-columns:none!important}
#view-overview .dashboard-node .resource-line{min-width:0}
#view-overview .dashboard-node .metric-row{display:flex;justify-content:space-between;align-items:center;font-size:11px!important;gap:8px;margin:0 0 5px!important}
#view-overview .dashboard-node .resource-numbers{display:flex;align-items:baseline;gap:8px;min-width:0}
#view-overview .dashboard-node .resource-numbers small{color:#92a9c4;font-size:10px;white-space:nowrap}
#view-overview .dashboard-node .resource-numbers b{color:#e7f1ff;font-size:12px;font-weight:700;min-width:27px;text-align:right}
#view-overview .dashboard-node .bar{height:6px!important;background:#293c54;border-radius:5px!important}
#view-overview .dashboard-node .fill{background:#68aaf2}#view-overview .dashboard-node .resource-line:nth-child(2) .fill{background:#9a84ef!important}#view-overview .dashboard-node .resource-line:nth-child(3) .fill{background:#41cbb5!important}
#view-overview .dashboard-node .fill.hot{background:#f4be64!important}#view-overview .dashboard-node .fill.critical{background:#ef7383!important}
#view-overview .dashboard-node .node-activity{display:flex;align-items:center;justify-content:space-between;gap:5px;margin:13px 0 0!important;padding:10px 0!important;border-top:1px solid #2b4058!important;color:#99b1cc;font-size:10px!important}
#view-overview .dashboard-node .node-activity span{min-width:0;white-space:nowrap}#view-overview .dashboard-node .node-activity .uptime{color:#9bb2cd}
#view-overview .dashboard-node .period-traffic{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:8px;margin:0!important}
#view-overview .dashboard-node .period-traffic .traffic-cell{background:#102035!important;border:1px solid #253b52;border-radius:8px;padding:9px 10px}
#view-overview .dashboard-node .traffic-cell small{font-size:10px!important;color:#93a9c3}
#view-overview .dashboard-node .traffic-cell strong{font-size:15px!important;color:#ebf4ff;line-height:1.4;font-weight:720;margin-top:4px}
#view-overview .dashboard-node .quota-caption{display:flex;justify-content:space-between;margin-top:9px;color:#9bb5cf;font-size:10px}
#view-overview .dashboard-node .quota-track{height:5px;margin-top:4px}
#view-overview .dashboard-node .node-badges{display:flex;flex-wrap:wrap;gap:5px;margin-top:8px}
#view-overview .dashboard-node .alert-chip{display:inline-block;padding:3px 7px;border-radius:5px;background:#463827;color:#ffcb7d;font-size:10px}
#view-overview .dashboard-node .alert-chip.urgent{background:#4d2c39;color:#ff9dab}
#view-overview .dashboard-node .offline-hint{padding:12px 0;color:#8da2bc;font-size:11px}
#view-overview .dashboard-empty{display:flex;flex-direction:column;gap:5px;align-items:center;justify-content:center;grid-column:1/-1;border:1px dashed #39536d;border-radius:13px;padding:36px;color:#a0b5cf}
#view-overview .dashboard-empty strong{font-size:14px}#view-overview .dashboard-empty span{font-size:11px}
@media(min-width:1250px){#view-overview .server-grid[data-layout=compact],#view-overview .server-grid[data-layout=cards]{grid-template-columns:repeat(4,minmax(0,1fr))!important}}
@media(min-width:780px) and (max-width:1249px){#view-overview .server-grid[data-layout=compact],#view-overview .server-grid[data-layout=cards]{grid-template-columns:repeat(2,minmax(0,1fr))!important}}
@media(max-width:779px){#view-overview .server-grid[data-layout=compact],#view-overview .server-grid[data-layout=cards]{grid-template-columns:1fr!important}#view-overview .dashboard-toolbar{align-items:stretch}#view-overview .dashboard-filters{justify-content:flex-start;flex-wrap:wrap}#view-overview .dashboard-filters input{flex:1 1 100%;max-width:none;width:100%}#view-overview .dashboard-filters select{flex:1 1 43%;max-width:none}#view-overview .dashboard-filters .layout-switch{flex:1 1 auto;width:auto!important}#view-overview .dashboard-filters .primary-btn{flex:0 0 auto}}
@media(max-width:520px){#view-overview .summary-strip{gap:6px!important;grid-template-columns:repeat(2,minmax(0,1fr))!important}#view-overview .summary-strip>div{min-height:55px!important;padding:8px!important;gap:8px}#view-overview .summary-strip .summary-icon{width:25px;height:25px;flex-basis:25px}#view-overview .summary-strip strong{font-size:19px!important}#view-overview .dashboard-heading h1{font-size:18px!important}#view-overview .dashboard-updated{font-size:10px}#view-overview .dashboard-node .node-activity{font-size:10px!important}}

`
