const assert = require("node:assert/strict");
const test = require("node:test");
const { initializeApp, loadSettings, saveSettings } = require("./app.js");

test("initializeApp wires the engine and view", () => {
  const document = { title: "", getElementById: () => null };
  const events = {};
  const engine = {
    status: "idle",
    subscribe: (listener) => listener({ mode: "work", remainingSec: 10, totalSec: 10, completedWorkSessions: 0 }),
    tick: () => {},
    start: () => { events.started = true; },
    pause: () => {},
    reset: () => {},
  };
  const view = {
    render: () => { events.rendered = true; },
    setProgress: () => {},
    setMessage: () => {},
    elements: {
      startButton: { addEventListener: () => {} },
      resetButton: { addEventListener: () => {} },
    },
  };

  const app = initializeApp(document, {
    Engine: function Engine() { return engine; },
    viewFactory: () => view,
  });

  assert.equal(document.title, "ポモドーロタイマー");
  assert.equal(app.engine, engine);
  assert.equal(events.rendered, true);
  clearInterval(app.interval);
});

test("settings are saved and loaded from storage", () => {
  const values = new Map();
  const storage = {
    getItem: (key) => values.get(key) || null,
    setItem: (key, value) => values.set(key, value),
  };
  const settings = { workDurationSec: 1200, autoStart: true };

  saveSettings(storage, settings);

  assert.deepEqual(loadSettings(storage), settings);
});
