package main

// Monitor stylesheet segment; concatenated in original cascade order.
const styleCSSPersonal=`/* V0.9.14 - unified customizable card field and node-management controls */
#view-overview .overview-quick-filters{display:flex;flex-wrap:wrap;gap:7px;margin:-1px 0 7px}
#view-overview .overview-quick-filters button{border:1px solid var(--t-border,#38536c);background:var(--t-card,#14253c);color:var(--t-text,#c9daf1);border-radius:9px;padding:7px 15px;font-size:12px;min-height:34px}
#view-overview .overview-quick-filters button.selected{background:#28547c!important;border-color:#4e90c4!important;color:#fff!important}
#view-overview .summary-strip>div:last-child{position:relative}
#view-overview .alert-filter-shortcut{position:absolute;right:6px;top:5px;border:1px solid #dfac5855;border-radius:6px;padding:2px 5px;background:transparent;color:#d3a24f;font-size:10px;cursor:pointer}
#view-overview .node-group-heading{grid-column:1/-1!important;display:flex!important;align-items:center!important;justify-content:space-between;width:100%;gap:10px;padding:9px 12px!important;min-height:37px;background:var(--t-card,#13253a);border:1px solid var(--t-border,#344a65);border-radius:9px;color:var(--t-text,#dfeaf8);font-size:12px;font-weight:650;text-align:left;box-shadow:none!important}
#view-overview .node-group-heading small{font-size:11px;color:var(--node-muted,#9eb0c9);font-weight:400}
#view-overview .node-group-heading:focus-visible{outline:2px solid #60a5fa;outline-offset:2px}
#view-overview .server-grid .dashboard-node.glass-node{position:relative}
#view-overview .glass-node .node-favorite-trigger{position:absolute;z-index:4;right:38px;top:9px;width:26px;height:26px;min-width:26px;padding:0;border:1px solid #67809d66;background:var(--node-inlay,#23344b);color:var(--node-muted,#9daec7);border-radius:50%;display:grid;place-items:center;font-size:18px;line-height:1;cursor:pointer}
#view-overview .glass-node .node-favorite-trigger.is-favorite{color:#e7ad3f!important;background:#fff0c329;border-color:#e7ad3f77}
#view-overview .glass-node .card-identity{padding-right:69px!important}
#view-overview .glass-node .metric-five:empty{display:none}
#view-settings .card-fields-head{display:flex;align-items:center;justify-content:space-between;gap:10px;flex-wrap:wrap;margin-bottom:12px}
#view-settings .card-fields-head label{display:flex;align-items:center;gap:9px;color:var(--t-text,#d3e2f4);font-size:12px}
#view-settings #card-fields-layout{min-height:38px;padding:7px 10px;border:1px solid var(--t-border,#344f6b);border-radius:8px;background:var(--t-card,#152741);color:var(--t-text,#eaf3ff)}
#view-settings .card-fields-options{display:grid;grid-template-columns:repeat(3,minmax(0,1fr));gap:9px}
#view-settings .card-fields-options label{display:flex;align-items:center;gap:9px;font-size:12px;padding:10px 11px;border:1px solid var(--t-border,#354a65);background:var(--t-card,#152640);border-radius:8px;cursor:pointer;min-width:0}
#view-settings .card-fields-options input{accent-color:#4f99d8;width:16px;height:16px;flex:none}
#view-settings .card-fields-preview{margin-top:12px;padding:11px 13px;border-radius:9px;background:var(--t-card,#14263d);border:1px dashed var(--t-border,#36516e);font-size:11px;line-height:1.6;color:var(--t-text,#adc3df);overflow-wrap:anywhere}
#view-settings .order-move-controls{display:flex!important;align-items:center;gap:2px;flex:none}
#view-settings .order-move-controls button{display:grid;place-items:center;width:23px;min-width:23px;height:27px;padding:0;border:1px solid var(--t-border,#38516b);border-radius:6px;background:var(--t-card,#1e334b);color:var(--t-text,#d5e5f9);font-size:13px;cursor:pointer}
#view-settings .order-move-controls .node-drag-handle{cursor:grab;font-size:17px;touch-action:none}
#view-settings .order-move-controls .node-drag-handle:active{cursor:grabbing}
#view-settings .node-edit-row .edit-node-identity{gap:5px!important}
#view-settings .node-edit-row .edit-node-label{min-width:0;overflow:hidden}
#view-settings .node-edit-row .edit-node-label strong{white-space:nowrap;text-overflow:ellipsis;overflow:hidden;display:block}
#view-settings #node-order-status{margin:0 0 8px;font-size:11px}
@media(max-width:800px){#view-settings .card-fields-options{grid-template-columns:repeat(2,minmax(0,1fr))}#view-settings .order-move-controls button{width:29px;min-width:29px;height:32px}#view-settings .order-move-controls .node-drag-handle{display:none!important}}
@media(max-width:420px){#view-overview .overview-quick-filters button{flex:1;padding:7px 8px}#view-settings .card-fields-options{grid-template-columns:1fr 1fr}#view-settings .card-fields-options label{font-size:11px;padding:8px 7px}#view-overview .glass-node .node-favorite-trigger{right:40px;top:8px;width:28px;height:28px;font-size:19px}#view-overview .node-group-heading{font-size:12px}#view-overview .alert-filter-shortcut{font-size:9px}}

#view-overview .server-grid[data-layout=compact] .glass-node .metric-five,#view-overview .server-grid[data-layout=cards] .glass-node .metric-five{grid-template-columns:repeat(var(--metric-count,5),minmax(0,1fr))!important}
`
