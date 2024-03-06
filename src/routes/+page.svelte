<script>
  const vapidPublicKey =
    "BA5GSOTHKfCUOtPMYFdht02DBDOqp-bVbl84JNk4AuKBhIs9MlyItwO38-VdZt4xiiYW5QJleDCr8NCQevEUXjw";
  function urlBase64ToUint8Array(base64String) {
    const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
    const base64 = (base64String + padding)
      .replace(/\-/g, "+")
      .replace(/_/g, "/");
    const rawData = window.atob(base64);
    return Uint8Array.from([...rawData].map((char) => char.charCodeAt(0)));
  }
  async function subscribe(registration) {
    let subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(vapidPublicKey),
    });
    return subscription;
  }
  async function subscribeClick() {
    console.log("button clicked");
    navigator.serviceWorker.register("service-worker.js");
    let registration = await navigator.serviceWorker.ready;
    let subscription = await registration.pushManager.getSubscription();
    if (!subscription) {
      subscription = await subscribe(registration);
    }
    if (!subscription) {
      window.alert("subscribe failed!");
      return;
    }
    let result = await fetch("/subscribe", {
      method: "POST",
      headers: {
        "Content-Type": "application/json", // JSON形式のデータのヘッダー
      },
      body: JSON.stringify(subscription),
    });
    if (!result.ok || result.status != 200) {
      window.alert("subscribe failed!");
      return;
    }
    let res = await result.json();
    if (res.error != null) {
      window.alert(res.error);
      return;
    }
  }
</script>

<div class="container h-full mx-auto flex justify-center items-center">
  <div class="space-y-5">
    <button class="btn variant-soft-primary" on:click={subscribeClick}
      >Subscribe!</button
    >
  </div>
</div>
