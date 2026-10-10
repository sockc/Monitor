package main

// Monitor stylesheet segment; concatenated in original cascade order.
const styleCSSResponsive=`/* V0.9.10: mobile overview hierarchy and genuinely different card modes */
#view-overview .dashboard-toolbar{display:flex!important;align-items:center!important;justify-content:space-between!important;gap:14px!important;flex-wrap:wrap!important;margin:12px 0!important}
#view-overview .dashboard-toolbar-head{display:flex;align-items:center;justify-content:space-between;gap:12px;flex:0 0 auto;min-width:0}
#view-overview .dashboard-toolbar-head .toolbar-label{display:flex;gap:8px;align-items:baseline;white-space:nowrap}
#view-overview .dashboard-toolbar-head #add-node{height:35px;padding:7px 12px!important;font-size:12px;white-space:nowrap}
#view-overview .dashboard-filters{display:flex;align-items:center;justify-content:flex-end;gap:8px;flex:1 1 auto;min-width:0}
#view-overview .filter-search{flex:1 1 170px;min-width:135px;max-width:240px}
#view-overview .filter-search #node-search{width:100%!important;max-width:none!important;min-width:0!important;margin:0!important}
#view-overview .filter-selects{display:flex;align-items:center;gap:8px;min-width:0}
#view-overview .filter-selects select{min-width:112px;max-width:160px;width:auto;flex:0 1 auto}
#view-overview .dashboard-filters .layout-switch{display:flex;width:auto!important;min-width:0;margin:0!important}
#view-overview .server-grid[data-layout=compact] .glass-node,
#view-overview .server-grid[data-layout=list] .glass-node{min-height:0!important}
#view-overview .server-grid[data-layout=compact] .glass-node .card-identity,
#view-overview .server-grid[data-layout=list] .glass-node .card-identity{
 justify-content:flex-start!important;text-align:left!important;gap:7px!important;
 padding:0 34px 0 0!important;min-height:32px!important;margin:0 0 10px!important
}
#view-overview .server-grid[data-layout=compact] .glass-node .identity-text,
#view-overview .server-grid[data-layout=list] .glass-node .identity-text{flex:1 1 auto;min-width:0;max-width:none!important;align-items:flex-start!important;text-align:left!important}
#view-overview .server-grid[data-layout=compact] .glass-node .identity-text strong,
#view-overview .server-grid[data-layout=list] .glass-node .identity-text strong{font-size:13px!important;line-height:1.35}
#view-overview .server-grid[data-layout=compact] .glass-node .state-text,
#view-overview .server-grid[data-layout=list] .glass-node .state-text{flex:0 0 auto;margin-left:auto!important;font-size:11px!important;white-space:nowrap}
#view-overview .server-grid[data-layout=compact] .dashboard-node{padding:13px 14px 12px!important}
#view-overview .server-grid[data-layout=compact] .glass-node .metric-five{margin:0 0 10px!important;gap:4px!important;grid-template-columns:repeat(5,minmax(0,1fr))!important}
#view-overview .server-grid[data-layout=compact] .glass-node .five-metric{gap:3px!important;min-width:0}
#view-overview .server-grid[data-layout=compact] .glass-node .five-metric small{font-size:10px!important}
#view-overview .server-grid[data-layout=compact] .glass-node .five-metric strong{font-size:12px!important;font-variant-numeric:tabular-nums}
#view-overview .server-grid[data-layout=compact] .glass-node .five-track{height:3px!important}
#view-overview .glass-node .compact-month{display:flex;justify-content:space-between;align-items:baseline;gap:10px;
 border-top:1px solid var(--node-stroke,#30445f);padding:9px 0 0;margin:0;
 color:var(--node-muted,#9eb1c8);font-size:11px;line-height:1.45}
#view-overview .glass-node .compact-month b{font-weight:750;color:var(--node-ink,#eaf2ff)}
#view-overview .glass-node .compact-month span:last-child{flex:0 0 auto}
#view-overview .glass-node .compact-quota{height:4px;border-radius:6px;background:var(--node-track,#294157);
 overflow:hidden;margin:6px 0 0}
#view-overview .glass-node .compact-quota i{display:block;height:100%;border-radius:6px;background:#36c997}
#view-overview .glass-node .compact-quota i.warn{background:#f5aa50}
#view-overview .glass-node .compact-quota i.danger{background:#f27186}
#view-overview .server-grid[data-layout=compact] .glass-node .card-alerts{margin-top:9px!important}
#view-overview .server-grid[data-layout=list] .dashboard-node{display:flex!important;flex-direction:column!important;align-items:stretch!important;
 justify-content:center!important;gap:0!important;padding:12px 14px!important;
 border-radius:12px!important;grid-template-columns:none!important}
#view-overview .server-grid[data-layout=list] .glass-node .card-identity{margin:0 0 4px!important}
#view-overview .glass-node .node-list-summary{font-size:11px!important;color:var(--node-muted,#9fb0c9);
 text-align:left;white-space:nowrap;overflow:hidden;text-overflow:ellipsis;padding-left:24px;font-variant-numeric:tabular-nums}
#view-overview .server-grid[data-layout=list] .glass-node .card-alerts{margin:7px 0 0!important}
#view-overview .server-grid[data-layout=list] .glass-node .node-info-trigger,
#view-overview .server-grid[data-layout=compact] .glass-node .node-info-trigger{top:8px;right:9px}
@media(max-width:779px){
 #view-overview .dashboard-heading{display:flex!important;align-items:flex-start!important;margin:1px 0 11px!important;gap:6px!important}
 #view-overview .dashboard-heading h1{font-size:19px!important;margin:0!important;line-height:1.35!important}
 #view-overview .dashboard-heading .eyebrow,#view-overview #overview-hint{display:none!important}
 #view-overview .dashboard-updated{font-size:10px!important;white-space:nowrap;margin-top:5px!important}
 #view-overview .summary-strip{display:grid!important;grid-template-columns:repeat(4,minmax(0,1fr))!important;gap:6px!important;margin-bottom:13px!important}
 #view-overview .summary-strip>div{min-width:0;min-height:63px!important;display:flex!important;
  flex-direction:column!important;align-items:center!important;justify-content:center!important;
  padding:8px 3px!important;gap:0!important;text-align:center!important;border-radius:10px!important}
 #view-overview .summary-strip>div>div{min-width:0}
 #view-overview .summary-strip .summary-icon{display:none!important}
 #view-overview .summary-strip small{display:block;font-size:10px!important;line-height:1.35;white-space:nowrap}
 #view-overview .summary-strip strong{display:block;font-size:21px!important;line-height:1.2!important;margin:3px 0 0!important;font-variant-numeric:tabular-nums}
 #view-overview .dashboard-toolbar{display:flex!important;flex-direction:column!important;align-items:stretch!important;gap:10px!important;margin:10px 0 11px!important}
 #view-overview .dashboard-toolbar-head{display:flex;width:100%;justify-content:space-between;align-items:center;gap:8px}
 #view-overview .dashboard-toolbar-head .toolbar-label h3{font-size:16px!important;margin:0}
 #view-overview .dashboard-toolbar-head #add-node{height:36px!important;padding:7px 11px!important;font-size:12px!important}
 #view-overview .dashboard-filters{display:grid!important;grid-template-columns:repeat(2,minmax(0,1fr))!important;gap:8px!important;
  width:100%!important;max-width:none!important;flex:none!important;align-items:stretch!important;justify-content:stretch!important}
 #view-overview .filter-search,#view-overview .filter-selects,#view-overview .dashboard-filters .layout-switch{
  grid-column:1/-1!important;min-width:0!important;max-width:none!important;width:100%!important;margin:0!important;flex:none!important}
 #view-overview .filter-search{display:block}
 #view-overview .filter-selects{display:grid!important;grid-template-columns:repeat(2,minmax(0,1fr));gap:8px!important}
 #view-overview .dashboard-filters input,#view-overview .dashboard-filters select{
  width:100%!important;max-width:none!important;min-width:0!important;box-sizing:border-box!important;
  height:40px!important;min-height:40px!important;margin:0!important;padding:8px 11px!important;font-size:12px!important}
 #view-overview .dashboard-filters .layout-switch{display:flex!important;justify-content:stretch!important;align-items:stretch!important;padding:3px!important;height:41px!important}
 #view-overview .dashboard-filters .layout-button{flex:1 1 0!important;min-width:0!important;padding:7px 2px!important;font-size:12px!important;height:33px!important}
 #view-overview .server-grid[data-layout=compact],
 #view-overview .server-grid[data-layout=cards],
 #view-overview .server-grid[data-layout=list]{display:grid!important;grid-template-columns:minmax(0,1fr)!important;gap:9px!important}
 #view-overview .server-grid[data-layout=compact] .dashboard-node{padding:12px 13px!important}
 #view-overview .server-grid[data-layout=list] .dashboard-node{padding:9px 13px!important}
 #view-overview .server-grid[data-layout=cards] .dashboard-node{padding:15px 13px!important}
}
@media(max-width:350px){
 #view-overview .summary-strip{grid-template-columns:repeat(2,minmax(0,1fr))!important}
 #view-overview .summary-strip>div{min-height:50px!important}
 #view-overview .dashboard-updated{font-size:9px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .five-metric strong{font-size:11px!important}
}



`
