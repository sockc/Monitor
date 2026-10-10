package main

import (
 "regexp"
 "strings"
 "testing"
)

// Dashboard markup and JS are embedded in the Server binary. Guard the
// interactive controls against accidental removal during UI-only releases.
func TestDashboardUIControls(t *testing.T) {
 ids := []string{
  "nodes", "add-node", "overview-hint", "node-visible-count", "node-search", "node-group", "node-order", "detail-back", "detail-title", "detail-subtitle",
  "detail-status", "detail-metrics", "detail-hours", "detail-chart",
  "detail-chart-status", "detail-traffic", "detail-info",
  "detail-network-info", "detail-performance", "detail-network", "detail-system",
  "settings-nodes", "settings-limits", "settings-alerts", "settings-account",
  "show-add-node", "hide-add-node", "add-panel", "theme-select", "settings-appearance", "create-node",
  "metadata-form", "metadata-name", "metadata-display", "metadata-group", "metadata-location", "metadata-provider", "metadata-country", "metadata-price", "metadata-currency", "metadata-cycle", "metadata-port", "metadata-ipv4", "metadata-ipv6", "metadata-start", "metadata-expires", "metadata-quota",
  "metadata-notes", "metadata-auto-location", "metadata-status", "node-editor-list", "metadata-advanced", "metadata-advanced-close",
  "node-rebind-panel", "node-rebind-command", "node-rebind-copy", "node-rebind-close", "node-rebind-status",
  "limit-name", "limit-timezone", "limit-quota", "limit-expires",
  "node-limits-form", "limit-status", "alerts-form", "alerts-status",
  "alert-events", "refresh-alerts", "history-chart", "history-name",
  "history-hours", "history-metric", "traffic-bars", "traffic-summary",
 }
 seen := make(map[string]bool)
 for _, match := range regexp.MustCompile(`id="([^"]+)"`).FindAllStringSubmatch(dashboardHTML, -1) {
  if seen[match[1]] { t.Errorf("duplicate HTML id: %s", match[1]) }
  seen[match[1]] = true
 }
 for _, id := range ids {
  if !seen[id] { t.Errorf("missing HTML control id %q", id) }
 }
 for _, v := range []string{"nodes","limits","alerts","appearance","account"} {
  if !strings.Contains(dashboardHTML, `data-settings-tab="`+v+`"`) {
   t.Errorf("missing settings tab %q", v)
  }
 }
 for _, v := range []string{"performance","network","system"} {
  if !strings.Contains(dashboardHTML, `data-detail-tab="`+v+`"`) {
   t.Errorf("missing detail tab %q", v)
  }
 }
 if !strings.Contains(styleCSS, "repeat(4,minmax(0,1fr))") { t.Error("desktop dashboard must support four equal columns") }
 if !strings.Contains(appJS, "n.location") { t.Error("dashboard is missing node location") }
 for _,theme:=range []string{"glass","dark","light"}{if !strings.Contains(dashboardHTML, `data-theme-choice="`+theme+`"`){t.Errorf("missing %s theme choice",theme)}}
 if !strings.Contains(appJS, "monitor-theme"){t.Error("theme persistence missing")}
 if !strings.Contains(dashboardHTML, "metadata-auto-location"){t.Error("automatic geo help text missing")}
 for _,piece:=range []string{"metric-five","card-cumulative","capability-tags","quota-linear","countryFlag","IPv4","IPv6","planPort"}{if !strings.Contains(appJS,piece){t.Errorf("missing VPS card component: %s",piece)}}
 for _,phrase:=range []string{"价格未设置","到期未设置","端口未设","IPv4 ?","IPv6 ?"}{
  if strings.Contains(appJS,phrase){t.Errorf("unconfigured VPS field must be hidden, found %q",phrase)}
 }
 if !strings.Contains(appJS, "/static/flags/"){t.Error("expected local SVG flag images")}
 if !strings.Contains(appJS,"n.public_ipv4")||!strings.Contains(appJS,"n.public_ipv6"){t.Error("expected both outbound IP families in dashboard")}
 if !strings.Contains(appJS, "card-cumulative"){t.Error("single-row total traffic missing")}

 if strings.Contains(dashboardHTML, `id="node-name"`) || strings.Contains(dashboardHTML, `id="rename-target"`){t.Error("manual node ID editing must not appear in settings")}
 if strings.Contains(dashboardHTML, `id="rename-node"`){t.Error("identity-renaming button must not appear")}
 for _,term:=range []string{"renderNodeEditorRows","data-node-action=","/api/v1/node-rebind","generatedInstallCommand"}{
  if !strings.Contains(appJS,term){t.Errorf("node inventory missing %s",term)}
 }

}
