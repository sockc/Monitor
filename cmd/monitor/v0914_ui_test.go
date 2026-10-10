package main

import (
 "path/filepath"
 "strings"
 "testing"
)

func TestNodeFavoritesSurviveDatabaseReopen(t *testing.T) {
 path:=filepath.Join(t.TempDir(),"monitor.db")
 db,e:=openDB(path);if e!=nil{t.Fatal(e)}
 if _,e=db.Exec("INSERT INTO node_favorites(name,favorite) VALUES(?,?)","n-7",1);e!=nil{t.Fatal(e)}
 if e=db.Close();e!=nil{t.Fatal(e)}
 db,e=openDB(path);if e!=nil{t.Fatal(e)};defer db.Close()
 var favorite int
 if e=db.QueryRow("SELECT favorite FROM node_favorites WHERE name=?","n-7").Scan(&favorite);e!=nil{t.Fatal(e)}
 if favorite!=1{t.Fatalf("favorite value lost: %d",favorite)}
}

func TestCustomDashboardControls(t *testing.T) {
 for _,id:=range []string{
  `id="overview-alert-filter"`,`data-node-filter="all"`,`data-node-filter="favorite"`,
  `data-node-filter="alert"`,`id="node-order-status"`,`id="card-fields-layout"`,
  `id="card-fields-options"`,`id="card-fields-reset"`,`id="card-fields-preview"`,
 } {
  if !strings.Contains(dashboardHTML,id){t.Errorf("missing dashboard control %s",id)}
 }
 for _,needle:=range []string{
  "function cardFieldOn(", "function saveCardFields(", "function saveNodeFavorite(",
  "function nodeHasAlert(", "function persistNodeOrder(", "data-group-header",
  "data-node-action=\"up\"", "data-node-action=\"down\"",
  "dragstart", "monitor-card-fields", "monitor-collapsed-groups",
  "monitor-quick-filter", "currentLayout==='compact'", "currentLayout==='list'",
 } {
  if !strings.Contains(appJS,needle){t.Errorf("missing JS behavior %q",needle)}
 }
 for _,needle:=range []string{
  "node-favorite-trigger", "node-group-heading", "card-fields-options",
  "--metric-count", "order-move-controls",
 }{
  if !strings.Contains(styleCSS,needle){t.Errorf("missing style %q",needle)}
 }
}
