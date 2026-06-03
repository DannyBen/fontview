const style = document.createElement("style");
style.textContent = fonts.map((font) => {
  return `@font-face{font-family:${JSON.stringify(font.id)};src:url(${JSON.stringify(font.path)});}`;
}).join("\n");
document.head.append(style);

const allUngroupedCategory = "__all_ungrouped__";

const state = {
  view: "map",
  font: fonts[0]?.id,
  category: allUngroupedCategory,
  size: 72,
  weight: "400",
  bold: false,
  italic: false,
  query: ""
};

const els = {
  viewButtons: document.querySelectorAll("[data-view]"),
  viewPanels: document.querySelectorAll("[data-view-panel]"),
  font: document.querySelector("#fontSelect"),
  category: document.querySelector("#categorySelect"),
  size: document.querySelector("#sizeInput"),
  weight: document.querySelector("#weightSelect"),
  bold: document.querySelector("#boldButton"),
  italic: document.querySelector("#italicButton"),
  search: document.querySelector("#searchInput"),
  sample: document.querySelector("#sampleText"),
  content: document.querySelector("#content"),
  meta: document.querySelector("#fontMeta"),
  toast: document.querySelector("#toast")
};

for (const font of fonts) {
  const option = document.createElement("option");
  option.value = font.id;
  option.textContent = font.name;
  els.font.append(option);
}

for (const button of els.viewButtons) {
  button.addEventListener("click", () => {
    state.view = button.dataset.view;
    renderView();
  });
}
els.font.addEventListener("change", () => {
  state.font = els.font.value;
  state.category = allUngroupedCategory;
  render();
});
els.category.addEventListener("change", () => {
  state.category = els.category.value;
  renderGlyphs();
});
els.size.addEventListener("change", () => {
  setSize(els.size.value);
});
els.size.addEventListener("blur", () => {
  els.size.value = String(state.size);
});
els.size.addEventListener("keydown", (event) => {
  handleSizeShortcut(event);
});
els.weight.addEventListener("change", () => {
  state.weight = els.weight.value;
  applyFontSettings();
});
els.bold.addEventListener("click", () => {
  state.bold = !state.bold;
  els.bold.setAttribute("aria-pressed", String(state.bold));
  applyFontSettings();
});
els.italic.addEventListener("click", () => {
  state.italic = !state.italic;
  els.italic.setAttribute("aria-pressed", String(state.italic));
  applyFontSettings();
});
els.search.addEventListener("input", () => {
  state.query = els.search.value.trim().toLowerCase();
  renderGlyphs();
});
document.addEventListener("keydown", (event) => {
  if (event.target === els.size) {
    return;
  }
  if (isEditableTarget(event.target)) {
    return;
  }
  handleSizeShortcut(event);
});

function activeFont() {
  return fonts.find((font) => font.id === state.font) ?? fonts[0];
}

function render() {
  const font = activeFont();
  document.documentElement.style.setProperty("--active-font", JSON.stringify(font.id));
  applyFontSettings();
  els.meta.textContent = `${font.glyphs.length.toLocaleString()} glyphs`;

  const categories = [
    { value: allUngroupedCategory, label: "All (Ungrouped)" },
    { value: "All", label: "All" },
    ...Array.from(new Set(font.glyphs.map((glyph) => glyph.category)), (category) => {
      return { value: category, label: category };
    })
  ];
  els.category.replaceChildren(...categories.map((category) => {
    const option = document.createElement("option");
    option.value = category.value;
    option.textContent = category.label;
    option.selected = category.value === state.category;
    return option;
  }));
  if (!categories.some((category) => category.value === state.category)) {
    state.category = allUngroupedCategory;
    els.category.value = allUngroupedCategory;
  }
  renderView();
  renderGlyphs();
}

function renderView() {
  for (const button of els.viewButtons) {
    button.classList.toggle("is-active", button.dataset.view === state.view);
  }
  for (const panel of els.viewPanels) {
    panel.classList.toggle("is-active", panel.dataset.viewPanel === state.view);
  }
  if (state.view === "type") {
    els.sample.focus();
  }
}

function applyFontSettings() {
  const weight = state.bold ? "700" : state.weight;
  const style = state.italic ? "italic" : "normal";
  document.documentElement.style.setProperty("--glyph-size", `${state.size}px`);
  document.documentElement.style.setProperty("--sample-size", `${state.size}px`);
  document.documentElement.style.setProperty("--glyph-weight", weight);
  document.documentElement.style.setProperty("--sample-weight", weight);
  document.documentElement.style.setProperty("--glyph-style", style);
  document.documentElement.style.setProperty("--sample-style", style);
}

function setSize(value) {
  const next = Number(value);
  if (!Number.isFinite(next)) {
    els.size.value = String(state.size);
    return;
  }
  state.size = clamp(Math.round(next), 8, 320);
  els.size.value = String(state.size);
  applyFontSettings();
}

function handleSizeShortcut(event) {
  if (event.altKey || event.ctrlKey || event.metaKey) {
    return;
  }
  if (event.key === "+" || event.code === "NumpadAdd") {
    event.preventDefault();
    setSize(state.size + 5);
  } else if (event.key === "-" || event.code === "Minus" || event.code === "NumpadSubtract") {
    event.preventDefault();
    setSize(state.size - 5);
  }
}

function isEditableTarget(target) {
  return target instanceof HTMLInputElement ||
    target instanceof HTMLTextAreaElement ||
    target instanceof HTMLSelectElement ||
    target?.isContentEditable;
}

function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max);
}

function renderGlyphs() {
  const font = activeFont();
  const query = state.query;
  const glyphs = font.glyphs.filter((glyph) => {
    const matchesCategory = state.category === allUngroupedCategory ||
      state.category === "All" ||
      glyph.category === state.category;
    const haystack = `${glyph.char} ${glyph.code} ${glyph.category} ${glyph.name}`.toLowerCase();
    return matchesCategory && (!query || haystack.includes(query));
  });

  if (state.category === allUngroupedCategory) {
    els.content.replaceChildren(section("All", glyphs));
    return;
  }

  const grouped = groupBy(glyphs, (glyph) => glyph.category);
  els.content.replaceChildren(...Array.from(grouped, ([category, items]) => section(category, items)));
}

function section(category, glyphs) {
  const root = document.createElement("section");
  root.className = "section";

  const title = document.createElement("div");
  title.className = "section-title";
  title.innerHTML = '<h2></h2><span class="meta"></span>';
  title.querySelector("h2").textContent = category;
  title.querySelector(".meta").textContent = `${glyphs.length.toLocaleString()} glyphs`;

  const grid = document.createElement("div");
  grid.className = "grid";
  for (const glyph of glyphs) {
    grid.append(glyphButton(glyph));
  }

  root.append(title, grid);
  return root;
}

function glyphButton(glyph) {
  const button = document.createElement("button");
  const label = glyphLabel(glyph);
  button.className = "glyph";
  button.type = "button";
  button.title = `Copy ${label}`;
  button.innerHTML = '<span class="glyph-symbol"></span>';
  button.querySelector(".glyph-symbol").textContent = glyph.char;
  button.addEventListener("click", async () => {
    await navigator.clipboard.writeText(glyph.char);
    showToast(`Copied ${label}`);
  });
  return button;
}

function glyphLabel(glyph) {
  const visible = glyph.name || glyph.char;
  return `${visible} (${glyph.code})`;
}

function groupBy(items, keyFn) {
  const result = new Map();
  for (const item of items) {
    const key = keyFn(item);
    const group = result.get(key) ?? [];
    group.push(item);
    result.set(key, group);
  }
  return result;
}

let toastTimer;
function showToast(message) {
  els.toast.textContent = message;
  els.toast.classList.add("is-visible");
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => els.toast.classList.remove("is-visible"), 1200);
}

render();
