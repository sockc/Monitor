package main

import (
 "strings"
 "testing"
)

func TestStatsPickerAndAlertFilters(t *testing.T) {
 for _, id := range []string{
  "history-node-pick", "history-node-dialog", "history-node-query",
  "history-node-options", "history-node-close", "history-point",
  "detail-point", "alert-summary", "node-info-backdrop",
 } {
  if !strings.Contains(dashboardHTML, `id="`+id+`"`) { t.Errorf("UI control missing: %s",id) }
 }
 for _, mode := range []string{"active","recovered","all"} {
  if !strings.Contains(dashboardHTML,`data-alert-mode="`+mode+`"`) {t.Errorf("missing alert mode %s",mode)}
 }
 for _, fn := range []string{"renderHistoryPicker()", "bindHistoryPoint(", "addEventListener('pointermove'","addEventListener('keydown'","function renderAlerts(items)", "historyChartPoints=data", "detailChartPoints=points"} {
  if !strings.Contains(appJS,fn){t.Errorf("missing interaction %s",fn)}
 }
 for _, css := range []string{".node-picker::backdrop",".chart-readout",".event-filters",".node-info-backdrop"} {
  if !strings.Contains(uiInteractionCSS,css){t.Errorf("missing component style %s",css)}
 }
}
