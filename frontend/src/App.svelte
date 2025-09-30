<script lang="ts">
  // import "./app.css";
  import "./style.css";
  import {
    GitHubAPIService,
    Notification,
  } from "../bindings/github.com/jcelaya775/github-notifier";
  import { Events } from "@wailsio/runtime";
  import { onMount } from "svelte";

  let notifications = $state<Notification[]>();

  function handleKeydown(event: KeyboardEvent) {
    console.log(event.key);
    if (event.key === "Escape") {
      Events.Emit("escape-pressed");
    }
  }

  onMount(async () => {
    notifications = await GitHubAPIService.GetNotifications();
  });

  $inspect(notifications)
</script>

<!--<svelte:document on:keydown={handleKeydown} />-->

<main class="p-4">
  <h2 class="text-2xl mb-8" style="font-family: 'Hubot Sans'">
    Github Notifier
  </h2>

  <div>
    {#each notifications as notification}
      <p>{notification.subject.title}</p>
    {/each}
  </div>
</main>

<style></style>
