package main

// Monitor stylesheet segment; concatenated in original cascade order.
const styleCSSThemes=`/* V0.9.5 selectable themes: default pastel glass reference, plus dark/light. */
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
