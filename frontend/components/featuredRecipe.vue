<script setup>
import { ref, onMounted } from 'vue'
import { useRuntimeConfig } from '#imports'
import vegs from '~/assets/css/images/vegs.png'
import { useRouter } from 'vue-router'

const loading = ref(true)
const errorMessage = ref('')
const recipes = ref([])
const router = useRouter()
const config = useRuntimeConfig()
const GRAPHQL_URL = config.public?.graphqlUrl || 'https://gebeta-app.hasura.app/v1/graphql'
const HASURA_ADMIN_SECRET = '5zghjagw4B9BsAy7grRxTOxrtrL6o5ivzSSPuXpwW4z1AXyS5hN3xo06T4HEBllw'

const QUERY = `
  query MyQuery {
    recipes(order_by: {created_at:desc}, limit: 3) {
      id
      title
      description
      cook_time_minutes
      featured_image
      recipe_ratings_aggregate {
        aggregate {
          avg {
            rating
          }
        }
      }
      user {
        full_name
      }
    }
  }
`

function truncate(text, max = 80) {
  if (!text) return ''
  return text.length > max ? text.slice(0, max) + '…' : text
}

onMounted(async () => {
  try {
    const headers = { 'Content-Type': 'application/json' }
    if (HASURA_ADMIN_SECRET) headers['x-hasura-admin-secret'] = HASURA_ADMIN_SECRET

    const res = await fetch(GRAPHQL_URL, {
      method: 'POST',
      headers,
      body: JSON.stringify({ query: QUERY })
    })
    const data = await res.json()
    if (data.errors) throw new Error(data.errors?.[0]?.message || 'Failed to load recipes')

    // Add independent like/bookmark states for each recipe
    recipes.value = (data.data?.recipes || []).map(r => ({
      ...r,
      isLiked: false,
      isBookmarked: false
    }))
  } catch (err) {
    errorMessage.value = err?.message || 'Unexpected error loading recipes'
  } finally {
    loading.value = false
  }
})

const goToRecipe = (id) => router.push(`/recipes/${id}`)
</script>

<template>
  <div class="w-full h-[97vh] bg-[#fae3cd]">
    <div class="mx-[30vw] pt-[40px]">
      <h1 class="mb-5 text-4xl ml-[7vw] font-serif font-bold">Featured Recipes</h1>
      <p class="text-lg font-semibold ml-5 text-gray-500 -mt-2">
        Discover the most loved recipe from the community
      </p>
    </div>

    <div class="ml-10 flex gap-10 justify-center mt-12">
      <div
        v-for="recipe in recipes"
        :key="recipe.id"
        class="relative card card-md bg-white w-80 h-[55vh] shadow-sm rounded-xl overflow-hidden cursor-pointer"
        @click="goToRecipe(recipe.id)"
      >
        <figure class="relative">
          <img
            class="w-full h-[180px] object-cover hover:scale-110 transition-all rounded-lg"
            :src="recipe.featured_image || vegs"
            :alt="recipe.title"
            @error="($event) => ($event.target.src = vegs)"
          />

          <div class="absolute top-2 left-2 z-10">
            <button class="w-7 h-7 transition-all" @click.stop="recipe.isLiked = !recipe.isLiked">
              <IconHeart :class="[recipe.isLiked ? 'text-red-500' : '']" class="w-5 h-5 hover:scale-110 transition-all drop-shadow"/>
            </button>
          </div>

          <div class="absolute top-2 w-8 h-8 right-2 z-10">
            <button @click.stop="recipe.isBookmarked = !recipe.isBookmarked">
              <IconBookmark :class="[recipe.isBookmarked ? 'text-blue-500' : '']" class="w-5 h-5 hover:scale-110 transition-all drop-shadow"/>
            </button>
          </div>
        </figure>

        <div class="card-body p-4 pt-2">
          <div class="flex items-center gap-1 text-sm text-yellow-500">
            <IconStar class="w-4 h-4" />
            <span class="text-xs font-medium">{{ Number(recipe.recipe_ratings_aggregate.aggregate.avg.rating ?? 0).toFixed(1) }}</span>
          </div>

          <div class="flex items-center justify-between mt-1">
            <h2 class="text-md font-bold text-gray-800">{{ recipe.title }}</h2>
            
          </div>

          <p class="text-sm text-gray-600 mt-1">{{ truncate(recipe.description, 80) }}</p>

          <div class="flex items-center justify-between mt-1 text-xs text-gray-500">
            <div class="flex items-center gap-1">
              <IconUsers class="w-4 h-4 text-gray-500" />
              <span>By {{ recipe.user?.full_name || 'Unknown' }}</span>
            </div>
            <div class="flex items-center gap-1">
              <IconClock9 class="w-4 h-4" />
              <span>{{ recipe.cook_time_minutes || 0 }} min</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="loading && recipes.length === 0" class="text-center text-orange-400 text-2xl py-10">
      Loading recipes…
    </div>
    <div v-else-if="errorMessage && recipes.length === 0" class="text-center text-red-500 py-10">
      {{ errorMessage }}
    </div>
  </div>
</template>
