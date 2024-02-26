<script>
  import { page } from '$app/stores'
  import GameEventsTable from '$src/components/game-events-table.svelte'
  import GameHistoryTable from '$src/components/game-history-table.svelte'
  import GameTeams from '$src/components/game-teams.svelte'
  import GameNodes from '$src/components/game-nodes.svelte'
  import GameActions from '$src/components/game-actions.svelte'
  import GameStats from '$src/components/game-stats.svelte'
  import Scoreboard from '$src/components/scoreboard.svelte'
  import { onMount } from 'svelte'
  import { currentGame, gamelist, pollGame, pollGames, fetchGame } from '$lib/api'
  import { List, Li, Heading, Table, TableBody, TableBodyCell, TableBodyRow, TableHead, TableHeadCell, Badge, Indicator, Button, ButtonGroup, GradientButton } from 'flowbite-svelte';
  import { ArrowRightOutline, CheckCircleSolid } from 'flowbite-svelte-icons';

  export let uuid

  onMount(() => {
    if (uuid == "current") {
      // nothing changes except current games
      pollGame(uuid)
      pollGames()
    } else {
      fetchGame(uuid)
      fetchGame("all")
    }
  })
</script>

<div class="grid grid-rows-3 grid-cols-12 gap-4">
  <div class="row-span-3 col-span-3">
    <GameTeams />
    <GameNodes />
    <GameStats />
    <GameHistoryTable uuid={uuid} />
  </div>

  <div class="col-span-6 grid-rows-3">
    <Scoreboard uuid={uuid} />
    <GameActions />
  </div>
  <div class="row-span-2 col-span-3">
    <GameEventsTable />
  </div>
</div>
