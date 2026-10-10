package main

// Extra UI surfaces split out from the legacy stylesheet.
const uiInteractionCSS=`
/* V0.9.13: statistical controls, events, and mobile details */
#history-node-pick{min-width:160px;max-width:100%;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;text-align:left}
#view-statistics .accessible-picker-source{position:absolute!important;width:1px!important;height:1px!important;padding:0!important;margin:-1px!important;overflow:hidden!important;clip-path:inset(50%)!important;white-space:nowrap!important;border:0!important}
.node-picker{width:min(510px,calc(100vw - 22px));max-height:min(74vh,720px);background:var(--ui-card,#15243a);color:var(--ui-text,#eef5ff);border:1px solid var(--ui-border,#385373);box-shadow:0 22px 66px #0008;border-radius:16px;padding:15px;box-sizing:border-box;overscroll-behavior:contain}
.node-picker::backdrop{background:#0713219c;backdrop-filter:blur(2px)}
.node-picker .node-picker-head{display:flex;justify-content:space-between;align-items:center;margin-bottom:11px;gap:10px}
.node-picker .node-picker-head strong{font-size:16px}
.node-picker #history-node-query{box-sizing:border-box;width:100%;min-height:43px;border-radius:8px;padding:11px;font-size:14px;background:var(--ui-card-2,#192a43);color:inherit;border:1px solid var(--ui-border,#385373)}
.node-picker .node-picker-options{display:flex;flex-direction:column;gap:5px;margin-top:12px;max-height:min(52vh,530px);overflow:auto;overscroll-behavior:contain}
.node-picker .node-picker-item{display:flex;flex-direction:column;align-items:flex-start;text-align:left;gap:4px;min-height:53px;border:1px solid var(--ui-border,#385373);background:var(--ui-card-2,#192a43);border-radius:9px;padding:9px 12px;color:inherit;overflow-wrap:anywhere}
.node-picker .node-picker-item strong{font-size:14px}
.node-picker .node-picker-item small{font-size:11px;opacity:.7}
.node-picker .node-picker-item[aria-pressed=true]{border-color:#62a9ee;background:#274464}
#view-statistics canvas,#view-detail canvas{touch-action:pan-y;cursor:crosshair;display:block;max-width:100%;height:auto}
.chart-readout{margin:8px 0 4px;padding:8px 10px;border-radius:8px;min-height:17px;background:var(--ui-card-2,#192a43);border:1px solid var(--ui-border,#385373);font-size:12px;line-height:1.5;font-variant-numeric:tabular-nums;color:var(--ui-text,#e5f1ff)}
.event-filters{display:flex;align-items:center;gap:6px;flex-wrap:wrap;margin:7px 0 14px}
.event-filters button{padding:7px 13px!important;border:1px solid var(--ui-border,#385373);background:transparent;border-radius:8px;font-size:12px;color:inherit}
.event-filters button.selected{background:#315f8c;border-color:#4488bd;color:#f5faff}
.event-filters #alert-summary{font-size:11px;margin-left:auto}
.node-info-backdrop[hidden]{display:none!important}
.node-info-backdrop{position:fixed;inset:0;background:#07121f8f;z-index:9998;backdrop-filter:blur(1.5px)}

:is(html[data-theme=glass],html[data-theme=light]) .node-picker{background:#fafcff;color:#2e425e;border-color:#c5d3e6}
:is(html[data-theme=glass],html[data-theme=light]) .node-picker #history-node-query,
:is(html[data-theme=glass],html[data-theme=light]) .node-picker .node-picker-item{background:#f4f8ff;color:#2e425e;border-color:#c5d3e6}
:is(html[data-theme=glass],html[data-theme=light]) .node-picker .node-picker-item[aria-pressed=true]{background:#dbeafa;border-color:#5796cd}
:is(html[data-theme=glass],html[data-theme=light]) .chart-readout{background:#eef4fc;color:#35506d;border-color:#c6d6e9}
:is(html[data-theme=glass],html[data-theme=light]) .event-filters button{color:#3e5670;border-color:#b2c7db}
:is(html[data-theme=glass],html[data-theme=light]) .event-filters button.selected{color:#fff;background:#37699b}
@media(max-width:779px){
 #view-statistics .stat-filters{display:grid!important;grid-template-columns:minmax(0,1.2fr) minmax(0,.75fr) minmax(0,.75fr)!important;gap:6px!important;align-items:center}
 #view-statistics #history-node-pick{grid-column:1!important;min-width:0!important;min-height:39px!important;width:100%;max-width:none;padding:8px!important;font-size:12px;text-align:left}
 #view-statistics #history-hours,#view-statistics #history-metric{min-width:0!important;width:100%!important;max-width:none}
 #view-overview #node-quick-info{position:fixed!important;top:auto!important;left:0!important;right:0!important;bottom:0!important;width:100%!important;max-width:none!important;max-height:min(76vh,750px)!important;overflow:auto!important;border-radius:18px 18px 0 0!important;padding:18px 15px max(18px,env(safe-area-inset-bottom))!important;box-sizing:border-box!important;z-index:9999!important;box-shadow:0 -12px 55px #0006!important}
 .node-picker{width:calc(100vw - 18px);padding:12px}
 .node-picker .node-picker-item{min-height:48px}
 .event-filters #alert-summary{flex-basis:100%;margin-left:0}
}
`
