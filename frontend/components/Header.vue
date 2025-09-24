<script setup>

const route= useRoute()

</script>


<template>
    <header>
      <div class="navbar fixed top-0 left-0 w-full bg-transparent backdrop-blur-md z-50 transition-all flex justify-between duration-300">
        <NuxtLink to="/" class="flex  gap-2 items-center hover:scale-105 md:ml-7 ml-3  text-center justify-center flex-row ">
          <div class="text-white font-bold w-8 h-8 rounded-sm pl-1 pt-1 bg-orange-500"><IconUtensils /></div>
          <div class=" font-mono text-xl ">ገበታ</div>
        </NuxtLink>
  <div class="flex gap-2">
    <!-- <input type="text" placeholder="Search" class="input input-bordered w-24 md:w-auto" /> -->
     <div class="flex gap-6 mr-20 md:mt-2 text-lg">
    <NuxtLink
      to="/recipes"
      class="nav-link md:mt-1"
      :class="{ 'active-link': route.path === '/recipes' }"
    >Recipes</NuxtLink>

    <NuxtLink to="/addrecipe"
    
      class="btn btn-dash btn-secondary border-orange-300 hover:bg-orange-500 hover:text-black text-orange-400"
    >
      Create Recipe
    </NuxtLink>

    <NuxtLink
      to="/creators"
      class="nav-link md:mt-1"
      :class="{ 'active-link': route.path === '/creators' }"
    >Creators</NuxtLink>
  </div>
    <div class="dropdown dropdown-end mr-5">
      <div tabindex="0" role="button" class="btn btn-ghost btn-circle avatar">
        <div class="w-10 rounded-full">
          <img
            alt="Tailwind CSS Navbar component"
            src="https://img.daisyui.com/images/stock/photo-1534528741775-53994a69daeb.webp" />
        </div>
      </div>
      <ul
        tabindex="0"
        class="menu menu-sm dropdown-content bg-orange-300  rounded-box z-1 mt-3 w-52 p-2 shadow">
        <li class="hover:bg-orange-600 rounded-box">
          <a class="justify-between">
            Profile
          </a>
        </li>
        
        <li class="hover:bg-orange-600 rounded-box"><NuxtLink to="/login">Logout</NuxtLink></li>
      </ul>
    </div>
  </div>
</div>
    </header>
</template>
<style scoped>
@reference "tailwindcss";
.nav-link {
  @apply relative text-sm font-bold transition-all duration-300;
}

.nav-link::after {
  content: '';
  @apply absolute left-0 bottom-[-1px] w-full h-[3px] bg-orange-400 opacity-0 transition-all duration-100;
}

.nav-link:hover::after {
  opacity: 100;
}

.active-link::after {
  @apply absolute left-0 bottom-[-1px] w-full h-[3px] bg-orange-600 blur-[1px] opacity-100;
}
</style>