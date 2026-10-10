package main

import (
 "regexp"
 "strings"
 "testing"
)

// Guard against regressions which previously made the three responsive
// dashboard layouts indistinguishable and stacked filters at mobile widths.
func TestMobileDashboardControls(t *testing.T) {
 elements:=[]string{"node-search","node-group","metadata-sort","add-node","total","online","offline","alert-count"}
 ids:=map[string]int{}
 for _,m:=range regexp.MustCompile(`\bid="([^"]+)"`).FindAllStringSubmatch(dashboardHTML,-1){ids[m[1]]++}
 for _,id:=range elements{if ids[id]!=1{t.Errorf("element %q should appear exactly once, got %d",id,ids[id])}}
 if ids["node-order"]!=0{t.Error("homepage sort selector should be removed")}
 if !strings.Contains(styleCSS+mobileRefinementCSS,".filter-selects{grid-column:2!important;grid-row:1!important}"){t.Error("search and group should share first filter row")}
 for _,markup:=range []string{`class="dashboard-toolbar-head"`,`class="filter-search"`,`class="filter-selects"`,`data-layout="cards"`,`data-layout="compact"`,`data-layout="list"`}{
  if !strings.Contains(dashboardHTML,markup){t.Errorf("dashboard markup missing %q",markup)}
 }
}

func TestCardLayoutsHaveDistinctContent(t *testing.T) {
 for _,js:=range []string{
  "if(currentLayout==='list')", "if(currentLayout==='compact')",
  "node-list-summary", "compact-month", "card-cumulative",
  "card-quota", "node-info-trigger", "percentLabel(ratio)",
  "localStorage.setItem('monitor-node-layout'",
 }{
  if !strings.Contains(appJS,js){t.Errorf("layout logic missing %q",js)}
 }
 for _,css:=range []string{
  "dashboard-toolbar-head", "filter-search", "filter-selects",
  "grid-template-columns:repeat(4,minmax(0,1fr))",
  "server-grid[data-layout=compact]", "server-grid[data-layout=list]",
  "node-list-summary", "compact-month", "@media(max-width:350px)",
 }{
  if !strings.Contains(styleCSS+mobileRefinementCSS,css){t.Errorf("responsive style missing %q",css)}
 }
}
