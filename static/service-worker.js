/**
 * Service Worker
 * Base implementation of a service worker for caching static assets.
 * 
 * The goal is to cache all the static assets of the application, so that the
 * application can be used offline. We do not want to cache other user's images though.
 * User assets are contained in static/users/<userId> and should only cache current user's assets or none.
 */

// Cache name
const CACHE_NAME = 'static-cache-v1';
const urlsToCache = [
    '/',
    '/static/icons/icon-192x192.png',
    '/static/js/*.js',
    '/static/css/*.css',
];

// Cache the static assets during the install phase
self.addEventListener('install', event => {
    event.waitUntil(
        caches.open(CACHE_NAME)
            .then(cache => {
                return cache.addAll(urlsToCache);
            })
    );
});

// Serve cached content when offline
self.addEventListener('fetch', event => {
    event.respondWith(
        caches.match(event.request)
            .then(response => {
                return response || fetch(event.request);
            })
    );
});
