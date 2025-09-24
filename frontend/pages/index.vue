<script setup>
import { Search,FeaturedRecipe,AvatarUser } from '#components';
</script>


<template>
  <div class="pt-[70px]">

    <!-- // hero section -->
    <div  class="hero min-h-screen" style="background-image: url('/image/hero2.png');">
        <div class="hero-overlay"></div>
        <div class="hero-content text-neutral-content text-center">
          <div class="max-w-[90vw]">
            <h1 class="mb-5 text-5xl font-bold">Find Your Next Favorite Recipe</h1>
            <p class="mb-5">
              Browse,cook and share amazing meals from everywhere
            </p>
            <NuxtLink to="/recipes" class="btn btn-primary bg-orange-500 border-orange-400">Get Started</NuxtLink>
          </div>
        </div>
    </div>

    <!-- search and catagories -->
    <div><search/></div>

    <!-- featured recipes -->
    <div><FeaturedRecipe/></div>

    <!-- user avatar  -->

    <div><AvatarUser/></div>
  </div>
   

  </template>
  