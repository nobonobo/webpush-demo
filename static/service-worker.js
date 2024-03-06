self.addEventListener("push", (event) => {
  if (!event.data) {
    return;
  }
  let data = event.data.json();
  console.log(`[Service Worker] Push Receive data: "${JSON.stringify(data)}"`);
  let options = {
    icon: data.icon || "./favicon.png",
    body: data.body,
    tag: data.tag,
    data: data.data,
    badge: data.badge,
    requireInteraction: data.requireInteraction || false,
  };
  if (data.actions) {
    options.actions = data.actions;
  }
  if (data.lang) {
    options.lang = data.lang;
  }
  if (data.image) {
    options.image = data.image;
  }
  if (data.silent) {
    options.silent = data.silent;
  } else {
    if (data.vibrate) {
      options.vibrate = data.vibrate;
    }
  }
  event.waitUntil(
    self.registration.showNotification(
      data.title || "Push Notification",
      options
    )
  );
});
self.addEventListener("notificationclick", (event) => {
  console.log(event);
  let url = event.notification.data || location.origin + "/info.html";
  event.notification.close(); // Android needs explicit close.
  event.waitUntil(
    clients.matchAll({ type: "window" }).then((windowClients) => {
      // Check if there is already a window/tab open with the target URL
      for (var i = 0; i < windowClients.length; i++) {
        var client = windowClients[i];
        // If so, just focus it.
        if (client.url === url && "focus" in client) {
          return client.focus();
        }
      }
      // If not, then open the target URL in a new window/tab.
      if (clients.openWindow) {
        return clients.openWindow(url);
      }
    })
  );
});
