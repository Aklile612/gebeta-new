<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRuntimeConfig } from '#imports'

const route = useRoute()
const userId = route.params.id

const loading = ref(true)
const errorMessage = ref('')
const recipes = ref([])
const config = useRuntimeConfig()

const GRAPHQL_URL = config.public?.graphqlUrl || 'https://gebeta-app.hasura.app/v1/graphql'
const HASURA_ADMIN_SECRET = '5zghjagw4B9BsAy7grRxTOxrtrL6o5ivzSSPuXpwW4z1AXyS5hN3xo06T4HEBllw'

const QUERY = `
  query RecipesByUser($userId: uuid!) {
    recipes(where: { user_id: { _eq: $userId } }, order_by: { created_at: desc }) {
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

onMounted(async () => {
  try {
    const headers = {
      'Content-Type': 'application/json',
      'x-hasura-admin-secret': HASURA_ADMIN_SECRET,
    }

    const res = await fetch(GRAPHQL_URL, {
      method: 'POST',
      headers,
      body: JSON.stringify({ query: QUERY, variables: { userId } }),
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
</script>

<template>
  <div class="bg-[#fae3cd] min-h-screen  mx-auto w-full  px-6  pt-[70px]">
    <div class="mb-6">
        <NuxtLink to="/" class="flex items-center text-gray-700 hover:text-orange-500 mb-2">
          <IconChevronLeft class="mr-1" />
          <span class="text-sm font-medium">Back</span>
        </NuxtLink>
        <h1 class="text-3xl font-bold mb-6">
        Recipes by {{ recipes[0]?.user?.full_name || 'User' }}
        </h1>
      </div>
    

    <!-- Loading & Error states -->
    <div v-if="loading" class="text-center py-10">
      <IconLoaderCircle class="w-12 h-12 animate-spin mx-auto text-orange-400" />
    </div>
    <div v-else-if="errorMessage" class="text-center text-red-500 py-10">
      {{ errorMessage }}
    </div>

    <!-- Recipes Grid -->
    <div v-else class="grid grid-cols-1  justify-center sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="recipe in recipes"
        :key="recipe.id"
        @click="$router.push(`/recipes/${recipe.id}`)"
        class="card bg-white shadow-sm hover:shadow-md transition-shadow cursor-pointer"
      >
        <figure class="relative">
          <img
            class="w-full h-40 object-cover hover:scale-110 transition-all rounded-lg"
            :src="recipe.featured_image || '/fallback.jpg'"
            :alt="recipe.title"
          />
          <span
            class="badge badge-warning absolute left-2 top-2 text-xs"
          >{{ recipe.cook_time_minutes || 0 }} min</span>
        </figure>
        <div class="card-body p-3">
          <h2
            class="card-title text-sm hover:text-orange-400 transition-all font-bold"
          >
            {{ recipe.title }}
          </h2>
          <p class="text-xs text-muted">
            {{ recipe.description?.slice(0, 80) }}
          </p>
          <div class="flex items-center justify-between mt-1">
            <div class="flex items-center gap-1">
              <span class="text-yellow-500 text-sm">★★★★★</span>
              <span class="text-xs">{{
                Number(
                  recipe.recipe_ratings_aggregate.aggregate.avg.rating ?? 0
                ).toFixed(1)
              }}</span>
            </div>
            <div class="text-xs text-muted font-semibold">
              {{ recipe.user?.full_name || 'Unknown' }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
