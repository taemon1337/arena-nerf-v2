<script>
  import { gamelist } from '$lib/api'
  import { Table, TableHead, TableBody, TableBodyCell, TableBodyRow, TableHeadCell, Li, List } from 'flowbite-svelte';
  import { CheckCircleSolid, CheckCircleOutline } from 'flowbite-svelte-icons';
  export let uuid
</script>

<Table noborder={true} shadow striped={true}>
  <TableHead>
    <TableHeadCell>Game History</TableHeadCell>
  </TableHead>
  <TableBody>
    {#each $gamelist as game}
      <List> 
        <Li icon>
          {#if game == uuid}
          <CheckCircleSolid class="w-3.5 h-3.5 me-2 text-green-500 dark:text-green-400" />
          {:else}
          <CheckCircleOutline class="w-3.5 h-3.5 me-2" />
          {/if}
          <a href="/games/{game}">{game}</a>
        </Li>
      </List>
    {/each}
  </TableBody>
</Table>
