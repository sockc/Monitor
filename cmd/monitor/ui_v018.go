package main

// Responsive TCP quality and heartbeat coverage display.
const uiV018CSS = `#view-detail .network-quality-panel{min-width:0}
#view-detail .network-quality-kpis{display:grid;grid-template-columns:repeat(4,minmax(0,1fr));gap:9px;margin:9px 0 10px}
#view-detail .network-quality-kpis>div{min-width:0;border:1px solid var(--t-border,#334c69);background:var(--t-card,#13253a);border-radius:10px;padding:13px 14px;display:flex;flex-direction:column;gap:7px}
#view-detail .network-quality-kpis small{font-size:11px;color:var(--t-muted,#92abc8)}
#view-detail .network-quality-kpis strong{font-size:20px;color:var(--t-text,#eaf2ff);font-variant-numeric:tabular-nums;overflow-wrap:anywhere}
#view-detail .network-quality-panel .field-hint{font-size:11px;line-height:1.55;margin:5px 0 16px}
#view-detail .network-quality-panel .panel-title{margin-top:8px}
#view-detail #uptime-windows{font-variant-numeric:tabular-nums}
@media(max-width:779px){
 #view-detail .network-quality-panel{padding:12px!important}
 #view-detail .network-quality-kpis{grid-template-columns:repeat(2,minmax(0,1fr))!important;gap:6px!important;margin:8px 0!important}
 #view-detail .network-quality-kpis>div{padding:9px 7px!important;gap:5px!important;min-width:0}
 #view-detail .network-quality-kpis small{font-size:10px!important}
 #view-detail .network-quality-kpis strong{font-size:15px!important}
 #view-detail #network-time{font-size:12px!important}
 #view-overview .dashboard-toolbar{margin-bottom:8px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .router-meta{margin-top:4px!important}
}
@media(max-width:350px){
 #view-detail .network-quality-kpis>div{padding:8px 4px!important}
 #view-detail .network-quality-kpis strong{font-size:13px!important}
}
`
