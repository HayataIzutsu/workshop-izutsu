(function init(globalScope) {
  const DEFAULT_SETTINGS = {
    workDuration: 25,
    breakDuration: 5,
    theme: "light",
    startSound: true,
    endSound: true,
    tickSound: false,
  };

  const ALLOWED_WORK_DURATIONS = [15, 25, 35, 45];
  const ALLOWED_BREAK_DURATIONS = [5, 10, 15];
  const ALLOWED_THEMES = ["light", "dark", "focus"];

  function sanitizeSettings(input = {}) {
    const settings = { ...DEFAULT_SETTINGS, ...input };
    return {
      workDuration: ALLOWED_WORK_DURATIONS.includes(Number(settings.workDuration))
        ? Number(settings.workDuration)
        : DEFAULT_SETTINGS.workDuration,
      breakDuration: ALLOWED_BREAK_DURATIONS.includes(Number(settings.breakDuration))
        ? Number(settings.breakDuration)
        : DEFAULT_SETTINGS.breakDuration,
      theme: ALLOWED_THEMES.includes(settings.theme) ? settings.theme : DEFAULT_SETTINGS.theme,
      startSound: Boolean(settings.startSound),
      endSound: Boolean(settings.endSound),
      tickSound: Boolean(settings.tickSound),
    };
  }

  function createUsageStats() {
    return {
      customizationCount: 0,
      optionUsage: {
        workDuration: {},
        breakDuration: {},
        theme: {},
      },
    };
  }

  function recordOptionUsage(stats, settings) {
    stats.optionUsage.workDuration[settings.workDuration] =
      (stats.optionUsage.workDuration[settings.workDuration] || 0) + 1;
    stats.optionUsage.breakDuration[settings.breakDuration] =
      (stats.optionUsage.breakDuration[settings.breakDuration] || 0) + 1;
    stats.optionUsage.theme[settings.theme] = (stats.optionUsage.theme[settings.theme] || 0) + 1;
  }

  function recordCustomization(stats, previous, next) {
    if (JSON.stringify(previous) !== JSON.stringify(next)) {
      stats.customizationCount += 1;
      recordOptionUsage(stats, next);
    }
  }

  function toSeconds(minutes) {
    return minutes * 60;
  }

  function formatRemaining(seconds) {
    const minute = String(Math.floor(seconds / 60)).padStart(2, "0");
    const second = String(seconds % 60).padStart(2, "0");
    return `${minute}:${second}`;
  }

  function nextPhase(currentPhase) {
    return currentPhase === "work" ? "break" : "work";
  }

  function createBeepPlayer() {
    let audioContext;
    return function playBeep(durationMs = 100, frequency = 660) {
      if (typeof window === "undefined" || typeof window.AudioContext === "undefined") {
        return;
      }
      audioContext = audioContext || new window.AudioContext();
      const oscillator = audioContext.createOscillator();
      const gain = audioContext.createGain();
      oscillator.frequency.value = frequency;
      gain.gain.value = 0.05;
      oscillator.connect(gain);
      gain.connect(audioContext.destination);
      oscillator.start();
      setTimeout(() => oscillator.stop(), durationMs);
    };
  }

  function createTimerApp(documentRoot) {
    const workDurationSelect = documentRoot.getElementById("workDuration");
    const breakDurationSelect = documentRoot.getElementById("breakDuration");
    const themeSelect = documentRoot.getElementById("theme");
    const startSoundInput = documentRoot.getElementById("startSound");
    const endSoundInput = documentRoot.getElementById("endSound");
    const tickSoundInput = documentRoot.getElementById("tickSound");
    const phaseLabel = documentRoot.getElementById("phaseLabel");
    const timeDisplay = documentRoot.getElementById("timeDisplay");
    const startButton = documentRoot.getElementById("startButton");
    const resetButton = documentRoot.getElementById("resetButton");
    const statsDisplay = documentRoot.getElementById("stats");

    let currentSettings = sanitizeSettings(DEFAULT_SETTINGS);
    let phase = "work";
    let remainingSeconds = toSeconds(currentSettings.workDuration);
    let intervalId;
    const stats = createUsageStats();
    const playBeep = createBeepPlayer();

    function syncTheme(theme) {
      documentRoot.body.classList.remove("theme-light", "theme-dark", "theme-focus");
      documentRoot.body.classList.add(`theme-${theme}`);
    }

    function render() {
      phaseLabel.textContent = phase === "work" ? "作業" : "休憩";
      timeDisplay.textContent = formatRemaining(remainingSeconds);
      statsDisplay.textContent = JSON.stringify(stats, null, 2);
    }

    function readSettingsFromUI() {
      return sanitizeSettings({
        workDuration: Number(workDurationSelect.value),
        breakDuration: Number(breakDurationSelect.value),
        theme: themeSelect.value,
        startSound: startSoundInput.checked,
        endSound: endSoundInput.checked,
        tickSound: tickSoundInput.checked,
      });
    }

    function applySettings(nextSettings) {
      const previous = currentSettings;
      currentSettings = nextSettings;
      syncTheme(currentSettings.theme);
      recordCustomization(stats, previous, currentSettings);
      if (!intervalId) {
        remainingSeconds = toSeconds(phase === "work" ? currentSettings.workDuration : currentSettings.breakDuration);
      }
      render();
    }

    function stopTimer() {
      if (intervalId) {
        clearInterval(intervalId);
        intervalId = undefined;
      }
    }

    function resetTimer() {
      stopTimer();
      phase = "work";
      remainingSeconds = toSeconds(currentSettings.workDuration);
      render();
    }

    function advancePhase() {
      phase = nextPhase(phase);
      remainingSeconds = toSeconds(phase === "work" ? currentSettings.workDuration : currentSettings.breakDuration);
      if (currentSettings.endSound) {
        playBeep(180, 520);
      }
    }

    function startTimer() {
      if (intervalId) {
        return;
      }
      if (currentSettings.startSound) {
        playBeep(140, 880);
      }
      intervalId = setInterval(() => {
        remainingSeconds -= 1;
        if (currentSettings.tickSound) {
          playBeep(40, 420);
        }
        if (remainingSeconds <= 0) {
          advancePhase();
        }
        render();
      }, 1000);
    }

    workDurationSelect.addEventListener("change", () => applySettings(readSettingsFromUI()));
    breakDurationSelect.addEventListener("change", () => applySettings(readSettingsFromUI()));
    themeSelect.addEventListener("change", () => applySettings(readSettingsFromUI()));
    startSoundInput.addEventListener("change", () => applySettings(readSettingsFromUI()));
    endSoundInput.addEventListener("change", () => applySettings(readSettingsFromUI()));
    tickSoundInput.addEventListener("change", () => applySettings(readSettingsFromUI()));
    startButton.addEventListener("click", startTimer);
    resetButton.addEventListener("click", resetTimer);

    applySettings(currentSettings);

    return {
      sanitizeSettings,
      createUsageStats,
      recordCustomization,
      recordOptionUsage,
      nextPhase,
      formatRemaining,
    };
  }

  const publicAPI = {
    DEFAULT_SETTINGS,
    ALLOWED_WORK_DURATIONS,
    ALLOWED_BREAK_DURATIONS,
    ALLOWED_THEMES,
    sanitizeSettings,
    createUsageStats,
    recordCustomization,
    recordOptionUsage,
    nextPhase,
    formatRemaining,
  };

  if (typeof module !== "undefined" && module.exports) {
    module.exports = publicAPI;
  } else {
    globalScope.Pomodoro = publicAPI;
    if (typeof document !== "undefined") {
      createTimerApp(document);
    }
  }
})(typeof globalThis !== "undefined" ? globalThis : window);
