<script>
  import { scoreboard } from '$lib/api'
  import { Heading, Badge, Indicator } from 'flowbite-svelte';
</script>

<div>
  <Heading tag="h3" class="dark:text-gray-400">Game Teams</Heading>
  {#each Object.entries($scoreboard || {}) as [team, count]}
    <Badge class="relative m-4 p-4" color={team}>
      {team}
      <Indicator color="{team}" border size="xl" placement="top-right">
        <span class="text-white text-xs font-bold">{count}</span>
      </Indicator>
    </Badge>
  {/each}
</div>
