function createView(document) {
  const elements = {
    modeLabel: document.getElementById("mode-label"),
    remainingTime: document.getElementById("remaining-time"),
    timerRing: document.getElementById("timer-ring"),
    startButton: document.getElementById("start-button"),
    resetButton: document.getElementById("reset-button"),
    completedSessions: document.getElementById("completed-sessions"),
    focusTime: document.getElementById("focus-time"),
    statusMessage: document.getElementById("status-message"),
    taskName: document.getElementById("task-name"),
    taskNote: document.getElementById("task-note"),
    sessionHistory: document.getElementById("session-history"),
  };
  const modeLabels = { work: "作業中", shortBreak: "短い休憩", longBreak: "長い休憩" };

  function render(state) {
    elements.modeLabel.textContent = modeLabels[state.mode];
    elements.remainingTime.textContent = formatTime(state.remainingSec);
    elements.timerRing.style.setProperty("--progress", `${progressPercent(state.remainingSec, state.totalSec)}%`);
    elements.startButton.textContent = state.status === "running" ? "一時停止" : state.status === "paused" ? "再開" : "開始";
    elements.completedSessions.textContent = state.completedWorkSessions;
    return state;
  }

  function setProgress(completedSessions, focusTimeSec) {
    elements.completedSessions.textContent = completedSessions;
    elements.focusTime.textContent = formatDuration(focusTimeSec);
  }

  function setMessage(message) {
    elements.statusMessage.textContent = message;
  }

  function setHistory(sessions) {
    elements.sessionHistory.replaceChildren();
    if (sessions.length === 0) {
      elements.sessionHistory.append(Object.assign(document.createElement("li"), { textContent: "まだセッションはありません。" }));
      return;
    }
    sessions.slice(0, 10).forEach((session) => {
      const item = document.createElement("li");
      const task = session.task ? `: ${session.task}` : "";
      item.textContent = `${session.type}${task} (${formatDuration(session.durationSec)})`;
      elements.sessionHistory.append(item);
    });
  }

  return { render, setProgress, setMessage, setHistory, elements };
}

if (typeof globalThis !== "undefined") globalThis.createView = createView;
if (typeof module !== "undefined" && module.exports) module.exports = { createView };