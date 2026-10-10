package main

// Monitor stylesheet segment; concatenated in original cascade order.
const styleCSSAdmin=`/* V0.9.8 immutable node identities and editable multi-node inventory */
#metadata-advanced[hidden],#node-rebind-panel[hidden],#add-panel[hidden],#view-settings .node-identity-readonly[hidden]{display:none!important}
#view-settings .node-editor-head,#view-settings .node-edit-row{display:grid;grid-template-columns:minmax(130px,1fr) minmax(135px,1.1fr) minmax(90px,.7fr) minmax(125px,1fr) minmax(230px,1.7fr);align-items:center;gap:9px}
#view-settings .node-editor-head{font-size:10px;letter-spacing:.2px;color:#829bb8;padding:0 11px 10px;border-bottom:1px solid #29425f}
#view-settings .node-editor-list{display:flex;flex-direction:column;gap:6px;margin-top:10px;min-width:0}
#view-settings .node-edit-row{padding:10px 10px 12px;background:#102136;border:1px solid #2b4362;border-radius:11px;min-width:0}
#view-settings .edit-node-identity{display:flex;align-items:center;gap:9px;min-width:0}
#view-settings .edit-node-identity>span:last-child{display:flex;flex-direction:column;gap:4px;min-width:0}
#view-settings .edit-node-identity strong{font-size:12px;font-weight:650;color:#e5f1ff;overflow:hidden;white-space:nowrap;text-overflow:ellipsis}
#view-settings .edit-node-identity small{font-size:9px;color:#7794b5;overflow:hidden;white-space:nowrap;text-overflow:ellipsis}
#view-settings .edit-node-state{height:8px;width:8px;flex:0 0 8px;border-radius:50%;background:#ef8290}
#view-settings .edit-node-state.on{background:#35d6a1}
#view-settings .node-edit-row label{min-width:0;display:block}
#view-settings .node-edit-row input{box-sizing:border-box;width:100%;min-width:0;border-radius:8px;background:#0c1a2e;color:#eaf4ff;border:1px solid #344e6c;padding:8px 9px;font-size:12px}
#view-settings .node-edit-actions{display:flex;align-items:center;justify-content:flex-end;gap:4px;flex-wrap:wrap;min-width:0}
#view-settings .node-edit-actions button{padding:7px 8px;border-radius:7px;font-size:10px;line-height:1.3}
#view-settings .node-edit-actions .danger-btn{background:transparent;color:#e6a2ab;border-color:#694152}
#view-settings .node-edit-mobile-label{display:none;font-size:10px;color:#839cb8;margin-bottom:4px}
#view-settings .row-feedback{grid-column:1/-1;font-size:11px;color:#72d9a9;min-height:0}
#view-settings .row-feedback:empty{display:none}
#view-settings #node-rebind-command,#view-settings #install-command{width:100%;box-sizing:border-box;background:#0c1a2e;color:#e6f2ff;border:1px solid #354c69;border-radius:8px;padding:11px;font-size:11px;resize:vertical;line-height:1.5}
#view-settings .settings-card:has(.node-editor-list){padding:15px!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .node-edit-row{background:#ffffff93;border-color:#ffffffbc}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .node-editor-head{border-bottom-color:#ffffffbd;color:#687a91}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .node-edit-row input,
:is(html[data-theme=glass],html[data-theme=light]) #view-settings #node-rebind-command{background:#fffc!important;color:#2d4160!important;border:1px solid #d5deea!important}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .edit-node-identity strong{color:#2d405f}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .edit-node-identity small{color:#7f8ea5}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .node-edit-actions .danger-btn{color:#b94b66;border-color:#dfacb8}
:is(html[data-theme=glass],html[data-theme=light]) #view-settings .node-edit-mobile-label{color:#72829b}
@media(max-width:1380px){#view-settings .node-editor-head{display:none}#view-settings .node-edit-row{grid-template-columns:repeat(3,minmax(0,1fr));gap:10px}#view-settings .node-edit-row .edit-node-identity{grid-column:1/-1}#view-settings .node-edit-actions{grid-column:1/-1;justify-content:flex-start}#view-settings .node-edit-mobile-label{display:block}}
@media(max-width:720px){#view-settings .node-edit-row{grid-template-columns:repeat(2,minmax(0,1fr));gap:10px}#view-settings .node-edit-row .edit-node-identity{grid-column:1/-1}#view-settings .node-edit-row label:has(.row-display){grid-column:1/-1}#view-settings .node-edit-actions{grid-column:1/-1}}
@media(max-width:420px){#view-settings .node-edit-row{grid-template-columns:1fr}#view-settings .node-edit-row>*{grid-column:1/-1}#view-settings .node-edit-actions button{flex:1 1 auto}}

/* V0.9.9 compact boot-time stamp and floating machine-info popover. */
#view-overview .server-grid .dashboard-node.glass-node{position:relative!important}
#view-overview .glass-node .card-identity{box-sizing:border-box;padding:0 23px 0 7px;min-width:0}
#view-overview .glass-node .node-info-trigger{position:absolute;top:9px;right:9px;z-index:2;display:grid;place-items:center;width:24px;height:24px;min-width:24px;padding:0;border:1px solid #ffffff91;border-radius:50%;background:#ffffff7d;color:#496181;font-size:14px;font-weight:800;line-height:1;cursor:pointer;box-shadow:0 2px 7px #1b30401c;transition:background .15s,border-color .15s,transform .15s}
#view-overview .glass-node .node-info-trigger:hover,#view-overview .glass-node .node-info-trigger[aria-expanded=true],#view-overview .glass-node .node-info-trigger:focus-visible{background:#d8e9ff;border-color:#8bb4ed;color:#244d82;transform:scale(1.05);outline-offset:2px}
#view-overview .glass-node .node-boot-line{display:flex;align-items:center;justify-content:space-between;flex-wrap:wrap;gap:3px 8px;margin:10px 0 0;padding-top:8px;border-top:1px solid #ffffff61;color:var(--node-muted);font-size:10px;line-height:1.5;text-align:left}
#view-overview .glass-node .node-boot-line time{color:var(--node-ink);font-size:10px;font-variant-numeric:tabular-nums}
#view-overview .node-quick-info[hidden]{display:none!important}
#view-overview .node-quick-info{position:fixed;z-index:9999;box-sizing:border-box;width:min(355px,calc(100vw - 20px));max-height:calc(100vh - 20px);overflow:auto;top:0;left:0;padding:13px 14px 11px;background:#f8fafff2;border:1px solid #ffffffee;box-shadow:0 16px 40px #152c4952;border-radius:14px;backdrop-filter:blur(24px);-webkit-backdrop-filter:blur(24px);color:#2f415a;text-align:left}
#view-overview .node-quick-info .node-info-heading{display:flex;align-items:flex-start;justify-content:space-between;gap:9px;padding-bottom:10px;border-bottom:1px solid #dce4ee}
#view-overview .node-quick-info .node-info-heading>div{display:flex;flex-direction:column;gap:4px;min-width:0}
#view-overview .node-quick-info .node-info-heading strong{font-size:14px;color:#273b58;word-break:break-word}
#view-overview .node-quick-info .node-info-heading small{font-size:10px;color:#647b97}
#view-overview .node-quick-info .node-info-close{width:23px;height:23px;flex:0 0 23px;display:grid;place-items:center;padding:0;background:#e4ecf7;border:0;border-radius:7px;color:#4c6686;font-size:18px;line-height:1;cursor:pointer}
#view-overview .node-quick-info .node-info-close:hover{background:#d1e0f5}
#view-overview .node-quick-info .node-info-grid{margin:10px 0 8px;display:grid;gap:0}
#view-overview .node-quick-info .info-line{display:flex;justify-content:space-between;align-items:baseline;gap:9px;padding:6px 1px;border-bottom:1px solid #dfe6ef9c}
#view-overview .node-quick-info dt{font-size:11px;flex:0 0 auto;color:#708199}
#view-overview .node-quick-info dd{font-size:11px;flex:1;min-width:0;text-align:right;line-height:1.4;color:#2c425f;margin:0;overflow-wrap:anywhere}
#view-overview .node-quick-info .node-info-foot{font-size:10px;color:#7a8a9f;margin:8px 0 0;line-height:1.5}
html[data-theme=dark] #view-overview .glass-node .node-info-trigger{color:#b3d2f6;background:#2a3c58;border:1px solid #4e6284}
html[data-theme=dark] #view-overview .glass-node .node-info-trigger:hover,html[data-theme=dark] #view-overview .glass-node .node-info-trigger[aria-expanded=true]{background:#335b83;color:white}
html[data-theme=dark] #view-overview .glass-node .node-boot-line{border-top-color:#2c425e}
html[data-theme=dark] #view-overview .node-quick-info{background:#17273df5;border:1px solid #45607c;color:#ecf3ff;box-shadow:0 16px 40px #0009}
html[data-theme=dark] #view-overview .node-quick-info .node-info-heading{border-bottom-color:#344862}
html[data-theme=dark] #view-overview .node-quick-info .node-info-heading strong{color:#f0f6ff}
html[data-theme=dark] #view-overview .node-quick-info .node-info-heading small,html[data-theme=dark] #view-overview .node-quick-info dt,html[data-theme=dark] #view-overview .node-quick-info .node-info-foot{color:#a3b5ce}
html[data-theme=dark] #view-overview .node-quick-info dd{color:#e3effe}
html[data-theme=dark] #view-overview .node-quick-info .info-line{border-bottom-color:#31435b}
html[data-theme=dark] #view-overview .node-quick-info .node-info-close{background:#253a55;color:#d2e5ff}
@media(hover:none),(pointer:coarse){#view-overview .glass-node .node-info-trigger{height:29px;width:29px;min-width:29px;top:8px;right:8px;font-size:16px}#view-overview .glass-node .card-identity{padding-right:32px}}
@media(max-width:420px){#view-overview .glass-node .node-boot-line{font-size:9px}#view-overview .node-quick-info{padding:12px 11px}}

`
