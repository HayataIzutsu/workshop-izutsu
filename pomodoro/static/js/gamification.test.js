const test = require("node:test");
const assert = require("node:assert/strict");
const game = require("./gamification.js");

test("完了セッションでXP加算とレベルアップ判定ができる", function () {
  const profile = { xp: 90, level: 1, streak: 0, totalCompletedSessions: 0, badges: [] };
  const sessions = [];

  const result = game.recordCompletedSession(profile, sessions, {
    now: new Date("2026-09-09T00:00:00Z"),
    focusMinutes: 25,
    plannedMinutes: 25,
  });

  assert.equal(result.profile.xp, 115);
  assert.equal(result.profile.level, 2);
  assert.equal(result.leveledUp, true);
});

test("3日連続でthree-day-streakバッジを獲得できる", function () {
  let profile = { xp: 0, level: 1, streak: 0, totalCompletedSessions: 0, badges: [] };
  let sessions = [];

  ["2026-09-07T01:00:00Z", "2026-09-08T01:00:00Z", "2026-09-09T01:00:00Z"].forEach((date) => {
    const result = game.recordCompletedSession(profile, sessions, {
      now: new Date(date),
      focusMinutes: 25,
      plannedMinutes: 25,
    });
    profile = result.profile;
    sessions = result.sessions;
  });

  assert.ok(profile.badges.includes("three-day-streak"));
});

test("7日以内10回完了でweekly-10バッジを獲得できる", function () {
  let profile = { xp: 0, level: 1, streak: 0, totalCompletedSessions: 0, badges: [] };
  let sessions = [];

  const timestamps = [
    "2026-09-03T09:00:00Z",
    "2026-09-04T09:00:00Z",
    "2026-09-05T09:00:00Z",
    "2026-09-06T09:00:00Z",
    "2026-09-07T09:00:00Z",
    "2026-09-08T09:00:00Z",
    "2026-09-09T09:00:00Z",
    "2026-09-09T10:00:00Z",
    "2026-09-09T11:00:00Z",
    "2026-09-09T12:00:00Z",
  ];

  timestamps.forEach((timestamp) => {
    const result = game.recordCompletedSession(profile, sessions, {
      now: new Date(timestamp),
      focusMinutes: 25,
      plannedMinutes: 25,
    });
    profile = result.profile;
    sessions = result.sessions;
  });

  assert.ok(profile.badges.includes("weekly-10"));
});

test("週次/月次統計で完了率と平均集中時間を算出できる", function () {
  const now = new Date("2026-09-09T12:00:00Z");
  const sessions = [
    { completed: true, completedAt: "2026-09-09T01:00:00Z", focusMinutes: 25, plannedMinutes: 25 },
    { completed: true, completedAt: "2026-09-08T01:00:00Z", focusMinutes: 30, plannedMinutes: 30 },
    { completed: false, completedAt: "2026-09-08T02:00:00Z", focusMinutes: 0, plannedMinutes: 25 },
  ];

  const weekly = game.buildPeriodStats(sessions, 7, now);

  assert.equal(weekly.plannedSessions, 3);
  assert.equal(weekly.completedSessions, 2);
  assert.equal(weekly.completionRate, 66.7);
  assert.equal(weekly.averageFocusMinutes, 27.5);
});
