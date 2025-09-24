<script setup>
import { ref, onMounted } from 'vue'
import { useRuntimeConfig } from '#imports'
import { useRouter } from 'vue-router'

const router= useRouter()
const loading = ref(true)
const errorMessage = ref('')
const topUsers = ref([])
const config = useRuntimeConfig()

const GRAPHQL_URL = config.public?.graphqlUrl || 'https://gebeta-app.hasura.app/v1/graphql'
const HASURA_ADMIN_SECRET = '5zghjagw4B9BsAy7grRxTOxrtrL6o5ivzSSPuXpwW4z1AXyS5hN3xo06T4HEBllw'

// Query: top 4 contributors ordered by recipe count
const QUERY = `
  query TopRecipeCreators {
    users(order_by: {recipes_aggregate: {count: desc}}, limit: 4) {
      id
      full_name
      recipes_aggregate {
        aggregate {
          count
        }
      }
    }
  }

`

onMounted(async () => {
  try {
    const headers = { 
      'Content-Type': 'application/json',
      'x-hasura-admin-secret': HASURA_ADMIN_SECRET 
    }

    const res = await fetch(GRAPHQL_URL, {
      method: 'POST',
      headers,
      body: JSON.stringify({ query: QUERY })
    })
    
    const data = await res.json()
    if (data.errors) {
      throw new Error(data.errors?.[0]?.message || 'Failed to load top contributors')
    }
    topUsers.value = data.data?.users || []
  } catch (err) {
    errorMessage.value = err?.message || 'Unexpected error loading top contributors'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="bg-white w-full h-[70vh]">
    <div class="mx-[30vw] pt-[70px]">
      <h1 class="mb-5 text-4xl ml-[7vw] font-serif font-bold">Top Recipe Creators</h1>
    </div>

    <!-- Loading & Error states -->
    <div v-if="loading && topUsers.length === 0" class="text-center text-orange-400 text-2xl py-10">
      <IconLoaderCircle class="w-10 h-10 mx-auto mb-4 animate-spin" />
    </div>
    <div v-else-if="errorMessage && topUsers.length === 0" class="text-center text-red-500 py-10">
      {{ errorMessage }}
    </div>

    <!-- Contributors list -->
    <div v-else class="flex gap-9 mt-8 justify-center">
      <div
        v-for="user in topUsers"
        :key="user.id"
        class="w-[240px] p-4 bg-[#fae3cd] rounded-xl flex flex-col items-center shadow-sm"
      >
        <!-- Icon in circle -->
        <div class="w-16 h-16 bg-white hover:shadow-slate-300 hover:scale-110 rounded-full flex items-center justify-center shadow-md">
          <IconUsers class="text-orange-500 w-8 h-8" />
        </div>

        <!-- User info -->
        <div class="mt-4 text-center">
          <p class="font-bold text-sm text-gray-800">{{ user.full_name || 'Unknown User' }}</p>
          <p class="text-xs text-gray-600 mt-1">
            {{ user.recipes_aggregate?.aggregate?.count || 0 }} Recipes Submitted
          </p>

          <!-- View Button -->
          <button
            @click="$router.push(`/creators/${user.id}`)"
            class="mt-3 text-sm bg-orange-300 hover:bg-orange-500 text-gray-800 hover:text-white font-medium py-1 px-3 rounded-md transition-all duration-200"
          >
            View Recipes
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
