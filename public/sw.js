const CACHE_NAME = "alarm-color-v1";
const urlsToCache = [
  "/", 
  "/index.html",
  "/manifest.json",
  "/assets/index-[hash].js", 
  "/assets/index-[hash].css",
  "/icon-192.png",
  "/icon-512.png",
];

//install chace
self.addEventListener("install", (event) => {
  event.waitUntil(
    caches
      .open(CACHE_NAME)
      .then((cache) => {
        console.log("Opened cache");
        return cache.addAll(urlsToCache);
      })
      .catch((error) => {
        console.log("Cache addAll failed:", error);
      })
  );
  self.skipWaiting();
});

// activate new cache
self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches.keys().then((cacheNames) => {
      return Promise.all(
        cacheNames.map((cacheName) => {
          if (cacheName !== CACHE_NAME) {
            console.log("Deleting old cache:", cacheName);
            return caches.delete(cacheName);
          }
        })
      );
    })
  );
  // Force active new sw
  self.clients.claim();
});

// Fetch: cache first, then network
self.addEventListener("fetch", (event) => {
  // Get Requests
  if (event.request.method !== "GET") {
    return;
  }

  event.respondWith(
    caches.open(CACHE_NAME).then((cache) => {
      return cache.match(event.request).then((cachedResponse) => {
        // if cached, return
        if (cachedResponse) {
          return cachedResponse;
        }

        // if no network get from cache
        return fetch(event.request)
          .then((networkResponse) => {
            // on success clone and cache
            if (networkResponse && networkResponse.status === 200) {
              const responseToCache = networkResponse.clone();
              cache.put(event.request, responseToCache);
            }
            return networkResponse;
          })
          .catch(() => {
            //on fail show online page
            console.log("Fetch failed; returning offline content");
            return caches.match("/"); // fallback to index.html
          });
      });
    })
  );
});

// Push notifications
self.addEventListener("push", (event) => {
  const options = {
    body: event.data ? event.data.text() : "آلارم فعال شد!",
    icon: "/icon-192.png",
    badge: "/icon-192.png",
    vibrate: [100, 50, 100],
    data: { dateOfArrival: Date.now(), primaryKey: 1 },
  };
  event.waitUntil(self.registration.showNotification("آلارم!", options));
});
