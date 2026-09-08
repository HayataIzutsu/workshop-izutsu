const DEFAULT_TIMER_SETTINGS = Object.freeze({
  workDurationSec: 25 * 60,
  shortBreakDurationSec: 5 * 60,
  longBreakDurationSec: 15 * 60,
  longBreakInterval: 4,
  autoStart: false,
});

class PomodoroEngine {
  constructor(settings = {}) {
    this.settings = { ...DEFAULT_TIMER_SETTINGS, ...settings };
    this.completedWorkSessions = 0;
    this.status = "idle";
    this.mode = "work";
    this.remainingSec = this.durationFor(this.mode);
    this.listeners = new Set();
    this.onSessionComplete = null;
  }

  durationFor(mode) {
    if (mode === "shortBreak") return this.settings.shortBreakDurationSec;
    if (mode === "longBreak") return this.settings.longBreakDurationSec;
    return this.settings.workDurationSec;
  }

  subscribe(listener) {
    this.listeners.add(listener);
    listener(this.state());
    return () => this.listeners.delete(listener);
  }

  state() {
    return {
      status: this.status,
      mode: this.mode,
      remainingSec: this.remainingSec,
      totalSec: this.durationFor(this.mode),
      completedWorkSessions: this.completedWorkSessions,
    };
  }

  notify() {
    const currentState = this.state();
    this.listeners.forEach((listener) => listener(currentState));
  }

  start() {
    if (this.status === "idle" || this.status === "paused") {
      this.status = "running";
      this.notify();
    }
  }

  pause() {
    if (this.status === "running") {
      this.status = "paused";
      this.notify();
    }
  }

  reset() {
    this.status = "idle";
    this.mode = "work";
    this.remainingSec = this.durationFor(this.mode);
    this.notify();
  }

  tick(deltaSec) {
    if (this.status !== "running" || deltaSec <= 0) return;
    let remainingDelta = deltaSec;
    while (remainingDelta >= this.remainingSec) {
      remainingDelta -= this.remainingSec;
      this.completeSession();
      if (this.status !== "running") return;
    }
    this.remainingSec -= remainingDelta;
    this.notify();
  }

  completeSession() {
    const completed = {
      type: this.mode === "work" ? "work" : this.mode,
      durationSec: this.durationFor(this.mode),
      completedAt: new Date().toISOString(),
    };
    if (this.mode === "work") this.completedWorkSessions += 1;
    if (typeof this.onSessionComplete === "function") this.onSessionComplete(completed);

    if (this.mode === "work") {
      this.mode = this.completedWorkSessions % this.settings.longBreakInterval === 0 ? "longBreak" : "shortBreak";
    } else {
      this.mode = "work";
    }
    this.remainingSec = this.durationFor(this.mode);
    this.status = this.settings.autoStart ? "running" : "idle";
    this.notify();
  }
}

if (typeof globalThis !== "undefined") globalThis.PomodoroEngine = PomodoroEngine;
if (typeof module !== "undefined" && module.exports) {
  module.exports = { PomodoroEngine, DEFAULT_TIMER_SETTINGS };
}