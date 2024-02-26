<script>
  import { currentGame } from '$lib/api'
  import { Badge } from 'flowbite-svelte'
</script>

<div class="p-2 m-2">
  <p>
    Started <Badge color="gray">{$currentGame.start_at}</Badge>
  </p>
  {#if $currentGame.completed}
  <p>
    Ended <Badge color="gray">{$currentGame.end_at}</Badge>
  </p>
  {:else}
    <p>
      <Badge color="green">In progress</Badge>
    </p>
  {/if}
</div>
