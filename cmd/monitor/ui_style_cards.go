package main

// Monitor stylesheet segment; concatenated in original cascade order.
const styleCSSCards=`/* 0.9.6 reference-inspired centered glass VPS cards; four per desktop row. */
html{--node-ink:#eaf3ff;--node-muted:#a4b6ce;--node-inlay:#0f2035;--node-track:#273d56;--node-stroke:#324862}
html[data-theme=glass]{--node-ink:#273348;--node-muted:#66738b;--node-inlay:#ffffffa6;--node-track:#ffffffa8;--node-stroke:#ffffffca}
html[data-theme=light]{--node-ink:#263852;--node-muted:#697f98;--node-inlay:#f9fbff;--node-track:#dfe7f2;--node-stroke:#d8e2ee}
#view-overview .server-grid .dashboard-node.glass-node{display:flex!important;flex-direction:column!important;gap:0!important;text-align:center;min-height:0;overflow:hidden;padding:14px 11px 13px!important;border-radius:20px!important;color:var(--node-ink)!important}
#view-overview .glass-node .card-identity{display:flex;align-items:center;justify-content:center;gap:6px;width:100%;min-height:43px;min-width:0;margin:0 0 6px}
#view-overview .glass-node .status-dot{width:8px;height:8px;border-radius:50%;background:#f07389;flex:0 0 8px;box-shadow:0 0 0 3px #ef738320}
#view-overview .glass-node .status-dot.on{background:#2acb82;box-shadow:0 0 0 3px #2acb8220}
#view-overview .glass-node .country-flag{font-size:20px;line-height:1;flex:0 0 auto;filter:saturate(1.1)}
#view-overview .glass-node .identity-text{display:flex;flex-direction:column;gap:3px;min-width:0;max-width:calc(100% - 92px)}
#view-overview .glass-node .identity-text strong{font-weight:770;line-height:1.27;font-size:13px;color:var(--node-ink);white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
#view-overview .glass-node .identity-text small{color:var(--node-muted);font-size:10px;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
#view-overview .glass-node .state-text{font-size:10px;color:#24ad73;flex:0 0 auto;white-space:nowrap}
#view-overview .glass-node.node-offline .state-text{color:#e07683}
#view-overview .glass-node .card-plan{display:flex;align-items:center;justify-content:center;flex-wrap:wrap;column-gap:9px;row-gap:4px;min-height:26px;margin:0 0 13px;color:var(--node-muted);font-size:10px}
#view-overview .glass-node .plan-price{max-width:100%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#view-overview .glass-node .plan-remaining{white-space:nowrap}
#view-overview .glass-node .expiry-progress{height:5px;width:41px;min-width:34px;max-width:42px;border-radius:10px;background:var(--node-track);display:inline-block;overflow:hidden}
#view-overview .glass-node .expiry-progress i{display:block;height:100%;background:#f3a257;border-radius:10px}
#view-overview .glass-node .metric-five{display:grid;grid-template-columns:repeat(5,minmax(0,1fr));gap:5px;margin:0 0 12px;width:100%}
#view-overview .glass-node .five-metric{display:flex;flex-direction:column;align-items:center;min-width:0;gap:4px;overflow:hidden}
#view-overview .glass-node .five-metric small{font-size:10px;line-height:1.1;color:var(--node-muted);white-space:nowrap}
#view-overview .glass-node .five-metric strong{font-size:12px;font-weight:770;line-height:1.4;color:var(--node-ink);white-space:nowrap;max-width:100%;overflow:hidden;text-overflow:ellipsis}
#view-overview .glass-node .five-track{height:4px;border-radius:6px;width:90%;background:var(--node-track);overflow:hidden}
#view-overview .glass-node .five-track.ghost{background:transparent}
#view-overview .glass-node .five-track i{display:block;height:100%;background:#56cc9a;border-radius:6px}
#view-overview .glass-node .five-metric:nth-child(2) .five-track i{background:#f1a25b}
#view-overview .glass-node .five-metric:nth-child(3) .five-track i{background:#4bcc8b}
#view-overview .glass-node .five-track i.warm{background:#f59e4d!important}
#view-overview .glass-node .five-track i.critical{background:#f16d7d!important}
#view-overview .glass-node .card-cumulative{display:grid;grid-template-columns:repeat(2,minmax(0,1fr));gap:7px;margin-bottom:12px;min-width:0}
#view-overview .glass-node .card-cumulative>div{background:var(--node-inlay);border:1px solid var(--node-stroke);padding:7px 5px;border-radius:11px;min-width:0;display:flex;flex-direction:column;gap:3px;align-items:center}
#view-overview .glass-node .card-cumulative small{font-size:9px;color:var(--node-muted);font-weight:600}
#view-overview .glass-node .card-cumulative strong{font-size:12px;font-weight:760;color:var(--node-ink);white-space:nowrap;overflow:hidden;text-overflow:ellipsis;max-width:100%}
#view-overview .glass-node .capability-tags{display:flex;align-items:center;justify-content:center;flex-wrap:wrap;gap:4px;min-height:25px;margin-bottom:12px}
#view-overview .glass-node .plan-tag{display:inline-block;padding:4px 6px;border-radius:7px;white-space:nowrap;line-height:1.25;font-size:9px;font-weight:640}
#view-overview .glass-node .tag-bandwidth{background:#2c81f0;color:white}
#view-overview .glass-node .tag-quota{background:#1eb782;color:white}
#view-overview .glass-node .tag-ip4{background:#9556db;color:white}
#view-overview .glass-node .tag-ip6{background:#d753a5;color:white}
#view-overview .glass-node .tag-muted{background:#98a6b87a;color:var(--node-ink)}
#view-overview .glass-node .tag-off{background:#d8a4a474;color:var(--node-ink)}
#view-overview .glass-node .card-quota{margin-top:auto;display:flex;flex-direction:column;gap:5px}
#view-overview .glass-node .quota-figures{display:flex;align-items:baseline;justify-content:space-between;gap:6px;min-width:0;color:var(--node-ink)}
#view-overview .glass-node .quota-figures>span{font-size:11px;line-height:1.25;white-space:nowrap;overflow:hidden;text-overflow:ellipsis}
#view-overview .glass-node .quota-figures em{font-style:normal;color:var(--node-muted)}
#view-overview .glass-node .quota-figures b{font-size:11px;color:var(--node-muted);white-space:nowrap;font-weight:550}
#view-overview .glass-node .quota-linear{height:7px;background:var(--node-track);overflow:hidden;border-radius:8px;box-shadow:inset 0 1px 2px #3333330d}
#view-overview .glass-node .quota-linear i{display:block;height:100%;background:#56cbae;border-radius:8px}
#view-overview .glass-node .quota-linear i.warn{background:#f1a154}
#view-overview .glass-node .quota-linear i.danger{background:#e86f83}
#view-overview .glass-node .card-alerts{display:flex;justify-content:center;gap:4px;flex-wrap:wrap;margin-top:6px}
#view-overview .glass-node .card-alerts span{font-size:9px;color:#bc5761;background:#ffdede6e;padding:3px 5px;border-radius:6px}
html[data-theme=glass] #view-overview .server-grid .dashboard-node.glass-node{background:#fffdfba7!important;backdrop-filter:blur(23px);-webkit-backdrop-filter:blur(23px);border:1px solid #ffffffd4!important;box-shadow:0 8px 26px #55455a23!important}
html[data-theme=light] #view-overview .server-grid .dashboard-node.glass-node{background:#fff!important;box-shadow:0 4px 20px #34455a14!important}
html[data-theme=dark] #view-overview .server-grid .dashboard-node.glass-node{background:#17283f!important;box-shadow:0 6px 22px #060f1c6b!important}
#view-settings .form-section-title{display:flex;flex-direction:column;gap:3px;border-top:1px solid #374963;padding-top:14px;margin-top:4px}
#view-settings .form-section-title strong{font-size:13px;color:#e8f1ff}
#view-settings .form-section-title small{font-size:10px;color:#9baec5}
html[data-theme=glass] #view-settings .form-section-title,html[data-theme=light] #view-settings .form-section-title{border-color:#ffffffd0}
html[data-theme=glass] #view-settings .form-section-title strong,html[data-theme=light] #view-settings .form-section-title strong{color:#344861}
html[data-theme=glass] #view-settings .form-section-title small,html[data-theme=light] #view-settings .form-section-title small{color:#65768d}
@media(max-width:420px){#view-overview .glass-node .five-metric strong{font-size:11px}#view-overview .glass-node .country-flag{font-size:18px}#view-overview .glass-node .plan-tag{font-size:9px;padding:4px 6px}}

/* 0.9.7: hide missing values, accurate self-hosted flags and slim upload/download. */
#view-overview .glass-node .country-flag{width:24px;height:18px;flex:0 0 24px;display:inline-flex;align-items:center;justify-content:center;font-size:0;overflow:hidden;border-radius:3px;box-shadow:0 1px 3px #192c4c35;background:#ffffffa8}
#view-overview .glass-node .country-flag .flag-image{display:block;object-fit:cover;object-position:center;width:24px;height:18px;border-radius:3px}
#view-overview .glass-node .card-cumulative{display:flex!important;flex-direction:row!important;align-items:center;justify-content:space-between;gap:7px!important;min-height:28px;margin:0 0 11px!important;padding:6px 9px;border-radius:10px;border:1px solid var(--node-stroke);background:var(--node-inlay)}
#view-overview .glass-node .card-cumulative>span{display:flex;align-items:center;justify-content:center;gap:4px;flex:1 1 0;min-width:0;color:var(--node-muted);font-size:10px;white-space:nowrap}
#view-overview .glass-node .card-cumulative b{color:var(--node-ink);font-size:11px;font-weight:750;overflow:hidden;text-overflow:ellipsis}
#view-overview .glass-node .card-plan{margin-bottom:10px;min-height:0}
#view-overview .glass-node .capability-tags{min-height:0;margin-bottom:11px}
#view-overview .glass-node .card-identity{margin-bottom:9px}
#view-overview .glass-node .quota-linear{margin-top:2px}
@media(max-width:390px){#view-overview .glass-node .card-cumulative{padding:6px}#view-overview .glass-node .card-cumulative>span{font-size:9px}#view-overview .glass-node .card-cumulative b{font-size:10px}}

`
