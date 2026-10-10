package main

import ("net/http";"strings")

func init() {
 // The main server uses this handler at /static/ to serve embedded assets.
}
func staticHandler() http.Handler {
 return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){
  if strings.HasPrefix(r.URL.Path,"/static/flags/"){serveFlagSVG(w,r);return}
  switch r.URL.Path{
  case "/static/style.css":w.Header().Set("Content-Type","text/css; charset=utf-8");w.Write([]byte(styleCSS));w.Write([]byte(mobileRefinementCSS));w.Write([]byte(uiInteractionCSS));w.Write([]byte(uiV015CSS))
  case "/static/app.js":w.Header().Set("Content-Type","application/javascript; charset=utf-8");w.Write([]byte(appJS))
  default:http.NotFound(w,r)
  }
 })
}
const styleCSS=styleCSSFoundation+styleCSSDashboard+styleCSSThemes+styleCSSCards+styleCSSAdmin+styleCSSResponsive+styleCSSPersonal
const appJS=`const $=x=>document.getElementById(x);
const validThemes=['glass','dark','light'];
let currentTheme='glass';
try{const stored=localStorage.getItem('monitor-theme');if(validThemes.includes(stored))currentTheme=stored}catch(e){}
function setTheme(next){if(!validThemes.includes(next))next='glass';currentTheme=next;document.documentElement.dataset.theme=next;const chooser=$('theme-select');if(chooser)chooser.value=next;document.querySelectorAll('[data-theme-choice]').forEach(b=>{const active=b.dataset.themeChoice===next;b.classList.toggle('selected',active);b.setAttribute('aria-pressed',String(active))});try{localStorage.setItem('monitor-theme',next)}catch(e){}}
setTheme(currentTheme);
$('theme-select').addEventListener('change',e=>setTheme(e.target.value));
document.querySelectorAll('[data-theme-choice]').forEach(b=>b.addEventListener('click',()=>{setTheme(b.dataset.themeChoice);$('theme-status').textContent='已切换为：'+({glass:'梦幻玻璃',dark:'经典深色',light:'简约浅色'}[currentTheme])}));

function esc(v){return String(v??'').replace(/[&<>"']/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]))}
function metric(label,n){n=Math.max(0,Math.min(100,Number(n)||0));return '<div><div class="metric-row"><span>'+label+'</span><b>'+n.toFixed(1)+'%</b></div><div class="bar"><div class="fill" style="width:'+n+'%"></div></div></div>'}
function formatBytes(n){if(!n)return '0 B';const units=['B','KB','MB','GB','TB'];let i=0;while(n>=1024&&i<units.length-1){n/=1024;i++}return n.toFixed(i?1:0)+' '+units[i]}
let prevNames="";let cachedNodes=[];let nodeFetchSequence=0,nodeMutationSequence=0,trafficFetchSequence=0,alertFetchSequence=0,historyFetchSequence=0,detailHistorySequence=0,nodeFetchBusy=false,trafficFetchBusy=false;
async function refresh(){if(nodeFetchBusy)return;nodeFetchBusy=true;const ticket=++nodeFetchSequence,mutation=nodeMutationSequence;try{const r=await fetch('/api/v1/nodes',{cache:'no-store'});if(r.status===401){location.href='/login';return}if(!r.ok)throw Error('HTTP '+r.status);const nodes=await r.json();if(ticket!==nodeFetchSequence||mutation!==nodeMutationSequence)return;nodes.sort(nodeOrder);cachedNodes=nodes;if(selectedNode)updateDetail();updateGroupOptions(nodes);const names=nodes.map(n=>n.name+'|'+(n.display_name||n.hostname||n.name)).join('|');if(names!==prevNames){prevNames=names;const ed=$("metadata-name"),old=ed.value;ed.replaceChildren(...nodes.map(n=>new Option(n.display_name||n.name,n.name)));if(nodes.some(n=>n.name===old))ed.value=old;loadMetadata();const sel=$('history-name');const limitSelect=$('limit-name'),previousLimit=limitSelect.value;limitSelect.replaceChildren(...nodes.map(n=>new Option(n.display_name||n.name,n.name)));if(nodes.some(n=>n.name===previousLimit))limitSelect.value=previousLimit;fillLimitsForm();const previousHistoryName=sel.value;sel.replaceChildren(...nodes.map(n=>new Option(n.display_name||n.hostname||n.name,n.name)));if(nodes.some(n=>n.name===previousHistoryName))sel.value=previousHistoryName;updateHistoryPickerLabel();drawHistory();}renderNodeEditorRows(nodes);const up=nodes.filter(n=>n.online).length;$('total').textContent=nodes.length;$('online').textContent=up;$('offline').textContent=nodes.length-up;$('overview-hint').textContent='节点实时状态与资源使用';renderOverviewNodes();$('updated').textContent='更新 '+new Date().toLocaleTimeString()}catch(e){if(ticket===nodeFetchSequence)$('updated').textContent='节点更新失败：'+e.message+'（保留上次数据）'}finally{nodeFetchBusy=false}}
let historyChartPoints=[],detailChartPoints=[],historyChartCursor=-1,detailChartCursor=-1;
function updateHistoryPickerLabel(){const select=$('history-name');const o=select.selectedOptions[0];$('history-node-pick').textContent=o?o.textContent+' ▾':'选择服务器 ▾'}
function renderHistoryPicker(){
 const search=$('history-node-query').value.trim().toLowerCase();
 const root=$('history-node-options');root.replaceChildren();
 for(const n of cachedNodes.filter(n=>!search||[n.name,n.display_name,n.hostname,n.group,n.location].some(x=>String(x||'').toLowerCase().includes(search)))){
  const button=document.createElement('button');button.type='button';button.className='node-picker-item';
  const name=document.createElement('strong');name.textContent=n.display_name||n.hostname||n.name;
  const sub=document.createElement('small');sub.textContent=[n.group,n.location].filter(Boolean).join(' · ')||'节点';
  button.append(name,sub);button.setAttribute('aria-pressed',String(n.name===$('history-name').value));
  button.addEventListener('click',()=>{$('history-name').value=n.name;updateHistoryPickerLabel();$('history-name').dispatchEvent(new Event('change',{bubbles:true}));$('history-node-dialog').close()});root.append(button)
 }
 if(!root.children.length){const t=document.createElement('p');t.className='muted';t.textContent='没有找到匹配的服务器';root.append(t)}
}
$('history-node-pick').addEventListener('click',()=>{$('history-node-query').value='';renderHistoryPicker();const dialog=$('history-node-dialog');if(typeof dialog.showModal==='function'){dialog.showModal();$('history-node-query').focus()}else{dialog.setAttribute('open','');$('history-node-query').focus()}});
$('history-node-close').addEventListener('click',()=>{$('history-node-dialog').close()});
$('history-node-dialog').addEventListener('click',e=>{if(e.target===$('history-node-dialog'))$('history-node-dialog').close()});
$('history-node-query').addEventListener('input',renderHistoryPicker);
function bindHistoryPoint(canvasId,statusId,kind){
 const canvas=$(canvasId),out=$(statusId);
 const getPoints=()=>kind==='detail'?detailChartPoints:historyChartPoints;
 const getCursor=()=>kind==='detail'?detailChartCursor:historyChartCursor;
 const setCursor=i=>{if(kind==='detail')detailChartCursor=i;else historyChartCursor=i};
 const show=i=>{
  const points=getPoints();if(i<0||i>=points.length)return;
  setCursor(i);const p=points[i],date=new Date(p.time*1000).toLocaleString();
  out.textContent=date+' · '+(kind==='detail'?'CPU '+Number(p.cpu).toFixed(1)+'% · 内存 '+Number(p.memory).toFixed(1)+'% · 磁盘 '+Number(p.disk).toFixed(1)+'%':({cpu:'CPU',memory:'内存',disk:'磁盘'})[$('history-metric').value]+' '+Number(p[$('history-metric').value]).toFixed(1)+'%')
 };
 const byPosition=clientX=>{
  const points=getPoints();if(!points.length)return;
  const rect=canvas.getBoundingClientRect();if(rect.width===0)return;
  const x=(clientX-rect.left)/rect.width*canvas.width;
  const min=points[0].time,max=Math.max(min+1,points[points.length-1].time);
  const fraction=Math.max(0,Math.min(1,(x-(kind==='detail'?42:40))/(kind==='detail'?842:canvas.width-55)));
  const at=min+fraction*(max-min);let index=0;
  for(let i=1;i<points.length;i++)if(Math.abs(points[i].time-at)<Math.abs(points[index].time-at))index=i;
  show(index)
 };
 canvas.addEventListener('pointermove',e=>byPosition(e.clientX));
 canvas.addEventListener('pointerdown',e=>byPosition(e.clientX));
 canvas.addEventListener('keydown',e=>{
  if(e.key!=='ArrowLeft'&&e.key!=='ArrowRight'&&e.key!=='Home'&&e.key!=='End')return;
  const pts=getPoints();if(!pts.length)return;e.preventDefault();
  const next=e.key==='Home'?0:e.key==='End'?pts.length-1:Math.min(pts.length-1,Math.max(0,getCursor()<0?0:getCursor()+(e.key==='ArrowRight'?1:-1)));
  show(next)
 });
}
bindHistoryPoint('history-chart','history-point','history');
bindHistoryPoint('detail-chart','detail-point','detail');
async function drawHistory(){const ticket=++historyFetchSequence,name=$('history-name').value,hours=$('history-hours').value,metricValue=$('history-metric').value;if(!name)return;try{const r=await fetch('/api/v1/history?name='+encodeURIComponent(name)+'&hours='+encodeURIComponent($('history-hours').value));if(!r.ok)throw Error('HTTP '+r.status);const data=await r.json();if(ticket!==historyFetchSequence||$('history-name').value!==name||$('history-hours').value!==hours||$('history-metric').value!==metricValue)return;historyChartPoints=data;historyChartCursor=-1;$('history-point').textContent=data.length?'触摸或移动到曲线上查看具体数值':'当前时间段没有监控采样';const canvas=$('history-chart');const ctx=canvas.getContext('2d');const w=canvas.width,h=canvas.height;ctx.clearRect(0,0,w,h);ctx.strokeStyle='#34455e';ctx.lineWidth=1;for(let y=0;y<=100;y+=25){const py=20+(100-y)/100*(h-50);ctx.beginPath();ctx.moveTo(40,py);ctx.lineTo(w-15,py);ctx.stroke();ctx.fillStyle='#a6b5c9';ctx.font='12px sans-serif';ctx.fillText(y+'%',4,py+4)}const metric=$('history-metric').value;if(data.length){const lo=data[0].time,hi=Math.max(lo+1,data[data.length-1].time);ctx.beginPath();ctx.strokeStyle='#47a0f5';ctx.lineWidth=2;data.forEach((p,i)=>{const x=40+(p.time-lo)/(hi-lo)*(w-55),y=20+(100-p[metric])/100*(h-50);i?ctx.lineTo(x,y):ctx.moveTo(x,y)});ctx.stroke()}$('history-status').textContent=data.length+' 个采样区间';drawTraffic()}catch(e){if(ticket===historyFetchSequence)$('history-status').textContent='历史数据获取失败：'+e.message}}
['history-name','history-hours','history-metric'].forEach(k=>document.addEventListener('change',e=>{if(e.target.id===k)drawHistory()}));async function drawTraffic(){const name=$('history-name').value;if(!name)return;try{const r=await fetch('/api/v1/traffic?name='+encodeURIComponent(name)+'&hours='+encodeURIComponent($('history-hours').value));if(!r.ok)throw Error('HTTP '+r.status);const d=await r.json();$('traffic-summary').textContent='接收 ↓ '+formatBytes(d.rx)+' · 发送 ↑ '+formatBytes(d.tx)+'（统计范围内的采样增量）';const max=Math.max(1,...d.buckets.map(b=>b.rx+b.tx));$('traffic-bars').innerHTML=d.buckets.slice(-24).map(b=>'<div class="item"><span>'+new Date(b.time*1000).toLocaleString()+'</span><div class="track"><div class="value" style="width:'+Math.round((b.rx+b.tx)/max*100)+'%"></div></div><span>'+formatBytes(b.rx+b.tx)+'</span></div>').join('')}catch(e){$('traffic-summary').textContent=e.message}}
function switchSettingsTab(tab){if(tab!=='nodes'&&document.getElementById('settings-nodes')&&!document.getElementById('settings-nodes').hidden&&!confirmPendingEdits())return;document.querySelectorAll('.settings-tab').forEach(b=>{const selected=b.dataset.settingsTab===tab;b.classList.toggle('selected',selected);b.setAttribute('aria-selected',String(selected));b.tabIndex=0});document.querySelectorAll('.settings-pane').forEach(p=>p.hidden=p.id!=='settings-'+tab)}
document.querySelectorAll('.settings-tab').forEach(b=>b.addEventListener('click',()=>switchSettingsTab(b.dataset.settingsTab)));
function toggleAddPanel(show){switchView('settings');switchSettingsTab('nodes');$('add-panel').hidden=!show;if(show){$('server-url').value=$('server-url').value||(location.origin.startsWith('https://')?location.origin:'');$('add-panel').scrollIntoView({behavior:'smooth',block:'nearest'})}}
$('add-node').addEventListener('click',()=>toggleAddPanel(true));
$('show-add-node').addEventListener('click',()=>toggleAddPanel(true));
$('hide-add-node').addEventListener('click',()=>{$('add-panel').hidden=true;$('show-add-node').focus()});
switchSettingsTab('nodes');
function generatedInstallCommand(url,id,token){
 return 'curl -fsSL https://raw.githubusercontent.com/sockc/Monitor/main/scripts/install-release.sh -o monitor-install.sh && sudo env MONITOR_SERVER='+url+' MONITOR_NODE_NAME='+id+' MONITOR_AGENT_TOKEN='+token+' bash monitor-install.sh agent';
}
function generatedReconnectCommand(url,id,token){
 return 'curl -fsSL https://raw.githubusercontent.com/sockc/Monitor/main/scripts/rebind-agent.sh -o /tmp/monitor-rebind.sh && sudo env MONITOR_SERVER='+url+' MONITOR_NODE_NAME='+id+' MONITOR_AGENT_TOKEN='+token+' bash /tmp/monitor-rebind.sh';
}
$('create-node').addEventListener('submit',async e=>{
 e.preventDefault();
 let parsed;try{parsed=new URL($('server-url').value.trim())}catch(error){$('create-status').textContent='Server 地址无效';return}
 const safeHost=[...parsed.hostname].every(c=>'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.-'.includes(c));
 if(parsed.protocol!=='https:'||!safeHost||!parsed.hostname||parsed.username||parsed.password||parsed.search||parsed.hash||parsed.pathname!=='/'){$('create-status').textContent='请输入完整的 HTTPS Server 域名';return}
 const url=parsed.origin;
 try{
  const r=await fetch('/api/v1/tokens',{method:'POST',headers:{'Content-Type':'application/json'},body:'{}'});
  if(!r.ok)throw Error('HTTP '+r.status);
  const data=await r.json();
  $('install-command').value=generatedInstallCommand(url,data.name,data.token);
  $('create-status').textContent='节点已创建（随机 ID '+data.name+'），复制命令到新 VPS 安装即可。';
  await refresh()
 }catch(err){$('create-status').textContent='创建失败：'+err.message}
});
$('copy-command').addEventListener('click',async()=>{
 try{await navigator.clipboard.writeText($('install-command').value);$('create-status').textContent='已复制到剪贴板'}
 catch(e){$('create-status').textContent='复制失败，请手动复制'}
});
let editorSignature='';
function renderNodeEditorRows(nodes,force=false){
 const area=$('node-editor-list');
 if(!area)return;
 const signature=nodes.map(n=>[n.name,n.display_name,n.group,n.manual_location,n.sort_order,n.online].join('|')).join('\\n');
 if(!force&&signature===editorSignature)return;
 if(!force&&(area.contains(document.activeElement)||area.dataset.dirty==='true'))return;
 editorSignature=signature;
 if(!nodes.length){area.innerHTML='<p class="field-hint">暂无节点，请点击右上角「添加节点」自动生成。</p>';return}
 area.innerHTML=nodes.map(n=>{
  const id=esc(n.name),label=esc(n.display_name||n.name),group=esc(n.group||''),loc=esc(n.manual_location||'');
  return '<div class="node-edit-row" data-edit-node="'+id+'">'+
   '<div class="edit-node-identity"><span class="order-move-controls"><button type="button" class="node-drag-handle" draggable="true" title="拖动调整服务器顺序" aria-label="拖动 '+label+' 排序">⠿</button><button type="button" data-node-action="up" class="node-move" title="上移" aria-label="将 '+label+' 上移">↑</button><button type="button" data-node-action="down" class="node-move" title="下移" aria-label="将 '+label+' 下移">↓</button></span><span class="edit-node-state '+(n.online?'on':'off')+'"></span><span class="edit-node-label"><strong>'+label+'</strong><small title="固定节点 ID">ID: '+id+'</small></span></div>'+
   '<label><span class="node-edit-mobile-label">首页名称</span><input class="row-display" maxlength="60" placeholder="首页显示名称" aria-label="'+id+' 的首页名称" value="'+esc(n.display_name||'')+'"></label>'+
   '<label><span class="node-edit-mobile-label">分组</span><input class="row-group" maxlength="40" placeholder="分组" aria-label="'+id+' 的分组" value="'+group+'"></label>'+
   '<label><span class="node-edit-mobile-label">位置</span><input class="row-location" maxlength="80" placeholder="自动 IP 位置" aria-label="'+id+' 的手动位置" value="'+loc+'"></label>'+
   '<label><span class="node-edit-mobile-label">排序</span><input class="row-sort" type="number" min="0" max="9999" step="1" aria-label="'+id+' 的首页顺序" title="1 最靠前；0 使用默认名称排序" value="'+(Number(n.sort_order)||0)+'"></label>'+
   '<button type="button" data-node-action="toggle" class="subtle-btn node-edit-toggle" aria-expanded="false">编辑节点</button>'+
   '<div class="node-edit-actions"><button type="button" data-node-action="save" class="primary-btn">保存</button><button type="button" data-node-action="advanced" class="subtle-btn">详细配置</button><details class="node-maintenance"><summary>节点维护</summary><div class="maintenance-actions"><button type="button" data-node-action="rebind" class="subtle-btn">重新连接</button><button type="button" data-node-action="delete" class="danger-btn">删除节点</button></div></details></div>'+
   '<div class="row-feedback" aria-live="polite"></div></div>'
 }).join('')
}
function showAdvancedNode(id){
 const sel=$('metadata-name');
 if(![...sel.options].some(o=>o.value===id))return;
 sel.value=id;loadMetadata();$('metadata-advanced').hidden=false;$('metadata-advanced').scrollIntoView({behavior:'smooth',block:'start'});$('metadata-display').focus()
}
$('metadata-advanced-close').addEventListener('click',()=>{$('metadata-advanced').hidden=true});
$('node-rebind-close').addEventListener('click',()=>{$('node-rebind-panel').hidden=true;$('node-rebind-command').value=''});
$('node-rebind-copy').addEventListener('click',async()=>{
 try{await navigator.clipboard.writeText($('node-rebind-command').value);$('node-rebind-status').textContent='已复制；请到对应 VPS 执行'}
 catch(e){$('node-rebind-status').textContent='复制失败，请手动复制'}
});
let nodeOrderSaving=false,nodeDraggedID='';
async function persistNodeOrder(){
 const area=$('node-editor-list');
 if(nodeOrderSaving)return;
 const names=[...area.querySelectorAll('[data-edit-node]')].map(row=>row.dataset.editNode);
 if(names.length!==cachedNodes.length)return;
 nodeOrderSaving=true;$('node-order-status').textContent='正在保存顺序…';
 try{
  const r=await fetch('/api/v1/node-ui',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({action:'reorder',names})});
  if(!r.ok)throw Error(await r.text());
  nodeMutationSequence++;
  names.forEach((name,i)=>{const n=cachedNodes.find(x=>x.name===name);if(n)n.sort_order=i+1});
  cachedNodes.sort(nodeOrder);editorSignature='';$('node-order-status').textContent='顺序已同步到服务器';
  renderNodeEditorRows(cachedNodes,true);renderOverviewNodes();
 }catch(e){$('node-order-status').textContent='排序失败：'+e.message;editorSignature='';renderNodeEditorRows(cachedNodes,true)}
 finally{nodeOrderSaving=false}
}
$('node-editor-list').addEventListener('dragstart',e=>{
 const handle=e.target.closest('.node-drag-handle');if(!handle)return;
 const area=$('node-editor-list');if(area.dataset.dirty==='true'||nodeOrderSaving){e.preventDefault();$('node-order-status').textContent='请先保存节点资料，再调整排序';return}
 const row=handle.closest('[data-edit-node]');if(!row)return;
 nodeDraggedID=row.dataset.editNode;e.dataTransfer.effectAllowed='move';e.dataTransfer.setData('text/plain',nodeDraggedID);
});
$('node-editor-list').addEventListener('dragover',e=>{
 if(!nodeDraggedID)return;const target=e.target.closest('[data-edit-node]');if(!target||target.dataset.editNode===nodeDraggedID)return;
 e.preventDefault();e.dataTransfer.dropEffect='move'
});
$('node-editor-list').addEventListener('drop',e=>{
 if(!nodeDraggedID)return;
 const target=e.target.closest('[data-edit-node]'),source=[...$('node-editor-list').children].find(x=>x.dataset.editNode===nodeDraggedID);
 if(target&&source&&target!==source){
  e.preventDefault();const before=[...$('node-editor-list').children].indexOf(source)<[...$('node-editor-list').children].indexOf(target);
  $('node-editor-list').insertBefore(source,before?target.nextSibling:target);persistNodeOrder();
 }
 nodeDraggedID=''
});
$('node-editor-list').addEventListener('dragend',()=>{nodeDraggedID=''});

function updateNodeEditorDirty(){const area=$('node-editor-list');area.dataset.dirty=[...area.querySelectorAll('[data-edit-node]')].some(row=>row.dataset.dirty==='true')?'true':'false'}
$('node-editor-list').addEventListener('input',e=>{if(e.target.matches('input')){const row=e.target.closest('[data-edit-node]');if(row)row.dataset.dirty='true';updateNodeEditorDirty()}});
$('node-editor-list').addEventListener('click',async e=>{
 const btn=e.target.closest('button[data-node-action]');if(!btn)return;
 const row=btn.closest('[data-edit-node]');if(!row)return;
 const id=row.dataset.editNode,action=btn.dataset.nodeAction,feedback=row.querySelector('.row-feedback');
 const node=cachedNodes.find(n=>n.name===id);if(!node)return;
 if(action==='up'||action==='down'){
  if($('node-editor-list').dataset.dirty==='true'||nodeOrderSaving){$('node-order-status').textContent='请先保存节点资料，再调整排序';return}
  const rows=[...$('node-editor-list').querySelectorAll('[data-edit-node]')],index=rows.indexOf(row),next=rows[index+(action==='up'?-1:1)];
  if(!next)return;
  if(action==='up')row.parentNode.insertBefore(row,next);else row.parentNode.insertBefore(next,row);
  persistNodeOrder();return
 }
 if(action==='toggle'){const expanded=row.classList.toggle('editing');btn.textContent=expanded?'收起编辑':'编辑节点';btn.setAttribute('aria-expanded',String(expanded));return}
 if(action==='advanced'){showAdvancedNode(id);return}
 if(action==='save'){
  btn.disabled=true;
  const payload={name:id,display_name:row.querySelector('.row-display').value.trim(),group:row.querySelector('.row-group').value.trim(),location:row.querySelector('.row-location').value.trim(),sort_order:Number(row.querySelector('.row-sort').value),notes:node.notes||''};
  try{
   const r=await fetch('/api/v1/node-metadata',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(payload)});
   if(!r.ok)throw Error(await r.text());
   nodeMutationSequence++;
   Object.assign(node,{display_name:payload.display_name,group:payload.group,manual_location:payload.location,location:payload.location||node.auto_location||'',sort_order:payload.sort_order});
   const label=row.querySelector('.edit-node-identity strong');if(label)label.textContent=payload.display_name||node.hostname||id;
   row.dataset.dirty='false';updateNodeEditorDirty();
   feedback.textContent='✓ 已保存至服务器';feedback.dataset.result='success';
   cachedNodes.sort(nodeOrder);prevNames='';renderOverviewNodes();
   if($('node-editor-list').dataset.dirty!=='true')refresh()
  }catch(error){feedback.textContent='保存失败，修改内容已保留：'+error.message;feedback.dataset.result='error'}finally{btn.disabled=false}
  return
 }
 if(action==='rebind'){
  if(!confirm('重新连接 '+(node.display_name||id)+' 会轮换令牌，旧 Agent 必须执行新命令。确定继续吗？'))return;
  btn.disabled=true;
  try{
   const r=await fetch('/api/v1/node-rebind',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({name:id})});
   if(!r.ok)throw Error(await r.text());
   const data=await r.json();
   $('node-rebind-command').value=generatedReconnectCommand(location.origin,data.name,data.token);
   $('node-rebind-status').textContent='节点 '+(node.display_name||id)+' 的命令已生成';
   $('node-rebind-panel').hidden=false;
   $('node-rebind-panel').scrollIntoView({behavior:'smooth',block:'nearest'})
  }catch(error){feedback.textContent='重新连接失败：'+error.message}finally{btn.disabled=false}
  return
 }
 if(action==='delete'){
  if(prompt('删除将永久清除节点历史。请输入固定节点 ID '+id+' 确认删除：')!==id)return;
  btn.disabled=true;
  try{
   const r=await fetch('/api/v1/node-manage',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({action:'delete',name:id})});
   if(!r.ok)throw Error(await r.text());prevNames='';editorSignature='';await refresh()
  }catch(error){feedback.textContent='删除失败：'+error.message}finally{btn.disabled=false}
 }
});
function fillPlanLimits(name){const v=nodeLimitsCache[name]||{};$("metadata-quota").value=v.quota_gb||0;$("metadata-expires").value=v.expires_on||""}
function loadMetadata(){const n=cachedNodes.find(x=>x.name===$("metadata-name").value);if(!n)return;
 const p=n.profile||{};
 $("metadata-display").value=n.display_name||"";
 $("metadata-group").value=n.group||"";
 $("metadata-sort").value=Number(n.sort_order)||0;
 $("metadata-location").value=n.manual_location||"";
 $("metadata-auto-location").textContent=n.auto_location?"IP 自动识别："+n.auto_location+(n.public_ip?"（"+n.public_ip+"）":"")+" · 清空手动位置即可恢复自动识别":n.public_ip?"公网 IP "+n.public_ip+" · 定位暂不可用":"等待 Agent 上报公网 IP 位置；旧 Agent 需升级";
 $("metadata-notes").value=n.notes||"";
 $("metadata-provider").value=p.provider||"";
 $("metadata-country").value=p.country_code||"";
 $("metadata-price").value=p.price_value||0;
 $("metadata-currency").value=p.price_currency||"USD";
 $("metadata-cycle").value=p.billing_cycle||"year";
 $("metadata-port").value=p.port_mbps||0;
 $("metadata-ipv4").value=String(p.has_ipv4??-1);
 $("metadata-ipv6").value=String(p.has_ipv6??-1);
 $("metadata-start").value=p.period_start||"";
 fillPlanLimits(n.name)
}
$("metadata-name").addEventListener("change",loadMetadata);
$("metadata-form").addEventListener("submit",async e=>{
 e.preventDefault();
 const selected=$("metadata-name").value;
 const payload={name:selected,display_name:$("metadata-display").value,group:$("metadata-group").value,sort_order:Number($("metadata-sort").value),location:$("metadata-location").value,notes:$("metadata-notes").value,
  profile:{provider:$("metadata-provider").value.trim(),country_code:$("metadata-country").value.trim().toUpperCase(),price_value:Number($("metadata-price").value),price_currency:$("metadata-currency").value,billing_cycle:$("metadata-cycle").value,period_start:$("metadata-start").value,port_mbps:Number($("metadata-port").value),has_ipv4:Number($("metadata-ipv4").value),has_ipv6:Number($("metadata-ipv6").value)},
  quota_gb:Number($("metadata-quota").value),expires_on:$("metadata-expires").value
 };
 try{
  const r=await fetch("/api/v1/node-metadata",{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify(payload)});
  if(!r.ok)throw Error(await r.text());
  $("metadata-status").textContent="保存成功，套餐与到期数据已同步";
  prevNames="";await loadNodeLimits();await refresh();await refreshPeriods()
 }catch(err){$("metadata-status").textContent="保存失败："+err.message}
});
let restoreGroup='';
try{const v=JSON.parse(sessionStorage.getItem('monitor-overview-filters')||'{}');if(typeof v.search==='string')$('node-search').value=v.search.slice(0,120);if(typeof v.group==='string')restoreGroup=v.group.slice(0,80)}catch(e){}
function saveOverviewFilters(){try{sessionStorage.setItem('monitor-overview-filters',JSON.stringify({search:$('node-search').value,group:$('node-group').value}))}catch(e){}}
function updateGroupOptions(nodes){const el=$('node-group'),value=el.value,groups=[...new Set(nodes.map(n=>n.group).filter(Boolean))].sort();const next=['',...groups];if([...el.options].map(x=>x.value).join('|')!==next.join('|')){el.replaceChildren(new Option('全部分组',''),...groups.map(x=>new Option(x,x)));el.value=(restoreGroup&&groups.includes(restoreGroup))?restoreGroup:value;restoreGroup=''}}
function nodeOrder(a,b){const aa=Number(a.sort_order)>0?Number(a.sort_order):10000,bb=Number(b.sort_order)>0?Number(b.sort_order):10000;return aa-bb||(a.display_name||a.hostname||a.name).localeCompare(b.display_name||b.hostname||b.name,'zh-Hans-CN')}
function nodeHasAlert(n){
 const p=trafficPeriods[n.name]||{},quota=Number(p.quota_gb)||0,used=(Number(p.month_rx)||0)+(Number(p.month_tx)||0);
 return !n.online||(n.online&&[n.cpu,n.memory,n.disk].some(v=>Number(v)>=90))||
  (quota>0&&p.has_samples&&used>=quota*800000000)||
  (p.expires_on&&daysUntil(p.expires_on,p.timezone)<=7)||
  lastAlertEvents.some(e=>!e.end&&e.node===n.name);
}
function filteredNodes(nodes){const q=$('node-search').value.trim().toLowerCase(),group=$('node-group').value;return nodes.filter(n=>(!group||n.group===group)&&(!q||[n.name,n.display_name,n.group,n.location,n.os,n.notes,n.hostname].some(x=>String(x||'').toLowerCase().includes(q)))&&(nodeFilter!=='favorite'||n.favorite)&&(nodeFilter!=='alert'||nodeHasAlert(n))).sort(nodeOrder)}
['node-search','node-group'].forEach(id=>{const control=$(id);control.addEventListener('input',()=>{saveOverviewFilters();renderOverviewNodes()});control.addEventListener('change',()=>{saveOverviewFilters();renderOverviewNodes()})});
let currentLayout='compact';let trafficPeriods={};let nodeLimitsCache={};
const cardFieldLabels={cpu:'CPU 使用率',memory:'内存使用率',disk:'磁盘使用率',speed:'实时上传／下载',traffic:'本月流量',price:'套餐与到期',cumulative:'累计网卡流量',tags:'IPv4／IPv6 与带宽',boot:'开机时间'};
const cardFieldDefaults={cards:['cpu','memory','disk','speed','traffic','price','cumulative','tags','boot'],compact:['cpu','memory','disk','speed','traffic'],list:['cpu','memory','disk']};
const cardFieldAllowed={cards:Object.keys(cardFieldLabels),compact:['cpu','memory','disk','speed','traffic','price','boot'],list:['cpu','memory','disk','speed','traffic']};
let cardFieldPreferences={};
try{const saved=JSON.parse(localStorage.getItem('monitor-card-fields')||'{}');for(const mode of Object.keys(cardFieldDefaults)){if(Array.isArray(saved[mode]))cardFieldPreferences[mode]=saved[mode].filter(x=>cardFieldAllowed[mode].includes(x))}}catch(e){}
function cardFieldOn(mode,key){return (cardFieldPreferences[mode]||cardFieldDefaults[mode]).includes(key)}
function saveCardFields(){try{localStorage.setItem('monitor-card-fields',JSON.stringify(cardFieldPreferences))}catch(e){}}
function updateCardFieldsEditor(){
 const mode=$('card-fields-layout').value,selected=cardFieldPreferences[mode]||cardFieldDefaults[mode],root=$('card-fields-options');
 root.replaceChildren();for(const key of cardFieldAllowed[mode]){
  const label=document.createElement('label'),check=document.createElement('input');check.type='checkbox';check.value=key;check.checked=selected.includes(key);check.addEventListener('change',()=>{const current=cardFieldPreferences[mode]||cardFieldDefaults[mode];cardFieldPreferences[mode]=check.checked?[...new Set([...current,key])]:current.filter(x=>x!==key);saveCardFields();renderOverviewNodes();updateCardFieldsEditor()});label.append(check,document.createTextNode(cardFieldLabels[key]));root.append(label);
 }
 $('card-fields-preview').textContent='预览：服务器名称 · 在线状态'+(selected.length?' · '+selected.map(x=>cardFieldLabels[x]).join(' · '):'（仅基础信息）')+' · 异常提醒';
}
$('card-fields-layout').addEventListener('change',updateCardFieldsEditor);
$('card-fields-reset').addEventListener('click',()=>{delete cardFieldPreferences[$('card-fields-layout').value];saveCardFields();updateCardFieldsEditor();renderOverviewNodes()});
updateCardFieldsEditor();
let nodeFilter='all',collapsedGroups={};
try{const f=localStorage.getItem('monitor-quick-filter');if(['all','favorite','alert'].includes(f))nodeFilter=f;const g=JSON.parse(localStorage.getItem('monitor-collapsed-groups')||'{}');if(g&&typeof g==='object'&&!Array.isArray(g))collapsedGroups=g}catch(e){}
function setNodeFilter(next){nodeFilter=next;try{localStorage.setItem('monitor-quick-filter',next)}catch(e){};document.querySelectorAll('[data-node-filter]').forEach(b=>{const on=b.dataset.nodeFilter===next;b.classList.toggle('selected',on);b.setAttribute('aria-pressed',String(on))});renderOverviewNodes()}
document.querySelectorAll('[data-node-filter]').forEach(b=>b.addEventListener('click',()=>setNodeFilter(b.dataset.nodeFilter)));
$('overview-alert-filter').addEventListener('click',()=>{setNodeFilter('alert');$('nodes').scrollIntoView({behavior:'smooth',block:'start'})});
document.querySelectorAll('[data-node-filter]').forEach(b=>{const on=b.dataset.nodeFilter===nodeFilter;b.classList.toggle('selected',on);b.setAttribute('aria-pressed',String(on))});

try{const saved=localStorage.getItem('monitor-node-layout');if(['cards','compact','list'].includes(saved))currentLayout=saved}catch(e){}
function updateLayout(){const container=$('nodes');container.classList.add('server-grid');container.dataset.layout=currentLayout;document.querySelectorAll('.layout-button').forEach(b=>{const active=b.dataset.layout===currentLayout;b.classList.toggle('active',active);b.setAttribute('aria-pressed',String(active))})}
document.querySelectorAll('.layout-button').forEach(b=>b.addEventListener('click',()=>{currentLayout=b.dataset.layout;try{localStorage.setItem('monitor-node-layout',currentLayout)}catch(e){}updateLayout();renderOverviewNodes()}));
updateLayout();
function capacity(n){return n>0?formatBytes(n):'—'}
function hardware(n){return '<div class="hardware-line"><span>'+((n.cpu_cores>0?n.cpu_cores+'C':'—')+' CPU')+'</span><span>内存 '+capacity(n.memory_total)+'</span><span>磁盘 '+capacity(n.disk_total)+'</span></div>'}
function resourceMetric(label,val,used,total){const v=Math.min(100,Math.max(0,Number(val)||0)),status=v>=90?'critical':v>=75?'hot':'',quantity=total>0?capacity(used)+' / '+capacity(total):'';return '<div class="resource-line"><div class="metric-row"><span>'+label+'</span><span class="resource-numbers">'+(quantity?'<small>'+quantity+'</small>':'')+'<b>'+Math.round(v)+'%</b></span></div><div class="bar"><div class="fill '+status+'" style="width:'+v+'%"></div></div></div>'}
function resourceWithCapacity(label,pct,used,total){return resourceMetric(label,pct,used,total)}
// Incremental patch: retain card DOM, pointer focus and open info popovers while
// telemetry refreshes every five seconds. Only changed attributes/text are updated.
function patchDashboardElement(oldNode,newNode){
 if(oldNode.nodeType!==newNode.nodeType||oldNode.nodeName!==newNode.nodeName){oldNode.replaceWith(newNode.cloneNode(true));return}
 if(oldNode.nodeType===3){if(oldNode.nodeValue!==newNode.nodeValue)oldNode.nodeValue=newNode.nodeValue;return}
 if(oldNode.nodeType!==1)return;
 for(const attr of [...oldNode.attributes])if(!newNode.hasAttribute(attr.name))oldNode.removeAttribute(attr.name);
 for(const attr of [...newNode.attributes])if(oldNode.getAttribute(attr.name)!==attr.value)oldNode.setAttribute(attr.name,attr.value);
 let a=oldNode.firstChild,b=newNode.firstChild;
 while(a||b){
  if(!a){oldNode.appendChild(b.cloneNode(true));b=b.nextSibling;continue}
  if(!b){const next=a.nextSibling;a.remove();a=next;continue}
  const nextA=a.nextSibling,nextB=b.nextSibling;
  patchDashboardElement(a,b);a=nextA;b=nextB;
 }
}
function renderOverviewNodes(){
 const visible=filteredNodes(cachedNodes),root=$('nodes'),layoutChanged=root.dataset.renderedLayout!==currentLayout;
 if(layoutChanged){root.replaceChildren();root.dataset.renderedLayout=currentLayout}
 const grouped=new Map();
 for(const n of visible){const g=n.group||'未分组';if(!grouped.has(g))grouped.set(g,[]);grouped.get(g).push(n)}
 const desired=[];
 for(const [group,nodes] of grouped){
  desired.push({group});
  if(!collapsedGroups[group])for(const n of nodes)desired.push({node:n});
 }
 const keys=new Set(desired.map(v=>v.node?'n:'+v.node.name:'g:'+v.group));
 for(const el of [...root.children]){const key=el.dataset.node?'n:'+el.dataset.node:el.dataset.groupHeader?'g:'+el.dataset.groupHeader:'';if(!keys.has(key))el.remove()}
 if(!visible.length){
  root.replaceChildren();const empty=document.createElement('div');empty.className='dashboard-empty';
  empty.textContent=nodeFilter==='favorite'?'还没有收藏的服务器':nodeFilter==='alert'?'当前筛选中没有异常服务器':'没有匹配的服务器';
  root.appendChild(empty)
 }else{
  const tpl=document.createElement('template');
  for(let i=0;i<desired.length;i++){
   const item=desired[i];let existing, fresh;
   if(item.node){
    const name=item.node.name;
    existing=[...root.children].find(el=>el.dataset.node===name);
    tpl.innerHTML=renderNode(item.node);fresh=tpl.content.firstElementChild;
   }else{
    existing=[...root.children].find(el=>el.dataset.groupHeader===item.group);
    fresh=document.createElement('button');fresh.type='button';fresh.className='node-group-heading';fresh.dataset.groupHeader=item.group;
    fresh.setAttribute('aria-expanded',String(!collapsedGroups[item.group]));
    const label=document.createElement('span');label.textContent=(collapsedGroups[item.group]?'▸ ':'▾ ')+item.group;
    const count=document.createElement('small');count.textContent=grouped.get(item.group).length+' 台';
    fresh.append(label,count);
   }
   if(!existing){root.insertBefore(fresh,root.children[i]||null)}
   else{if(root.children[i]!==existing)root.insertBefore(existing,root.children[i]||null);
    if(item.node)patchDashboardElement(existing,fresh);
    else{existing.setAttribute('aria-expanded',fresh.getAttribute('aria-expanded'));existing.replaceChildren(...fresh.childNodes)}
   }
  }
 }
 $('node-visible-count').textContent='显示 '+visible.length+' / '+cachedNodes.length+' 台';
 if(quickInfoNode){
  if(visible.some(n=>n.name===quickInfoNode)&&!collapsedGroups[(cachedNodes.find(n=>n.name===quickInfoNode)||{}).group||'未分组']){
   const trigger=[...document.querySelectorAll('.node-info-trigger')].find(b=>b.dataset.infoNode===quickInfoNode);
   if(trigger)trigger.setAttribute('aria-expanded','true');
   if(!layoutChanged)positionNodeInfo();else showNodeInfo(quickInfoNode,quickInfoPinned)
  }else hideNodeInfo()
 }
}
$('nodes').addEventListener('click',e=>{const b=e.target.closest('[data-group-header]');if(!b)return;const g=b.dataset.groupHeader;collapsedGroups[g]=!collapsedGroups[g];try{localStorage.setItem('monitor-collapsed-groups',JSON.stringify(collapsedGroups))}catch(err){};renderOverviewNodes()});
async function saveNodeFavorite(name,favorite){
 const res=await fetch('/api/v1/node-ui',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({action:'favorite',name,favorite})});
 if(!res.ok)throw Error(await res.text());
 nodeMutationSequence++;const node=cachedNodes.find(n=>n.name===name);if(node)node.favorite=favorite;renderOverviewNodes();
}
function dateInZone(zone){try{const d=new Intl.DateTimeFormat('en-GB',{timeZone:zone||'UTC',year:'numeric',month:'2-digit',day:'2-digit'}).formatToParts(new Date());const get=x=>d.find(y=>y.type===x).value;return get('year')+'-'+get('month')+'-'+get('day')}catch(e){return new Date().toISOString().slice(0,10)}}
function daysUntil(date,zone){if(!date)return null;return Math.round((Date.parse(date+'T00:00:00Z')-Date.parse(dateInZone(zone)+'T00:00:00Z'))/86400000)}
function countryFlag(code){const c=String(code||'').toUpperCase();if(!/^[A-Z]{2}$/.test(c))return '';return '<img class="flag-image" src="/static/flags/'+c.toLowerCase()+'.svg" alt="'+c+' 国旗" width="24" height="18" loading="lazy">'}
function binaryAmount(bytes){let n=Number(bytes);if(!Number.isFinite(n)||n<0)return '—';const units=['B','KiB','MiB','GiB','TiB'];let i=0;while(n>=1024&&i<units.length-1){n/=1024;i++}return n.toFixed(i?i<3?1:2:0)+' '+units[i]}
function planAmountGB(gb){const n=Number(gb)||0;return n>=1000?(n/1000).toFixed(n%1000===0?0:1)+' TB':n.toFixed(n%1===0?0:1)+' GB'}
function compactSpeed(bytes){const n=Number(bytes)||0;return formatBytes(Math.max(0,n)).replace(' KB',' K').replace(' MB',' M').replace(' GB',' G')+'/s'}
function planPort(mbps){const n=Number(mbps)||0;if(n<=0)return '';return n>=1000&&n%1000===0?(n/1000)+'G 口':n+'Mbps'}
function expiryInfo(p,profile){if(!p||!p.expires_on)return '';const d=daysUntil(p.expires_on,p.timezone);if(d===null)return '';const left=d<0?'已到期 '+(-d)+' 天':d===0?'今天到期':'剩余 '+d+' 天';let bar='';if(profile.period_start){const start=Date.parse(profile.period_start+'T00:00:00Z'),end=Date.parse(p.expires_on+'T00:00:00Z'),now=Date.parse(dateInZone(p.timezone)+'T00:00:00Z');if(Number.isFinite(start)&&Number.isFinite(end)&&end>start){const pct=Math.max(0,Math.min(100,(now-start)/(end-start)*100));bar='<span class="expiry-progress" title="本期已过 '+pct.toFixed(0)+'%"><i style="width:'+pct.toFixed(1)+'%"></i></span>'}}return '<span class="plan-remaining">'+left+'</span>'+bar}
function metricFive(label,value,percentage,kind){const isPercent=kind==='percent',pct=Math.max(0,Math.min(100,Number(percentage)||0));return '<div class="five-metric"><small>'+label+'</small><strong>'+value+'</strong>'+(isPercent?'<span class="five-track"><i class="'+(pct>=90?'critical':pct>=75?'warm':'')+'" style="width:'+pct+'%"></i></span>':'<span class="five-track ghost"></span>')+'</div>'}
function ipCapability(label,manualValue,publicIP,cls){const manual=Number(manualValue);if(manual===0||(manual!==1&&!publicIP))return '';const hint=publicIP?'出站公网地址：'+esc(publicIP):'手动标注支持；尚未获取出站地址';return '<span class="plan-tag '+cls+'" title="'+hint+'">'+label+'</span>'}
function nodeBootDate(value,withSeconds=false){
 if(!value)return '';
 const d=new Date(value);
 if(!Number.isFinite(d.getTime()))return '';
 const fmt={year:'numeric',month:'2-digit',day:'2-digit',hour:'2-digit',minute:'2-digit',hour12:false};
 if(withSeconds)fmt.second='2-digit';
 return new Intl.DateTimeFormat('zh-CN',fmt).format(d).replaceAll('/','-');
}
function nodeUptime(seconds){
 const n=Number(seconds)||0;if(n<=0)return '';
 const days=Math.floor(n/86400),hours=Math.floor(n%86400/3600),mins=Math.floor(n%3600/60);
 return days?days+'天 '+hours+'小时':hours?hours+'小时 '+mins+'分':mins+'分钟';
}
function nodeInfoContent(n){
 const rows=[];
 const add=(label,value)=>{if(value!==null&&value!==undefined&&String(value)!=='')rows.push('<div class="info-line"><dt>'+esc(label)+'</dt><dd>'+esc(value)+'</dd></div>')};
 add('固定节点 ID',n.name);
 add('主机名称',n.hostname);
 add('系统发行版',n.os);
 add('系统架构',n.arch);
 const cpu=[n.cpu_model,n.cpu_cores>0?n.cpu_cores+' 核':''].filter(Boolean).join(' · ');
 add('CPU',cpu);
 if(n.memory_total>0)add('物理内存',formatBytes(n.memory_total));
 if(n.disk_total>0)add('系统磁盘',formatBytes(n.disk_total));
 if(n.swap_total>0)add('Swap',formatBytes(n.swap_used)+' / '+formatBytes(n.swap_total));
 add('服务器位置',n.location);
 add('服务商',n.profile?.provider);
 add('公网 IPv4（出站）',n.public_ipv4);
 add('公网 IPv6（出站）',n.public_ipv6);
 add('开机时间（估算，本地）',nodeBootDate(n.boot_time,true));
 add('已运行',nodeUptime(n.uptime));
 add('Agent 版本',n.agent_version);
 if(n.last_seen){const d=nodeBootDate(n.last_seen,true);add('最后上报（本地）',d)}
 return '<div class="node-info-heading"><div><strong>'+esc(n.display_name||n.name)+'</strong><small>机器基本信息 · 不影响节点连接</small></div><button class="node-info-close" type="button" aria-label="关闭机器信息">×</button></div><dl class="node-info-grid">'+rows.join('')+'</dl><p class="node-info-foot">开机时间由 Agent 的系统运行时长推算；时间按当前浏览器所在时区显示。</p>';
}
let quickInfoNode='',quickInfoPinned=false,quickInfoTimer=0;
function hideNodeInfo(){
 clearTimeout(quickInfoTimer);quickInfoTimer=0;
 quickInfoNode='';quickInfoPinned=false;
 const panel=$('node-quick-info');panel.hidden=true;$('node-info-backdrop').hidden=true;
 document.querySelectorAll('.node-info-trigger[aria-expanded="true"]').forEach(b=>b.setAttribute('aria-expanded','false'));
}
function positionNodeInfo(){
 const panel=$('node-quick-info');if(panel.hidden||!quickInfoNode)return;
 const btn=[...document.querySelectorAll('#nodes .node-info-trigger')].find(b=>b.dataset.infoNode===quickInfoNode);
 if(!btn){hideNodeInfo();return}
 if(window.matchMedia('(max-width:779px)').matches){panel.style.left='';panel.style.top='';return}
 const rect=btn.getBoundingClientRect(),w=panel.offsetWidth,h=panel.offsetHeight;
 const left=Math.max(10,Math.min(rect.right-w,window.innerWidth-w-10));
 const below=rect.bottom+9,above=rect.top-h-9;
 const top=below+h<=window.innerHeight-10?below:Math.max(10,above);
 panel.style.left=Math.round(left)+'px';panel.style.top=Math.round(top)+'px';
}
function showNodeInfo(id,pinned=false){
 const n=cachedNodes.find(x=>x.name===id);if(!n)return;
 clearTimeout(quickInfoTimer);
 quickInfoNode=id;quickInfoPinned=pinned;
 const panel=$('node-quick-info');panel.innerHTML=nodeInfoContent(n);panel.hidden=false;$('node-info-backdrop').hidden=!window.matchMedia('(max-width:779px)').matches;
 document.querySelectorAll('.node-info-trigger').forEach(b=>b.setAttribute('aria-expanded',String(b.dataset.infoNode===id)));
 positionNodeInfo();
}
function scheduleHideNodeInfo(){clearTimeout(quickInfoTimer);if(!quickInfoPinned)quickInfoTimer=setTimeout(hideNodeInfo,180)}
function percentLabel(pct){return pct>0&&pct<0.1?'<0.1%':pct.toFixed(1)+'%'}
function renderNode(n){
 const name=esc(n.display_name||n.name),id=esc(n.name),online=Boolean(n.online),profile=n.profile||{},period=trafficPeriods[n.name]||{},code=n.country_code||'';
 const loc=n.location||'',provider=profile.provider||'';
 const context=[loc,provider].filter(Boolean).map(esc).join(' · ');
 const p=Number(profile.price_value)||0,cycle=({month:'月',quarter:'季',year:'年',one_time:'一次'})[profile.billing_cycle]||'年';
 const price=p>0?p.toLocaleString('en-US',{maximumFractionDigits:2})+' '+esc(profile.price_currency||'USD')+'/'+cycle:'';
 const expiration=expiryInfo(period,profile);
 const priceRow=price||expiration?'<div class="card-plan">'+(price?'<span class="plan-price">'+price+'</span>':'')+expiration+'</div>':'';
 const chosen=['cpu','memory','disk','speed'].filter(key=>cardFieldOn(currentLayout,key));
 const metrics='<div class="metric-five" style="--metric-count:'+Math.max(1,chosen.length+(chosen.includes('speed')?1:0))+'">'+(chosen.includes('cpu')?metricFive('CPU',online?Number(n.cpu).toFixed(1)+'%':'—',n.cpu,'percent'):'')+(chosen.includes('memory')?metricFive('内存',online?Number(n.memory).toFixed(1)+'%':'—',n.memory,'percent'):'')+(chosen.includes('disk')?metricFive('存储',online?Number(n.disk).toFixed(1)+'%':'—',n.disk,'percent'):'')+(chosen.includes('speed')?metricFive('上传',online?compactSpeed(n.tx_speed):'—',0,'speed')+metricFive('下载',online?compactSpeed(n.rx_speed):'—',0,'speed'):'')+'</div>';
 const cumulative='<div class="card-cumulative"><span>↑ 累计上传 <b>'+binaryAmount(n.tx_bytes||0)+'</b></span><span>↓ 累计下载 <b>'+binaryAmount(n.rx_bytes||0)+'</b></span></div>';
 const port=planPort(profile.port_mbps);
 const tags=(port?'<span class="plan-tag tag-bandwidth">'+port+'</span>':'')+(Number(period.quota_gb)>0?'<span class="plan-tag tag-quota">'+planAmountGB(period.quota_gb)+'/月</span>':'')+ipCapability('IPv4',profile.has_ipv4??-1,n.public_ipv4||'','tag-ip4')+ipCapability('IPv6',profile.has_ipv6??-1,n.public_ipv6||'','tag-ip6');
 const labels=tags?'<div class="capability-tags">'+tags+'</div>':'';
 const used=(Number(period.month_rx)||0)+(Number(period.month_tx)||0),quota=(Number(period.quota_gb)||0)*1000000000,has=Boolean(period.has_samples),ratio=quota>0?used/quota*100:0,clamp=Math.min(100,Math.max(0,ratio));
 const quotaBar=quota>0&&has?'<div class="card-quota"><div class="quota-figures"><span>'+planAmountGB(used/1000000000)+' <em>/ '+planAmountGB(period.quota_gb)+'</em></span><b>'+percentLabel(ratio)+'</b></div><div class="quota-linear"><i class="'+(ratio>=100?'danger':ratio>=80?'warn':'')+'" style="width:'+clamp.toFixed(1)+'%"></i></div></div>':'';
 const days=daysUntil(period.expires_on,period.timezone),warnings=[];
 if(days!==null&&days<=7)warnings.push(days<0?'服务器已到期':'即将到期');
 if(quota>0&&has&&ratio>=80)warnings.push('流量接近额度');
 if(online&&[n.cpu,n.memory,n.disk].some(v=>Number(v)>=90))warnings.push('资源偏高');
 const alert=warnings.length?'<div class="card-alerts">'+warnings.map(w=>'<span>'+w+'</span>').join('')+'</div>':'';
 const flag=countryFlag(code);
 const header='<div class="card-identity"><span class="status-dot '+(online?'on':'off')+'"></span>'+(flag?'<span class="country-flag" title="'+esc(code)+'">'+flag+'</span>':'')+'<div class="identity-text"><strong>'+name+'</strong>'+(context?'<small>'+context+'</small>':'')+'</div><span class="state-text">'+(online?'在线':'离线')+'</span></div>';
 const boot=nodeBootDate(n.boot_time);
 const bootLine=boot?'<div class="node-boot-line"><span>开机 <time datetime="'+esc(n.boot_time)+'">'+esc(boot)+'</time></span>'+(nodeUptime(n.uptime)?'<span>'+esc(nodeUptime(n.uptime))+'</span>':'')+'</div>':'';
 const infoButton='<button type="button" class="node-info-trigger" data-info-node="'+id+'" aria-label="查看 '+name+' 的机器基本信息" aria-controls="node-quick-info" aria-expanded="false" title="机器基本信息">!</button>';
 const favButton='<button type="button" class="node-favorite-trigger '+(n.favorite?'is-favorite':'')+'" data-fav-node="'+id+'" aria-label="'+(n.favorite?'取消收藏 ':'收藏 ')+name+'" title="'+(n.favorite?'取消收藏':'收藏服务器')+'" aria-pressed="'+Boolean(n.favorite)+'">'+(n.favorite?'★':'☆')+'</button>';
 const cardStart='<article role="button" tabindex="0" data-node="'+id+'" class="node dashboard-node glass-node '+(online?'':'node-offline')+'" aria-label="查看 '+name+' 详情">';
 if(currentLayout==='list'){
  const parts=[cardFieldOn('list','cpu')?'CPU '+Math.round(Number(n.cpu)||0)+'%':'',cardFieldOn('list','memory')?'内存 '+Math.round(Number(n.memory)||0)+'%':'',cardFieldOn('list','disk')?'存储 '+Math.round(Number(n.disk)||0)+'%':'',cardFieldOn('list','speed')?'↓ '+compactSpeed(n.rx_speed)+' ↑ '+compactSpeed(n.tx_speed):''].filter(Boolean);
  const listStat=online?parts.join(' · '):'上次在线：'+(n.last_seen?esc(new Date(n.last_seen).toLocaleString()):'未知');
  return cardStart+favButton+infoButton+header+(listStat?'<div class="node-list-summary">'+listStat+'</div>':'')+(cardFieldOn('list','traffic')&&has?'<div class="compact-month">本月 '+planAmountGB(used/1000000000)+'</div>':'')+alert+'</article>';
 }
 if(currentLayout==='compact'){
  const compactMonthly=quota>0&&has?
   '<div class="compact-month"><span>本月 <b>'+planAmountGB(used/1000000000)+'</b> / '+planAmountGB(period.quota_gb)+'</span><span>'+percentLabel(ratio)+'</span></div><div class="compact-quota"><i class="'+(ratio>=100?'danger':ratio>=80?'warn':'')+'" style="width:'+clamp.toFixed(1)+'%"></i></div>':
   has?'<div class="compact-month"><span>本月流量</span><b>'+planAmountGB(used/1000000000)+'</b></div>':'';
  return cardStart+favButton+infoButton+header+(chosen.length?metrics:'')+(cardFieldOn('compact','price')?priceRow:'')+(cardFieldOn('compact','traffic')?compactMonthly:'')+(cardFieldOn('compact','boot')?bootLine:'')+alert+'</article>';
 }
 return cardStart+favButton+infoButton+header+(cardFieldOn('cards','price')?priceRow:'')+(chosen.length?metrics:'')+(cardFieldOn('cards','cumulative')?cumulative:'')+(cardFieldOn('cards','tags')?labels:'')+(cardFieldOn('cards','traffic')?quotaBar:'')+(cardFieldOn('cards','boot')?bootLine:'')+alert+'</article>';
}
async function alertsLoad(){try{const r=await fetch('/api/v1/alerts');if(!r.ok)return;const d=await r.json(),a=d.settings;$('a-offline').value=a.offline_seconds;$('a-cpu').value=a.cpu_threshold;$('a-memory').value=a.memory_threshold;$('a-disk').value=a.disk_threshold;$('a-duration').value=a.duration_seconds;$('a-webhook').value=a.webhook||'';renderAlerts(d.events)}catch(e){$('alerts-status').textContent=e.message}}
function alertName(k){return ({quota_80:'月流量达到 80%',quota_90:'月流量达到 90%',quota_100:'月流量达到 100%',expiry_30:'30 天内到期',expiry_15:'15 天内到期',expiry_7:'7 天内到期',expiry_overdue:'服务器已到期'})[k]||k}
let lastAlertEvents=[],alertFilter='active';
document.querySelectorAll('[data-alert-mode]').forEach(button=>button.addEventListener('click',()=>{
 alertFilter=button.dataset.alertMode;renderAlerts(lastAlertEvents);
}));
function renderAlerts(items){
 lastAlertEvents=Array.isArray(items)?items:[];
 const active=lastAlertEvents.filter(e=>!e.end).length;
 const recovered=lastAlertEvents.length-active;
 $('alert-count').textContent=active;
 $('alert-summary').textContent=active+' 条告警中 · '+recovered+' 条已恢复';
 document.querySelectorAll('[data-alert-mode]').forEach(b=>{const selected=b.dataset.alertMode===alertFilter;b.classList.toggle('selected',selected);b.setAttribute('aria-pressed',String(selected))});
 const root=$('alert-events');root.replaceChildren();
 const visible=lastAlertEvents.filter(e=>alertFilter==='all'||(alertFilter==='active'?!e.end:!!e.end));
 if(!visible.length){const empty=document.createElement('div');empty.className='event-empty';empty.textContent=alertFilter==='active'?'目前没有未恢复的告警':alertFilter==='recovered'?'暂无已恢复记录':'暂无告警记录';root.append(empty);return}
 for(const e of visible.slice(0,60)){
  const item=document.createElement('div');item.className='event-item'+(e.end?' recovered':'');
  const dot=document.createElement('span');dot.className='event-dot';
  const content=document.createElement('div');content.className='event-content';
  const title=document.createElement('strong');title.textContent=alertName(e.kind);
  const sub=document.createElement('small');const n=cachedNodes.find(n=>n.name===e.node);
  sub.textContent=(n?(n.display_name||n.hostname||n.name):e.node)+' · '+new Date(e.start*1000).toLocaleString();
  content.append(title,sub);
  const status=document.createElement('span');status.className='event-status';status.textContent=e.end?'已恢复':'告警中';
  item.append(dot,content,status);root.append(item)
 }
}

$('alerts-form').addEventListener('submit',async e=>{e.preventDefault();const payload={offline_seconds:Number($('a-offline').value),cpu_threshold:Number($('a-cpu').value),memory_threshold:Number($('a-memory').value),disk_threshold:Number($('a-disk').value),duration_seconds:Number($('a-duration').value),webhook:$('a-webhook').value.trim()};try{const r=await fetch('/api/v1/alerts',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(payload)});if(!r.ok)throw Error(await r.text());$('alerts-status').textContent='已保存'}catch(err){$('alerts-status').textContent=err.message}});
let selectedNode=null;
function openNode(name){selectedNode=name;switchView('detail');switchDetailTab('performance');updateDetail();window.scrollTo(0,0)}
function detailRow(label,value){return '<div><span class="muted">'+esc(label)+'</span><br><strong>'+esc(value)+'</strong></div>'}
function updateDetail(){if(!selectedNode)return;const n=cachedNodes.find(x=>x.name===selectedNode);if(!n){$('detail-title').textContent='节点不存在';$('detail-status').textContent='此节点可能已删除';return}
$('detail-title').textContent=n.display_name||n.name;$('detail-subtitle').textContent=(n.group||'未分组')+' · '+n.name;
$('detail-status').textContent=n.online?'● 在线 · 监控数据实时更新':'● 离线 · 最后上报 '+(n.last_seen?new Date(n.last_seen).toLocaleString():'未知');
$('detail-metrics').innerHTML=[['CPU',n.cpu],['内存',n.memory],['磁盘',n.disk]].map(x=>'<div><small>'+x[0]+'</small><strong>'+(n.online?Math.round(x[1])+'%':'—')+'</strong></div>').join('');
$('detail-info').innerHTML=detailRow('服务器位置',n.location||'尚未定位')+detailRow('位置来源',n.location_source==='manual'?'手动覆盖':n.location_source==='ip'?'公网 IP 自动识别':'未识别')+detailRow('定位 IP',n.public_ip||'—')+detailRow('公网 IPv4（出站）',n.public_ipv4||'未检测到')+detailRow('公网 IPv6（出站）',n.public_ipv6||'未检测到')+detailRow('主机名',n.hostname||'—')+detailRow('系统发行版',n.os||'—')+detailRow('Agent 版本',n.agent_version||'未上报')+detailRow('CPU 逻辑核心',n.cpu_cores>0?n.cpu_cores+' C':'—')+detailRow('CPU 型号',n.cpu_model||'—')+detailRow('负载 1/5/15 分钟',[n.load1,n.load5,n.load15].every(x=>x!==undefined)?[n.load1,n.load5,n.load15].map(x=>Number(x).toFixed(2)).join(' / '):'—')+detailRow('Swap',n.swap_total>0?capacity(n.swap_used)+' / '+capacity(n.swap_total):'无或未上报')+detailRow('内存使用',n.memory_total>0?capacity(n.memory_used)+' / '+capacity(n.memory_total):'—')+detailRow('磁盘使用',n.disk_total>0?capacity(n.disk_used)+' / '+capacity(n.disk_total):'—')+detailRow('架构',n.arch||'—')+detailRow('运行时长',n.uptime?Math.floor(n.uptime/86400)+' 天':'—')+detailRow('最后上报',n.last_seen?new Date(n.last_seen).toLocaleString():'—')+(n.notes?detailRow('备注',n.notes):'');
const p=trafficPeriods[n.name];$('detail-network-info').innerHTML=detailRow('实时下载',n.online?formatBytes(n.rx_speed)+'/s':'—')+detailRow('实时上传',n.online?formatBytes(n.tx_speed)+'/s':'—')+detailRow('累计接收（网卡计数）',formatBytes(n.rx_bytes))+detailRow('累计发送（网卡计数）',formatBytes(n.tx_bytes))+detailRow('今日流量',p&&p.has_samples?formatBytes(p.today_rx+p.today_tx):'—')+detailRow('本月流量',p&&p.has_samples?formatBytes(p.month_rx+p.month_tx):'—')+detailRow('统计时区',p?p.timezone||'UTC':'UTC')+detailRow('月流量额度',p&&p.quota_gb>0?p.quota_gb+' GB':'未设置')+detailRow('磁盘读取',n.online?formatBytes(n.disk_read_speed)+'/s':'—')+detailRow('磁盘写入',n.online?formatBytes(n.disk_write_speed)+'/s':'—');
}
async function updateDetailHistory(){if(!selectedNode)return;const ticket=++detailHistorySequence,name=selectedNode,hours=$('detail-hours').value;try{const rs=await Promise.all([fetch('/api/v1/history?name='+encodeURIComponent(name)+'&hours='+hours),fetch('/api/v1/traffic?name='+encodeURIComponent(name)+'&hours='+hours)]);if(!rs.every(r=>r.ok))throw Error('数据接口异常');const points=await rs[0].json(),traffic=await rs[1].json();if(ticket!==detailHistorySequence||selectedNode!==name||$('detail-hours').value!==hours)return;detailChartPoints=points;detailChartCursor=-1;$('detail-point').textContent=points.length?'触摸图表查看 CPU、内存和磁盘':'当前时间段没有监控采样';const ctx=$('detail-chart').getContext('2d'),w=900,h=240;ctx.clearRect(0,0,w,h);const colors=['#60a5fa','#a78bfa','#2dd4bf'];ctx.strokeStyle='#34455e';ctx.lineWidth=1;for(let value=0;value<=100;value+=25){const y=18+(100-value)/100*196;ctx.beginPath();ctx.moveTo(42,y);ctx.lineTo(884,y);ctx.stroke();ctx.fillStyle='#9eb0c9';ctx.font='12px sans-serif';ctx.fillText(value+'%',6,y+4)}
if(points.length>0){const min=points[0].time,max=Math.max(min+1,points[points.length-1].time);['cpu','memory','disk'].forEach((key,k)=>{ctx.beginPath();ctx.strokeStyle=colors[k];ctx.lineWidth=2;points.forEach((p,i)=>{const x=42+(p.time-min)/(max-min)*842,y=18+(100-Math.max(0,Math.min(100,p[key])))/100*196;if(!i)ctx.moveTo(x,y);else ctx.lineTo(x,y)});ctx.stroke()})}
$('detail-chart-status').textContent='CPU（蓝） · 内存（紫） · 磁盘（青） · '+points.length+' 个区间';
$('detail-traffic').textContent='接收 ↓ '+formatBytes(traffic.rx)+'　发送 ↑ '+formatBytes(traffic.tx)+'（采样估算）'}catch(e){if(ticket===detailHistorySequence)$('detail-chart-status').textContent='历史数据加载失败：'+e.message}}
$('nodes').addEventListener('click',e=>{
 const fav=e.target.closest('[data-fav-node]');
 if(fav){e.stopPropagation();const id=fav.dataset.favNode,node=cachedNodes.find(n=>n.name===id);if(node){fav.disabled=true;saveNodeFavorite(id,!node.favorite).catch(err=>{fav.disabled=false;alert('收藏失败：'+err.message)})}return}
 const info=e.target.closest('.node-info-trigger');
 if(info){
  e.stopPropagation();const id=info.dataset.infoNode;
  if(quickInfoNode===id&&quickInfoPinned)hideNodeInfo();else showNodeInfo(id,true);
  return;
 }
 const card=e.target.closest('[data-node]');if(card){hideNodeInfo();openNode(card.dataset.node)}
});
$('nodes').addEventListener('keydown',e=>{
 if(e.target.closest('.node-info-trigger')||e.target.closest('[data-fav-node]')||e.target.closest('[data-group-header]'))return;
 if(e.key==='Enter'||e.key===' '){const card=e.target.closest('[data-node]');if(card){e.preventDefault();hideNodeInfo();openNode(card.dataset.node)}}
});
$('nodes').addEventListener('pointerover',e=>{
 const button=e.target.closest('.node-info-trigger');
 if(button&&e.pointerType==='mouse'&&!quickInfoPinned)showNodeInfo(button.dataset.infoNode,false);
});
$('nodes').addEventListener('pointerout',e=>{
 const button=e.target.closest('.node-info-trigger');
 if(button&&!button.contains(e.relatedTarget))scheduleHideNodeInfo();
});
$('node-quick-info').addEventListener('pointerenter',()=>clearTimeout(quickInfoTimer));
$('node-quick-info').addEventListener('pointerleave',scheduleHideNodeInfo);
$('node-info-backdrop').addEventListener('click',hideNodeInfo);
$('node-quick-info').addEventListener('click',e=>{if(e.target.closest('.node-info-close'))hideNodeInfo();e.stopPropagation()});
document.addEventListener('click',e=>{
 if(!e.target.closest('.node-info-trigger')&&!e.target.closest('#node-quick-info'))hideNodeInfo();
});
document.addEventListener('keydown',e=>{if(e.key==='Escape'&&quickInfoNode)hideNodeInfo()});
window.addEventListener('scroll',()=>{if(quickInfoNode)positionNodeInfo()},{passive:true,capture:true});
window.addEventListener('resize',()=>{if(quickInfoNode)positionNodeInfo()});

function switchDetailTab(tab){document.querySelectorAll('.detail-tab').forEach(b=>{const on=b.dataset.detailTab===tab;b.classList.toggle('selected',on);b.setAttribute('aria-selected',String(on));b.tabIndex=0});document.querySelectorAll('.detail-pane').forEach(p=>p.hidden=p.id!=='detail-'+tab)}
document.querySelectorAll('.detail-tab').forEach(b=>b.addEventListener('click',()=>switchDetailTab(b.dataset.detailTab)));
switchDetailTab('performance');
$('refresh-alerts').addEventListener('click',refreshAlertEvents);
$('detail-back').addEventListener('click',()=>switchView('overview'));
$('detail-hours').addEventListener('change',updateDetailHistory);
async function loadNodeLimits(){try{const r=await fetch('/api/v1/node-limits',{cache:'no-store'});if(!r.ok)throw Error('HTTP '+r.status);nodeLimitsCache=await r.json();fillLimitsForm();if($('metadata-name').value)fillPlanLimits($('metadata-name').value)}catch(e){$('limit-status').textContent='加载失败：'+e.message}}
function fillLimitsForm(){const name=$('limit-name').value;if(!name)return;const v=nodeLimitsCache[name]||{};$('limit-timezone').value=v.timezone||'UTC';$('limit-quota').value=v.quota_gb||0;$('limit-expires').value=v.expires_on||''}
$('limit-name').addEventListener('change',fillLimitsForm);
$('node-limits-form').addEventListener('submit',async e=>{e.preventDefault();const payload={name:$('limit-name').value,timezone:$('limit-timezone').value.trim(),quota_gb:Number($('limit-quota').value),expires_on:$('limit-expires').value};try{const r=await fetch('/api/v1/node-limits',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(payload)});if(!r.ok)throw Error(await r.text());$('limit-status').textContent='已保存';await loadNodeLimits();await refreshPeriods()}catch(err){$('limit-status').textContent='保存失败：'+err.message}});
function confirmPendingEdits(){const rows=$('node-editor-list');if(rows.dataset.dirty==='true'){if(!confirm('节点资料尚未保存，确定放弃修改吗？'))return false;rows.dataset.dirty='false'}return true}
window.addEventListener('beforeunload',e=>{if($('node-editor-list').dataset.dirty==='true'){e.preventDefault();e.returnValue=''}});
let currentView='overview',overviewScrollY=0;
function switchView(view){if(view!=='settings'&&!$('view-settings').hidden&&!confirmPendingEdits())return;if(currentView==='overview'&&view!=='overview')overviewScrollY=window.scrollY;if(view!=='overview'&&quickInfoNode)hideNodeInfo();for(const el of document.querySelectorAll('.view'))el.hidden=el.id!=='view-'+view;for(const b of document.querySelectorAll('.nav-btn')){const active=b.dataset.view===view;b.classList.toggle('selected',active);b.setAttribute('aria-current',active?'page':'false')}if(view==='statistics')drawHistory();if(view==='alerts')refreshAlertEvents();if(view==='detail')updateDetailHistory();if(view==='overview'&&currentView!=='overview')requestAnimationFrame(()=>window.scrollTo({top:overviewScrollY,behavior:'instant'}));currentView=view}
document.querySelectorAll('.nav-btn').forEach(b=>b.addEventListener('click',()=>switchView(b.dataset.view)));
async function refreshAlertEvents(){const ticket=++alertFetchSequence;try{const r=await fetch('/api/v1/alerts');if(!r.ok)throw Error('HTTP '+r.status);const data=await r.json();if(ticket!==alertFetchSequence)return;renderAlerts(data.events)}catch(e){if(ticket===alertFetchSequence&&!$('view-alerts').hidden)$('alert-summary').textContent='刷新失败 · 请稍后重试'}}
switchView('overview');async function refreshPeriods(){if(trafficFetchBusy)return;trafficFetchBusy=true;const ticket=++trafficFetchSequence;try{const r=await fetch('/api/v1/traffic-summary',{cache:'no-store'});if(!r.ok)throw Error('HTTP '+r.status);const data=await r.json();if(ticket!==trafficFetchSequence)return;trafficPeriods=data;renderOverviewNodes();if(selectedNode)updateDetail()}catch(e){if(ticket===trafficFetchSequence&&$('updated').textContent.includes('加载中'))$('updated').textContent='流量数据暂不可用'}finally{trafficFetchBusy=false}}
refresh();alertsLoad();loadNodeLimits();refreshPeriods();setInterval(refresh,5000);setInterval(refreshPeriods,15000);setInterval(drawHistory,30000);setInterval(refreshAlertEvents,30000);`
