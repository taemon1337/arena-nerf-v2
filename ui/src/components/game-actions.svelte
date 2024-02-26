<script>
  import { currentGame, sendaction } from '$lib/api'
  import { Heading, Button, GradientButton } from 'flowbite-svelte';
  import { CirclePlusSolid, CirclePauseSolid, CogOutline } from 'flowbite-svelte-icons';

  let newgame = () => sendaction("ui:game:begin", "The admin has started a new game!")
  let stopgame = () => sendaction("ui:game:end", "The admin has stopped the game!")

</script>

<div class="p-2 m-2">
  <Heading class="mb-2" tag="h3">Actions</Heading>
  {#if $currentGame.completed}
  <GradientButton color="greenToBlue" on:click={newgame}>
    <CirclePlusSolid class="w-3.5 h-3.5 me-2" />
    New Game
  </GradientButton>
  {:else}
  <GradientButton color="greenToBlue" on:click={stopgame}>
    <CirclePauseSolid class="w-3.5 h-3.5 me-2" />
    Stop Game
  </GradientButton>
  {/if}
  <GradientButton color="greenToBlue">
    <CogOutline class="w-3.5 h-3.5 me-2" />
    Settings
  </GradientButton>
</div>
