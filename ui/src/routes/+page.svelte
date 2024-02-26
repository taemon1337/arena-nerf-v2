<script>
  import { uuid } from '$lib/api'
  import GamePage from '../components/gamepage.svelte'
</script>

<GamePage uuid={$uuid}></GamePage>
