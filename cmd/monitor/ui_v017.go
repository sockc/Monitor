package main

// The V0.9.17 refinements apply after the shared desktop and mobile styles.
// OpenWrt telemetry adds one optional line only to router cards.
const uiV017CSS = `
#view-overview .glass-node .router-meta{margin-top:5px;color:var(--node-muted,#93a8c2);font-size:10px;line-height:1.4;overflow:hidden;text-overflow:ellipsis;white-space:nowrap}
@media(max-width:779px){
 #view-overview .overview-headline{margin:0 0 7px!important}
 #view-overview .dashboard-heading{margin-bottom:7px!important}
 #view-overview .summary-strip{margin-bottom:7px!important;gap:6px!important}
 #view-overview .summary-strip>div{min-height:48px!important;padding:6px 2px!important}
 #view-overview .summary-strip small{font-size:10px!important}
 #view-overview .summary-strip strong{font-size:18px!important}
 #view-overview .dashboard-toolbar{margin:4px 0 8px!important}
 #view-overview .node-group-heading{min-height:34px!important;padding:7px 10px!important}
}
`
