<script>
	import '../app.postcss';
  import { fetchGame } from '$lib/api'
  import { beforeNavigate } from '$app/navigation'
  import { Navbar, NavBrand, NavLi, NavUl, NavHamburger, DarkMode} from 'flowbite-svelte'

  beforeNavigate(async ({ to, cancel }) => {
    if (to?.params?.uuid) {
      console.log('fetching ', to.params.uuid)
      fetchGame(to.params.uuid)
    }
  })

  let btnClass="text-gray-500 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-700 focus:outline-none focus:ring-gray-200 dark:focus:ring-gray-700 rounded-lg text-sm p-2"
</script>

<Navbar let:hidden let:toggle>
  <NavBrand href="/">
    <img
      src="https://flowbite.com/docs/images/logo.svg"
      class="mr-3 h-6 sm:h-9"
      alt="Flowbite Logo"
    />
    <span class="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
      Arena Nerf Console
    </span>
  </NavBrand>
  <NavHamburger />
  <NavUl>
    <NavLi>
      <DarkMode title="Toggle dark and light mode" class={btnClass} />
    </NavLi>
  </NavUl>
</Navbar>

<div class="container mx-auto bg-white rounded-xl overflow-y-auto h-full">
  <slot />
</div>
