package main

import (
 "strings"
 "testing"
)

// The layered CSS must keep the old cascade and all three selectable themes.
// These checks also protect the UX fixes against accidental regressions.
func TestV015StylesheetLayers(t *testing.T) {
 parts:=[]struct{name,value string}{
  {"foundation",styleCSSFoundation},{"dashboard",styleCSSDashboard},
  {"themes",styleCSSThemes},{"cards",styleCSSCards},
  {"admin",styleCSSAdmin},{"responsive",styleCSSResponsive},
  {"personal",styleCSSPersonal},
 }
 var combined strings.Builder
 for _,p:=range parts{
  if len(p.value)<1000{t.Errorf("%s stylesheet is unexpectedly empty",p.name)}
  combined.WriteString(p.value)
 }
 if got:=combined.String();got!=styleCSS{t.Fatal("refactored CSS lost or reordered existing rules")}
 for _,selector:=range []string{"data-theme=glass","data-theme=dark","data-theme=light","server-grid[data-layout=compact]","server-grid[data-layout=cards]","server-grid[data-layout=list]"}{
  if !strings.Contains(styleCSS+mobileRefinementCSS+uiInteractionCSS+uiV015CSS,selector){t.Errorf("missing theme/layout styling %s",selector)}
 }
 if !strings.Contains(uiV015CSS,"@media(max-width:779px)") { t.Fatal("new mobile rules missing") }
}

func TestV015DashboardInteractivity(t *testing.T) {
 for _,piece:=range []string{
  "nodeFetchBusy","trafficFetchBusy","nodeFetchSequence","nodeMutationSequence",
  "historyFetchSequence","detailHistorySequence",
  "updateNodeEditorDirty","row.dataset.dirty='true'",
  "保存失败，修改内容已保留","feedback.dataset.result='success'",
  "saveOverviewFilters","monitor-overview-filters","overviewScrollY",
  "requestAnimationFrame(()=>window.scrollTo",
 }{
  if !strings.Contains(appJS,piece){t.Errorf("missing stable UI behavior: %s",piece)}
 }
 for _,piece:=range []string{"uiV015CSS","mobileRefinementCSS","uiInteractionCSS"}{
  if !strings.Contains(strings.Join([]string{uiV015CSS,mobileRefinementCSS,uiInteractionCSS},""),"@media"){t.Fatal("no responsive styles")}
  if piece==""{t.Fatal("invalid stylesheet name")}
 }
}
