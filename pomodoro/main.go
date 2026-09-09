// Pomodoro timer app
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const pageHTML = `<!doctype html>
<html lang="ja">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title>Pomodoro Pattern A</title>
  <style>
    :root {
      color-scheme: dark;
      --bg: #0b1220;
      --panel: #121a2b;
      --text: #e6edf7;
      --muted: #9fb0ce;
      --ring-bg: #22314f;
    }

    * { box-sizing: border-box; }

    body {
      margin: 0;
      min-height: 100vh;
      display: grid;
      place-items: center;
      background: radial-gradient(circle at 30% 20%, #1a2b47 0%, var(--bg) 45%, #070b12 100%);
      color: var(--text);
      font-family: "Inter", "Noto Sans JP", system-ui, sans-serif;
      overflow: hidden;
    }

    #particleCanvas {
      position: fixed;
      inset: 0;
      width: 100%;
      height: 100%;
      pointer-events: none;
      opacity: 0;
      transition: opacity 0.5s ease;
    }

    body.focus-running #particleCanvas {
      opacity: 0.55;
    }

    #rippleLayer {
      position: fixed;
      inset: 0;
      overflow: hidden;
      pointer-events: none;
    }

    .ripple {
      position: absolute;
      width: 12px;
      height: 12px;
      border-radius: 50%;
      border: 2px solid rgba(129, 225, 255, 0.5);
      transform: translate(-50%, -50%) scale(0);
      animation: ripple 2.8s ease-out forwards;
    }

    @keyframes ripple {
      0% {
        opacity: 0.9;
        transform: translate(-50%, -50%) scale(0.2);
      }
      100% {
        opacity: 0;
        transform: translate(-50%, -50%) scale(14);
      }
    }

    .panel {
      position: relative;
      z-index: 1;
      width: min(92vw, 400px);
      padding: 24px;
      border-radius: 20px;
      background: color-mix(in srgb, var(--panel), transparent 8%);
      box-shadow: 0 16px 40px rgba(0, 0, 0, 0.35);
      backdrop-filter: blur(6px);
      text-align: center;
    }

    .mode {
      margin: 0 0 16px;
      color: var(--muted);
      font-weight: 600;
      letter-spacing: 0.08em;
      text-transform: uppercase;
    }

    .timer-wrap {
      position: relative;
      width: 250px;
      height: 250px;
      margin: 0 auto 16px;
      display: grid;
      place-items: center;
    }

    .time {
      position: absolute;
      font-size: 44px;
      font-weight: 700;
      letter-spacing: 0.03em;
      font-variant-numeric: tabular-nums;
    }

    svg {
      width: 100%;
      height: 100%;
      transform: rotate(-90deg);
    }

    .track {
      fill: none;
      stroke: var(--ring-bg);
      stroke-width: 12;
    }

    .progress {
      fill: none;
      stroke: #4aa4ff;
      stroke-width: 12;
      stroke-linecap: round;
      transition: stroke 0.25s linear;
    }

    .controls {
      display: flex;
      gap: 8px;
      justify-content: center;
      flex-wrap: wrap;
    }

    button {
      border: 0;
      padding: 10px 14px;
      border-radius: 10px;
      font-weight: 600;
      cursor: pointer;
      color: #0a1424;
      background: #9bb9ff;
    }

    button.secondary {
      background: #2a3d60;
      color: var(--text);
    }
  </style>
</head>
<body>
  <canvas id="particleCanvas" aria-hidden="true"></canvas>
  <div id="rippleLayer" aria-hidden="true"></div>

  <main class="panel">
    <p class="mode" id="modeLabel">Focus Session</p>
    <div class="timer-wrap">
      <svg viewBox="0 0 250 250" aria-hidden="true">
        <circle class="track" cx="125" cy="125" r="108"></circle>
        <circle class="progress" id="progressRing" cx="125" cy="125" r="108"></circle>
      </svg>
      <div class="time" id="timeText">25:00</div>
    </div>
    <div class="controls">
      <button id="startPauseButton">開始</button>
      <button id="resetButton" class="secondary">リセット</button>
      <button id="modeButton" class="secondary">休憩モードへ</button>
    </div>
  </main>

  <script>
    const FOCUS_SECONDS = 25 * 60;
    const BREAK_SECONDS = 5 * 60;

    const modeLabel = document.getElementById('modeLabel');
    const timeText = document.getElementById('timeText');
    const progressRing = document.getElementById('progressRing');
    const startPauseButton = document.getElementById('startPauseButton');
    const resetButton = document.getElementById('resetButton');
    const modeButton = document.getElementById('modeButton');
    const rippleLayer = document.getElementById('rippleLayer');

    const canvas = document.getElementById('particleCanvas');
    const ctx = canvas.getContext('2d');

    let isFocusMode = true;
    let isRunning = false;
    let sessionSeconds = FOCUS_SECONDS;
    let remainingSeconds = sessionSeconds;
    let startedAt = 0;
    let previousRemaining = remainingSeconds;
    let displayRatio = 1;
    let rippleIntervalId = null;

    const radius = 108;
    const circumference = Math.PI * 2 * radius;
    progressRing.style.strokeDasharray = String(circumference);

    const particles = [];

    function resizeCanvas() {
      canvas.width = window.innerWidth;
      canvas.height = window.innerHeight;
      if (particles.length === 0) {
        for (let i = 0; i < 56; i += 1) {
          particles.push({
            x: Math.random() * canvas.width,
            y: Math.random() * canvas.height,
            size: 1 + Math.random() * 2.4,
            alpha: 0.12 + Math.random() * 0.26,
            velocityX: (Math.random() - 0.5) * 0.18,
            velocityY: -0.05 - Math.random() * 0.22,
          });
        }
      }
    }

    function formatTime(seconds) {
      const value = Math.max(0, Math.ceil(seconds));
      const minute = String(Math.floor(value / 60)).padStart(2, '0');
      const second = String(value % 60).padStart(2, '0');
      return minute + ':' + second;
    }

    function lerp(start, end, t) {
      return start + (end - start) * t;
    }

    function getProgressColor(ratio) {
      const clamped = Math.max(0, Math.min(1, ratio));
      if (clamped > 0.5) {
        const t = (1 - clamped) / 0.5;
        const r = Math.round(lerp(74, 255, t));
        const g = Math.round(lerp(164, 224, t));
        const b = Math.round(lerp(255, 80, t));
        return 'rgb(' + r + ', ' + g + ', ' + b + ')';
      }

      const t = clamped / 0.5;
      const r = 255;
      const g = Math.round(lerp(224, 82, 1 - t));
      const b = Math.round(lerp(80, 82, 1 - t));
      return 'rgb(' + r + ', ' + g + ', ' + b + ')';
    }

    function updateModeVisual() {
      modeLabel.textContent = isFocusMode ? 'Focus Session' : 'Break Session';
      modeButton.textContent = isFocusMode ? '休憩モードへ' : '集中モードへ';
      document.body.classList.toggle('focus-running', isFocusMode && isRunning);
    }

    function stopRipple() {
      if (rippleIntervalId !== null) {
        clearInterval(rippleIntervalId);
        rippleIntervalId = null;
      }
    }

    function startRipple() {
      stopRipple();
      if (!isFocusMode || !isRunning) {
        return;
      }

      const createRipple = () => {
        const ripple = document.createElement('div');
        ripple.className = 'ripple';
        ripple.style.left = String(30 + Math.random() * 40) + '%';
        ripple.style.top = String(30 + Math.random() * 40) + '%';
        rippleLayer.appendChild(ripple);
        setTimeout(() => ripple.remove(), 2800);
      };

      createRipple();
      rippleIntervalId = window.setInterval(createRipple, 1300);
    }

    function resetSession(keepMode = true) {
      isRunning = false;
      startPauseButton.textContent = '開始';
      if (!keepMode) {
        isFocusMode = true;
      }
      sessionSeconds = isFocusMode ? FOCUS_SECONDS : BREAK_SECONDS;
      remainingSeconds = sessionSeconds;
      previousRemaining = remainingSeconds;
      displayRatio = 1;
      updateModeVisual();
      stopRipple();
      render();
    }

    function toggleMode() {
      isFocusMode = !isFocusMode;
      resetSession(true);
    }

    function toggleRunning() {
      if (isRunning) {
        isRunning = false;
        startPauseButton.textContent = '再開';
        previousRemaining = remainingSeconds;
        stopRipple();
      } else {
        isRunning = true;
        startPauseButton.textContent = '一時停止';
        startedAt = performance.now();
        previousRemaining = remainingSeconds;
        startRipple();
      }
      updateModeVisual();
    }

    function updateTimer(now) {
      if (!isRunning) {
        return;
      }

      const elapsed = (now - startedAt) / 1000;
      remainingSeconds = Math.max(0, previousRemaining - elapsed);

      if (remainingSeconds <= 0) {
        isRunning = false;
        startPauseButton.textContent = '開始';
        stopRipple();
      }
    }

    function renderParticles() {
      ctx.clearRect(0, 0, canvas.width, canvas.height);
      if (!(isFocusMode && (isRunning || remainingSeconds < sessionSeconds))) {
        return;
      }

      particles.forEach((particle) => {
        particle.x += particle.velocityX;
        particle.y += particle.velocityY;

        if (particle.y < -10) {
          particle.y = canvas.height + 10;
          particle.x = Math.random() * canvas.width;
        }
        if (particle.x < -10) {
          particle.x = canvas.width + 10;
        }
        if (particle.x > canvas.width + 10) {
          particle.x = -10;
        }

        const gradient = ctx.createRadialGradient(
          particle.x,
          particle.y,
          0,
          particle.x,
          particle.y,
          particle.size * 6,
        );
        gradient.addColorStop(0, 'rgba(130, 225, 255, ' + particle.alpha + ')');
        gradient.addColorStop(1, 'rgba(130, 225, 255, 0)');

        ctx.fillStyle = gradient;
        ctx.beginPath();
        ctx.arc(particle.x, particle.y, particle.size * 5.8, 0, Math.PI * 2);
        ctx.fill();
      });
    }

    function render() {
      const targetRatio = remainingSeconds / sessionSeconds;
      displayRatio += (targetRatio - displayRatio) * 0.16;

      const color = getProgressColor(displayRatio);
      const offset = circumference * (1 - displayRatio);
      progressRing.style.strokeDashoffset = String(offset);
      progressRing.style.stroke = color;
      timeText.textContent = formatTime(remainingSeconds);
      renderParticles();
    }

    function tick(now) {
      updateTimer(now);
      render();
      requestAnimationFrame(tick);
    }

    startPauseButton.addEventListener('click', toggleRunning);
    resetButton.addEventListener('click', () => resetSession(true));
    modeButton.addEventListener('click', toggleMode);
    window.addEventListener('resize', resizeCanvas);

    resizeCanvas();
    resetSession(true);
    requestAnimationFrame(tick);
  </script>
</body>
</html>
`

func rootHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, pageHTML)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Pomodoro timer running at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
