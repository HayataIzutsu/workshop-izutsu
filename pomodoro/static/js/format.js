function formatTime(totalSeconds) {
  const safeSeconds = Math.max(0, Math.floor(totalSeconds));
  const minutes = Math.floor(safeSeconds / 60);
  const seconds = safeSeconds % 60;
  return `${String(minutes).padStart(2, "0")}:${String(seconds).padStart(2, "0")}`;
}

function progressPercent(remainingSeconds, totalSeconds) {
  if (totalSeconds <= 0) return 0;
  return Math.max(0, Math.min(100, (remainingSeconds / totalSeconds) * 100));
}

function formatDuration(totalSeconds) {
  const minutes = Math.floor(Math.max(0, totalSeconds) / 60);
  if (minutes < 60) return `${minutes}分`;
  const hours = Math.floor(minutes / 60);
  const restMinutes = minutes % 60;
  return restMinutes === 0 ? `${hours}時間` : `${hours}時間${restMinutes}分`;
}

if (typeof module !== "undefined" && module.exports) {
  module.exports = { formatTime, progressPercent, formatDuration };
}