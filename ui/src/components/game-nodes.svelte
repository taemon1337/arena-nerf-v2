<script>
  import { nodeboard } from '$lib/api'
  import { Heading, Badge, Indicator } from 'flowbite-svelte';
</script>

<div>
  <Heading tag="h3" class="dark:text-gray-400">Game Nodes</Heading>
  {#each Object.entries($nodeboard || {}) as [node, count]}
    <Badge class="p-4 m-4 relative">
      {node}
      <Indicator border size="xl" placement="top-right">
        <span class="text-white text-xs font-bold">{count}</span>
      </Indicator>
    </Badge>
  {/each}
</div>
