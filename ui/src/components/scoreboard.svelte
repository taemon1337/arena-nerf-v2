<script>
  import { currentGame, scoreboard } from '$lib/api'
  import { Heading, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Badge, Indicator } from 'flowbite-svelte';
  export let uuid;
</script>

<div>
  <div>
    <Heading tag="h1" class="mb-4 dark:text-gray-400" customSize="text-4xl font-extrabold  md:text-5xl lg:text-6xl">
      {#if uuid == "current"}
        Current Game Stats
      {:else}
        Game Stats
      {/if}
      <Badge color="green">{$currentGame.status}</Badge>
    {#if $currentGame.winner}
      <Badge class="relative m-2 p-2" color={$currentGame.winner}>
        winner: {$currentGame.winner}
        <Indicator color="{$currentGame.winner}" border size="xl" placement="top-right">
          <span class="text-white text-xs font-bold">{$currentGame.highscore}</span>
        </Indicator>
      </Badge>
      {/if}
    </Heading>
  </div>

  <Table striped={true}>
    <TableHead>
      <TableHeadCell></TableHeadCell>
      <TableHeadCell>Rank</TableHeadCell>
      <TableHeadCell>Team name</TableHeadCell>
      <TableHeadCell>Points</TableHeadCell>
    </TableHead>
    <TableBody class="divide-y">
      {#each Object.entries($scoreboard || {}) as [team, count]}
      <TableBodyRow>
        <TableBodyCell><Indicator color={team} /></TableBodyCell>
        <TableBodyCell>{team}</TableBodyCell>
        <TableBodyCell>{count}</TableBodyCell>
        <TableBodyCell></TableBodyCell>
      </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
</div>
