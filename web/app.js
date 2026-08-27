// Greenshields web UI. Every number shown is returned by the Go backend via
// the /api/* JSON endpoints; nothing is hard-coded or pre-computed in JS.

function num(id) {
  const v = parseFloat(document.getElementById(id).value);
  return isNaN(v) ? NaN : v;
}

function showError(el, msg) {
  el.innerHTML = '<div class="err">后端错误：' + msg + '</div>';
}

async function postJSON(url, body) {
  const resp = await fetch(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  });
  const data = await resp.json();
  if (!resp.ok) {
    throw new Error(data.error || ('HTTP ' + resp.status));
  }
  return data;
}

document.getElementById('load').addEventListener('click', async () => {
  try {
    const data = await (await fetch('/api/example')).json();
    document.getElementById('vf').value = data.vf;
    document.getElementById('kj').value = data.kj;
    document.getElementById('curveOut').innerHTML =
      '<span class="muted">已加载示例「' + (data.name || 'unnamed') + '」</span>';
  } catch (e) {
    showError(document.getElementById('curveOut'), e.message);
  }
});

document.getElementById('compute').addEventListener('click', async () => {
  const out = document.getElementById('curveOut');
  const svgOut = document.getElementById('svgOut');
  try {
    const vf = num('vf'), kj = num('kj');
    const data = await postJSON('/api/curve', { vf, kj });
    out.innerHTML =
      '<b>通行能力点</b>：km = ' + data.km.toFixed(3) +
      '，qm = ' + data.qm.toFixed(3) +
      '，vm = ' + data.vm.toFixed(3) +
      '<br><span class="muted">曲线点列（来自 /api/curve，共 ' + data.points.length + ' 点）</span>';
    // Render the SVG from the backend point list (no half-circle hard-coded).
    svgOut.innerHTML = renderCurve(data);
  } catch (e) {
    svgOut.innerHTML = '';
    showError(out, e.message);
  }
});

document.getElementById('qk').addEventListener('click', async () => {
  const out = document.getElementById('qkOut');
  try {
    const data = await postJSON('/api/qk', { vf: num('vf'), kj: num('kj'), k: num('k') });
    out.innerHTML =
      'k = ' + data.k.toFixed(3) + ' → v = ' + data.v.toFixed(3) +
      '，q = ' + data.q.toFixed(3) +
      '，侧：<b>' + data.side + '</b>' + (data.congested ? '（拥堵侧）' : '');
  } catch (e) {
    showError(out, e.message);
  }
});

document.getElementById('wave').addEventListener('click', async () => {
  const out = document.getElementById('waveOut');
  try {
    const data = await postJSON('/api/wave', {
      vf: num('vf'), kj: num('kj'), k1: num('k1'), k2: num('k2'),
    });
    out.innerHTML =
      'w = ' + data.w.toFixed(4) + '，传播方向：<b>' + data.direction + '</b>' +
      '（q₁=' + data.q1.toFixed(3) + ', q₂=' + data.q2.toFixed(3) + '）';
  } catch (e) {
    showError(out, e.message);
  }
});

// Build an SVG of the q(k) parabola purely from backend-supplied points.
function renderCurve(data) {
  const W = 520, H = 360, pad = 44;
  const qm = data.qm || 1;
  const kj = data.kj || 1;
  const plotW = W - 2 * pad, plotH = H - 2 * pad;
  const x = (k) => pad + (k / kj) * plotW;
  const y = (q) => (H - pad) - (q / qm) * plotH;
  let pts = data.points.map(p => x(p.k).toFixed(1) + ',' + y(p.q).toFixed(1)).join(' ');
  const cx = x(data.km), cy = y(data.qm);
  return '<svg xmlns="http://www.w3.org/2000/svg" width="' + W + '" height="' + H +
    '" viewBox="0 0 ' + W + ' ' + H + '">' +
    '<rect width="100%" height="100%" fill="#fff"/>' +
    '<line x1="' + pad + '" y1="' + (H - pad) + '" x2="' + (W - pad) + '" y2="' + (H - pad) + '" stroke="#444"/>' +
    '<line x1="' + pad + '" y1="' + (H - pad) + '" x2="' + pad + '" y2="' + pad + '" stroke="#444"/>' +
    '<polyline fill="none" stroke="#1f77b4" stroke-width="2" points="' + pts + '"/>' +
    '<circle cx="' + cx.toFixed(1) + '" cy="' + cy.toFixed(1) + '" r="4" fill="#d62728"/>' +
    '<text x="' + (cx + 6).toFixed(1) + '" y="' + (cy - 6).toFixed(1) + '" font-size="11" fill="#d62728">qm=' + data.qm.toFixed(2) + '</text>' +
    '<text x="' + (W - pad - 14) + '" y="' + (H - pad + 18) + '" font-size="12" fill="#444">k</text>' +
    '<text x="' + (pad - 30) + '" y="' + (pad + 4) + '" font-size="12" fill="#444">q</text>' +
    '</svg>';
}
