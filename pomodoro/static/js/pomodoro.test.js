const test = require("node:test");
const assert = require("node:assert/strict");
const pomodoro = require("./pomodoro");

test("sanitizeSettings は許可されない値をデフォルトに戻す", () => {
  const sanitized = pomodoro.sanitizeSettings({
    workDuration: 99,
    breakDuration: 99,
    theme: "unknown",
    startSound: 1,
    endSound: 0,
    tickSound: "yes",
  });

  assert.equal(sanitized.workDuration, 25);
  assert.equal(sanitized.breakDuration, 5);
  assert.equal(sanitized.theme, "light");
  assert.equal(sanitized.startSound, true);
  assert.equal(sanitized.endSound, false);
  assert.equal(sanitized.tickSound, true);
});

test("recordCustomization は設定変更時のみカウントする", () => {
  const stats = pomodoro.createUsageStats();
  const before = pomodoro.sanitizeSettings();
  const same = pomodoro.sanitizeSettings();
  const changed = pomodoro.sanitizeSettings({ workDuration: 35, breakDuration: 10, theme: "dark" });

  pomodoro.recordCustomization(stats, before, same);
  assert.equal(stats.customizationCount, 0);

  pomodoro.recordCustomization(stats, same, changed);
  assert.equal(stats.customizationCount, 1);
  assert.equal(stats.optionUsage.workDuration[35], 1);
  assert.equal(stats.optionUsage.breakDuration[10], 1);
  assert.equal(stats.optionUsage.theme.dark, 1);
});

test("nextPhase と formatRemaining の基本動作", () => {
  assert.equal(pomodoro.nextPhase("work"), "break");
  assert.equal(pomodoro.nextPhase("break"), "work");
  assert.equal(pomodoro.formatRemaining(1500), "25:00");
  assert.equal(pomodoro.formatRemaining(59), "00:59");
});
