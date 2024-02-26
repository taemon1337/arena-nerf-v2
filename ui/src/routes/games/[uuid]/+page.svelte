<script>
  export let data;
  import GamePage from '$src/components/gamepage.svelte'
</script>

<GamePage uuid={data.uuid}></GamePage>
