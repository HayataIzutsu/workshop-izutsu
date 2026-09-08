const CACHE_NAME = "pomodoro-static-v1";
const APP_SHELL = ["/", "/style.css", "/js/format.js", "/js/pomodoroEngine.js", "/js/view.js", "/js/api.js", "/js/app.js"];

self.addEventListener("install", (event) => {
  event.waitUntil(caches.open(CACHE_NAME).then((cache) => cache.addAll(APP_SHELL)));
});

self.addEventListener("fetch", (event) => {
  if (event.request.method !== "GET") return;
  event.respondWith(caches.match(event.request).then((cached) => cached || fetch(event.request)));
});