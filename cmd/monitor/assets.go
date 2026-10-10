package main

import "net/http"

func init() {
 // The main server uses this handler at /static/ to serve embedded assets.
}
func staticHandler() http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  switch r.URL.Path{
  case "/static/style.css":w.Header().Set("Content-Type","text/css; charset=utf-8");w.Write([]byte(styleCSS))
  case "/static/app.js":w.Header().Set("Content-Type","application/javascript; charset=utf-8");w.Write([]byte(appJS))
  default:http.NotFound(w,r)
  }
 })
}
const styleCSS=`.period-traffic{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:7px;margin-top:8px}.period-traffic>div{background:#10243a;border:1px solid #29405a;border-radius:7px;padding:6px 8px}.period-traffic small{display:block;font-size:10px;color:#9eb1cb}.period-traffic strong{font-size:12px;color:#e2efff}.period-traffic em{display:block;font-style:normal;color:#9eb1cb;font-size:9px;margin-top:2px}.server-grid[data-layout=compact]{grid-template-columns:repeat(auto-fill,minmax(270px,1fr))!important}.server-grid[data-layout=compact] .meters{gap:6px!important}.server-grid[data-layout=compact] .node{padding:12px!important}@media(max-width:700px){.server-grid[data-layout=compact]{grid-template-columns:1fr!important}}.hardware-line{display:flex;flex-wrap:wrap;gap:5px;margin:7px 0 8px}.hardware-line span{font-size:10px;padding:3px 6px;border-radius:5px;background:#263a53;color:#d3e5f8}.resource-info{font-size:10px;color:#adc0d7;margin-top:3px;text-align:right}.server-grid[data-layout=compact] .hardware-line{gap:3px;margin:5px 0}.server-grid[data-layout=compact] .hardware-line span{font-size:9px;padding:2px 4px}.server-grid[data-layout=list] .hardware-line{grid-column:1/-1;margin:1px 0}.server-grid[data-layout=list] .resource-info{display:none}@media(min-width:780px){#view-detail .detail-heading{gap:10px!important}#view-detail .detail-heading h1{font-size:15px!important}#view-detail .detail-heading p,#view-detail #detail-status{font-size:11px!important}#view-detail .detail-heading button,#view-detail .detail-heading select{font-size:11px!important;padding:6px 9px!important}#view-detail .detail-metrics{gap:7px!important;margin-top:9px!important}#view-detail .detail-metrics>div{padding:8px 10px!important}#view-detail .detail-metrics strong{font-size:17px!important}#view-detail .detail-metrics small{font-size:10px!important;margin-bottom:2px!important}#view-detail .history-panel,#view-statistics .history-panel{padding:10px 12px!important;margin-top:9px!important}#view-detail .history-panel h3,#view-statistics .history-panel h3{font-size:13px!important;margin:2px 0 8px!important}#detail-info{grid-template-columns:repeat(4,minmax(0,1fr))!important;gap:6px!important}#detail-info>div{padding:7px 9px!important;font-size:11px!important;line-height:1.4!important}#detail-info strong{font-size:12px!important;font-weight:600!important}#view-statistics>h1{font-size:16px!important;margin:5px 0 10px!important}#view-statistics .history-controls{gap:7px!important;margin-bottom:6px!important}#view-statistics .history-controls select{font-size:11px!important;padding:6px!important}#detail-chart,#history-chart{max-height:240px;object-fit:contain}#traffic-bars .item{font-size:11px!important;margin:5px 0!important}}@media(min-width:900px){body{font-size:12px!important}header{padding:7px 2.4%!important;min-height:48px}header h2{font-size:17px;margin:0}.top-nav{padding:3px}.top-nav button{padding:6px 12px;font-size:12px}header form button{padding:6px 10px;font-size:12px}main{max-width:1580px!important;margin:10px auto!important;padding:0 16px!important}.page-heading{margin:0 0 8px!important}.page-heading h1{font-size:16px!important}.stats{gap:7px!important}.stats>div{padding:7px 11px!important;min-height:0!important}.stats small{font-size:11px!important}.stats strong{font-size:19px!important}.toolbar{margin-top:12px!important;margin-bottom:7px!important}.toolbar h3{font-size:13px}.toolbar button,.history-controls select,.history-controls input{font-size:11px;padding:6px 9px}.history-controls{gap:6px!important;margin-bottom:8px!important}.layout-button{padding:5px 8px!important;font-size:11px!important}.server-grid{gap:8px!important}.server-grid[data-layout=cards]{grid-template-columns:repeat(auto-fill,minmax(260px,1fr))!important}.server-grid[data-layout=compact]{grid-template-columns:repeat(auto-fill,minmax(205px,1fr))!important}.server-grid .node{padding:10px 11px!important}.server-grid .node-head{margin-bottom:8px!important}.server-grid .node-title{font-size:12px!important}.server-grid .node-subtitle{font-size:10px!important}.server-grid .node .meters{gap:7px!important}.server-grid .node .metric-row{font-size:11px!important;margin-bottom:3px!important}.server-grid .node .meta{font-size:10px!important;padding-top:7px!important;margin-top:8px!important}.server-grid[data-layout=compact] .node{padding:9px!important}.server-grid[data-layout=list] .node{padding:8px 10px!important}.view{min-height:0!important}footer{padding:13px 0!important;font-size:10px}}.traffic-totals{display:flex;flex-wrap:wrap;gap:5px 12px;margin-top:6px;font-size:11px;color:#a9bed5}.traffic-totals b{color:#e1edff}.server-grid[data-layout=compact] .traffic-totals{font-size:10px;gap:3px 8px}.server-grid[data-layout=list] .traffic-totals{margin-top:3px}@media(max-width:600px){.traffic-totals{font-size:10px}}.server-grid .node{cursor:pointer;transition:border-color .15s,transform .15s,background .15s}.server-grid .node:hover{background:#1b2c45;transform:translateY(-1px);border-color:#5685b9!important}.server-grid .node:focus-visible{outline:2px solid #74b9ff;outline-offset:3px}.detail-heading{display:flex;align-items:center;gap:18px;flex-wrap:wrap}.detail-heading h1{font-size:21px;margin:0}.detail-heading p{margin:4px 0}.detail-heading button,.detail-heading select{background:#1b314c;border:1px solid #385579;color:white;border-radius:9px;padding:10px}.detail-metrics{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:12px;margin-top:18px}.detail-metrics>div{background:#17243a;border:1px solid #2a3b53;border-radius:12px;padding:16px}.detail-metrics strong{font-size:23px}.detail-metrics small{display:block;color:#9fb0c9;margin-bottom:5px}#detail-info{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:9px}#detail-info>div{padding:12px;background:#101d30;border-radius:8px;overflow-wrap:anywhere}#detail-chart{width:100%;height:auto}@media(max-width:600px){.detail-metrics{gap:6px}.detail-metrics>div{padding:9px}.detail-metrics strong{font-size:17px}#detail-info{grid-template-columns:1fr}}.stats{grid-template-columns:repeat(4,minmax(0,1fr))!important;gap:9px!important}.stats>div{min-height:0!important;padding:11px 13px!important;border-radius:11px!important;border-color:#2a4059!important;position:relative}.stats small{font-size:12px!important;color:#a3b5cc!important}.stats strong{font-size:23px!important;margin-top:3px!important;line-height:1.25}.stats>div:nth-child(1){border-left:3px solid #60a5fa!important}.stats>div:nth-child(2){border-left:3px solid #34d399!important}.stats>div:nth-child(3){border-left:3px solid #fb7185!important}.stats>div:nth-child(4){border-left:3px solid #fbbf24!important}.stats>div:nth-child(1) strong{color:#93c5fd}.stats>div:nth-child(2) strong{color:#34d399}.stats>div:nth-child(3) strong{color:#fb7185}.stats>div:nth-child(4) strong{color:#fbbf24}.page-heading{margin:0 0 12px!important}.page-heading h1{font-size:19px!important;margin:0!important}.toolbar{margin-top:17px!important;margin-bottom:10px!important}.toolbar h3{margin:0;font-size:15px}.history-controls{gap:8px!important;margin-bottom:12px!important}.server-grid .node{border-radius:12px!important;padding:14px!important;border-left:3px solid #34d399!important}.server-grid .node.node-offline{border-left-color:#fb7185!important}.server-grid .node-head{margin-bottom:13px!important;align-items:center}.server-grid .node-title{font-weight:650;font-size:14px;min-width:0;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}.server-grid .node-subtitle{font-size:11px;color:#91a5be;margin-top:3px}.server-grid .node .meters{gap:11px}.server-grid .node .metric-row{font-size:12px;margin-bottom:5px}.server-grid .node .metric-row b{font-weight:650}.server-grid .node .bar{height:5px}.server-grid .node .fill{background:#60a5fa}.server-grid .node .meters>div:nth-child(2) .fill{background:#a78bfa}.server-grid .node .meters>div:nth-child(3) .fill{background:#2dd4bf}.server-grid .node .fill.hot{background:#fb923c!important}.server-grid .node .fill.critical{background:#fb7185!important}.server-grid .node .meta{border-top:1px solid #2a4059;padding-top:10px;margin-top:12px!important;display:flex;justify-content:space-between;gap:8px;font-size:11px;flex-wrap:wrap}.server-grid .node .rx{color:#34d399}.server-grid .node .tx{color:#60a5fa}.server-grid[data-layout=compact] .node{padding:11px!important}.server-grid[data-layout=compact] .node-head{margin-bottom:10px!important}.server-grid[data-layout=compact] .node .meters{gap:6px!important}.server-grid[data-layout=compact] .node .meta{padding-top:8px;margin-top:9px!important}.server-grid[data-layout=list] .node{grid-template-columns:minmax(160px,1fr) minmax(300px,2fr);padding:11px 14px!important}.server-grid[data-layout=list] .node .meta{grid-column:1/-1;margin-top:0!important;padding-top:7px}.server-grid[data-layout=list] .node-head{margin:0!important}@media(max-width:600px){.stats{gap:5px!important}.stats>div{padding:8px 5px!important;text-align:center}.stats small{font-size:10px!important}.stats strong{font-size:19px!important}.server-grid[data-layout=list] .node{display:block}.page-heading h1{font-size:17px!important}}.layout-switch{display:inline-flex;gap:4px;margin-left:auto;background:#101c2f;padding:4px;border-radius:10px;border:1px solid #30425b}.layout-button{white-space:nowrap;background:transparent;color:#9fb0c9;padding:8px 11px;font-size:12px}.layout-button.active{background:#345782;color:#fff}.server-grid{display:grid;gap:12px}.server-grid[data-layout=cards]{grid-template-columns:repeat(auto-fit,minmax(min(100%,320px),1fr))}.server-grid[data-layout=compact]{grid-template-columns:repeat(auto-fit,minmax(min(100%,245px),1fr))}.server-grid[data-layout=list]{grid-template-columns:1fr}.server-grid .node{margin:0!important;min-width:0}.server-grid[data-layout=compact] .node{padding:13px}.server-grid[data-layout=compact] .meters{grid-template-columns:1fr;gap:8px}.server-grid[data-layout=list] .node{display:grid;grid-template-columns:minmax(160px,1fr) minmax(280px,2fr);align-items:center;gap:10px 20px}.server-grid[data-layout=list] .node-head{margin:0}.server-grid[data-layout=list] .meta{grid-column:1/-1;margin:0}@media(max-width:650px){.layout-switch{margin-left:0;width:100%;justify-content:space-around}.layout-button{flex:1;padding:8px 4px}.server-grid[data-layout=list] .node{display:block}.server-grid[data-layout=cards],.server-grid[data-layout=compact]{grid-template-columns:1fr}}.view[hidden]{display:none!important}.top-nav{display:flex;gap:6px;background:#0d192c;border-radius:12px;padding:5px}.top-nav button{background:transparent;color:#9aabc5;padding:10px 17px}.top-nav button.selected{background:#29486c;color:#fff}header{gap:16px;flex-wrap:wrap}main{max-width:1200px!important}.view{min-height:62vh}.view>h1,.page-heading h1{font-size:22px;font-weight:650;margin:14px 0 23px}.page-heading{display:flex;align-items:center;justify-content:space-between}.history-panel{margin-top:16px!important}.stats>div{min-height:120px}footer{text-align:center;padding:30px 0}body{background:#0c1423!important}@media(max-width:650px){.top-nav{order:3;width:100%;justify-content:space-around}.top-nav button{padding:9px 12px;flex:1}.page-heading h1{font-size:20px}header h2{margin:0}}#alerts-form{display:flex;flex-wrap:wrap;gap:12px;align-items:center}#alerts-form input,#node-search,#node-group,#node-order{background:#0b1220;color:white;padding:8px;border:1px solid #34455e;border-radius:7px}#alerts-form input[type=number]{width:75px}#alert-events{font-size:13px;margin-top:12px}#metadata-form{display:flex;flex-wrap:wrap;gap:8px;margin:10px 0}#metadata-form input,#metadata-form select,#metadata-form textarea{background:#0b1220;color:#fff;border:1px solid #34455e;border-radius:7px;padding:9px;max-width:100%}#manage-name,#rename-target{background:#0b1220;border:1px solid #34455e;border-radius:7px;padding:9px;color:white;max-width:100%}#add-panel[hidden]{display:none}#create-node{display:flex;gap:9px;flex-wrap:wrap}#create-node input{padding:10px;color:#fff;background:#0b1220;border:1px solid #34455e;border-radius:8px;flex:1;min-width:170px}#install-command{width:100%;background:#0b1220;color:#cee4ff;padding:10px;border:1px solid #34455e;border-radius:8px}#traffic-bars .item{display:flex;gap:8px;align-items:center;margin:8px 0;font-size:12px}#traffic-bars .track{background:#2b3b50;height:9px;flex:1;border-radius:4px}#traffic-bars .value{height:100%;border-radius:4px;background:#298be8}*{box-sizing:border-box}body{background:#0b1220;color:#eaf1fc;font:15px system-ui,sans-serif;margin:0}header{display:flex;justify-content:space-between;align-items:center;background:#121f33;padding:12px 5%}main{max-width:1100px;margin:28px auto;padding:0 18px}button{background:#293b58;border:0;color:white;padding:8px 14px;border-radius:8px;cursor:pointer}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:12px}.stats>div,.node{background:#17243a;border:1px solid #2a3b53;border-radius:12px;padding:18px}.stats small,.stats strong{display:block}.stats small,.muted,.toolbar span{color:#97a6bc}.stats strong{font-size:27px;margin-top:8px}.toolbar{display:flex;justify-content:space-between;align-items:center;margin-top:26px}.node{margin:10px 0}.node-head,.metric-row{display:flex;justify-content:space-between;gap:12px}.node-head{font-weight:bold;margin-bottom:15px}.pill{font-size:12px}.pill.ok{color:#5ad6a3}.pill.bad{color:#fc8892}.meters{display:grid;grid-template-columns:repeat(3,1fr);gap:16px}.metric-row{font-size:13px;margin-bottom:6px}.bar{height:7px;border-radius:7px;background:#2b3b50;overflow:hidden}.fill{height:100%;background:#388add}.meta{font-size:12px;color:#a3b3c9;margin-top:14px;overflow-wrap:anywhere}.history-panel{background:#17243a;border:1px solid #2a3b53;border-radius:12px;padding:18px;margin-top:20px}.history-controls{display:flex;gap:8px;flex-wrap:wrap;margin-bottom:16px}.history-controls select{background:#0b1220;color:white;border:1px solid #34455e;padding:8px;border-radius:7px}canvas{width:100%;height:auto}@media(max-width:600px){.meters{grid-template-columns:1fr}.stats strong{font-size:23px}.stats>div{padding:12px}}@media(min-width:780px){.server-grid[data-layout=compact]{grid-template-columns:repeat(auto-fill,minmax(288px,1fr))!important}.server-grid[data-layout=compact] .node{padding:12px!important}}@media(max-width:779px){.server-grid[data-layout=compact]{grid-template-columns:1fr!important}}.agent-tag{display:inline-block;padding:2px 5px;border-radius:5px;background:#1b3a42;color:#7ddac5;font-size:10px;margin-left:4px}.agent-tag.old{background:#433628;color:#ffcb7f}.quota-track{height:5px;border-radius:5px;background:#2b4058;overflow:hidden;margin-top:6px}.quota-fill{height:100%;background:#34d399}.quota-fill.warn{background:#f59e0b}.quota-fill.high{background:#fb7185}.node-badges{display:flex;gap:6px;flex-wrap:wrap;font-size:10px;margin-top:6px}.node-badges .expired{color:#fb7185}.node-badges .soon{color:#fbbf24}.period-traffic .coverage{font-size:9px;color:#96a7bd;margin-top:4px}#node-limits-form{display:flex;align-items:center;flex-wrap:wrap;gap:9px}#node-limits-form input,#node-limits-form select{padding:8px;background:#0b1220;border:1px solid #34455e;border-radius:7px;color:#f0f5ff;font-size:12px;max-width:220px}#node-limits-form label{font-size:12px}#node-limits-form label input{width:145px}
/* V0.9.2: unified dense desktop dashboard and redesigned management UI */
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

/* V0.9.5 selectable themes: default pastel glass reference, plus dark/light. */
.header-tools{display:flex;align-items:center;gap:10px;white-space:nowrap}
.theme-switch{display:flex;align-items:center;gap:5px;border:1px solid #46617e;border-radius:10px;padding:4px 8px;background:#20344d;color:#e6f4ff;font-size:13px}
.theme-switch select{appearance:none;-webkit-appearance:none;min-width:85px;padding:4px 0;color:inherit;background:transparent;border:0;outline:0;cursor:pointer;font-size:12px;font-weight:650}
.theme-switch option{color:#202d42;background:#fafbff}
.theme-options{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:11px}
.theme-choice{display:flex;flex-direction:column;align-items:flex-start;gap:7px;padding:12px;border:1px solid #3c5270;background:#1a2d43;border-radius:13px;text-align:left;min-width:0;color:#d8e9ff}
.theme-choice strong{font-size:13px}.theme-choice small{font-size:11px;color:#9eb1ca;line-height:1.4}
.theme-choice.selected{border-color:#74acf1;box-shadow:0 0 0 2px #74acf12f}
.theme-thumb{display:block;width:100%;height:68px;border-radius:9px;border:1px solid #ffffff55}
.theme-thumb-glass{background:radial-gradient(ellipse at 70% 15%,#fffbf3 0%,transparent 45%),radial-gradient(ellipse at 38% 63%,#c39dcd 0%,transparent 51%),radial-gradient(ellipse at 80% 84%,#e7b6d6 0%,transparent 60%),#cfd5e8}
.theme-thumb-dark{background:linear-gradient(145deg,#0b1525 20%,#213854 50%,#15243a 100%)}
.theme-thumb-light{background:linear-gradient(135deg,#f9fbff 0%,#e4ecf5 100%)}
html[data-theme=glass]{--t-text:#28324b;--t-muted:#5b6780;--t-card:#fffdfb9c;--t-card-deep:#fffeffb8;--t-border:#ffffffc9;--t-input:#ffffffc4;--t-shadow:0 10px 26px #38315419}
html[data-theme=glass] body{background:radial-gradient(ellipse at 67% 12%,#fffcf4c9 0%,transparent 40%),radial-gradient(ellipse at 48% 31%,#85799db5 0%,transparent 46%),radial-gradient(ellipse at 22% 69%,#e4bddacd 0%,transparent 52%),radial-gradient(ellipse at 82% 80%,#becde9e0 0%,transparent 53%),linear-gradient(115deg,#cbd3e5,#f1e6df 48%,#e6d4e8)!important;color:var(--t-text)!important;background-attachment:fixed!important;min-height:100vh}
html[data-theme=glass] body:before{content:"";position:fixed;inset:-12%;z-index:-1;pointer-events:none;opacity:.54;filter:blur(56px);background:radial-gradient(ellipse 22% 55% at 46% 32%,#746380 0%,transparent 90%),radial-gradient(ellipse 22% 42% at 59% 66%,#d09cc0 0%,transparent 80%),radial-gradient(ellipse 26% 44% at 82% 34%,#fffbf5 0%,transparent 84%);background-attachment:fixed}
html[data-theme=light]{--t-text:#202e47;--t-muted:#607087;--t-card:#fffffff2;--t-card-deep:#ffffff;--t-border:#d9e2ee;--t-input:#f8fbff;--t-shadow:0 4px 16px #102b4d0e}
html[data-theme=light] body{background:linear-gradient(150deg,#f8fbff,#eaf0f8 54%,#f5f8fc)!important;color:var(--t-text)!important;min-height:100vh}
:is(html[data-theme=glass],html[data-theme=light]) header{background:var(--t-card)!important;backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px);border-bottom:1px solid var(--t-border)!important}
:is(html[data-theme=glass],html[data-theme=light]) header h2{color:#27344c!important}
:is(html[data-theme=glass],html[data-theme=light]) header .top-nav{background:#ffffff71!important;border:1px solid #ffffffa8!important;box-shadow:none!important}
:is(html[data-theme=glass],html[data-theme=light]) .top-nav button{color:#5d6c81!important}
:is(html[data-theme=glass],html[data-theme=light]) .top-nav button.selected{background:#f8fbfff2!important;color:#2b3856!important;box-shadow:0 3px 8px #3b4b6c1c!important}
:is(html[data-theme=glass],html[data-theme=light]) .theme-switch{background:#fff9;border-color:#ffffffcd;color:#384968}
:is(html[data-theme=glass],html[data-theme=light]) .logout-btn{background:#ffffff9c!important;border:1px solid #ffffffc9!important;color:#3a4d69!important}
:is(html[data-theme=glass],html[data-theme=light]) h1,
:is(html[data-theme=glass],html[data-theme=light]) h2,
:is(html[data-theme=glass],html[data-theme=light]) h3{color:#29344f!important}
:is(html[data-theme=glass],html[data-theme=light]) .muted,
:is(html[data-theme=glass],html[data-theme=light]) .field-hint{color:var(--t-muted)!important}
:is(html[data-theme=glass],html[data-theme=light]) .eyebrow{color:#796e9a!important}
:is(html[data-theme=glass],html[data-theme=light]) .primary-btn{background:#eef0ffdf!important;border-color:#ffffffce!important;color:#38476a!important;box-shadow:0 3px 11px #6c678d20!important}
:is(html[data-theme=glass],html[data-theme=light]) .subtle-btn,
:is(html[data-theme=glass],html[data-theme=light]) .layout-switch{background:#ffffff9c!important;border-color:#ffffffe0!important;color:#354669!important}
:is(html[data-theme=glass],html[data-theme=light]) .layout-button{color:#65738d!important}
:is(html[data-theme=glass],html[data-theme=light]) .layout-button.active{background:#fffefa!important;color:#293957!important}
:is(html[data-theme=glass],html[data-theme=light]) .stats>div,
:is(html[data-theme=glass],html[data-theme=light]) .history-panel,
:is(html[data-theme=glass],html[data-theme=light]) .settings-card,
:is(html[data-theme=glass],html[data-theme=light]) .settings-menu,
:is(html[data-theme=glass],html[data-theme=light]) .settings-note{background:var(--t-card)!important;border:1px solid var(--t-border)!important;color:var(--t-text)!important;box-shadow:var(--t-shadow)!important;backdrop-filter:blur(20px);-webkit-backdrop-filter:blur(20px)}
:is(html[data-theme=glass],html[data-theme=light]) .settings-tab{color:#58677f!important}
:is(html[data-theme=glass],html[data-theme=light]) .settings-tab.selected{background:#fffefa!important;border-color:#fff!important;color:#304464!important}
:is(html[data-theme=glass],html[data-theme=light]) .menu-icon{background:#e0e3f1!important;color:#566f94!important}
:is(html[data-theme=glass],html[data-theme=light]) .settings-card-heading p,
:is(html[data-theme=glass],html[data-theme=light]) .settings-pane-heading p,
:is(html[data-theme=glass],html[data-theme=light]) .menu-copy small{color:#6b7c92!important}
:is(html[data-theme=glass],html[data-theme=light]) .form-field{color:#4b5b75!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .settings-form input,
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .settings-form select,
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .settings-form textarea,
:is(html[data-theme=glass],html[data-theme=light]) .history-controls select,
:is(html[data-theme=glass],html[data-theme=light]) .history-controls input{background:var(--t-input)!important;border:1px solid var(--t-border)!important;color:var(--t-text)!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .settings-form input::placeholder,
:is(html[data-theme=glass],html[data-theme=light]) input::placeholder{color:#8694a8!important}
:is(html[data-theme=glass],html[data-theme=light]) .theme-choice{background:#ffffff91!important;border-color:#ffffffd5!important;color:#293957!important}
:is(html[data-theme=glass],html[data-theme=light]) .theme-choice strong{color:#293957!important}
:is(html[data-theme=glass],html[data-theme=light]) .theme-choice.selected{border-color:#94aadd!important;box-shadow:0 0 0 2px #7295c429!important}
:is(html[data-theme=glass],html[data-theme=light]) .theme-choice small{color:#65758a!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-heading h1{color:#29334b!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview #overview-hint,
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-updated{color:#65758e!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .summary-strip>div{background:var(--t-card)!important;border:1px solid var(--t-border)!important;box-shadow:var(--t-shadow)!important;border-left:2px solid #a8add4!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .summary-strip small{color:#576d86!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .summary-strip strong{color:#233b60!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .summary-strip .summary-icon{background:#ffffff83!important;color:#527ab5!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .toolbar-label h3,
:is(html[data-theme=glass],html[data-theme=light]) #view-overview #node-visible-count{color:#43556d!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-filters input,
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-filters select,
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-filters .layout-switch{background:#fffafccc!important;border:1px solid #fffdfbdd!important;color:#354968!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-filters .primary-btn{background:#ffffffc9!important;color:#3a4e70!important;border:1px solid #ffffffdf!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .server-grid .dashboard-node{background:var(--t-card)!important;color:var(--t-text)!important;border:1px solid var(--t-border)!important;box-shadow:var(--t-shadow)!important;backdrop-filter:blur(21px);-webkit-backdrop-filter:blur(21px)}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .server-grid .dashboard-node:hover{background:var(--t-card-deep)!important;border-color:#ffffffed!important;box-shadow:0 13px 30px #3f365021!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .node-title{color:#233750!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .node-subtitle{color:#61748d!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .node-location{color:#466591!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .pill.ok{background:#daf6e8b5;color:#177856!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .pill.bad{background:#f9dce4b5;color:#bc4669!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .hardware-line span{background:#e7eefa9e!important;color:#4c6281!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .metric-row{color:#3c4b66!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .resource-numbers small{color:#73849c!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .resource-numbers b{color:#273e5c!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .bar{background:#aebad154!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .node-activity{border-top:1px solid #c9ccdbba!important;color:#4e6a8b!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .node-activity .uptime{color:#64758d!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .period-traffic .traffic-cell{background:#ffffff7b!important;border:1px solid #f3f3fac9!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .traffic-cell small{color:#6c7892!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .traffic-cell strong{color:#2e4264!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .quota-caption{color:#5b7291!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-overview .dashboard-node .quota-track{background:#b6b7d54c!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-detail .detail-live{background:#e7f5ecbf;color:#20734e;border-color:#bde6d7}
:is(html[data-theme=glass],html[data-theme=light]) #view-detail .detail-data-grid>div,
:is(html[data-theme=glass],html[data-theme=light]) #view-detail .detail-metrics>div{background:#ffffffa5!important;color:#4e6180!important;border:1px solid #ffffffb8!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-detail .detail-data-grid strong,
:is(html[data-theme=glass],html[data-theme=light]) #view-detail .detail-metrics strong{color:#2d415e!important}
:is(html[data-theme=glass],html[data-theme=light]) .section-tabs{background:#ffffff8a!important;border-color:#ffffffbe!important}
:is(html[data-theme=glass],html[data-theme=light]) .detail-tab{color:#67748c}
:is(html[data-theme=glass],html[data-theme=light]) .detail-tab.selected{background:#fff9!important;color:#324463!important}
:is(html[data-theme=glass],html[data-theme=light]) .event-item{background:#ffffffa5!important;border-color:#ffffffb8!important}
:is(html[data-theme=glass],html[data-theme=light]) .event-content strong{color:#2d405f!important}
:is(html[data-theme=glass],html[data-theme=light]) .event-content small{color:#6f8195!important}
:is(html[data-theme=glass],html[data-theme=light]) .event-item.recovered{background:#f4fffbba!important}
:is(html[data-theme=glass],html[data-theme=light]) #traffic-bars .track{background:#9faecc66!important}
:is(html[data-theme=glass],html[data-theme=light]) .settings-note strong{color:#415879!important}
:is(html[data-theme=glass],html[data-theme=light]) .settings-note p{color:#586d85!important}
@media(max-width:680px){.header-tools{margin-left:auto}.theme-switch{padding:3px 7px}.theme-switch select{font-size:11px;min-width:72px}.theme-options{grid-template-columns:1fr}.theme-thumb{height:50px}}
@media(prefers-reduced-motion:reduce){button,.node{transition:none!important;animation:none!important}}
`
const appJS=`const $=x=>document.getElementById(x);
const validThemes=['glass','dark','light'];
let currentTheme='glass';
try{const stored=localStorage.getItem('monitor-theme');if(validThemes.includes(stored))currentTheme=stored}catch(e){}
function setTheme(next){if(!validThemes.includes(next))next='glass';currentTheme=next;document.documentElement.dataset.theme=next;const chooser=$('theme-select');if(chooser)chooser.value=next;document.querySelectorAll('[data-theme-choice]').forEach(b=>{const active=b.dataset.themeChoice===next;b.classList.toggle('selected',active);b.setAttribute('aria-pressed',String(active))});try{localStorage.setItem('monitor-theme',next)}catch(e){}}
setTheme(currentTheme);
$('theme-select').addEventListener('change',e=>setTheme(e.target.value));
document.querySelectorAll('[data-theme-choice]').forEach(b=>b.addEventListener('click',()=>{setTheme(b.dataset.themeChoice);$('theme-status').textContent='已切换为：'+({glass:'梦幻玻璃',dark:'经典深色',light:'简约浅色'}[currentTheme])}));

function esc(v){return String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]))}
function metric(label,n){n=Math.max(0,Math.min(100,Number(n)||0));return '<div><div class="metric-row"><span>'+label+'</span><b>'+n.toFixed(1)+'%</b></div><div class="bar"><div class="fill" style="width:'+n+'%"></div></div></div>'}
function formatBytes(n){if(!n)return '0 B';const units=['B','KB','MB','GB','TB'];let i=0;while(n>=1024&&i<units.length-1){n/=1024;i++}return n.toFixed(i?1:0)+' '+units[i]}
let prevNames="";let cachedNodes=[];async function refresh(){try{const r=await fetch('/api/v1/nodes',{cache:'no-store'});if(r.status===401){location.href='/login';return}if(!r.ok)throw Error('HTTP '+r.status);const nodes=await r.json();nodes.sort((a,b)=>(a.group||"").localeCompare(b.group||"")||(a.display_name||a.name).localeCompare(b.display_name||b.name));cachedNodes=nodes;if(selectedNode)updateDetail();updateGroupOptions(nodes);const names=nodes.map(n=>n.name).join('|');if(names!==prevNames){prevNames=names;const ed=$("metadata-name"),old=ed.value;ed.replaceChildren(...nodes.map(n=>new Option(n.display_name||n.name,n.name)));if(nodes.some(n=>n.name===old))ed.value=old;loadMetadata();const sel=$('history-name');const manage=$('manage-name');const was=manage.value;manage.replaceChildren(...nodes.map(n=>new Option(n.name,n.name)));if(nodes.some(n=>n.name===was))manage.value=was;const limitSelect=$('limit-name'),previousLimit=limitSelect.value;limitSelect.replaceChildren(...nodes.map(n=>new Option(n.display_name||n.name,n.name)));if(nodes.some(n=>n.name===previousLimit))limitSelect.value=previousLimit;fillLimitsForm();const previousHistoryName=sel.value;sel.replaceChildren(...nodes.map(n=>new Option(n.name,n.name)));if(nodes.some(n=>n.name===previousHistoryName))sel.value=previousHistoryName;drawHistory();}const up=nodes.filter(n=>n.online).length;$('total').textContent=nodes.length;$('online').textContent=up;$('offline').textContent=nodes.length-up;$('overview-hint').textContent='共 '+nodes.length+' 台服务器 · '+up+' 在线 · '+(nodes.length-up)+' 离线';renderOverviewNodes();$('updated').textContent='最后刷新 '+new Date().toLocaleTimeString()}catch(e){$('updated').textContent='获取失败：'+e.message}}
async function drawHistory(){const name=$('history-name').value;if(!name)return;try{const r=await fetch('/api/v1/history?name='+encodeURIComponent(name)+'&hours='+encodeURIComponent($('history-hours').value));if(!r.ok)throw Error('HTTP '+r.status);const data=await r.json();const canvas=$('history-chart');const ctx=canvas.getContext('2d');const w=canvas.width,h=canvas.height;ctx.clearRect(0,0,w,h);ctx.strokeStyle='#34455e';ctx.lineWidth=1;for(let y=0;y<=100;y+=25){const py=20+(100-y)/100*(h-50);ctx.beginPath();ctx.moveTo(40,py);ctx.lineTo(w-15,py);ctx.stroke();ctx.fillStyle='#a6b5c9';ctx.font='12px sans-serif';ctx.fillText(y+'%',4,py+4)}const metric=$('history-metric').value;if(data.length){const lo=data[0].time,hi=Math.max(lo+1,data[data.length-1].time);ctx.beginPath();ctx.strokeStyle='#47a0f5';ctx.lineWidth=2;data.forEach((p,i)=>{const x=40+(p.time-lo)/(hi-lo)*(w-55),y=20+(100-p[metric])/100*(h-50);i?ctx.lineTo(x,y):ctx.moveTo(x,y)});ctx.stroke()}$('history-status').textContent=data.length+' 个采样区间';drawTraffic()}catch(e){$('history-status').textContent=e.message}}
['history-name','history-hours','history-metric'].forEach(k=>document.addEventListener('change',e=>{if(e.target.id===k)drawHistory()}));async function drawTraffic(){const name=$('history-name').value;if(!name)return;try{const r=await fetch('/api/v1/traffic?name='+encodeURIComponent(name)+'&hours='+encodeURIComponent($('history-hours').value));if(!r.ok)throw Error('HTTP '+r.status);const d=await r.json();$('traffic-summary').textContent='接收 ↓ '+formatBytes(d.rx)+' · 发送 ↑ '+formatBytes(d.tx)+'（统计范围内的采样增量）';const max=Math.max(1,...d.buckets.map(b=>b.rx+b.tx));$('traffic-bars').innerHTML=d.buckets.slice(-24).map(b=>'<div class="item"><span>'+new Date(b.time*1000).toLocaleString()+'</span><div class="track"><div class="value" style="width:'+Math.round((b.rx+b.tx)/max*100)+'%"></div></div><span>'+formatBytes(b.rx+b.tx)+'</span></div>').join('')}catch(e){$('traffic-summary').textContent=e.message}}
function switchSettingsTab(tab){document.querySelectorAll('.settings-tab').forEach(b=>{const selected=b.dataset.settingsTab===tab;b.classList.toggle('selected',selected);b.setAttribute('aria-selected',String(selected));b.tabIndex=0});document.querySelectorAll('.settings-pane').forEach(p=>p.hidden=p.id!=='settings-'+tab)}
document.querySelectorAll('.settings-tab').forEach(b=>b.addEventListener('click',()=>switchSettingsTab(b.dataset.settingsTab)));
function toggleAddPanel(show){switchView('settings');switchSettingsTab('nodes');$('add-panel').hidden=!show;if(show){$('server-url').value=$('server-url').value||(location.origin.startsWith('https://')?location.origin:'');$('add-panel').scrollIntoView({behavior:'smooth',block:'nearest'})}}
$('add-node').addEventListener('click',()=>toggleAddPanel(true));
$('show-add-node').addEventListener('click',()=>toggleAddPanel(true));
$('hide-add-node').addEventListener('click',()=>{$('add-panel').hidden=true;$('show-add-node').focus()});
switchSettingsTab('nodes');
$('create-node').addEventListener('submit',async e=>{e.preventDefault();const name=$('node-name').value.trim(),url=$('server-url').value.trim().replace(/\/$/,'');if(!/^[a-zA-Z0-9_-]{1,64}$/.test(name)||!/^https:\/\/[a-zA-Z0-9.-]+(?::[0-9]+)?$/.test(url)){$('create-status').textContent='需要合法节点名称和 HTTPS 域名';return}if(!confirm('确定创建或轮换节点 '+name+' 的令牌吗？旧令牌会立即失效。'))return;try{const r=await fetch('/api/v1/tokens',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name})});if(!r.ok)throw Error('HTTP '+r.status);const data=await r.json();const installer='https://raw.githubusercontent.com/sockc/Monitor/main/scripts/install-release.sh';$('install-command').value='curl -fsSL '+installer+' -o monitor-install.sh && sudo env MONITOR_SERVER='+url+' MONITOR_NODE_NAME='+name+' MONITOR_AGENT_TOKEN='+data.token+' bash monitor-install.sh agent';$('create-status').textContent='专属命令已生成，请妥善保管（令牌仅显示一次）';refresh()}catch(err){$('create-status').textContent=err.message}});
$('copy-command').addEventListener('click',async()=>{try{await navigator.clipboard.writeText($('install-command').value);$('create-status').textContent='已复制到剪贴板'}catch(e){$('create-status').textContent='复制失败，请手动复制'}});
async function manage(action){const name=$('manage-name').value;if(!name)return;let new_name='';if(action==='rename'){new_name=$('rename-target').value.trim();if(!/^[a-zA-Z0-9_-]{1,64}$/.test(new_name)){$('manage-status').textContent='新名称无效';return}}if(!confirm('确认对 '+name+' 执行 '+action+'？此操作可能造成 Agent 离线。'))return;try{const r=await fetch('/api/v1/node-manage',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({action,name,new_name})});if(!r.ok)throw Error(await r.text());$('manage-status').textContent='操作成功。重命名后需重新签发令牌并更新 Agent。';prevNames='';refresh()}catch(e){$('manage-status').textContent=e.message}}
$('rename-node').addEventListener('click',()=>manage('rename'));$('revoke-node').addEventListener('click',()=>manage('revoke'));$('delete-node').addEventListener('click',()=>manage('delete'));
function fillPlanLimits(name){const v=nodeLimitsCache[name]||{};$("metadata-quota").value=v.quota_gb||0;$("metadata-expires").value=v.expires_on||""}
function loadMetadata(){const n=cachedNodes.find(x=>x.name===$("metadata-name").value);if(!n)return;
 const p=n.profile||{};
 $("metadata-display").value=n.display_name||"";
 $("metadata-group").value=n.group||"";
 $("metadata-location").value=n.manual_location||"";
 $("metadata-auto-location").textContent=n.auto_location?"IP 自动识别："+n.auto_location+(n.public_ip?"（"+n.public_ip+"）":"")+" · 清空手动位置即可恢复自动识别":n.public_ip?"公网 IP "+n.public_ip+" · 定位暂不可用":"等待 Agent 上报公网 IP 位置；旧 Agent 需升级";
 $("metadata-notes").value=n.notes||"";
 $("metadata-provider").value=p.provider||"";
 $("metadata-country").value=p.country_code||"";
 $("metadata-price").value=p.price_value||0;
 $("metadata-currency").value=p.price_currency||"USD";
 $("metadata-cycle").value=p.billing_cycle||"year";
 $("metadata-port").value=p.port_mbps||0;
 $("metadata-ipv4").value=String(p.has_ipv4??-1);
 $("metadata-ipv6").value=String(p.has_ipv6??-1);
 $("metadata-start").value=p.period_start||"";
 fillPlanLimits(n.name)
}
$("metadata-name").addEventListener("change",loadMetadata);
$("metadata-form").addEventListener("submit",async e=>{
 e.preventDefault();
 const selected=$("metadata-name").value;
 const payload={name:selected,display_name:$("metadata-display").value,group:$("metadata-group").value,location:$("metadata-location").value,notes:$("metadata-notes").value,
  profile:{provider:$("metadata-provider").value.trim(),country_code:$("metadata-country").value.trim().toUpperCase(),price_value:Number($("metadata-price").value),price_currency:$("metadata-currency").value,billing_cycle:$("metadata-cycle").value,period_start:$("metadata-start").value,port_mbps:Number($("metadata-port").value),has_ipv4:Number($("metadata-ipv4").value),has_ipv6:Number($("metadata-ipv6").value)},
  quota_gb:Number($("metadata-quota").value),expires_on:$("metadata-expires").value
 };
 try{
  const r=await fetch("/api/v1/node-metadata",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify(payload)});
  if(!r.ok)throw Error(await r.text());
  $("metadata-status").textContent="保存成功，套餐与到期数据已同步";
  prevNames="";await loadNodeLimits();await refresh();await refreshPeriods()
 }catch(err){$("metadata-status").textContent="保存失败："+err.message}
});
function updateGroupOptions(nodes){const el=$('node-group'),value=el.value,groups=[...new Set(nodes.map(n=>n.group).filter(Boolean))].sort();const next=['',...groups];if([...el.options].map(x=>x.value).join('|')!==next.join('|')){el.replaceChildren(new Option('全部分组',''),...groups.map(x=>new Option(x,x)));el.value=value}}
function filteredNodes(nodes){const q=$('node-search').value.trim().toLowerCase(),group=$('node-group').value,order=$('node-order').value;const out=nodes.filter(n=>(!group||n.group===group)&&(!q||[n.name,n.display_name,n.group,n.location,n.os,n.notes,n.hostname].some(x=>String(x||'').toLowerCase().includes(q))));if(order==='offline')out.sort((a,b)=>Number(a.online)-Number(b.online));if(order==='cpu')out.sort((a,b)=>b.cpu-a.cpu);return out}
['node-search','node-group','node-order'].forEach(id=>$(id).addEventListener('input',()=>{renderOverviewNodes()}));
let currentLayout='compact';let trafficPeriods={};let nodeLimitsCache={};
try{const saved=localStorage.getItem('monitor-node-layout');if(['cards','compact','list'].includes(saved))currentLayout=saved}catch(e){}
function updateLayout(){const container=$('nodes');container.classList.add('server-grid');container.dataset.layout=currentLayout;document.querySelectorAll('.layout-button').forEach(b=>{const active=b.dataset.layout===currentLayout;b.classList.toggle('active',active);b.setAttribute('aria-pressed',String(active))})}
document.querySelectorAll('.layout-button').forEach(b=>b.addEventListener('click',()=>{currentLayout=b.dataset.layout;try{localStorage.setItem('monitor-node-layout',currentLayout)}catch(e){}updateLayout();renderOverviewNodes()}));
updateLayout();
function capacity(n){return n>0?formatBytes(n):'—'}
function hardware(n){return '<div class="hardware-line"><span>'+((n.cpu_cores>0?n.cpu_cores+'C':'—')+' CPU')+'</span><span>内存 '+capacity(n.memory_total)+'</span><span>磁盘 '+capacity(n.disk_total)+'</span></div>'}
function resourceMetric(label,val,used,total){const v=Math.min(100,Math.max(0,Number(val)||0)),status=v>=90?'critical':v>=75?'hot':'',quantity=total>0?capacity(used)+' / '+capacity(total):'';return '<div class="resource-line"><div class="metric-row"><span>'+label+'</span><span class="resource-numbers">'+(quantity?'<small>'+quantity+'</small>':'')+'<b>'+Math.round(v)+'%</b></span></div><div class="bar"><div class="fill '+status+'" style="width:'+v+'%"></div></div></div>'}
function resourceWithCapacity(label,pct,used,total){return resourceMetric(label,pct,used,total)}
function renderOverviewNodes(){const visible=filteredNodes(cachedNodes);$('nodes').innerHTML=visible.map(renderNode).join('')||'<div class="dashboard-empty"><strong>没有匹配的服务器</strong><span>试试修改搜索关键词或切换分组</span></div>';$('node-visible-count').textContent='显示 '+visible.length+' / '+cachedNodes.length+' 台'}

function dateInZone(zone){try{const d=new Intl.DateTimeFormat('en-GB',{timeZone:zone||'UTC',year:'numeric',month:'2-digit',day:'2-digit'}).formatToParts(new Date());const get=x=>d.find(y=>y.type===x).value;return get('year')+'-'+get('month')+'-'+get('day')}catch(e){return new Date().toISOString().slice(0,10)}}
function daysUntil(date,zone){if(!date)return null;return Math.round((Date.parse(date+'T00:00:00Z')-Date.parse(dateInZone(zone)+'T00:00:00Z'))/86400000)}
function periodBox(n){const p=trafficPeriods[n.name];const has=Boolean(p&&p.has_samples);const mk=(label,rx,tx)=>'<div class="traffic-cell"><small>'+label+'</small><strong>'+(has?formatBytes(rx+tx):'—')+'</strong></div>';let out='<div class="period-traffic">'+mk('今日流量',Number(p?.today_rx)||0,Number(p?.today_tx)||0)+mk('本月流量',Number(p?.month_rx)||0,Number(p?.month_tx)||0)+'</div>';
if(p&&p.quota_gb>0){const used=Number(p.month_rx)+Number(p.month_tx),ratio=used/(p.quota_gb*1000000000)*100;out+='<div class="quota-caption"><span>月流量额度</span><b>'+Math.min(999,Math.max(0,ratio)).toFixed(0)+'%</b></div><div class="quota-track"><div class="quota-fill '+(ratio>=100?'high':ratio>=80?'warn':'')+'" style="width:'+Math.max(0,Math.min(ratio,100))+'%"></div></div>'}
const badges=[];const left=daysUntil(p?.expires_on,p?.timezone);if(left!==null&&left<=30)badges.push('<span class="alert-chip '+(left<0?'urgent':'')+'">'+(left<0?'已过期 '+(-left)+' 天':left===0?'今天到期':'到期剩余 '+left+' 天')+'</span>');
if(has&&p.quota_gb>0&&(Number(p.month_rx)+Number(p.month_tx))/(p.quota_gb*1000000000)>=0.8)badges.push('<span class="alert-chip">流量接近额度</span>');
if(n.online&&(n.cpu>=90||n.memory>=90||n.disk>=90))badges.push('<span class="alert-chip urgent">资源使用偏高</span>');
if(badges.length)out+='<div class="node-badges">'+badges.join('')+'</div>';
return out}
function renderNode(n){const name=esc(n.display_name||n.name),id=esc(n.name),online=Boolean(n.online),loc=n.location?esc(n.location):'等待 IP 定位',os=esc(n.os||n.arch||'系统未知');
const rx=online?formatBytes(n.rx_speed)+'/s':'—',tx=online?formatBytes(n.tx_speed)+'/s':'—';
return '<article role="button" tabindex="0" data-node="'+id+'" class="node dashboard-node '+(online?'':'node-offline')+'" aria-label="查看 '+name+' 详情"><div class="node-head"><div class="node-name-wrap"><div class="node-title">'+name+'</div><div class="node-location" title="'+loc+'"><span class="location-symbol">⌖</span>'+loc+'</div><div class="node-subtitle">'+os+'</div></div><span class="pill '+(online?'ok':'bad')+'">'+(online?'● 在线':'● 离线')+'</span></div>'+hardware(n)+'<div class="meters">'+(online?resourceMetric('CPU',n.cpu)+resourceMetric('内存',n.memory,n.memory_used,n.memory_total)+resourceMetric('磁盘',n.disk,n.disk_used,n.disk_total):'<div class="offline-hint">最近上报 '+(n.last_seen?esc(new Date(n.last_seen).toLocaleString()):'未知')+'</div>')+'</div><div class="meta node-activity"><span class="rx">↓ '+rx+'</span><span class="tx">↑ '+tx+'</span><span class="uptime">运行 '+(online?Math.floor(n.uptime/86400)+' 天':'—')+'</span></div>'+periodBox(n)+'</article>'}

async function alertsLoad(){try{const r=await fetch('/api/v1/alerts');if(!r.ok)return;const d=await r.json(),a=d.settings;$('a-offline').value=a.offline_seconds;$('a-cpu').value=a.cpu_threshold;$('a-memory').value=a.memory_threshold;$('a-disk').value=a.disk_threshold;$('a-duration').value=a.duration_seconds;$('a-webhook').value=a.webhook||'';renderAlerts(d.events)}catch(e){$('alerts-status').textContent=e.message}}
function alertName(k){return ({quota_80:'月流量达到 80%',quota_90:'月流量达到 90%',quota_100:'月流量达到 100%',expiry_30:'30 天内到期',expiry_15:'15 天内到期',expiry_7:'7 天内到期',expiry_overdue:'服务器已到期'})[k]||k}
function renderAlerts(items){const active=items.filter(e=>!e.end).length;$('alert-count').textContent=active;const root=$('alert-events');root.replaceChildren();if(!items.length){const empty=document.createElement('div');empty.className='event-empty';empty.textContent='暂无告警记录 · 服务器状态正常时无需处理';root.append(empty);return}for(const e of items.slice(0,60)){const item=document.createElement('div');item.className='event-item'+(e.end?' recovered':'');const dot=document.createElement('span');dot.className='event-dot';const content=document.createElement('div');content.className='event-content';const title=document.createElement('strong');title.textContent=alertName(e.kind);const sub=document.createElement('small');sub.textContent=e.node+' · '+new Date(e.start*1000).toLocaleString();content.append(title,sub);const status=document.createElement('span');status.className='event-status';status.textContent=e.end?'已恢复':'告警中';item.append(dot,content,status);root.append(item)}}
$('alerts-form').addEventListener('submit',async e=>{e.preventDefault();const payload={offline_seconds:Number($('a-offline').value),cpu_threshold:Number($('a-cpu').value),memory_threshold:Number($('a-memory').value),disk_threshold:Number($('a-disk').value),duration_seconds:Number($('a-duration').value),webhook:$('a-webhook').value.trim()};try{const r=await fetch('/api/v1/alerts',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(payload)});if(!r.ok)throw Error(await r.text());$('alerts-status').textContent='已保存'}catch(err){$('alerts-status').textContent=err.message}});
let selectedNode=null;
function openNode(name){selectedNode=name;switchView('detail');switchDetailTab('performance');updateDetail();window.scrollTo(0,0)}
function detailRow(label,value){return '<div><span class="muted">'+esc(label)+'</span><br><strong>'+esc(value)+'</strong></div>'}
function updateDetail(){if(!selectedNode)return;const n=cachedNodes.find(x=>x.name===selectedNode);if(!n){$('detail-title').textContent='节点不存在';$('detail-status').textContent='此节点可能已删除';return}
$('detail-title').textContent=n.display_name||n.name;$('detail-subtitle').textContent=(n.group||'未分组')+' · '+n.name;
$('detail-status').textContent=n.online?'● 在线 · 监控数据实时更新':'● 离线 · 最后上报 '+(n.last_seen?new Date(n.last_seen).toLocaleString():'未知');
$('detail-metrics').innerHTML=[['CPU',n.cpu],['内存',n.memory],['磁盘',n.disk]].map(x=>'<div><small>'+x[0]+'</small><strong>'+(n.online?Math.round(x[1])+'%':'—')+'</strong></div>').join('');
$('detail-info').innerHTML=detailRow('服务器位置',n.location||'尚未定位')+detailRow('位置来源',n.location_source==='manual'?'手动覆盖':n.location_source==='ip'?'公网 IP 自动识别':'未识别')+detailRow('公网出口 IP',n.public_ip||'—')+detailRow('主机名',n.hostname||'—')+detailRow('系统发行版',n.os||'—')+detailRow('Agent 版本',n.agent_version||'未上报')+detailRow('CPU 逻辑核心',n.cpu_cores>0?n.cpu_cores+' C':'—')+detailRow('CPU 型号',n.cpu_model||'—')+detailRow('负载 1/5/15 分钟',[n.load1,n.load5,n.load15].every(x=>x!==undefined)?[n.load1,n.load5,n.load15].map(x=>Number(x).toFixed(2)).join(' / '):'—')+detailRow('Swap',n.swap_total>0?capacity(n.swap_used)+' / '+capacity(n.swap_total):'无或未上报')+detailRow('内存使用',n.memory_total>0?capacity(n.memory_used)+' / '+capacity(n.memory_total):'—')+detailRow('磁盘使用',n.disk_total>0?capacity(n.disk_used)+' / '+capacity(n.disk_total):'—')+detailRow('架构',n.arch||'—')+detailRow('运行时长',n.uptime?Math.floor(n.uptime/86400)+' 天':'—')+detailRow('最后上报',n.last_seen?new Date(n.last_seen).toLocaleString():'—')+(n.notes?detailRow('备注',n.notes):'');
const p=trafficPeriods[n.name];$('detail-network-info').innerHTML=detailRow('实时下载',n.online?formatBytes(n.rx_speed)+'/s':'—')+detailRow('实时上传',n.online?formatBytes(n.tx_speed)+'/s':'—')+detailRow('累计接收（网卡计数）',formatBytes(n.rx_bytes))+detailRow('累计发送（网卡计数）',formatBytes(n.tx_bytes))+detailRow('今日流量',p&&p.has_samples?formatBytes(p.today_rx+p.today_tx):'—')+detailRow('本月流量',p&&p.has_samples?formatBytes(p.month_rx+p.month_tx):'—')+detailRow('统计时区',p?p.timezone||'UTC':'UTC')+detailRow('月流量额度',p&&p.quota_gb>0?p.quota_gb+' GB':'未设置')+detailRow('磁盘读取',n.online?formatBytes(n.disk_read_speed)+'/s':'—')+detailRow('磁盘写入',n.online?formatBytes(n.disk_write_speed)+'/s':'—');
}
async function updateDetailHistory(){if(!selectedNode)return;const name=selectedNode,hours=$('detail-hours').value;try{const rs=await Promise.all([fetch('/api/v1/history?name='+encodeURIComponent(name)+'&hours='+hours),fetch('/api/v1/traffic?name='+encodeURIComponent(name)+'&hours='+hours)]);if(!rs.every(r=>r.ok))throw Error('数据接口异常');const points=await rs[0].json(),traffic=await rs[1].json();if(selectedNode!==name)return;const ctx=$('detail-chart').getContext('2d'),w=900,h=240;ctx.clearRect(0,0,w,h);const colors=['#60a5fa','#a78bfa','#2dd4bf'];ctx.strokeStyle='#34455e';ctx.lineWidth=1;for(let value=0;value<=100;value+=25){const y=18+(100-value)/100*196;ctx.beginPath();ctx.moveTo(42,y);ctx.lineTo(884,y);ctx.stroke();ctx.fillStyle='#9eb0c9';ctx.font='12px sans-serif';ctx.fillText(value+'%',6,y+4)}
if(points.length>0){const min=points[0].time,max=Math.max(min+1,points[points.length-1].time);['cpu','memory','disk'].forEach((key,k)=>{ctx.beginPath();ctx.strokeStyle=colors[k];ctx.lineWidth=2;points.forEach((p,i)=>{const x=42+(p.time-min)/(max-min)*842,y=18+(100-Math.max(0,Math.min(100,p[key])))/100*196;if(!i)ctx.moveTo(x,y);else ctx.lineTo(x,y)});ctx.stroke()})}
$('detail-chart-status').textContent='CPU（蓝） · 内存（紫） · 磁盘（青） · '+points.length+' 个区间';
$('detail-traffic').textContent='接收 ↓ '+formatBytes(traffic.rx)+'　发送 ↑ '+formatBytes(traffic.tx)+'（采样估算）'}catch(e){$('detail-chart-status').textContent='历史数据加载失败：'+e.message}}
$('nodes').addEventListener('click',e=>{const card=e.target.closest('[data-node]');if(card)openNode(card.dataset.node)});
$('nodes').addEventListener('keydown',e=>{if(e.key==='Enter'||e.key===' '){const card=e.target.closest('[data-node]');if(card){e.preventDefault();openNode(card.dataset.node)}}});
function switchDetailTab(tab){document.querySelectorAll('.detail-tab').forEach(b=>{const on=b.dataset.detailTab===tab;b.classList.toggle('selected',on);b.setAttribute('aria-selected',String(on));b.tabIndex=0});document.querySelectorAll('.detail-pane').forEach(p=>p.hidden=p.id!=='detail-'+tab)}
document.querySelectorAll('.detail-tab').forEach(b=>b.addEventListener('click',()=>switchDetailTab(b.dataset.detailTab)));
switchDetailTab('performance');
$('refresh-alerts').addEventListener('click',refreshAlertEvents);
$('detail-back').addEventListener('click',()=>switchView('overview'));
$('detail-hours').addEventListener('change',updateDetailHistory);
async function loadNodeLimits(){try{const r=await fetch('/api/v1/node-limits',{cache:'no-store'});if(!r.ok)throw Error('HTTP '+r.status);nodeLimitsCache=await r.json();fillLimitsForm();if($('metadata-name').value)fillPlanLimits($('metadata-name').value)}catch(e){$('limit-status').textContent='加载失败：'+e.message}}
function fillLimitsForm(){const name=$('limit-name').value;if(!name)return;const v=nodeLimitsCache[name]||{};$('limit-timezone').value=v.timezone||'UTC';$('limit-quota').value=v.quota_gb||0;$('limit-expires').value=v.expires_on||''}
$('limit-name').addEventListener('change',fillLimitsForm);
$('node-limits-form').addEventListener('submit',async e=>{e.preventDefault();const payload={name:$('limit-name').value,timezone:$('limit-timezone').value.trim(),quota_gb:Number($('limit-quota').value),expires_on:$('limit-expires').value};try{const r=await fetch('/api/v1/node-limits',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(payload)});if(!r.ok)throw Error(await r.text());$('limit-status').textContent='已保存';await loadNodeLimits();await refreshPeriods()}catch(err){$('limit-status').textContent='保存失败：'+err.message}});
function switchView(view){for(const el of document.querySelectorAll('.view'))el.hidden=el.id!=='view-'+view;for(const b of document.querySelectorAll('.nav-btn')){const active=b.dataset.view===view;b.classList.toggle('selected',active);b.setAttribute('aria-current',active?'page':'false')}if(view==='statistics')drawHistory();if(view==='alerts')refreshAlertEvents();if(view==='detail')updateDetailHistory()}
document.querySelectorAll('.nav-btn').forEach(b=>b.addEventListener('click',()=>switchView(b.dataset.view)));
async function refreshAlertEvents(){try{const r=await fetch('/api/v1/alerts');if(r.ok)renderAlerts((await r.json()).events)}catch(e){}}
switchView('overview');async function refreshPeriods(){try{const r=await fetch('/api/v1/traffic-summary',{cache:'no-store'});if(!r.ok)return;trafficPeriods=await r.json();renderOverviewNodes();if(selectedNode)updateDetail()}catch(e){}}
refresh();alertsLoad();loadNodeLimits();refreshPeriods();setInterval(refresh,5000);setInterval(refreshPeriods,15000);setInterval(drawHistory,30000);setInterval(async()=>{try{const r=await fetch('/api/v1/alerts');if(r.ok)renderAlerts((await r.json()).events)}catch(e){}},30000);`
