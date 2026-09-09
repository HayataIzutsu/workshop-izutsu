(function () {
  var storageProfileKey = "pomodoro.profile";
  var storageSessionsKey = "pomodoro.sessions";
  var game = window.PomodoroGame;

  var profileEl = {
    level: document.getElementById("levelValue"),
    xp: document.getElementById("xpValue"),
    streak: document.getElementById("streakValue"),
    levelUp: document.getElementById("levelUpMessage"),
  };

  var badgeList = document.getElementById("badgeList");
  var weeklyStats = document.getElementById("weeklyStats");
  var monthlyStats = document.getElementById("monthlyStats");

  var timerDisplay = document.getElementById("timerDisplay");
  var startButton = document.getElementById("startTimerButton");
  var completeButton = document.getElementById("completeSessionButton");

  var timerSeconds = 25 * 60;
  var intervalId = null;

  function defaultProfile() {
    return {
      xp: 0,
      level: 1,
      streak: 0,
      totalCompletedSessions: 0,
      badges: [],
    };
  }

  function loadProfile() {
    try {
      return JSON.parse(localStorage.getItem(storageProfileKey)) || defaultProfile();
    } catch (error) {
      return defaultProfile();
    }
  }

  function loadSessions() {
    try {
      return JSON.parse(localStorage.getItem(storageSessionsKey)) || [];
    } catch (error) {
      return [];
    }
  }

  function save(profile, sessions) {
    localStorage.setItem(storageProfileKey, JSON.stringify(profile));
    localStorage.setItem(storageSessionsKey, JSON.stringify(sessions));
  }

  function formatTimeline(timeline) {
    var days = Object.keys(timeline).sort();
    if (days.length === 0) {
      return "<p>データなし</p>";
    }

    return (
      "<ul>" +
      days
        .map(function (day) {
          return "<li>" + day + ": " + timeline[day].completed + "回 / " + timeline[day].focusMinutes + "分</li>";
        })
        .join("") +
      "</ul>"
    );
  }

  function renderStats(name, targetEl, stats) {
    targetEl.innerHTML =
      "<p>完了数: " +
      stats.completedSessions +
      " / " +
      stats.plannedSessions +
      "</p>" +
      "<p>完了率: " +
      stats.completionRate +
      "%</p>" +
      "<p>平均集中時間: " +
      stats.averageFocusMinutes +
      "分</p>" +
      "<h3>時系列(" +
      name +
      ")</h3>" +
      formatTimeline(stats.timeline);
  }

  function render(profile, sessions) {
    profileEl.level.textContent = String(profile.level);
    profileEl.xp.textContent = String(profile.xp);
    profileEl.streak.textContent = String(profile.streak);

    var badgeNames = {
      "three-day-streak": "3日連続",
      "weekly-10": "今週10回完了",
      "thirty-sessions": "30セッション達成",
      "seven-day-streak-reward": "7日ストリークリワード",
    };

    badgeList.innerHTML = profile.badges.length
      ? profile.badges
          .map(function (badgeKey) {
            return '<li class="badge-item">🏅 ' + (badgeNames[badgeKey] || badgeKey) + "</li>";
          })
          .join("")
      : "<li>まだバッジはありません</li>";

    var weekly = game.buildPeriodStats(sessions, 7, new Date());
    var monthly = game.buildPeriodStats(sessions, 30, new Date());
    renderStats("週", weeklyStats, weekly);
    renderStats("月", monthlyStats, monthly);
  }

  function refreshTimerDisplay() {
    var minutes = Math.floor(timerSeconds / 60)
      .toString()
      .padStart(2, "0");
    var seconds = (timerSeconds % 60).toString().padStart(2, "0");
    timerDisplay.textContent = minutes + ":" + seconds;
  }

  function startTimer() {
    if (intervalId) {
      return;
    }

    timerSeconds = 25 * 60;
    refreshTimerDisplay();
    completeButton.disabled = true;

    intervalId = setInterval(function () {
      timerSeconds -= 1;
      refreshTimerDisplay();

      if (timerSeconds <= 0) {
        clearInterval(intervalId);
        intervalId = null;
        completeButton.disabled = false;
      }
    }, 1000);
  }

  function completeSession() {
    var profile = loadProfile();
    var sessions = loadSessions();
    var result = game.recordCompletedSession(profile, sessions, {
      focusMinutes: 25,
      plannedMinutes: 25,
      now: new Date(),
    });

    save(result.profile, result.sessions);
    render(result.profile, result.sessions);

    if (result.leveledUp) {
      profileEl.levelUp.textContent = "レベルアップ！ Lv." + result.profile.level;
      profileEl.levelUp.classList.remove("hidden");
    } else {
      profileEl.levelUp.classList.add("hidden");
    }

    if (result.unlockedBadges.length > 0) {
      alert("新しいバッジを獲得: " + result.unlockedBadges.join(" / "));
    }

    completeButton.disabled = true;
  }

  startButton.addEventListener("click", startTimer);
  completeButton.addEventListener("click", completeSession);
  completeButton.disabled = true;

  var currentProfile = loadProfile();
  var currentSessions = loadSessions();
  render(currentProfile, currentSessions);
  refreshTimerDisplay();
})();
