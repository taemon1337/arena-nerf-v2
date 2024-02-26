<script>
  import { gameevents } from '$lib/api'
  import { Table, TableHead, TableBody, TableBodyCell, TableBodyRow, TableHeadCell } from 'flowbite-svelte';
</script>

<Table noborder={true} shadow striped={true}>
  <TableHead>
    <TableHeadCell>Game Event Stream</TableHeadCell>
  </TableHead>
  <TableBody>
    {#each [...($gameevents || [])].reverse() as event}
      <TableBodyRow>
        <TableBodyCell tdClass="px-6 py-0 font-small">{event}</TableBodyCell>
      </TableBodyRow>
    {/each}
  </TableBody>
</Table>
