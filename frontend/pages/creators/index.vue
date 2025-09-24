<script setup>
import { ref, onMounted } from 'vue'
import { useRuntimeConfig } from '#imports'

const loading = ref(true)
const errorMessage = ref('')
const users = ref([])
const config = useRuntimeConfig()

const GRAPHQL_URL = config.public?.graphqlUrl || 'https://gebeta-app.hasura.app/v1/graphql'
const HASURA_ADMIN_SECRET = '5zghjagw4B9BsAy7grRxTOxrtrL6o5ivzSSPuXpwW4z1AXyS5hN3xo06T4HEBllw'

const QUERY = `
  query AllRecipeCreators {
    users(order_by: { full_name: asc }) {
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
      throw new Error(data.errors?.[0]?.message || 'Failed to load users')
    }
    users.value = data.data?.users || []
  } catch (err) {
    errorMessage.value = err?.message || 'Unexpected error loading users'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="pt-[70px] pb-[20px] min-h-screen" style="background-color:#FAE3CD;">
    <div class="max-w-[1200px] mx-auto px-4">
      <div class="text-center mb-6">
        <h1 class="text-3xl font-bold">All Recipe Creators</h1>
        <p class="text-muted">Browse contributors from our community</p>
      </div>

      <div v-if="loading && users.length === 0" class="text-center text-orange-400 text-2xl py-10">
        Loading creators…
      </div>
      <div v-else-if="errorMessage && users.length === 0" class="text-center text-red-500 py-10">
        {{ errorMessage }}
      </div>

      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <div 
          v-for="user in users" 
          :key="user.id" 
          class="card bg-white shadow-sm hover:shadow-md transition-shadow cursor-pointer"
        >
          <div class="card-body p-4 text-center">
            <div class="w-16 h-16 bg-orange-500 rounded-full flex items-center justify-center mx-auto mb-3 text-white font-bold text-xl">
              {{ user.full_name ? user.full_name.charAt(0).toUpperCase() : 'U' }}
            </div>
            
            <h2 class="card-title justify-center text-sm font-bold">{{ user.full_name || 'Unknown User' }}</h2>
            <p class="text-xs text-muted">
              {{ user.recipes_aggregate?.aggregate?.count || 0 }} Recipes
            </p>
            
            <button class="mt-2 text-sm bg-orange-300 hover:bg-orange-500 text-gray-800 hover:text-white font-medium py-1 px-3 rounded-md transition-all">
              View Recipes
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>