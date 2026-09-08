async function fetchTodayStats(fetcher = fetch) {
  const response = await fetcher("/api/stats/today");
  if (!response.ok) throw new Error("今日の進捗を取得できませんでした");
  return response.json();
}

async function fetchHistory(fetcher = fetch, task = "") {
  const query = task ? `?task=${encodeURIComponent(task)}` : "";
  const response = await fetcher(`/api/history${query}`);
  if (!response.ok) throw new Error("履歴を取得できませんでした");
  return response.json();
}

async function recordSession(session, fetcher = fetch) {
  const response = await fetcher("/api/sessions", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(session),
  });
  if (!response.ok) throw new Error("セッションを記録できませんでした");
  return response.json();
}

const OFFLINE_QUEUE_KEY = "pomodoro-offline-sessions";

function readQueuedSessions(storage = globalThis.localStorage) {
  if (!storage) return [];
  try {
    return JSON.parse(storage.getItem(OFFLINE_QUEUE_KEY)) || [];
  } catch {
    return [];
  }
}

function queueSession(session, storage = globalThis.localStorage) {
  if (!storage) return;
  const queue = readQueuedSessions(storage);
  queue.push(session);
  storage.setItem(OFFLINE_QUEUE_KEY, JSON.stringify(queue));
}

async function recordSessionWithOfflineQueue(session, fetcher = fetch, storage = globalThis.localStorage) {
  try {
    return await recordSession(session, fetcher);
  } catch (error) {
    queueSession(session, storage);
    return { queued: true, error };
  }
}

async function flushQueuedSessions(fetcher = fetch, storage = globalThis.localStorage) {
  const queue = readQueuedSessions(storage);
  const remaining = [];
  for (const session of queue) {
    try {
      await recordSession(session, fetcher);
    } catch {
      remaining.push(session);
    }
  }
  if (storage) storage.setItem(OFFLINE_QUEUE_KEY, JSON.stringify(remaining));
  return { sent: queue.length - remaining.length, remaining: remaining.length };
}

function notifyCompletion(mode, notification = globalThis.Notification) {
  if (typeof notification !== "function") return false;
  if (notification.permission === "granted") {
    new notification(mode === "work" ? "作業セッション完了" : "休憩終了", { body: "次のセッションを始められます" });
    return true;
  }
  if (notification.permission === "default" && typeof notification.requestPermission === "function") {
    notification.requestPermission();
  }
  return false;
}

if (typeof globalThis !== "undefined") {
  globalThis.fetchTodayStats = fetchTodayStats;
  globalThis.fetchHistory = fetchHistory;
  globalThis.recordSession = recordSession;
  globalThis.recordSessionWithOfflineQueue = recordSessionWithOfflineQueue;
  globalThis.flushQueuedSessions = flushQueuedSessions;
  globalThis.notifyCompletion = notifyCompletion;
}
if (typeof module !== "undefined" && module.exports) module.exports = {
  fetchTodayStats,
  fetchHistory,
  recordSession,
  readQueuedSessions,
  queueSession,
  recordSessionWithOfflineQueue,
  flushQueuedSessions,
  notifyCompletion,
};