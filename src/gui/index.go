package gui

var indexHTML = `<!DOCTYPE html>
<html lang="id">
<head>
<meta charset="UTF-8">
<title>Ice Sliding Puzzle Solver</title>
<style>
*{box-sizing:border-box;margin:0;padding:0}
body{background:#0f172a;color:#e2e8f0;font-family:'Segoe UI',system-ui,sans-serif;height:100vh;display:flex;flex-direction:column;overflow:hidden}
header{background:#1e293b;padding:12px 20px;border-bottom:1px solid #334155;display:flex;align-items:center;gap:12px;flex-shrink:0}
header h1{font-size:18px;font-weight:700;color:#38bdf8}
header span{color:#64748b;font-size:13px}
.main{display:flex;flex:1;overflow:hidden}
.left{width:300px;background:#1e293b;border-right:1px solid #334155;overflow-y:auto;flex-shrink:0;padding:16px;display:flex;flex-direction:column;gap:12px}
.right{flex:1;display:flex;flex-direction:column;overflow:hidden}
.board-area{flex:1;overflow:auto;padding:16px;display:flex;align-items:flex-start;justify-content:flex-start}
.controls{background:#1e293b;border-top:1px solid #334155;padding:12px 16px;flex-shrink:0}

label{font-size:12px;font-weight:600;color:#94a3b8;text-transform:uppercase;letter-spacing:.05em;display:block;margin-bottom:4px}
select,input[type=file]{width:100%;background:#0f172a;border:1px solid #334155;color:#e2e8f0;padding:7px 10px;border-radius:6px;font-size:13px}
select:focus{outline:none;border-color:#38bdf8}
input[type=range]{width:100%;accent-color:#38bdf8}

.btn{width:100%;padding:9px;border:none;border-radius:6px;font-size:14px;font-weight:600;cursor:pointer;transition:opacity .15s}
.btn-primary{background:#0284c7;color:#fff}
.btn-primary:hover{opacity:.85}
.btn-primary:disabled{background:#334155;color:#64748b;cursor:default;opacity:1}
.btn-secondary{background:#1e3a5f;color:#38bdf8;border:1px solid #38bdf8}
.btn-secondary:hover{opacity:.85}
.btn-secondary:disabled{background:#1e293b;color:#64748b;border-color:#334155;cursor:default}

.section-title{font-size:12px;font-weight:700;color:#38bdf8;text-transform:uppercase;letter-spacing:.08em;padding-bottom:6px;border-bottom:1px solid #334155}

.stats-grid{display:grid;grid-template-columns:auto 1fr;gap:4px 10px;font-size:12px}
.stats-grid .k{color:#64748b}
.stats-grid .v{color:#e2e8f0;font-weight:500;word-break:break-all}

.legend{display:grid;grid-template-columns:1fr 1fr;gap:4px;font-size:11px}
.leg-item{display:flex;align-items:center;gap:5px}
.swatch{width:14px;height:14px;border-radius:2px;flex-shrink:0}

.status{padding:6px 10px;border-radius:6px;font-size:12px;font-weight:600;text-align:center}
.status-idle{background:#1e293b;color:#64748b;border:1px solid #334155}
.status-running{background:#1e3a5f;color:#38bdf8;border:1px solid #0284c7}
.status-ok{background:#14532d;color:#4ade80;border:1px solid #16a34a}
.status-fail{background:#450a0a;color:#f87171;border:1px solid #dc2626}

#board{display:grid;gap:1px;background:#0f172a}
.cell{width:46px;height:46px;display:flex;align-items:center;justify-content:center;font-weight:700;font-size:16px;border-radius:3px;transition:background .2s}
.cell-ice{background:#bfdbfe}
.cell-rock{background:#374151;color:#9ca3af}
.cell-lava{background:#dc2626;color:#fff}
.cell-exit{background:#16a34a;color:#fff}
.cell-number{background:#fbbf24;color:#1c1917}
.cell-collected{background:#6b7280;color:#e5e7eb}
.cell-player{background:#2563eb;color:#fff;box-shadow:0 0 0 2px #93c5fd}
.cell-path{background:#93c5fd}

.pb-row{display:flex;align-items:center;gap:10px;margin-bottom:8px}
.pb-btn{background:#1e293b;border:1px solid #334155;color:#e2e8f0;width:34px;height:34px;border-radius:6px;cursor:pointer;display:flex;align-items:center;justify-content:center;font-size:14px;transition:background .15s}
.pb-btn:hover{background:#334155}
.pb-btn:disabled{opacity:.4;cursor:default}
#step-info{font-size:12px;color:#94a3b8;flex:1;text-align:center}
input[type=range].step-range{flex:1;height:6px}
.speed-wrap{display:flex;align-items:center;gap:8px;font-size:12px;color:#94a3b8}
#speed-val{color:#38bdf8;min-width:36px;font-weight:600}

.file-drop{border:2px dashed #334155;border-radius:8px;padding:14px;text-align:center;font-size:12px;color:#64748b;cursor:pointer;transition:border-color .2s}
.file-drop:hover,.file-drop.drag-over{border-color:#38bdf8;color:#38bdf8}
</style>
</head>
<body>
<header>
  <h1>Ice Sliding Puzzle Solver</h1>
  <span>IF2211 Strategi Algoritma</span>
</header>
<div class="main">
  <div class="left">
    <div class="section-title">Input</div>
    <div>
      <label>File Puzzle (.txt)</label>
      <div class="file-drop" id="drop-zone" onclick="document.getElementById('file-input').click()">
        Klik atau drag &amp; drop file .txt
        <br><span id="file-name" style="color:#38bdf8;font-weight:600"></span>
      </div>
      <input type="file" id="file-input" accept=".txt" style="display:none" onchange="onFileChange(this)">
    </div>
    <div>
      <label>Algoritma</label>
      <select id="alg-select" onchange="onAlgChange()">
        <option value="astar">A*</option>
        <option value="ucs">UCS</option>
        <option value="gbfs">GBFS</option>
        <option value="idastar">IDA*</option>
      </select>
    </div>
    <div>
      <label>Heuristik</label>
      <select id="heur-select">
        <option value="manhattan">Manhattan</option>
        <option value="remaining">Remaining Digits</option>
        <option value="euclidean">Euclidean</option>
        <option value="maxcombined">Max Combined</option>
      </select>
    </div>
    <button class="btn btn-primary" id="run-btn" onclick="runSolve()" disabled>Cari Solusi</button>
    <div class="status status-idle" id="status-box">Pilih file input</div>

    <div class="section-title">Hasil</div>
    <div class="stats-grid">
      <span class="k">Algoritma</span><span class="v" id="s-alg">-</span>
      <span class="k">Heuristik</span><span class="v" id="s-heur">-</span>
      <span class="k">Cost</span><span class="v" id="s-cost">-</span>
      <span class="k">Steps</span><span class="v" id="s-steps">-</span>
      <span class="k">Iterasi</span><span class="v" id="s-iter">-</span>
      <span class="k">Nodes</span><span class="v" id="s-nodes">-</span>
      <span class="k">Durasi (ms)</span><span class="v" id="s-dur">-</span>
      <span class="k">Path</span><span class="v" id="s-path" style="font-family:monospace">-</span>
    </div>

    <button class="btn btn-secondary" id="save-btn" onclick="saveSolution()" disabled>Simpan Solusi (.txt)</button>

    <div class="section-title">Legenda</div>
    <div class="legend">
      <div class="leg-item"><div class="swatch" style="background:#2563eb"></div><span>Z = Pemain</span></div>
      <div class="leg-item"><div class="swatch" style="background:#16a34a"></div><span>O = Tujuan</span></div>
      <div class="leg-item"><div class="swatch" style="background:#374151"></div><span>X = Batu</span></div>
      <div class="leg-item"><div class="swatch" style="background:#dc2626"></div><span>L = Lava</span></div>
      <div class="leg-item"><div class="swatch" style="background:#fbbf24"></div><span>0-9 = Digit</span></div>
      <div class="leg-item"><div class="swatch" style="background:#6b7280"></div><span>* = Terlewati</span></div>
      <div class="leg-item"><div class="swatch" style="background:#bfdbfe"></div><span>Es</span></div>
      <div class="leg-item"><div class="swatch" style="background:#93c5fd"></div><span>Jejak</span></div>
    </div>
  </div>

  <div class="right">
    <div class="board-area">
      <div id="board"></div>
    </div>
    <div class="controls">
      <div class="pb-row">
        <button class="pb-btn" id="pb-first" onclick="pbFirst()" title="Awal" disabled>&#x23EE;</button>
        <button class="pb-btn" id="pb-prev"  onclick="pbPrev()"  title="Mundur" disabled>&#x25C0;</button>
        <button class="pb-btn" id="pb-play"  onclick="pbPlay()"  title="Play/Pause" disabled>&#x25B6;</button>
        <button class="pb-btn" id="pb-next"  onclick="pbNext()"  title="Maju" disabled>&#x25B6;&#x7C;</button>
        <button class="pb-btn" id="pb-last"  onclick="pbLast()"  title="Akhir" disabled>&#x23ED;</button>
        <div class="speed-wrap">
          <span>Kecepatan:</span>
          <input type="range" id="speed-range" min="0.25" max="4" step="0.25" value="1" oninput="onSpeedChange(this.value)" style="width:100px">
          <span id="speed-val">1.00x</span>
        </div>
      </div>
      <div class="pb-row">
        <input type="range" class="step-range" id="step-range" min="0" max="0" value="0" oninput="pbJump(parseInt(this.value))">
      </div>
      <div id="step-info">Step: - / -</div>
    </div>
  </div>
</div>

<script>
let state = {
  file: null,
  data: null,
  curFrame: 0,
  playing: false,
  timer: null,
  speed: 1.0,
  visitedPos: []
};

function onFileChange(input) {
  if (input.files.length === 0) return;
  state.file = input.files[0];
  document.getElementById('file-name').textContent = state.file.name;
  document.getElementById('run-btn').disabled = false;
  setStatus('idle', 'File dipilih: ' + state.file.name);
}

const dz = document.getElementById('drop-zone');
dz.addEventListener('dragover', e => { e.preventDefault(); dz.classList.add('drag-over'); });
dz.addEventListener('dragleave', () => dz.classList.remove('drag-over'));
dz.addEventListener('drop', e => {
  e.preventDefault(); dz.classList.remove('drag-over');
  const f = e.dataTransfer.files[0];
  if (!f) return;
  state.file = f;
  document.getElementById('file-name').textContent = f.name;
  document.getElementById('run-btn').disabled = false;
  setStatus('idle', 'File dipilih: ' + f.name);
});

function onAlgChange() {
  const alg = document.getElementById('alg-select').value;
  document.getElementById('heur-select').disabled = alg === 'ucs';
}

function setStatus(type, msg) {
  const el = document.getElementById('status-box');
  el.className = 'status status-' + type;
  el.textContent = msg;
}

async function runSolve() {
  if (!state.file) return;
  stopPlay();
  setStatus('running', 'Memproses...');
  document.getElementById('run-btn').disabled = true;
  document.getElementById('save-btn').disabled = true;

  const fd = new FormData();
  fd.append('file', state.file);
  fd.append('alg', document.getElementById('alg-select').value);
  fd.append('heur', document.getElementById('heur-select').value);

  try {
    const resp = await fetch('/api/solve', { method: 'POST', body: fd });
    const data = await resp.json();
    if (data.error) { setStatus('fail', 'Error: ' + data.error); return; }
    state.data = data;
    state.curFrame = 0;
    state.visitedPos = [];
    updateStats(data);
    if (!data.found) {
      setStatus('fail', 'Tidak ada solusi ditemukan.');
      renderFrame(0);
    } else {
      setStatus('ok', 'Solusi ditemukan! ' + data.steps + ' langkah, cost: ' + data.cost);
      initPlayback(data);
      renderFrame(0);
      document.getElementById('save-btn').disabled = false;
    }
  } catch(e) {
    setStatus('fail', 'Network error: ' + e.message);
  } finally {
    document.getElementById('run-btn').disabled = false;
  }
}

function updateStats(d) {
  document.getElementById('s-alg').textContent = d.algorithm;
  document.getElementById('s-heur').textContent = d.heuristic;
  document.getElementById('s-cost').textContent = d.found ? d.cost : '-';
  document.getElementById('s-steps').textContent = d.found ? d.steps : '-';
  document.getElementById('s-iter').textContent = d.iterations;
  document.getElementById('s-nodes').textContent = d.nodesGenerated;
  document.getElementById('s-dur').textContent = d.durationMs.toFixed(3);
  document.getElementById('s-path').textContent = d.found ? d.path : '-';
}

function initPlayback(data) {
  const n = data.frames.length - 1;
  const sr = document.getElementById('step-range');
  sr.max = n; sr.value = 0;
  setPbBtns(true);
}

function setPbBtns(enabled) {
  ['pb-first','pb-prev','pb-play','pb-next','pb-last'].forEach(id => {
    document.getElementById(id).disabled = !enabled;
  });
}

function renderFrame(idx) {
  const d = state.data;
  if (!d) return;
  if (idx < 0) idx = 0;
  if (idx >= d.frames.length) idx = d.frames.length - 1;
  state.curFrame = idx;

  const fr = d.frames[idx];
  document.getElementById('step-range').value = idx;

  const action = fr.action === 'Initial' ? 'Initial' : fr.action;
  document.getElementById('step-info').textContent =
    'Step ' + idx + ' / ' + (d.frames.length - 1) +
    '  |  Aksi: ' + action +
    (idx > 0 ? '  |  Step Cost: ' + fr.stepCost + '  |  Total Cost: ' + fr.totalCost : '');

  state.visitedPos = [];
  for (let i = 0; i <= idx; i++) {
    state.visitedPos.push([d.frames[i].px, d.frames[i].py]);
  }

  const grid = document.getElementById('board');
  grid.style.gridTemplateColumns = 'repeat(' + d.width + ', 46px)';
  grid.innerHTML = '';

  for (let y = 0; y < d.height; y++) {
    for (let x = 0; x < d.width; x++) {
      const cell = d.cells[y][x];
      const div = document.createElement('div');
      div.className = 'cell';

      let cls, txt = '';
      if (x === fr.px && y === fr.py) {
        cls = 'cell-player'; txt = 'Z';
      } else {
        switch (cell.t) {
          case 'rock':  cls = 'cell-rock';  txt = 'X'; break;
          case 'lava':  cls = 'cell-lava';  txt = 'L'; break;
          case 'exit':  cls = 'cell-exit';  txt = 'O'; break;
          case 'number':
            const col = fr.collected && cell.d < fr.collected.length && fr.collected[cell.d];
            cls = col ? 'cell-collected' : 'cell-number';
            txt = col ? '*' : String(cell.d);
            break;
          default:
            const isVisited = state.visitedPos.some(p => p[0] === x && p[1] === y);
            cls = isVisited ? 'cell-path' : 'cell-ice';
        }
      }
      div.classList.add(cls);
      div.textContent = txt;
      div.title = cell.t + (cell.t === 'number' ? ' ' + cell.d : '') + ' (cost:' + cell.c + ')';
      grid.appendChild(div);
    }
  }
}

function pbFirst() { stopPlay(); renderFrame(0); }
function pbPrev()  { stopPlay(); renderFrame(state.curFrame - 1); }
function pbNext()  { stopPlay(); renderFrame(state.curFrame + 1); }
function pbLast()  { stopPlay(); renderFrame(state.data.frames.length - 1); }
function pbJump(v) { stopPlay(); renderFrame(v); }

function pbPlay() {
  if (state.playing) { stopPlay(); return; }
  if (!state.data || !state.data.found) return;
  if (state.curFrame >= state.data.frames.length - 1) renderFrame(0);
  state.playing = true;
  document.getElementById('pb-play').textContent = '⏸';
  const ms = 1000 / state.speed;
  state.timer = setInterval(() => {
    if (state.curFrame < state.data.frames.length - 1) {
      renderFrame(state.curFrame + 1);
    } else {
      stopPlay();
    }
  }, ms);
}

function stopPlay() {
  if (!state.playing) return;
  clearInterval(state.timer);
  state.playing = false;
  document.getElementById('pb-play').innerHTML = '&#x25B6;';
}

function onSpeedChange(v) {
  state.speed = parseFloat(v);
  document.getElementById('speed-val').textContent = parseFloat(v).toFixed(2) + 'x';
  if (state.playing) { stopPlay(); pbPlay(); }
}

function saveSolution() {
  const d = state.data;
  if (!d || !d.found) return;
  let txt = '';
  txt += 'Algoritma       : ' + d.algorithm + '\n';
  txt += 'Heuristik       : ' + d.heuristic + '\n';
  txt += 'Path            : ' + d.path + '\n';
  txt += 'Cost            : ' + d.cost + '\n';
  txt += 'Steps           : ' + d.steps + '\n';
  txt += 'Iterations      : ' + d.iterations + '\n';
  txt += 'Nodes generated : ' + d.nodesGenerated + '\n';
  txt += 'Duration        : ' + d.durationMs.toFixed(3) + ' ms\n';
  txt += '\n== Step-by-step ==\n';
  d.frames.forEach(fr => {
    txt += '\n--- Step ' + fr.index + ' / ' + (d.frames.length-1) + '  (' + fr.action + ') ---\n';
    for (let y = 0; y < d.height; y++) {
      for (let x = 0; x < d.width; x++) {
        if (x === fr.px && y === fr.py) { txt += 'Z'; continue; }
        const c = d.cells[y][x];
        if (c.t === 'rock')  { txt += 'X'; continue; }
        if (c.t === 'lava')  { txt += 'L'; continue; }
        if (c.t === 'exit')  { txt += 'O'; continue; }
        if (c.t === 'number') {
          const col = fr.collected && c.d < fr.collected.length && fr.collected[c.d];
          txt += col ? '*' : String(c.d); continue;
        }
        txt += '*';
      }
      txt += '\n';
    }
    if (fr.index > 0) txt += 'Move: ' + fr.action + '  step cost: ' + fr.stepCost + '  total cost: ' + fr.totalCost + '\n';
  });

  const blob = new Blob([txt], {type: 'text/plain'});
  const a = document.createElement('a');
  a.href = URL.createObjectURL(blob);
  a.download = 'solusi.txt';
  a.click();
}
</script>
</body>
</html>`
