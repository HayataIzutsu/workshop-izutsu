const assert = require("node:assert/strict");
const test = require("node:test");
const {
  fetchTodayStats,
  fetchHistory,
  recordSession,
  recordSessionWithOfflineQueue,
  flushQueuedSessions,
  notifyCompletion,
} = require("./api.js");

test("fetchTodayStats returns JSON data", async () => {
  const fetcher = async (path) => {
    assert.equal(path, "/api/stats/today");
    return { ok: true, json: async () => ({ completedWorkSessions: 2, focusTimeSec: 3000 }) };
  };
  assert.deepEqual(await fetchTodayStats(fetcher), { completedWorkSessions: 2, focusTimeSec: 3000 });
});

test("fetchHistory requests the history endpoint", async () => {
  const fetcher = async (path) => {
    assert.equal(path, "/api/history?task=%E5%AE%9F%E8%A3%85");
    return { ok: true, json: async () => [] };
  };
  assert.deepEqual(await fetchHistory(fetcher, "実装"), []);
});

test("recordSession sends a JSON POST", async () => {
  const session = { type: "work", durationSec: 1500 };
  const fetcher = async (path, options) => {
    assert.equal(path, "/api/sessions");
    assert.equal(options.method, "POST");
    assert.equal(options.headers["Content-Type"], "application/json");
    assert.deepEqual(JSON.parse(options.body), session);
    return { ok: true, json: async () => ({ id: "session-1" }) };
  };
  assert.deepEqual(await recordSession(session, fetcher), { id: "session-1" });
});

test("API helpers reject failed responses", async () => {
  const fetcher = async () => ({ ok: false });
  await assert.rejects(() => fetchTodayStats(fetcher), /今日の進捗/);
  await assert.rejects(() => recordSession({}, fetcher), /セッション/);
});

test("failed session requests are queued and can be flushed later", async () => {
  const values = new Map();
  const storage = {
    getItem: (key) => values.get(key) || null,
    setItem: (key, value) => values.set(key, value),
  };
  const session = { type: "work", durationSec: 1500 };
  const failed = async () => ({ ok: false });
  const queued = await recordSessionWithOfflineQueue(session, failed, storage);
  assert.equal(queued.queued, true);

  const sent = [];
  const successful = async (_path, options) => {
    sent.push(JSON.parse(options.body));
    return { ok: true, json: async () => ({}) };
  };
  assert.deepEqual(await flushQueuedSessions(successful, storage), { sent: 1, remaining: 0 });
  assert.deepEqual(sent, [session]);
});

test("notifyCompletion uses granted browser notifications", () => {
  const notifications = [];
  function Notification(title, options) {
    notifications.push({ title, options });
  }
  Notification.permission = "granted";

  assert.equal(notifyCompletion("work", Notification), true);
  assert.equal(notifications[0].title, "作業セッション完了");
});