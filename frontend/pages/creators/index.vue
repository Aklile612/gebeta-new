<script setup>
import { ref, onMounted } from 'vue'
import { useRuntimeConfig } from '#imports'
import { useRouter } from 'vue-router'


const loading = ref(true)
const errorMessage = ref('')
const users = ref([])
const config = useRuntimeConfig()
const router = useRouter()
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
  <div class="bg-white w-full min-h-screen">
    <div class="mx-[30vw] pt-[70px]">
      <h1 class="mb-5 text-4xl ml-[7vw] font-serif font-bold">All Recipe Creators</h1>
    </div>

    <!-- Loading & Error states -->
    <div v-if="loading && users.length === 0" class="text-center text-orange-400 text-2xl py-10">
      <IconLoaderCircle class="w-16 h-16 mx-auto mb-4 animate-spin" />
    </div>
    <div v-else-if="errorMessage && users.length === 0" class="text-center text-red-500 py-10">
      {{ errorMessage }}
    </div>

    <!-- Cards container -->
    <div class="max-w-6xl mx-auto px-6">
      <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-9 mt-8">
        <div
          v-for="user in users"
          :key="user.id"
          class="w-[240px] p-4 bg-[#fae3cd] hover:bg-[#f7d5b6] transition-all rounded-xl flex flex-col items-center shadow-sm"
        >
          <!-- Icon in circle -->
          <div class="w-16 h-16 bg-white rounded-full flex items-center justify-center shadow-md hover:shadow-slate-300 hover:scale-110 transition-all">
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
  </div>
</template>

