const assert = require("node:assert/strict");
const test = require("node:test");
const { PomodoroEngine } = require("./pomodoroEngine.js");

function createEngine(overrides = {}) {
  return new PomodoroEngine({
    workDurationSec: 10,
    shortBreakDurationSec: 5,
    longBreakDurationSec: 8,
    longBreakInterval: 2,
    ...overrides,
  });
}

test("engine starts, pauses, resumes, and resets", () => {
  const engine = createEngine();
  assert.equal(engine.state().status, "idle");
  engine.start();
  engine.tick(3);
  assert.equal(engine.state().remainingSec, 7);
  engine.pause();
  engine.tick(2);
  assert.equal(engine.state().remainingSec, 7);
  engine.start();
  engine.tick(2);
  assert.equal(engine.state().remainingSec, 5);
  engine.reset();
  assert.deepEqual(engine.state(), {
    status: "idle", mode: "work", remainingSec: 10, totalSec: 10, completedWorkSessions: 0,
  });
});

test("work completion transitions to a short break", () => {
  const engine = createEngine();
  const completed = [];
  engine.onSessionComplete = (session) => completed.push(session);
  engine.start();
  engine.tick(10);

  assert.equal(engine.state().mode, "shortBreak");
  assert.equal(engine.state().remainingSec, 5);
  assert.equal(engine.state().status, "idle");
  assert.equal(engine.state().completedWorkSessions, 1);
  assert.equal(completed[0].type, "work");
});

test("every configured number of work sessions transitions to a long break", () => {
  const engine = createEngine();
  engine.start();
  engine.tick(10);
  engine.start();
  engine.tick(5);
  engine.start();
  engine.tick(10);

  assert.equal(engine.state().mode, "longBreak");
  assert.equal(engine.state().remainingSec, 8);
});

test("autoStart keeps the engine running after a session", () => {
  const engine = createEngine({ autoStart: true });
  engine.start();
  engine.tick(10);
  assert.equal(engine.state().status, "running");
});

test("invalid duration settings are normalized to positive integers", () => {
  const engine = createEngine({ workDurationSec: 0, shortBreakDurationSec: -3, longBreakDurationSec: NaN, longBreakInterval: 0 });
  engine.start();
  engine.tick(3);
  assert.ok(engine.state().remainingSec >= 1);
});