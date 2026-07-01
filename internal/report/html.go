package report

// since im not so good in working with graphics...
// THIS FILE IN MODULE IS CREATED BY AI (CLAUDE)

import (
	"afv/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type htmlNode struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Sub   string `json:"sub,omitempty"`
	Type  string `json:"type"`
}

type htmlEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"`
}

type htmlGraph struct {
	Nodes []htmlNode `json:"nodes"`
	Edges []htmlEdge `json:"edges"`
}

type htmlSummary struct {
	Functions int `json:"functions"`
	Endpoints int `json:"endpoints"`
	Calls     int `json:"calls"`
}

type htmlPayload struct {
	Target  string      `json:"target"`
	Program htmlGraph   `json:"program"`
	Auth    htmlGraph   `json:"auth"`
	Summary htmlSummary `json:"summary"`
}

func buildHTMLGraph(g model.Graph) htmlGraph {
	nodes := []htmlNode{}
	edges := []htmlEdge{}

	for _, n := range g.Nodes {
		switch n.Type {
		case "function":
			parts := strings.SplitN(n.ID, ":", 3)
			label := n.ID
			sub := ""
			if len(parts) == 3 {
				label = parts[2] + "()"
				sub = parts[1]
			}
			nodes = append(nodes, htmlNode{ID: n.ID, Label: label, Sub: sub, Type: "function"})
		case "endpoint":
			parts := strings.SplitN(n.ID, ":", 3)
			label := n.ID
			if len(parts) == 3 {
				label = parts[1] + " " + parts[2]
			}
			nodes = append(nodes, htmlNode{ID: n.ID, Label: label, Type: "endpoint"})
		case "external":
			parts := strings.SplitN(n.ID, ":", 2)
			label := n.ID
			if len(parts) == 2 {
				label = parts[1] + "()"
			}
			nodes = append(nodes, htmlNode{ID: n.ID, Label: label, Type: "external"})
		default:
			nodes = append(nodes, htmlNode{ID: n.ID, Label: n.ID, Type: n.Type})
		}
	}

	for _, e := range g.Edges {
		edges = append(edges, htmlEdge{From: e.From, To: e.To, Type: e.Type})
	}

	return htmlGraph{Nodes: nodes, Edges: edges}
}

func htmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return r.Replace(s)
}

func GenerateHTMLReport(target string, graph model.Graph, authGraph model.Graph) error {
	outDir := filepath.Join("output", target)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}

	functions, endpoints := 0, 0
	for _, n := range graph.Nodes {
		switch n.Type {
		case "function":
			functions++
		case "endpoint":
			endpoints++
		}
	}

	payload := htmlPayload{
		Target:  target,
		Program: buildHTMLGraph(graph),
		Auth:    buildHTMLGraph(authGraph),
		Summary: htmlSummary{
			Functions: functions,
			Endpoints: endpoints,
			Calls:     len(graph.Edges),
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	safeJSON := strings.ReplaceAll(string(payloadJSON), "</", `<\/`)

	page := strings.Replace(htmlPageTemplate, "__AFV_DATA__", safeJSON, 1)
	page = strings.Replace(page, "__AFV_TARGET__", htmlEscape(target), -1)

	reportPath := filepath.Join(outDir, "report.html")
	if err := os.WriteFile(reportPath, []byte(page), 0644); err != nil {
		return err
	}

	fmt.Println("[*]HTML report generated:", reportPath)
	return nil
}

const htmlPageTemplate = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Auth Flow Visualizer &mdash; __AFV_TARGET__</title>
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700;800&display=swap" rel="stylesheet">
<style>
:root {
  --bg: #090c12;
  --bg-raised: #0e131c;
  --panel: #10141d;
  --border: #1d2430;
  --text: #e7ecf3;
  --text-dim: #6f7A8f;
  --function: #4fd1c5;
  --endpoint: #f5a623;
  --external: #7b869c;
  --auth: #f0455c;
}

* { box-sizing: border-box; }
html { scroll-behavior: smooth; }

body {
  margin: 0;
  background: var(--bg);
  color: var(--text);
  font-family: 'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, monospace;
  font-size: 15px;
  line-height: 1.5;
}

body::before {
  content: "";
  position: fixed;
  inset: 0;
  background-image:
    radial-gradient(circle at 15% 20%, rgba(79,209,197,0.06), transparent 40%),
    radial-gradient(circle at 85% 75%, rgba(240,69,92,0.05), transparent 45%);
  pointer-events: none;
  z-index: 0;
}

main.afv-scroller {
  height: 100vh;
  overflow-y: auto;
  scroll-snap-type: y proximity;
}

.afv-section {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 8vh 7vw;
  scroll-snap-align: start;
  position: relative;
  opacity: 0;
  transform: translateY(14px);
  transition: opacity 0.7s ease, transform 0.7s ease;
}
.afv-section.afv-in-view { opacity: 1; transform: translateY(0); }

.afv-eyebrow {
  font-size: 12px;
  letter-spacing: 0.18em;
  text-transform: uppercase;
  color: var(--function);
  margin: 0 0 10px;
  display: flex;
  align-items: center;
  gap: 10px;
}
.afv-eyebrow::before { content: ""; width: 22px; height: 1px; background: var(--function); display: inline-block; }
.afv-section-auth .afv-eyebrow { color: var(--auth); }
.afv-section-auth .afv-eyebrow::before { background: var(--auth); }

.afv-section h2 {
  font-size: clamp(22px, 3vw, 32px);
  font-weight: 700;
  margin: 0 0 6px;
  letter-spacing: -0.01em;
}

.afv-section p.afv-desc {
  color: var(--text-dim);
  margin: 0 0 28px;
  max-width: 60ch;
  font-size: 13.5px;
}

/* ---- Hero ---- */
.afv-hero { align-items: flex-start; }
.afv-hero pre.afv-banner {
  font-size: clamp(8px, 1.35vw, 13px);
  line-height: 1.15;
  color: var(--function);
  margin: 0 0 20px;
  text-shadow: 0 0 24px rgba(79,209,197,0.25);
  white-space: pre;
}
.afv-hero h1 {
  font-size: clamp(28px, 5vw, 46px);
  margin: 0 0 10px;
  letter-spacing: -0.02em;
}
.afv-hero .afv-sub {
  color: var(--text-dim);
  margin: 0 0 30px;
  font-size: 14px;
}
.afv-hero .afv-sub code {
  color: var(--endpoint);
  background: rgba(245,166,35,0.08);
  padding: 2px 7px;
  border-radius: 4px;
  border: 1px solid rgba(245,166,35,0.2);
}
.afv-chip-row { display: flex; gap: 12px; flex-wrap: wrap; margin-bottom: 40px; }
.afv-chip {
  border: 1px solid var(--border);
  background: var(--panel);
  border-radius: 8px;
  padding: 12px 18px;
  min-width: 108px;
}
.afv-chip .afv-chip-num { font-size: 22px; font-weight: 700; color: var(--function); }
.afv-chip .afv-chip-label { font-size: 10.5px; color: var(--text-dim); text-transform: uppercase; letter-spacing: 0.1em; margin-top: 2px; }

.afv-scrolldown {
  color: var(--text-dim);
  font-size: 11px;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  display: flex;
  align-items: center;
  gap: 8px;
  animation: afv-bob 2s ease-in-out infinite;
}
@keyframes afv-bob { 0%,100% { transform: translateY(0); } 50% { transform: translateY(5px); } }

/* ---- Graph panel ---- */
.afv-panel {
  border: 1px solid var(--border);
  background: linear-gradient(180deg, var(--bg-raised), var(--panel));
  border-radius: 12px;
  padding: 4px;
  max-height: 62vh;
  overflow: auto;
  position: relative;
}
.afv-graph-svg { display: block; }
.afv-empty {
  padding: 60px 30px;
  text-align: center;
  color: var(--text-dim);
  font-size: 13px;
  max-width: 46ch;
  margin: 0 auto;
}

.afv-legend { display: flex; gap: 22px; margin-top: 16px; flex-wrap: wrap; }
.afv-legend-item { display: flex; align-items: center; gap: 8px; font-size: 11.5px; color: var(--text-dim); }
.afv-legend-swatch { width: 10px; height: 10px; border-radius: 50%; display: inline-block; }
.afv-legend-swatch.sw-function { background: var(--function); }
.afv-legend-swatch.sw-endpoint { background: var(--endpoint); border-radius: 2px; transform: rotate(45deg); }
.afv-legend-swatch.sw-external { background: transparent; border: 1.5px dashed var(--external); }
.afv-legend-swatch.sw-auth { background: var(--auth); }

/* node/edge styling (svg) */
.afv-shape { stroke-width: 1.6; transition: opacity 0.2s ease; cursor: default; }
.afv-shape-function { fill: rgba(79,209,197,0.10); stroke: var(--function); }
.afv-shape-endpoint { fill: rgba(245,166,35,0.12); stroke: var(--endpoint); }
.afv-shape-external { fill: rgba(123,134,156,0.08); stroke: var(--external); stroke-dasharray: 3 3; }

.afv-label { fill: var(--text); font-size: 11.5px; font-family: inherit; }
.afv-sublabel { fill: var(--text-dim); font-size: 9.5px; font-family: inherit; }

.afv-edge { fill: none; stroke: var(--function); stroke-opacity: 0.38; stroke-width: 1.4; transition: stroke-opacity 0.2s ease, stroke-width 0.2s ease; }
.afv-edge-endpoint { stroke: var(--endpoint); stroke-opacity: 0.5; }
.afv-edge-auth_flow { stroke: var(--auth); stroke-opacity: 0.55; }
marker[id$="-arrow-call"] path { fill: var(--function); }
marker[id$="-arrow-endpoint"] path { fill: var(--endpoint); }
marker[id$="-arrow-auth_flow"] path { fill: var(--auth); }

.afv-pulse { fill: var(--auth); filter: drop-shadow(0 0 4px var(--auth)); }

.afv-node.afv-active .afv-shape { stroke-width: 2.4; }
.afv-edge.afv-active { stroke-opacity: 0.95; stroke-width: 2.2; }
.afv-dim { opacity: 0.15; }

/* ---- Summary ---- */
.afv-summary-grid { display: flex; gap: 20px; flex-wrap: wrap; margin-bottom: 30px; }
.afv-summary-card {
  border: 1px solid var(--border);
  background: var(--panel);
  border-radius: 12px;
  padding: 26px 30px;
  min-width: 160px;
}
.afv-summary-num { font-size: 40px; font-weight: 800; color: var(--function); line-height: 1; }
.afv-summary-card:nth-child(2) .afv-summary-num { color: var(--endpoint); }
.afv-summary-card:nth-child(3) .afv-summary-num { color: var(--auth); }
.afv-summary-label { color: var(--text-dim); font-size: 11px; text-transform: uppercase; letter-spacing: 0.1em; margin-top: 8px; }
.afv-footer-line { color: var(--text-dim); font-size: 12px; margin-top: 10px; }

/* ---- Nav dots ---- */
.afv-navdots {
  position: fixed;
  right: 26px;
  top: 50%;
  transform: translateY(-50%);
  display: flex;
  flex-direction: column;
  gap: 14px;
  z-index: 50;
}
.afv-navdot {
  width: 8px; height: 8px; border-radius: 50%;
  background: var(--border);
  cursor: pointer;
  transition: all 0.25s ease;
  border: none;
  padding: 0;
}
.afv-navdot.afv-active { background: var(--function); transform: scale(1.5); box-shadow: 0 0 0 4px rgba(79,209,197,0.14); }

#afv-tooltip {
  position: fixed;
  pointer-events: none;
  background: var(--panel);
  border: 1px solid var(--border);
  padding: 6px 10px;
  border-radius: 6px;
  font-size: 11px;
  color: var(--text);
  opacity: 0;
  z-index: 100;
}

@media (max-width: 720px) {
  .afv-navdots { display: none; }
  .afv-panel { max-height: 50vh; }
}

</style>
</head>
<body>

<div class="afv-navdots">
  <button class="afv-navdot" data-target="afv-sec-hero" aria-label="Overview"></button>
  <button class="afv-navdot" data-target="afv-sec-program" aria-label="Program flow"></button>
  <button class="afv-navdot" data-target="afv-sec-auth" aria-label="Auth flow"></button>
  <button class="afv-navdot" data-target="afv-sec-summary" aria-label="Summary"></button>
</div>

<div id="afv-tooltip"></div>

<main class="afv-scroller">

  <section class="afv-section afv-hero" id="afv-sec-hero">
<pre class="afv-banner">   █████████   ███████████ █████   █████
  ███░░░░░███ ░░███░░░░░░█░░███   ░░███ 
 ░███    ░███  ░███   █ ░  ░███    ░███ 
 ░███████████  ░███████    ░███    ░███ 
 ░███░░░░░███  ░███░░░█    ░░███   ███  
 ░███    ░███  ░███  ░      ░░░█████░   
 █████   █████ █████          ░░███     
░░░░░   ░░░░░ ░░░░░            ░░░      </pre>
    <h1>Auth Flow Visualizer</h1>
    <h5>&lt;by abysseraphim github/&gt;<h5>
    <p class="afv-sub">Target: <code id="afv-hero-target">&nbsp;</code></p>

    <div class="afv-chip-row">
      <div class="afv-chip"><div class="afv-chip-num" id="afv-hero-functions">0</div><div class="afv-chip-label">Functions</div></div>
      <div class="afv-chip"><div class="afv-chip-num" id="afv-hero-endpoints">0</div><div class="afv-chip-label">Endpoints</div></div>
      <div class="afv-chip"><div class="afv-chip-num" id="afv-hero-calls">0</div><div class="afv-chip-label">Calls</div></div>
    </div>

    <div class="afv-scrolldown">&darr; scroll to explore</div>
  </section>

  <section class="afv-section" id="afv-sec-program">
    <div class="afv-eyebrow">01 &mdash; Program Flow</div>
    <h2>Every function, every call.</h2>
    <p class="afv-desc">The full call graph extracted from the codebase &mdash; who calls whom, and which requests leave the program entirely.</p>
    <div class="afv-panel"><div id="afv-program-graph"></div></div>
    <div class="afv-legend">
      <div class="afv-legend-item"><span class="afv-legend-swatch sw-function"></span>Function</div>
      <div class="afv-legend-item"><span class="afv-legend-swatch sw-endpoint"></span>Endpoint</div>
      <div class="afv-legend-item"><span class="afv-legend-swatch sw-external"></span>External call</div>
    </div>
  </section>

  <section class="afv-section afv-section-auth" id="afv-sec-auth">
    <div class="afv-eyebrow">02 &mdash; Auth Flow</div>
    <h2>Where authentication actually goes.</h2>
    <p class="afv-desc">Traced from every entry point whose name suggests a login, forward through everything it touches.</p>
    <div class="afv-panel"><div id="afv-auth-graph"></div></div>
    <div class="afv-legend">
      <div class="afv-legend-item"><span class="afv-legend-swatch sw-auth"></span>Auth-flow call</div>
      <div class="afv-legend-item"><span class="afv-legend-swatch sw-external"></span>External call</div>
    </div>
  </section>

  <section class="afv-section" id="afv-sec-summary">
    <div class="afv-eyebrow">03 &mdash; Summary</div>
    <h2>Analysis complete.</h2>
    <div class="afv-summary-grid">
      <div class="afv-summary-card"><div class="afv-summary-num" id="afv-stat-functions">0</div><div class="afv-summary-label">Functions</div></div>
      <div class="afv-summary-card"><div class="afv-summary-num" id="afv-stat-endpoints">0</div><div class="afv-summary-label">Endpoints</div></div>
      <div class="afv-summary-card"><div class="afv-summary-num" id="afv-stat-calls">0</div><div class="afv-summary-label">Calls</div></div>
    </div>
    <p class="afv-footer-line">Generated by Auth Flow Visualizer v1.0.0</p>
  </section>

</main>

<script type="application/json" id="afv-data">__AFV_DATA__</script>
<script>
(function () {
  "use strict";

  var DATA = JSON.parse(document.getElementById("afv-data").textContent);

  var RADIUS = { function: 30, endpoint: 20, external: 22 };
  var COLOR = { function: "var(--function)", endpoint: "var(--endpoint)", external: "var(--external)" };

  function truncate(s, n) {
    if (!s) return "";
    return s.length > n ? s.slice(0, n - 1) + "\u2026" : s;
  }

  function layoutGraph(nodes, edges) {
    var incoming = {};
    var i, n, e;
    for (i = 0; i < nodes.length; i++) incoming[nodes[i].id] = 0;
    for (i = 0; i < edges.length; i++) {
      e = edges[i];
      if (incoming[e.to] === undefined) incoming[e.to] = 0;
      incoming[e.to] += 1;
    }

    var layer = {};
    var hasRoot = false;
    for (i = 0; i < nodes.length; i++) {
      n = nodes[i];
      if (!incoming[n.id]) {
        layer[n.id] = 0;
        hasRoot = true;
      }
    }
    if (!hasRoot) {
      for (i = 0; i < nodes.length; i++) layer[nodes[i].id] = 0;
    }

    var changed = true;
    var guard = 0;
    while (changed && guard < nodes.length + 5) {
      changed = false;
      guard += 1;
      for (i = 0; i < edges.length; i++) {
        e = edges[i];
        if (layer[e.from] === undefined) continue;
        var next = layer[e.from] + 1;
        if (layer[e.to] === undefined || layer[e.to] < next) {
          layer[e.to] = next;
          changed = true;
        }
      }
    }
    for (i = 0; i < nodes.length; i++) {
      if (layer[nodes[i].id] === undefined) layer[nodes[i].id] = 0;
    }

    var byLayer = {};
    for (i = 0; i < nodes.length; i++) {
      n = nodes[i];
      var l = layer[n.id];
      if (!byLayer[l]) byLayer[l] = [];
      byLayer[l].push(n);
    }

    var typeOrder = { "function": 0, endpoint: 1, external: 2 };
    var layerKeys = Object.keys(byLayer).map(Number).sort(function (a, b) { return a - b; });
    layerKeys.forEach(function (l) {
      byLayer[l].sort(function (a, b) {
        var ta = typeOrder[a.type] === undefined ? 9 : typeOrder[a.type];
        var tb = typeOrder[b.type] === undefined ? 9 : typeOrder[b.type];
        if (ta !== tb) return ta - tb;
        return a.label.localeCompare(b.label);
      });
    });

    var layerGapX = 210;
    var nodeGapY = 96;
    var paddingX = 90;
    var paddingY = 70;

    var maxCount = 1;
    layerKeys.forEach(function (l) {
      if (byLayer[l].length > maxCount) maxCount = byLayer[l].length;
    });
    var height = Math.max(maxCount * nodeGapY + paddingY * 2, 340);

    var pos = {};
    layerKeys.forEach(function (l, li) {
      var items = byLayer[l];
      var totalH = items.length * nodeGapY;
      var startY = (height - totalH) / 2 + nodeGapY / 2;
      items.forEach(function (node, idx) {
        pos[node.id] = { x: paddingX + li * layerGapX, y: startY + idx * nodeGapY };
      });
    });

    var width = paddingX * 2 + (layerKeys.length - 1) * layerGapX + 60;
    return { pos: pos, width: Math.max(width, 460), height: height };
  }

  function svgEl(tag, attrs) {
    var el = document.createElementNS("http://www.w3.org/2000/svg", tag);
    for (var k in attrs) {
      if (Object.prototype.hasOwnProperty.call(attrs, k)) {
        el.setAttribute(k, attrs[k]);
      }
    }
    return el;
  }

  function nodeShape(n, x, y) {
    var g = svgEl("g", { class: "afv-node", "data-id": n.id, "data-type": n.type, transform: "translate(" + x + "," + y + ")" });

    if (n.type === "endpoint") {
      var r = RADIUS.endpoint;
      var d = svgEl("rect", {
        x: -r, y: -r, width: r * 2, height: r * 2, rx: 6,
        transform: "rotate(45)",
        class: "afv-shape afv-shape-endpoint"
      });
      g.appendChild(d);
    } else {
      var radius = n.type === "function" ? RADIUS.function : RADIUS.external;
      var c = svgEl("circle", { r: radius, class: "afv-shape afv-shape-" + n.type });
      g.appendChild(c);
    }

    var label = svgEl("text", { class: "afv-label", y: (n.type === "endpoint" ? RADIUS.endpoint : (n.type === "function" ? RADIUS.function : RADIUS.external)) + 20, "text-anchor": "middle" });
    label.textContent = truncate(n.label, 16);
    g.appendChild(label);

    if (n.sub) {
      var sub = svgEl("text", { class: "afv-sublabel", y: (n.type === "endpoint" ? RADIUS.endpoint : (n.type === "function" ? RADIUS.function : RADIUS.external)) + 34, "text-anchor": "middle" });
      sub.textContent = truncate(n.sub, 22);
      g.appendChild(sub);
    }

    var title = svgEl("title", {});
    title.textContent = n.label + (n.sub ? " \u2014 " + n.sub : "");
    g.appendChild(title);

    return g;
  }

  function edgePath(x1, y1, x2, y2) {
    var mx = (x1 + x2) / 2;
    return "M " + x1 + " " + y1 + " C " + mx + " " + y1 + ", " + mx + " " + y2 + ", " + x2 + " " + y2;
  }

  function renderGraph(container, nodes, edges, opts) {
    container.innerHTML = "";

    if (!nodes.length) {
      var empty = document.createElement("div");
      empty.className = "afv-empty";
      empty.textContent = opts.emptyText || "Nothing to show here.";
      container.appendChild(empty);
      return;
    }

    var layout = layoutGraph(nodes, edges);
    var svg = svgEl("svg", {
      viewBox: "0 0 " + layout.width + " " + layout.height,
      width: layout.width,
      height: layout.height,
      class: "afv-graph-svg"
    });

    var defs = svgEl("defs", {});
    var edgeTypes = {};
    edges.forEach(function (e) { edgeTypes[e.type] = true; });
    Object.keys(edgeTypes).forEach(function (t) {
      var marker = svgEl("marker", {
        id: opts.idPrefix + "-arrow-" + t, viewBox: "0 0 10 10", refX: "9", refY: "5",
        markerWidth: "7", markerHeight: "7", orient: "auto-start-reverse"
      });
      var arrowPath = svgEl("path", { d: "M 0 0 L 10 5 L 0 10 z" });
      marker.appendChild(arrowPath);
      defs.appendChild(marker);
    });
    svg.appendChild(defs);

    var edgeLayer = svgEl("g", { class: "afv-edges" });
    var nodeLayer = svgEl("g", { class: "afv-nodes" });

    var byId = {};
    nodes.forEach(function (n) { byId[n.id] = n; });

    edges.forEach(function (e, idx) {
      var from = layout.pos[e.from];
      var to = layout.pos[e.to];
      if (!from || !to) return;
      var fromType = byId[e.from] ? byId[e.from].type : "function";
      var toType = byId[e.to] ? byId[e.to].type : "function";
      var r1 = RADIUS[fromType] || 26;
      var r2 = RADIUS[toType] || 26;
      var x1 = from.x + r1, y1 = from.y;
      var x2 = to.x - r2, y2 = to.y;
      var d = edgePath(x1, y1, x2, y2);
      var pathId = opts.idPrefix + "-edge-" + idx;

      var path = svgEl("path", {
        d: d, id: pathId, class: "afv-edge afv-edge-" + e.type,
        "marker-end": "url(#" + opts.idPrefix + "-arrow-" + e.type + ")",
        "data-from": e.from, "data-to": e.to
      });
      edgeLayer.appendChild(path);

      if (opts.pulse) {
        var dot = svgEl("circle", { r: 3, class: "afv-pulse" });
        var animate = svgEl("animateMotion", { dur: (2.4 + (idx % 3) * 0.6) + "s", repeatCount: "indefinite" });
        var mpath = svgEl("mpath", {});
        mpath.setAttributeNS("http://www.w3.org/1999/xlink", "href", "#" + pathId);
        animate.appendChild(mpath);
        dot.appendChild(animate);
        edgeLayer.appendChild(dot);
      }
    });

    nodes.forEach(function (n) {
      var p = layout.pos[n.id];
      if (!p) return;
      nodeLayer.appendChild(nodeShape(n, p.x, p.y));
    });

    svg.appendChild(edgeLayer);
    svg.appendChild(nodeLayer);
    container.appendChild(svg);

    var tooltip = document.getElementById("afv-tooltip");

    nodeLayer.querySelectorAll(".afv-node").forEach(function (nodeEl) {
      nodeEl.addEventListener("mouseenter", function () {
        var id = nodeEl.getAttribute("data-id");
        container.querySelectorAll(".afv-node, .afv-edge, .afv-pulse").forEach(function (el) {
          el.classList.add("afv-dim");
        });
        nodeEl.classList.remove("afv-dim");
        nodeEl.classList.add("afv-active");
        container.querySelectorAll('.afv-edge[data-from="' + CSS.escape(id) + '"], .afv-edge[data-to="' + CSS.escape(id) + '"]').forEach(function (edgeEl) {
          edgeEl.classList.remove("afv-dim");
          edgeEl.classList.add("afv-active");
        });
      });
      nodeEl.addEventListener("mouseleave", function () {
        container.querySelectorAll(".afv-node, .afv-edge, .afv-pulse").forEach(function (el) {
          el.classList.remove("afv-dim");
          el.classList.remove("afv-active");
        });
      });
    });
  }

  function renderSummary() {
    var s = DATA.summary || {};
    var setText = function (id, val) {
      var el = document.getElementById(id);
      if (el) el.textContent = val;
    };
    setText("afv-stat-functions", s.functions || 0);
    setText("afv-stat-endpoints", s.endpoints || 0);
    setText("afv-stat-calls", s.calls || 0);
    setText("afv-hero-functions", s.functions || 0);
    setText("afv-hero-endpoints", s.endpoints || 0);
    setText("afv-hero-calls", s.calls || 0);
    setText("afv-hero-target", DATA.target || "");
  }

  function setupScrollSpy() {
    var sections = document.querySelectorAll(".afv-section");
    var dots = document.querySelectorAll(".afv-navdot");
    if (!sections.length) return;

    if (typeof IntersectionObserver === "undefined") {
      sections.forEach(function (s) { s.classList.add("afv-in-view"); });
      dots.forEach(function (dot) {
        dot.addEventListener("click", function () {
          var target = document.getElementById(dot.getAttribute("data-target"));
          if (target) target.scrollIntoView({ behavior: "smooth" });
        });
      });
      return;
    }

    var observer = new IntersectionObserver(function (entries) {
      entries.forEach(function (entry) {
        if (entry.isIntersecting) {
          entry.target.classList.add("afv-in-view");
          var id = entry.target.getAttribute("id");
          dots.forEach(function (dot) {
            dot.classList.toggle("afv-active", dot.getAttribute("data-target") === id);
          });
        }
      });
    }, { threshold: 0.5 });

    sections.forEach(function (s) { observer.observe(s); });

    dots.forEach(function (dot) {
      dot.addEventListener("click", function () {
        var target = document.getElementById(dot.getAttribute("data-target"));
        if (target) target.scrollIntoView({ behavior: "smooth" });
      });
    });
  }

  document.addEventListener("DOMContentLoaded", function () {
    renderSummary();

    renderGraph(document.getElementById("afv-program-graph"), DATA.program.nodes, DATA.program.edges, {
      idPrefix: "program",
      pulse: false,
      emptyText: "No function calls detected in this codebase."
    });

    renderGraph(document.getElementById("afv-auth-graph"), DATA.auth.nodes, DATA.auth.edges, {
      idPrefix: "auth",
      pulse: true,
      emptyText: "No authentication entry point detected (looking for a function with \u2018login\u2019 in its name)."
    });

    setupScrollSpy();
  });
})();

</script>
</body>
</html>
`
