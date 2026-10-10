package main

// mobileRefinementCSS owns the most recent mobile dashboard and admin overrides.
// Keeping it separate prevents further cascading patches to the legacy styles.
const mobileRefinementCSS=`/* V0.9.10: search and group share one mobile row, distinct node views. */
#view-overview .dashboard-toolbar{display:flex;flex-direction:column;align-items:stretch;gap:9px}
#view-overview .dashboard-toolbar-head{display:flex;align-items:center;justify-content:space-between;width:100%;gap:8px}
#view-overview .dashboard-filters{display:grid!important;grid-template-columns:minmax(0,1.5fr) minmax(110px,.85fr)!important;align-items:center!important;gap:8px!important;justify-content:stretch!important;width:100%!important;flex:none!important}
#view-overview .dashboard-filters .filter-search{grid-column:1!important;grid-row:1!important;min-width:0!important;max-width:none!important;width:100%!important;flex:none!important}
#view-overview .dashboard-filters .filter-selects{grid-column:2!important;grid-row:1!important;display:block!important;min-width:0!important;max-width:none!important;width:100%!important;flex:none!important}
#view-overview .dashboard-filters #node-search,#view-overview .dashboard-filters #node-group{box-sizing:border-box;width:100%!important;min-width:0!important;max-width:none!important;height:40px!important;min-height:40px!important;font-size:12px!important;padding:8px 9px!important;margin:0!important}
#view-overview .dashboard-filters .layout-switch{grid-column:1/-1!important;grid-row:2!important;width:100%!important;max-width:none!important;display:flex!important;margin:0!important;box-sizing:border-box}
#view-overview .dashboard-filters .layout-button{flex:1!important}
#view-overview .server-grid[data-layout=compact] .glass-node .card-identity{margin-bottom:9px!important}
#view-overview .server-grid[data-layout=compact] .glass-node .metric-five{margin:0 0 8px!important}
#view-overview .server-grid[data-layout=compact] .glass-node .compact-month{font-size:11px!important}
#view-overview .server-grid[data-layout=compact] .glass-node .compact-quota{margin:3px 0 0!important}
#view-overview .server-grid[data-layout=list] .dashboard-node.glass-node{min-height:0!important}
#view-overview .server-grid[data-layout=list] .glass-node .node-list-summary{font-size:11px!important}
#view-settings .node-editor-list .row-sort{max-width:100%!important;font-variant-numeric:tabular-nums}
@media(max-width:779px){
 #view-overview .dashboard-toolbar{gap:8px!important;margin:9px 0 10px!important}
 #view-overview .dashboard-filters{grid-template-columns:minmax(0,1.5fr) minmax(105px,.85fr)!important;gap:7px!important}
 #view-overview .dashboard-filters .filter-search{grid-column:1!important;grid-row:1!important}
 #view-overview .dashboard-filters .filter-selects{grid-column:2!important;grid-row:1!important}
 #view-overview .dashboard-filters .layout-switch{grid-column:1/-1!important;grid-row:2!important}
 #view-overview .summary-strip{grid-template-columns:repeat(4,minmax(0,1fr))!important;gap:6px!important;margin:0 0 10px!important}
 #view-overview .summary-strip>div{min-height:53px!important;padding:7px 3px!important}
 #view-overview .summary-strip .summary-icon{display:none!important}
 #view-overview .summary-strip small{font-size:10px!important;white-space:nowrap!important}
 #view-overview .summary-strip strong{font-size:19px!important;line-height:1.2!important;margin:2px 0 0!important}
 #view-overview .server-grid[data-layout=compact] .dashboard-node.glass-node{padding:11px 12px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .card-identity{min-height:35px!important;margin-bottom:8px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .five-metric{gap:3px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .five-metric strong{font-size:11px!important}
 #view-overview .server-grid[data-layout=compact] .glass-node .metric-five{margin:0 0 6px!important}
 #view-overview .server-grid[data-layout=list] .dashboard-node.glass-node{padding:9px 12px!important}
 #view-overview .server-grid[data-layout=list] .glass-node .identity-text small{display:none}
}
@media(max-width:335px){
 #view-overview .summary-strip{grid-template-columns:repeat(2,minmax(0,1fr))!important}
 #view-overview .dashboard-filters{grid-template-columns:minmax(0,1.1fr) minmax(96px,1fr)!important}
 #view-overview .dashboard-filters input,#view-overview .dashboard-filters select{font-size:11px!important}
}

@media(min-width:1381px){#view-settings .node-editor-head,#view-settings .node-edit-row{grid-template-columns:minmax(115px,1fr) minmax(125px,1.1fr) minmax(85px,.65fr) minmax(120px,.9fr) minmax(62px,.55fr) minmax(245px,1.8fr)!important}}

/* V0.9.12: one responsive mobile layout entry point; avoid layering further
   @media overrides into the legacy global stylesheet. */
#view-settings .node-edit-toggle{display:none}
#view-settings .node-maintenance{min-width:0;flex:1 1 auto;font-size:12px}
#view-settings .node-maintenance summary{cursor:pointer;border:1px solid #52617d;border-radius:8px;padding:8px 12px;color:#bcd0e5}
#view-settings .node-maintenance .maintenance-actions{display:flex;flex-wrap:wrap;gap:7px;margin-top:9px}
@media(max-width:779px){
 body>header{display:grid!important;grid-template-columns:minmax(0,1fr) auto;gap:8px 10px!important;padding:10px 12px!important;align-items:center!important}
 body>header h2{grid-column:1!important;grid-row:1!important;font-size:17px!important;margin:0!important;min-width:0}
 body>header .header-tools{grid-column:2!important;grid-row:1!important;display:flex!important;align-items:center!important;justify-content:flex-end!important;gap:7px!important;min-width:0}
 body>header .theme-switch{margin:0!important}
 body>header #theme-select{width:90px!important;max-width:90px!important;height:34px!important;font-size:11px!important;padding:5px!important}
 body>header .logout-btn{height:34px!important;min-height:34px!important;padding:5px 9px!important;font-size:12px!important}
 body>header .top-nav{grid-column:1/-1!important;grid-row:2!important;order:0!important;width:100%!important;margin:0!important;padding:3px!important;min-width:0}
 body>header .top-nav button{flex:1 1 0!important;padding:8px 4px!important;font-size:12px!important;white-space:nowrap}
 #view-settings .node-edit-row{grid-template-columns:minmax(0,1fr) auto!important;gap:9px!important;align-items:center!important}
 #view-settings .node-edit-row .edit-node-identity{grid-column:1!important;min-width:0}
 #view-settings .node-edit-row .node-edit-toggle{display:inline-flex!important;grid-column:2!important;grid-row:1!important;padding:8px 12px!important}
 #view-settings .node-edit-row:not(.editing)>label,
 #view-settings .node-edit-row:not(.editing)>.node-edit-actions,
 #view-settings .node-edit-row:not(.editing)>.row-feedback{display:none!important}
 #view-settings .node-edit-row.editing>label{grid-column:1/-1!important;min-width:0}
 #view-settings .node-edit-row.editing .node-edit-actions,
 #view-settings .node-edit-row.editing .row-feedback{grid-column:1/-1!important}
 #view-settings .node-edit-row .node-edit-actions{display:flex!important;flex-wrap:wrap;gap:8px;align-items:center}
 #view-overview .glass-node{scroll-margin-top:95px}
}
@media(max-width:365px){
 body>header h2{font-size:15px!important}
 body>header #theme-select{width:76px!important}
 #view-settings .node-edit-row .node-edit-toggle{padding:7px 8px!important;font-size:11px!important}
}
`
