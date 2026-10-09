package main

import "net/http"

func init() {
 // The main server uses this handler at /static/ to serve embedded assets.
}
func staticHandler() http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  switch r.URL.Path{
  case "/static/style.css":w.Header().Set("Content-Type","text/css; charset=utf-8");w.Write([]byte(styleCSS))
  case "/static/app.js":w.Header().Set("Content-Type","application/javascript; charset=utf-8");w.Write([]byte(appJS))
  default:http.NotFound(w,r)
  }
 })
}
const styleCSS=`*{box-sizing:border-box}body{background:#0b1220;color:#eaf1fc;font:15px system-ui,sans-serif;margin:0}header{display:flex;justify-content:space-between;align-items:center;background:#121f33;padding:12px 5%}main{max-width:1100px;margin:28px auto;padding:0 18px}button{background:#293b58;border:0;color:white;padding:8px 14px;border-radius:8px;cursor:pointer}.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:12px}.stats>div,.node{background:#17243a;border:1px solid #2a3b53;border-radius:12px;padding:18px}.stats small,.stats strong{display:block}.stats small,.muted,.toolbar span{color:#97a6bc}.stats strong{font-size:27px;margin-top:8px}.toolbar{display:flex;justify-content:space-between;align-items:center;margin-top:26px}.node{margin:10px 0}.node-head,.metric-row{display:flex;justify-content:space-between;gap:12px}.node-head{font-weight:bold;margin-bottom:15px}.pill{font-size:12px}.pill.ok{color:#5ad6a3}.pill.bad{color:#fc8892}.meters{display:grid;grid-template-columns:repeat(3,1fr);gap:16px}.metric-row{font-size:13px;margin-bottom:6px}.bar{height:7px;border-radius:7px;background:#2b3b50;overflow:hidden}.fill{height:100%;background:#388add}.meta{font-size:12px;color:#a3b3c9;margin-top:14px;overflow-wrap:anywhere}@media(max-width:600px){.meters{grid-template-columns:1fr}.stats strong{font-size:23px}.stats>div{padding:12px}}`
const appJS=`const $=x=>document.getElementById(x);
function esc(v){return String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]))}
function metric(label,n){n=Math.max(0,Math.min(100,Number(n)||0));return '<div><div class="metric-row"><span>'+label+'</span><b>'+n.toFixed(1)+'%</b></div><div class="bar"><div class="fill" style="width:'+n+'%"></div></div></div>'}
function formatBytes(n){if(!n)return '0 B';const units=['B','KB','MB','GB','TB'];let i=0;while(n>=1024&&i<units.length-1){n/=1024;i++}return n.toFixed(i?1:0)+' '+units[i]}
async function refresh(){try{const r=await fetch('/api/v1/nodes',{cache:'no-store'});if(r.status===401){location.href='/login';return}if(!r.ok)throw Error('HTTP '+r.status);const nodes=await r.json();nodes.sort((a,b)=>a.name.localeCompare(b.name));const up=nodes.filter(n=>n.online).length;$('total').textContent=nodes.length;$('online').textContent=up;$('offline').textContent=nodes.length-up;$('nodes').innerHTML=nodes.map(n=>'<article class="node"><div class="node-head"><span>'+esc(n.name)+'</span><span class="pill '+(n.online?'ok':'bad')+'">'+(n.online?'● 在线':'● 离线')+'</span></div><div class="meters">'+metric('CPU',n.cpu)+metric('内存',n.memory)+metric('磁盘 /',n.disk)+'</div><div class="meta">'+esc(n.hostname)+' · '+esc(n.os)+'/'+esc(n.arch)+' · 流量 ↓'+formatBytes(n.rx_bytes)+' ↑'+formatBytes(n.tx_bytes)+' · 运行 '+Math.floor(n.uptime/3600)+'h</div></article>').join('')||'<p>尚无节点。请安装 Agent 并连接到服务端。</p>';$('updated').textContent='更新：'+new Date().toLocaleTimeString()}catch(e){$('updated').textContent='获取失败：'+e.message}}
refresh();setInterval(refresh,5000);`
