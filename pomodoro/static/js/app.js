function loadSettings(storage) {
  try {
    return JSON.parse(storage.getItem("pomodoro-settings")) || {};
  } catch {
    return {};
  }
}

function saveSettings(storage, settings) {
  storage.setItem("pomodoro-settings", JSON.stringify(settings));
}

function registerServiceWorker(navigatorObject = globalThis.navigator) {
  if (!navigatorObject?.serviceWorker) return Promise.resolve(null);
  return navigatorObject.serviceWorker.register("/sw.js");
}

function initializeApp(document, dependencies = {}) {
  document.title = "ポモドーロタイマー";
  const Engine = dependencies.Engine || globalThis.PomodoroEngine;
  const viewFactory = dependencies.viewFactory || globalThis.createView;
  const storage = dependencies.storage || globalThis.localStorage;
  const engine = new Engine(storage ? loadSettings(storage) : {});
  const view = viewFactory(document);
  const interval = setInterval(() => engine.tick(1), 1000);

  engine.subscribe((state) => view.render(state));
  engine.onSessionComplete = async (session) => {
    session.task = view.elements.taskName.value.trim();
    session.note = view.elements.taskNote.value.trim();
    try {
      const result = typeof globalThis.recordSessionWithOfflineQueue === "function"
        ? await globalThis.recordSessionWithOfflineQueue(session)
        : await globalThis.recordSession(session);
      view.setMessage(result?.queued ? "通信復旧後にセッションを同期します" : "セッションを記録しました");
      if (typeof globalThis.notifyCompletion === "function") globalThis.notifyCompletion(session.type);
    } catch (error) {
      view.setMessage(error.message);
    }
  };
  view.elements.startButton.addEventListener("click", () => {
    if (engine.status === "running") engine.pause();
    else engine.start();
  });
  view.elements.resetButton.addEventListener("click", () => engine.reset());

  const settingsButton = document.getElementById("settings-button");
  const dialog = document.getElementById("settings-dialog");
  if (settingsButton && dialog) settingsButton.addEventListener("click", () => dialog.showModal());
  const settingsForm = document.getElementById("settings-form");
  if (settingsForm) settingsForm.addEventListener("submit", (event) => {
    if (event.submitter?.id !== "save-settings") return;
    const nextSettings = {
      workDurationSec: Number(document.getElementById("work-minutes").value) * 60,
      shortBreakDurationSec: Number(document.getElementById("short-break-minutes").value) * 60,
      longBreakDurationSec: Number(document.getElementById("long-break-minutes").value) * 60,
      longBreakInterval: Number(document.getElementById("long-break-interval").value),
      autoStart: document.getElementById("auto-start").checked,
    };
    if (storage) saveSettings(storage, nextSettings);
    view.setMessage("設定を保存しました。リロード後に反映されます");
  });

  if (typeof globalThis.fetchTodayStats === "function") {
    globalThis.fetchTodayStats().then((stats) => view.setProgress(stats.completedWorkSessions, stats.focusTimeSec)).catch((error) => view.setMessage(error.message));
  }
  if (typeof globalThis.fetchHistory === "function") {
    globalThis.fetchHistory().then((sessions) => view.setHistory(sessions)).catch((error) => view.setMessage(error.message));
  }
  return { engine, view, interval };
}

if (typeof document !== "undefined") document.addEventListener("DOMContentLoaded", () => initializeApp(document));
if (typeof document !== "undefined") registerServiceWorker().catch(() => {});

if (typeof module !== "undefined" && module.exports) module.exports = { initializeApp, loadSettings, saveSettings, registerServiceWorker };
