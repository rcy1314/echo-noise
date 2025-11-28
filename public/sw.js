const CACHE_NAME = 'noise-pwa-v3';
const ASSETS = [
  '/',
  '/favicon.ico',
  '/manifest.webmanifest'
];

self.addEventListener('install', (event) => {
  event.waitUntil(
    caches.open(CACHE_NAME).then((cache) => cache.addAll(ASSETS))
  );
});

self.addEventListener('activate', (event) => {
  event.waitUntil(
    caches.keys().then((keys) => Promise.all(keys.map((k) => {
      if (k !== CACHE_NAME) return caches.delete(k);
    })))
  );
});

self.addEventListener('fetch', (event) => {
  const req = event.request;
  const url = new URL(req.url);

  if (req.mode === 'navigate') {
    event.respondWith((async () => {
      try {
        const network = await fetch(req);
        const cache = await caches.open(CACHE_NAME);
        cache.put('/', network.clone());
        return network;
      } catch (_) {
        const cache = await caches.open(CACHE_NAME);
        const cached = await cache.match('/') || await cache.match('/index.html');
        return cached || fetch(req);
      }
    })());
    return;
  }

  if (url.pathname.startsWith('/_nuxt/')) {
    event.respondWith((async () => {
      const cache = await caches.open(CACHE_NAME);
      const cached = await cache.match(req);
      if (cached) {
        fetch(req).then((resp) => cache.put(req, resp.clone()));
        return cached;
      }
      try {
        const resp = await fetch(req);
        cache.put(req, resp.clone());
        return resp;
      } catch (_) {
        return cached || Response.error();
      }
    })());
    return;
  }

  if (url.pathname.startsWith('/api/')) {
    event.respondWith(fetch(req));
    return;
  }

  if (req.method === 'GET' && url.origin === self.location.origin) {
    event.respondWith((async () => {
      const cache = await caches.open(CACHE_NAME);
      const cached = await cache.match(req);
      if (cached) {
        fetch(req).then((resp) => cache.put(req, resp.clone()));
        return cached;
      }
      try {
        const resp = await fetch(req);
        cache.put(req, resp.clone());
        return resp;
      } catch (_) {
        return cached || Response.error();
      }
    })());
    return;
  }
});
