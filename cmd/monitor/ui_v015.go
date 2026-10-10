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
`;
