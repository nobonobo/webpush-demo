<script>
  import { Toast, initializeStores } from "@skeletonlabs/skeleton";
  import { AppBar, LightSwitch } from "@skeletonlabs/skeleton";
  import "../app.postcss";
  initializeStores();
</script>

<Toast />

<AppBar
  class="mb-4"
  gridColumns="grid-cols-3"
  slotDefault="place-self-center"
  slotTrail="place-content-end"
>
  <svelte:fragment slot="lead">
    <a href="/">WebPush Demo</a></svelte:fragment
  >
  <svelte:fragment slot="trail"><LightSwitch /></svelte:fragment>
</AppBar>
<slot />
