<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Recipe Header -->
    <div class="relative h-96 bg-white">
      <img v-if="recipe.featured_image" 
      :src="recipe.featured_image" 
      :alt="recipe.title" 
      class="w-full h-full object-cover" />
      <div class=""></div>
      <div class="absolute bottom-0 left-0 right-0 p-6 text-white">
        <div class="max-w-4xl mx-auto">
          <nav class="mb-4">
            <NuxtLink to="/recipes" class="flex items-center text-white hover:text-orange-300">
              <IconArrowLeft  class="w-5 h-5 text-white mr-2" />
              Back to recipes
            </NuxtLink>
          </nav>
          <h1 class="text-4xl font-bold mb-2">{{ recipe.title }}</h1>
          <p class="text-xl mb-4">{{ recipe.description }}</p>
          <div class="flex flex-wrap gap-6 text-sm">
            <div class="flex items-center">
              <IconClock  class="w-5 h-5 mr-2" />
              <span>{{ recipe.cook_time_minutes }} mins</span>
            </div>
            <div class="flex items-center">
              <IconUser  class="w-5 h-5 mr-2" />
              <span>By {{ recipe.user?.full_name }}</span>
            </div>
            <div v-if="averageRating" class="flex items-center">
              <IconStar  class="w-5 h-5 mr-1 text-yellow-400" />
              <span>{{ averageRating.toFixed(1) }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Recipe Content -->
    <div class="max-w-4xl mx-auto px-4 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column - Ingredients -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow-md p-6 sticky top-4">
            <h2 class="text-2xl font-bold mb-4">Ingredients</h2>
            <div class="space-y-3">
              <div 
                v-for="(ingredient, index) in recipe.ingredients" 
                :key="index"
                class="flex items-center p-2 rounded hover:bg-gray-50"
              >
                <input type="checkbox" :id="`ingredient-${index}`" class="mr-3">
                <label :for="`ingredient-${index}`" class="flex-1 cursor-pointer">
                  <span class="font-medium">{{ ingredient.quantity }}</span>
                  <span class="ml-2">{{ ingredient.name }}</span>
                </label>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="mt-6 space-y-3">
              <button 
                class="w-full bg-orange-500 text-white py-2 px-4 rounded-lg font-medium  transition"
              >
                Start Cooking
              </button>
              <button class="w-full border border-gray-300 py-2 px-4 rounded-lg font-medium hover:bg-gray-50 transition">
                Save to Favorites
              </button>
            </div>
          </div>
        </div>

        <!-- Right Column - Instructions -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-2xl font-bold mb-6">Instructions</h2>
            <div class="space-y-6">
              <div 
                v-for="step in sortedSteps" 
                :key="step.step_number"
                class="flex gap-4"
              >
                <div class="flex-shrink-0 w-8 h-8 bg-orange-500 text-white rounded-full flex items-center justify-center font-bold">
                  {{ step.step_number }}
                </div>
                <div class="flex-1">
                  <p class="mb-3">{{ step.description }}</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Additional Images -->
          <div v-if="recipe.recipe_images?.length" class="bg-white rounded-lg shadow-md p-6 mb-6">
            <h3 class="text-xl font-bold mb-4">Gallery</h3>
            <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
              <img 
                v-for="(image, index) in recipe.recipe_images" 
                :key="index"
                :src="image.image_url" 
                :alt="`${recipe.title} image ${index + 1}`"
                class="rounded-lg w-full h-32 object-cover cursor-pointer hover:opacity-80 transition"
                @click="openImageModal(image.image_url)"
              />
            </div>
          </div>

          <!-- Chef Info -->
          <div class="bg-white rounded-lg shadow-md p-6">
            <h3 class="text-xl font-bold mb-4">About the Chef</h3>
            <div class="flex items-center gap-4">
              <div class="w-16 h-16 bg-orange-100 rounded-full flex items-center justify-center">
                <User  class="w-8 h-8 text-orange-500" />
              </div>
              <div>
                <h4 class="font-bold text-lg">{{ recipe.user?.full_name }}</h4>
                <p class="text-gray-600">Recipe creator</p>
              </div>
            </div>
          </div>

          <!-- Comments -->
          <div v-if="sortedComments.length" class="bg-white rounded-lg shadow-md p-6 mt-6">
            <h3 class="text-xl font-bold mb-4">Comments</h3>
            <div class="space-y-4">
                <div v-for="comment in sortedComments" :key="comment.id" class="border-b pb-4 last:border-b-0">
                <div class="flex justify-between items-start mb-2">
                    <div class="flex items-center gap-2">
                    <div class="w-8 h-8 bg-gray-200 rounded-full flex items-center justify-center">
                        <IconUser class="w-4 h-4 text-gray-600" />
                    </div>
                    <span class="font-medium">{{ comment.user?.full_name || 'Anonymous' }}</span>
                    </div>
                    <span class="text-sm text-gray-500">{{ formatDate(comment.created_at) }}</span>
                </div>
                <p class="text-gray-700">{{ comment.comment }}</p>
                </div>
            </div>
          </div>
          <!-- add comment -->
           <div class="bg-white rounded-lg shadow-md p-6 mt-6">
      <h3 class="text-xl font-bold mb-4">Add a Comment</h3>
      <form @submit.prevent="onSubmit" class="space-y-4">
        <div>
          <textarea
            v-model="comment"
            v-bind="commentAttrs"
            placeholder="Share your thoughts about this recipe..."
            rows="4"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500 resize-none"
            :class="{ 'border-red-500': errors.comment }"
          ></textarea>
          <div v-if="errors.comment" class="text-red-500 text-xs mt-1">{{ errors.comment }}</div>
        </div>
        
        <button
          type="submit"
          :disabled="isSubmitting"
          class="bg-orange-500 hover:bg-orange-600 text-white py-2 px-6 rounded-md font-medium transition disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isSubmitting ? 'Posting...' : 'Post Comment' }}
        </button>
      </form>
    </div>
        </div>
      </div>
    </div>

    <!-- Image Modal -->
    <div v-if="selectedImage" class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50 p-4" @click="selectedImage = null">
      <div class="max-w-4xl max-h-full">
        <img :src="selectedImage" :alt="recipe.title" class="max-w-full max-h-full object-contain">
      </div>
    </div>
  </div>
</template>

<script setup>
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import { useAuthStore } from '~/stores/auth'
// Define route parameter
const route = useRoute()
const recipeId = route.params.id
const authStore = useAuthStore()
const config = useRuntimeConfig()
const GRAPHQL_URL = config.public?.graphqlUrl || 'https://gebeta-app.hasura.app/v1/graphql'
const HASURA_ADMIN_SECRET = '5zghjagw4B9BsAy7grRxTOxrtrL6o5ivzSSPuXpwW4z1AXyS5hN3xo06T4HEBllw'

// Recipe data
const recipe = ref({})
const loading = ref(true)
const errorMessage = ref('')
const selectedImage = ref(null)


const schema = yup.object({
  comment: yup.string().required('Comment is required').min(5, 'Comment must be at least 5 characters').max(500, 'Comment too long')
})
// Comment form setup
const { defineField, handleSubmit, errors, isSubmitting, resetForm } = useForm({
  validationSchema: schema
})
const [comment, commentAttrs] = defineField('comment')
const fetchRecipeData = async () => {
  if (!isValidRecipeId.value) {
    errorMessage.value = 'Invalid recipe ID'
    loading.value = false
    return
  }

  loading.value = true
  try {
    const headers = { 
      'Content-Type': 'application/json',
      'x-hasura-admin-secret': HASURA_ADMIN_SECRET
    }

    const res = await fetch(GRAPHQL_URL, {
      method: 'POST',
      headers,
      body: JSON.stringify({ 
        query: RECIPE_QUERY,
        variables: { id: recipeId  }
      })
    })

    const data = await res.json()

    if (data.errors) {
      throw new Error(data.errors?.[0]?.message || 'Failed to load recipe')
    }

    if (!data.data?.recipes_by_pk) {
      throw new Error('Recipe not found')
    }

    recipe.value = data.data.recipes_by_pk

  } catch (err) {
    errorMessage.value = err?.message || 'Unexpected error loading recipe'
    console.error('Error loading recipe:', err)
  } finally {
    loading.value = false
  }
}
onMounted(() => {
  authStore.loadToken()
})
// Call it onMounted
onMounted(fetchRecipeData)

// Comment submission handler
const onSubmit = handleSubmit(async (values) => {
  try {
    // Check if user is authenticated
    if (!authStore.isAuthenticated) {
      errorMessage.value = 'Please log in to post a comment'
      return
    }
    const formData = new FormData()
    formData.append('comment', values.comment)


    const response = await $fetch(`https://gebeta-new.onrender.com/comment_recipes/${recipeId}`, {
      method: 'POST',
      headers: {
        
        'Authorization': `Bearer ${authStore.token}`
      },
      body:  formData
    })
    
    console.log('Comment posted successfully:', response)
    
    // Refresh the recipe data to show the new comment
    await fetchRecipeData()
    
    // Reset the form
    resetForm()
    
  } catch (error) {
    console.error('Failed to post comment:', error)
    errorMessage.value = error.data?.message || 'Failed to post comment'
  }
})
// Validate recipeId
const isValidRecipeId = computed(() => !!recipeId)

// GraphQL query for single recipe
const RECIPE_QUERY = `
  query GetRecipe($id: uuid!) {
    recipes_by_pk(id: $id) {
      id
      title
      description
      cook_time_minutes
      prep_time_minutes
      featured_image
      created_at
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
      recipe_images {
        image_url
      }
      ingredients {
        name
        quantity
      }
      recipe_steps {
        description
        step_number
      }
      recipe_comments(order_by: {created_at: desc}) {
        id
        comment
        created_at
        user {
          full_name
        }
      }
    }
  }
`

// Computed properties
const averageRating = computed(() => {
  const rating = recipe.value.recipe_ratings_aggregate?.aggregate?.avg?.rating
  return rating ? parseFloat(rating) : 0
})

const sortedSteps = computed(() => {
  if (!recipe.value.recipe_steps) return []
  return [...recipe.value.recipe_steps].sort((a, b) => a.step_number - b.step_number)
})

const sortedComments = computed(() => {
  return (recipe.value.recipe_comments || []).map(c => ({
    ...c,
    id: c.id || crypto.randomUUID()  // fallback if id is missing
  }))
})


// Methods
const openImageModal = (imageUrl) => {
  selectedImage.value = imageUrl
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

// Fetch recipe data
onMounted(async () => {
  // Validate recipeId first
  if (!isValidRecipeId.value) {
    errorMessage.value = 'Invalid recipe ID'
    loading.value = false
    return
  }

  try {
    const headers = { 
      'Content-Type': 'application/json',
      'x-hasura-admin-secret': HASURA_ADMIN_SECRET
    }

    const res = await fetch(GRAPHQL_URL, {
      method: 'POST',
      headers,
      body: JSON.stringify({ 
        query: RECIPE_QUERY,
        variables: { id: recipeId  }
      })
    })

    const data = await res.json()
    
    if (data.errors) {
      throw new Error(data.errors?.[0]?.message || 'Failed to load recipe')
    }

    if (!data.data?.recipes_by_pk) {
      throw new Error('Recipe not found')
    }

    recipe.value = data.data.recipes_by_pk
    
  } catch (err) {
    errorMessage.value = err?.message || 'Unexpected error loading recipe'
    console.error('Error loading recipe:', err)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.sticky {
  position: sticky;
}
</style>