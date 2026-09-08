const assert = require("node:assert/strict");
const test = require("node:test");
const { formatTime, progressPercent, formatDuration } = require("./format.js");

test("formatTime formats minutes and seconds", () => {
  assert.equal(formatTime(1500), "25:00");
  assert.equal(formatTime(65), "01:05");
  assert.equal(formatTime(-1), "00:00");
});

test("progressPercent is bounded between zero and one hundred", () => {
  assert.equal(progressPercent(50, 100), 50);
  assert.equal(progressPercent(150, 100), 100);
  assert.equal(progressPercent(-1, 100), 0);
});

test("formatDuration displays minutes and hours", () => {
  assert.equal(formatDuration(1500), "25分");
  assert.equal(formatDuration(3600), "1時間");
  assert.equal(formatDuration(5400), "1時間30分");
});