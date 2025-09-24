<script setup>
import { ref, onMounted } from 'vue'
import { useRuntimeConfig } from '#imports'
// Use existing auto-registered components: <search />
import vegs from '~/assets/css/images/vegs.png'
import { useRouter } from 'vue-router'


const loading = ref(true)
const errorMessage = ref('')
const recipes = ref([])
const router= useRouter()
const config = useRuntimeConfig()
const GRAPHQL_URL = config.public?.graphqlUrl || 'https://gebeta-app.hasura.app/v1/graphql'
const HASURA_ADMIN_SECRET = '5zghjagw4B9BsAy7grRxTOxrtrL6o5ivzSSPuXpwW4z1AXyS5hN3xo06T4HEBllw'
const QUERY = `
  query MyQuery {
    recipes(order_by: {created_at:desc}) {
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
    if (HASURA_ADMIN_SECRET) {
      headers['x-hasura-admin-secret'] = HASURA_ADMIN_SECRET
    }
    const res = await fetch(GRAPHQL_URL, {
      method: 'POST',
      headers,
      body: JSON.stringify({ query: QUERY })
    })
    const data = await res.json()
    if (data.errors) {
      throw new Error(data.errors?.[0]?.message || 'Failed to load recipes')
    }
    recipes.value = data.data?.recipes || []
  } catch (err) {
    errorMessage.value = err?.message || 'Unexpected error loading recipes'
  } finally {
    loading.value = false
  }
})
const goToRecipe = (id) => {
  router.push(`/recipes/${id}`)  // assuming your single recipe page is at /recipes/[id].vue
}
</script>

<template>
  <div class="pt-[70px] pb-[20px] min-h-screen" style="background-color:#FAE3CD;">
    <div class="max-w-[1200px] mx-auto px-4">
      <div class="text-center mb-6">
        <h1 class="text-3xl font-bold">All Recipes</h1>
        <p class="text-muted">Browse meals shared by food lovers around the world</p>
      </div>

      <!-- Broad search: full width on page container -->
      <div class="mb-6">
        <div class="max-w-none w-full">
          <search />
        </div>
      </div>

      <!-- Loading / Error states -->
      <div v-if="loading && recipes.length === 0" class="text-center text-orange-400 text-2xl  py-10 "><IconLoaderCircle class="w-16 h-16 mx-auto mb-4 animate-spin" /></div>
      <div v-else-if="errorMessage && recipes.length === 0" class="text-center text-red-500 py-10">{{ errorMessage }}</div>

      <!-- Smaller Cards Grid -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div v-for="recipe in recipes" 
        :key="recipe.id" 
        @click="goToRecipe(recipe.id)"
        class="card bg-white shadow-sm hover:shadow-md transition-shadow">
          <figure class="relative">
            <img
              class="w-full h-40 object-cover hover:scale-110  transition-all rounded-lg" 
              :src="recipe.featured_image || vegs"
              :alt="recipe.title"
              @error="($event) => ($event.target.src = vegs)"
            />
            <span class="badge badge-warning absolute left-2 top-2 text-xs">{{ recipe.cook_time_minutes || 0 }} min</span>
          </figure>
          <div class="card-body p-3">
            <h2 class="card-title text-sm hover:text-orange-400 transition-all font-bold">{{ recipe.title }}</h2>
            <p class="text-xs text-muted">{{ truncate(recipe.description, 80) }}</p>
            <div class="flex items-center justify-between mt-1">
              <div class="flex items-center gap-1">
                <span class="text-yellow-500 text-sm">★★★★★</span>
                <span class="text-xs">{{ Number(recipe.recipe_ratings_aggregate.aggregate.avg.rating ?? 0).toFixed(1) }}</span>
              </div>
              <div class="text-xs text-muted font-semibold">{{ recipe.user?.full_name || 'Unknown' }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
