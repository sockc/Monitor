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
  "nodes", "add-node", "detail-back", "detail-title", "detail-subtitle",
  "detail-status", "detail-metrics", "detail-hours", "detail-chart",
  "detail-chart-status", "detail-traffic", "detail-info",
  "detail-network-info", "detail-performance", "detail-network", "detail-system",
  "settings-nodes", "settings-limits", "settings-alerts", "settings-account",
  "show-add-node", "hide-add-node", "add-panel", "create-node",
  "metadata-form", "metadata-name", "metadata-display", "metadata-group",
  "metadata-notes", "metadata-status", "manage-name", "rename-target",
  "rename-node", "revoke-node", "delete-node", "manage-status",
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
 for _, v := range []string{"nodes","limits","alerts","account"} {
  if !strings.Contains(dashboardHTML, `data-settings-tab="`+v+`"`) {
   t.Errorf("missing settings tab %q", v)
  }
 }
 for _, v := range []string{"performance","network","system"} {
  if !strings.Contains(dashboardHTML, `data-detail-tab="`+v+`"`) {
   t.Errorf("missing detail tab %q", v)
  }
 }
}
