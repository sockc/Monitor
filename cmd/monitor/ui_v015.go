package main

// Focused V0.9.15 interface polish. Uses existing theme variables rather
// than adding another theme, and deliberately leaves collected data alone.
const uiV015CSS=`
/* Consistent control hierarchy and accessible feedback */
#view-settings .row-feedback[data-result=success]{color:#21ae7c!important;font-weight:650}
#view-settings .row-feedback[data-result=error]{color:#ee6679!important;font-weight:650}
#view-settings .node-edit-row[data-dirty=true]{outline:1px solid #d8a74d82;outline-offset:0}
#view-settings .node-edit-row[data-dirty=true] .node-edit-toggle::after{content:" · 未保存";color:#d7a44d;font-size:10px}
#view-settings .node-edit-row button:disabled{opacity:.54;cursor:wait}
#view-settings .node-edit-actions button{min-height:34px}
#view-settings .node-edit-row .node-edit-mobile-label{line-height:1.25}
#view-settings #node-order-status{font-size:12px;line-height:1.5}
#view-settings .settings-card-heading>div{min-width:0}
#view-settings .settings-pane-heading{min-width:0}
#view-settings .form-feedback{overflow-wrap:anywhere}
#view-settings .node-edit-row input{box-sizing:border-box;max-width:100%}
#view-settings .node-edit-row input:focus-visible{outline:2px solid #459af4;outline-offset:1px}
#view-overview .server-grid .dashboard-node.glass-node{min-width:0;box-sizing:border-box}
#view-overview .server-grid .dashboard-node.glass-node .card-identity{min-width:0;max-width:100%}
#view-overview .server-grid .dashboard-node.glass-node .identity-text{min-width:0;flex:1 1 auto}
#view-overview .server-grid .dashboard-node.glass-node .identity-text strong{max-width:100%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
#view-overview .server-grid .dashboard-node.glass-node .state-text{flex-shrink:0;white-space:nowrap}
#view-overview .server-grid .dashboard-node.glass-node .metric-five{min-width:0}
#view-overview .server-grid .dashboard-node.glass-node .five-metric{min-width:0}
#view-overview .server-grid .dashboard-node.glass-node .five-metric strong{font-variant-numeric:tabular-nums}
#view-overview .server-grid .dashboard-node.glass-node .card-alerts{flex-wrap:wrap}
#view-overview .server-grid .dashboard-node.glass-node .card-alerts span{line-height:1.35}
#view-overview .group-heading button:focus-visible,#view-overview .node-favorite:focus-visible{outline:2px solid #459af4;outline-offset:2px}
#view-statistics .history-controls{min-width:0}
#view-statistics #history-point,#view-detail #detail-point{min-height:20px;font-variant-numeric:tabular-nums}
#view-alerts .event-item{min-width:0}
#view-alerts .event-content{min-width:0;overflow-wrap:anywhere}
@media(max-width:779px){
 body>header{box-sizing:border-box;max-width:100%;gap:7px}
 body>header .top-nav{min-width:0}
 body>header .top-nav .nav-btn{min-width:0}
 #view-overview .dashboard-heading{min-width:0}
 #view-overview .dashboard-toolbar-head{min-width:0}
 #view-overview .dashboard-toolbar-head .toolbar-label{min-width:0;overflow:hidden}
 #view-overview .dashboard-toolbar-head #node-visible-count{white-space:nowrap}
 #view-overview .server-grid .dashboard-node.glass-node .card-identity{gap:6px}
 #view-overview .server-grid .dashboard-node.glass-node .identity-text strong{font-size:13px;line-height:1.35}
 #view-overview .server-grid .dashboard-node.glass-node .state-text{font-size:10px}
 #view-overview .server-grid .dashboard-node.glass-node .five-metric strong{font-size:11px}
 #view-overview .server-grid .dashboard-node.glass-node .card-alerts{margin-top:8px}
 #view-settings .node-edit-row{min-width:0}
 #view-settings .node-edit-row .edit-node-identity{min-width:0}
 #view-settings .node-edit-row .row-feedback{font-size:11px}
 #view-statistics .stat-filters{flex-wrap:wrap}
 #view-detail .detail-top{min-width:0}
}
@media(prefers-reduced-motion:reduce){
 #view-overview .dashboard-node,#view-settings .settings-tab,#view-overview .layout-button{transition:none!important;animation:none!important}
}

/* Desktop dashboard density: title and four status counters share one row. */
@media(min-width:1050px){
 body>header{min-height:58px!important;padding-top:9px!important;padding-bottom:9px!important}
 #view-overview .overview-headline{display:grid;grid-template-columns:minmax(210px,1fr) minmax(520px,auto);align-items:center;gap:18px;margin:0 0 15px}
 #view-overview .overview-headline .dashboard-heading{display:flex!important;align-items:center!important;gap:12px!important;min-width:0;margin:0!important}
 #view-overview .overview-headline .dashboard-heading>div:first-child{min-width:0}
 #view-overview .overview-headline .dashboard-heading h1{font-size:21px!important;margin:0!important}
 #view-overview .overview-headline #overview-hint{display:none!important}
 #view-overview .overview-headline .dashboard-updated{font-size:10px!important;white-space:nowrap}
 #view-overview .overview-headline .summary-strip{display:grid!important;grid-template-columns:repeat(4,minmax(0,1fr))!important;gap:7px!important;margin:0!important;min-width:0}
 #view-overview .overview-headline .summary-strip>div{min-height:44px!important;padding:6px 10px!important;gap:5px!important;min-width:0!important;border-radius:9px!important}
 #view-overview .overview-headline .summary-strip .summary-icon{display:none!important}
 #view-overview .overview-headline .summary-strip>div>div{display:flex;align-items:center;justify-content:center;gap:7px;min-width:0}
 #view-overview .overview-headline .summary-strip small{font-size:11px!important;white-space:nowrap!important}
 #view-overview .overview-headline .summary-strip strong{font-size:19px!important;margin:0!important;font-variant-numeric:tabular-nums}
 #view-overview .overview-headline #overview-alert-filter{position:static!important;display:inline-flex;align-items:center;justify-content:center;min-width:0;padding:2px 5px!important;margin:0!important;font-size:10px!important}
 #view-overview .dashboard-toolbar{display:grid!important;grid-template-columns:auto minmax(0,1fr);align-items:center!important;gap:9px 12px!important;margin:0 0 10px!important}
 #view-overview .dashboard-toolbar-head{grid-column:1;grid-row:1;display:flex;align-items:center;gap:12px}
 #view-overview .dashboard-toolbar-head #add-node{order:2}
 #view-overview .overview-quick-filters{grid-column:2;grid-row:1;justify-self:end;margin:0!important}
 #view-overview .dashboard-filters{grid-column:1/-1;grid-row:2;display:grid!important;grid-template-columns:minmax(180px,1fr) minmax(140px,230px) auto;gap:9px!important;justify-content:stretch!important;width:100%}
 #view-overview .dashboard-filters .filter-search{min-width:0!important;max-width:none!important;width:100%!important}
 #view-overview .dashboard-filters #node-search{max-width:none!important;width:100%!important}
 #view-overview .dashboard-filters .filter-selects{min-width:0!important;max-width:none!important;width:100%!important}
 #view-overview .dashboard-filters #node-group{max-width:none!important;width:100%!important}
 #view-overview .dashboard-filters .layout-switch{justify-self:end!important}
 #view-overview .group-heading{margin-top:8px!important;margin-bottom:8px!important}
}
@media(min-width:1050px) and (max-width:1350px){
 #view-overview .overview-headline{grid-template-columns:minmax(175px,1fr) minmax(480px,1.8fr)}
 #view-overview .overview-headline .dashboard-updated{display:none!important}
 #view-overview .overview-headline .summary-strip>div{padding:6px!important}
 #view-overview .overview-headline .summary-strip small{font-size:10px!important}
}
@media(max-width:1049px){
 #view-overview .overview-headline{display:block}
 #view-overview .overview-headline .dashboard-heading{margin-bottom:11px!important}
}

/* Compact node cards: prioritize metrics, keep actionable warnings visible. */
@media(min-width:780px){
 #view-overview .server-grid[data-layout=compact] .dashboard-node.glass-node{padding:12px 13px 11px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .card-identity{margin-bottom:8px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .metric-five{margin-bottom:7px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .compact-month{margin-top:2px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .card-alerts{margin-top:7px!important}
}
#view-overview .glass-node .card-alerts{display:flex;flex-wrap:wrap;gap:5px}
#view-overview .glass-node .card-alerts span{max-width:100%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}

/* Overview toolbar: all group and layout controls directly follow Add. */
#view-overview .dashboard-toolbar-head .dashboard-toolbar-actions{display:flex;align-items:center;gap:8px;flex-wrap:wrap;min-width:0}
#view-overview .dashboard-toolbar-head .filter-selects{display:block;min-width:110px}
#view-overview .dashboard-toolbar-head #node-group{height:35px!important;min-height:35px!important;max-width:180px;width:100%;padding:7px 9px!important;font-size:11px!important;border-radius:8px!important;background:#14243a!important;border:1px solid #304863!important;color:#d5e4f8}
#view-overview .dashboard-toolbar-head .layout-switch{display:flex!important;width:auto!important;margin:0!important;flex:0 0 auto}
#view-overview .dashboard-toolbar-head .layout-button{height:30px!important;padding:5px 9px!important;font-size:11px!important}
@media(min-width:1050px){
 #view-overview .dashboard-toolbar{display:grid!important;grid-template-columns:minmax(0,1fr) auto;gap:8px 12px!important}
 #view-overview .dashboard-toolbar-head{grid-column:1;grid-row:1;display:flex;flex-wrap:wrap;gap:10px!important;min-width:0}
 #view-overview .dashboard-toolbar-head #add-node{order:0}
 #view-overview .overview-quick-filters{grid-column:2;grid-row:1}
 #view-overview .dashboard-filters{grid-column:1/-1;grid-row:2;display:flex!important;justify-content:stretch!important;width:100%}
 #view-overview .dashboard-filters .filter-search{flex:1 1 100%;max-width:none!important;width:100%!important}
 #view-overview .dashboard-filters #node-search{max-width:none!important;width:100%!important}
}
@media(max-width:1049px){
 #view-overview .dashboard-toolbar-head{flex-wrap:wrap!important;height:auto!important}
 #view-overview .dashboard-toolbar-head .dashboard-toolbar-actions{width:100%;gap:7px}
 #view-overview .dashboard-toolbar-head .filter-selects{flex:1 1 110px;min-width:100px}
 #view-overview .dashboard-toolbar-head #node-group{max-width:none}
 #view-overview .dashboard-filters{display:flex!important;width:100%!important}
 #view-overview .dashboard-filters .filter-search{flex:1!important;max-width:none!important;width:100%!important}
 #view-overview .dashboard-filters #node-search{width:100%!important;max-width:none!important}
}
`;
