(function (root, factory) {
  if (typeof module === "object" && module.exports) {
    module.exports = factory();
    return;
  }
  root.PomodoroGame = factory();
})(typeof self !== "undefined" ? self : this, function () {
  var XP_PER_COMPLETED_SESSION = 25;
  var LEVEL_XP_STEP = 100;

  function toDayKey(date) {
    return date.toISOString().slice(0, 10);
  }

  function parseDate(value) {
    return value instanceof Date ? value : new Date(value);
  }

  function calculateLevel(xp) {
    return Math.floor(Math.max(xp, 0) / LEVEL_XP_STEP) + 1;
  }

  function calculateStreak(sessions) {
    var completedDaySet = new Set();
    sessions
      .filter(function (session) {
        return session.completed;
      })
      .forEach(function (session) {
        completedDaySet.add(toDayKey(parseDate(session.completedAt)));
      });

    var days = Array.from(completedDaySet).sort();
    if (days.length === 0) {
      return 0;
    }

    var streak = 1;
    for (var i = days.length - 1; i > 0; i--) {
      var current = new Date(days[i] + "T00:00:00Z");
      var previous = new Date(days[i - 1] + "T00:00:00Z");
      var diffDays = Math.round((current - previous) / 86400000);
      if (diffDays === 1) {
        streak += 1;
        continue;
      }
      break;
    }

    return streak;
  }

  function countCompletedInLastDays(sessions, days, now) {
    var baseDate = parseDate(now);
    var start = new Date(baseDate);
    start.setUTCHours(0, 0, 0, 0);
    start.setUTCDate(start.getUTCDate() - (days - 1));
    var startKey = toDayKey(start);
    var endKey = toDayKey(baseDate);

    return sessions.filter(function (session) {
      if (!session.completed) {
        return false;
      }
      var key = toDayKey(parseDate(session.completedAt));
      return key >= startKey && key <= endKey;
    }).length;
  }

  function evaluateBadges(profile, sessions, now) {
    var existing = new Set(profile.badges || []);
    var unlocked = [];

    if (profile.streak >= 3 && !existing.has("three-day-streak")) {
      existing.add("three-day-streak");
      unlocked.push("3日連続バッジ");
    }

    if (countCompletedInLastDays(sessions, 7, now) >= 10 && !existing.has("weekly-10")) {
      existing.add("weekly-10");
      unlocked.push("今週10回完了バッジ");
    }

    if (profile.totalCompletedSessions >= 30 && !existing.has("thirty-sessions")) {
      existing.add("thirty-sessions");
      unlocked.push("30セッション達成バッジ");
    }

    if (profile.streak >= 7 && !existing.has("seven-day-streak-reward")) {
      existing.add("seven-day-streak-reward");
      unlocked.push("7日ストリークリワード");
    }

    return {
      badges: Array.from(existing),
      unlocked: unlocked,
    };
  }

  function buildPeriodStats(sessions, days, now) {
    var baseDate = parseDate(now);
    var start = new Date(baseDate);
    start.setUTCHours(0, 0, 0, 0);
    start.setUTCDate(start.getUTCDate() - (days - 1));
    var startKey = toDayKey(start);
    var endKey = toDayKey(baseDate);

    var periodSessions = sessions.filter(function (session) {
      var key = toDayKey(parseDate(session.completedAt));
      return key >= startKey && key <= endKey;
    });

    var completedSessions = periodSessions.filter(function (session) {
      return session.completed;
    });

    var totalFocusMinutes = completedSessions.reduce(function (sum, session) {
      return sum + (session.focusMinutes || 0);
    }, 0);

    var byDay = {};
    completedSessions.forEach(function (session) {
      var key = toDayKey(parseDate(session.completedAt));
      byDay[key] = byDay[key] || { completed: 0, focusMinutes: 0 };
      byDay[key].completed += 1;
      byDay[key].focusMinutes += session.focusMinutes || 0;
    });

    return {
      plannedSessions: periodSessions.length,
      completedSessions: completedSessions.length,
      completionRate:
        periodSessions.length === 0 ? 0 : Math.round((completedSessions.length / periodSessions.length) * 1000) / 10,
      averageFocusMinutes:
        completedSessions.length === 0 ? 0 : Math.round((totalFocusMinutes / completedSessions.length) * 10) / 10,
      timeline: byDay,
    };
  }

  function recordCompletedSession(profile, sessions, options) {
    var now = (options && options.now) || new Date();
    var focusMinutes = (options && options.focusMinutes) || 25;
    var plannedMinutes = (options && options.plannedMinutes) || 25;

    var createdSession = {
      completed: true,
      completedAt: parseDate(now).toISOString(),
      focusMinutes: focusMinutes,
      plannedMinutes: plannedMinutes,
    };

    var nextSessions = sessions.concat(createdSession);

    var previousLevel = profile.level || calculateLevel(profile.xp || 0);
    var nextXP = (profile.xp || 0) + XP_PER_COMPLETED_SESSION;
    var nextLevel = calculateLevel(nextXP);
    var totalCompletedSessions = (profile.totalCompletedSessions || 0) + 1;
    var streak = calculateStreak(nextSessions);

    var nextProfile = {
      xp: nextXP,
      level: nextLevel,
      streak: streak,
      totalCompletedSessions: totalCompletedSessions,
      badges: profile.badges || [],
    };

    var badgeResult = evaluateBadges(nextProfile, nextSessions, now);
    nextProfile.badges = badgeResult.badges;

    return {
      profile: nextProfile,
      sessions: nextSessions,
      leveledUp: nextLevel > previousLevel,
      unlockedBadges: badgeResult.unlocked,
      weeklyStats: buildPeriodStats(nextSessions, 7, now),
      monthlyStats: buildPeriodStats(nextSessions, 30, now),
    };
  }

  return {
    XP_PER_COMPLETED_SESSION: XP_PER_COMPLETED_SESSION,
    calculateLevel: calculateLevel,
    calculateStreak: calculateStreak,
    evaluateBadges: evaluateBadges,
    buildPeriodStats: buildPeriodStats,
    recordCompletedSession: recordCompletedSession,
  };
});
