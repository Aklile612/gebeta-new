<script setup>
import { useForm, useFieldArray } from 'vee-validate'
import * as yup from 'yup'
import { ref, watch, computed, onBeforeUnmount, onMounted, nextTick } from 'vue'
import { useToast } from '#imports'
import { useAuthStore } from '../stores/auth'
import { useRouter } from '#app'
// In your component
const toast = useToast();
const authStore = useAuthStore();
const router = useRouter();

// Load token from localStorage on component mount
onMounted(() => {
  authStore.loadToken();
})
// Enhanced schema with better validation messages
const schema = yup.object({
  title: yup.string().required('Recipe title is required').max(100, 'Title too long'),
  description: yup.string().required('Description is required').max(500, 'Description too long'),
  prepTime: yup.number()
    .min(1, 'Must be at least 1 minute')
    .max(600, 'Prep time too long')
    .required('Prep time is required'),
  cookTime: yup.number()
    .min(1, 'Must be at least 1 minute')
    .max(600, 'Cook time too long')
    .required('Cook time is required'),
  images: yup.array()
    .min(1, 'At least one image is required')
    .max(5, 'Maximum 5 images allowed'),
  category: yup.string().required('Category is required'),
  difficulty: yup.string().required('Difficulty is required'),
  ingredients: yup.array()
    .of(yup.object({
      name: yup.string().required('Ingredient name is required').max(50),
      quantity: yup.string().required('Quantity is required').max(20)
    }))
    .min(1, 'At least one ingredient is required'),
  steps: yup.array()
    .of(yup.object({
      description: yup.string().required('Step description is required').max(500)
    }))
    .min(1, 'At least one step is required'),
    price: yup.number()
  .when('isPremium', {
    is: true,
    then: (schema) => 
      schema
        .typeError('Price must be a number')
        .required('Price is required for premium recipes')
        .min(1, 'Price must be at least 1')
        .max(1000, 'Price too high'),
    otherwise: (schema) => schema.nullable().notRequired()
  })
})

// Form setup with enhanced initial values
const { handleSubmit, errors, values, setFieldValue, meta, touched, resetForm } = useForm({
  validationSchema: schema,
  initialValues: {
    title: '',
    description: '',
    prepTime: '',
    cookTime: '',
    images: [],
    category: '',
    difficulty: '',
    ingredients: [{ name: '', quantity: '' }],
    steps: [{ description: '' }],
    isPremium: false,
    price: null
  }
})

const isFormValid = computed(() => meta.value.valid)
const isSubmitting = ref(false)

// Enhanced field arrays with auto-focus on new items
const { fields: ingredientFields, push: addIngredient, remove: removeIngredient } = useFieldArray('ingredients')
const { fields: stepFields, push: addStep, remove: removeStep } = useFieldArray('steps')

const addNewIngredient = () => {
  addIngredient({ name: '', quantity: '' })
  nextTick(() => {
    const inputs = document.querySelectorAll('.ingredient-name-input')
    inputs[inputs.length - 1]?.focus()
  })
}

const addNewStep = () => {
  addStep({ description: '' })
  nextTick(() => {
    const textareas = document.querySelectorAll('.step-description-input')
    textareas[textareas.length - 1]?.focus()
  })
}

// File handling with better error states
const imagePreviews = ref([])
const fileInput = ref(null)
const fileError = ref(null)

function handleImageUpload(event) {
  fileError.value = null
  const MAX_SIZE = 5 * 1024 * 1024 // 5MB
  const MAX_FILES = 5
  const files = event.target.files
  
  if (!files || files.length === 0) return

  // Check number of files
  if (files.length + imagePreviews.value.length > MAX_FILES) {
    fileError.value = `You can upload maximum ${MAX_FILES} files`
    return
  }

  // Check file sizes and types
  const validFiles = Array.from(files).filter(file => {
    if (file.size > MAX_SIZE) {
      fileError.value = `File ${file.name} exceeds 5MB limit`
      return false
    }
    if (!['image/jpeg', 'image/png'].includes(file.type)) {
      fileError.value = `File ${file.name} must be JPG or PNG`
      return false
    }
    return true
  })

  if (validFiles.length === 0) return

  // Create new array with existing and new files
  const newImages = [...values.images, ...validFiles]
  setFieldValue('images', newImages)

  // Generate previews only for new files
  const newPreviews = validFiles.map(file => ({
    url: URL.createObjectURL(file),
    name: file.name,
    size: file.size,
    type: file.type
  }))
  
  imagePreviews.value = [...imagePreviews.value, ...newPreviews]
  event.target.value = '' // Reset file input
}

// Clean up object URLs
onBeforeUnmount(() => {
  imagePreviews.value.forEach(preview => {
    URL.revokeObjectURL(preview.url)
  })
})

const removeImage = (index) => {
  // Revoke the object URL
  URL.revokeObjectURL(imagePreviews.value[index].url)
  
  // Remove from both form state and previews
  const newImages = [...values.images]
  newImages.splice(index, 1)
  setFieldValue('images', newImages)
  
  imagePreviews.value.splice(index, 1)
}
watch(() => values.isPremium, (newValue) => {
  if (!newValue) {
    setFieldValue('price', null)
  }
})

// Enhanced form submission with loading state and better error handling
const onSubmit = handleSubmit(
  async (values) => {
    isSubmitting.value = true
    try {
      const formData = new FormData()
      
      // Append individual form fields as expected by backend
      formData.append('title', values.title)
      formData.append('description', values.description)
      formData.append('prep_time_minutes', values.prepTime.toString())
      formData.append('cook_time_minutes', values.cookTime.toString())
      formData.append('difficulty', values.difficulty)
      formData.append('category_name', values.category)
      formData.append('is_paid', values.isPremium ? 'true' : 'false')
      formData.append('price', values.price ? values.price.toString() : '0')
      formData.append('steps', JSON.stringify(values.steps))
      formData.append('ingredients', JSON.stringify(values.ingredients))
      
      // Append each image file
      values.images.forEach((file, index) => {
        formData.append('image', file)
      })

      // Check localStorage directly
      const tokenFromStorage = localStorage.getItem('token')
      console.log('Token from localStorage:', tokenFromStorage)
      
      // Reload token from localStorage to ensure it's current
      authStore.loadToken()
      
      // Check if token exists, redirect to login if not
      console.log('Auth store token:', authStore.token)
      console.log('Is authenticated:', authStore.isAuthenticated)
      
      // Use token from store or localStorage as fallback
      const tokenToUse = authStore.token || tokenFromStorage
      
      if (!tokenToUse) {
        alert('Please login first to submit a recipe.')
        router.push('/login')
        return
      }

      console.log('Using token:', tokenToUse)

      // Make API call with authentication
      const response = await $fetch('https://gebeta-new.onrender.com/add_recipes', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${tokenToUse}`
        },
        body: formData
      })
      
      console.log('Recipe submitted successfully:', response)
      alert("Recipe submitted successfully!")
      // Show success toast
      toast.add({
        title: "Success!",
        description: "Recipe saved successfully!",
        color: "green",
        timeout: 5000
      })
      
      resetForm()
      router.push('/recipes')
    } catch (error) {
      // Extract meaningful error info from $fetch error
      const status = error?.response?.status || error?.statusCode || 0
      const apiError = error?.data?.error || error?.data?.message || error?.message || 'Unknown error'
      const detail = error?.data?.detail || error?.data?.details

      console.error('Submission error:', { status, apiError, detail, raw: error })

      // Show error toast with server detail if present
      toast.add({
        title: `Error${status ? ` ${status}` : ''}`,
        description: `${apiError}${detail ? `: ${detail}` : ''}`,
        color: 'red',
        timeout: 6000
      })
    } finally {
      isSubmitting.value = false
    }
  },
  (errors) => {
    console.log('Validation errors:', errors)
    // Enhanced error scrolling
    const firstErrorKey = Object.keys(errors)[0]
    if (firstErrorKey) {
      let selector = `[name="${firstErrorKey}"]`
      
      // Handle array fields
      if (firstErrorKey.includes('[')) {
        const [fieldName, index] = firstErrorKey.split(/[[\]]/).filter(Boolean)
        selector = `[name="${fieldName}_${index}"]`
      }
      
      const errorElement = document.querySelector(selector)
      if (errorElement) {
        errorElement.scrollIntoView({
          behavior: 'smooth',
          block: 'center'
        })
        errorElement.focus()
      }
    }
  }
)

// Debug helpers
watch(values, (newVal) => {
  console.log('Form values updated:', JSON.parse(JSON.stringify(newVal)))
}, { deep: true })

const logFormState = () => {
  console.log('Current form state:', {
    ...values,
    priceType: typeof values.price,
    isPremium: values.isPremium,
    images: values.images.map(f => f.name)
  })
}
const handlePriceInput = (event) => {
  const value = event.target.value
  // Convert to number if not empty, otherwise set to null
  setFieldValue('price', value === '' ? null : Number(value))
}
</script>

<template>
  <div class="min-h-screen pt-[70px] bg-[#fae3cd] py-8 px-4">
    <div class="max-w-3xl mx-auto">
      <!-- Header -->
      <div class="mb-6">
        <NuxtLink to="/" class="flex items-center text-gray-700 hover:text-orange-500 mb-2">
          <IconChevronLeft class="mr-1" />
          <span class="text-sm font-medium">Back</span>
        </NuxtLink>
        <h1 class="text-2xl font-bold text-gray-800">Create New Recipe</h1>
        <p class="text-gray-600">Share your culinary creation with the community</p>
      </div>

      <!-- Form -->
      <form @submit.prevent="onSubmit" class="bg-white rounded-lg shadow-md p-6">
        <!-- Basic Information Section -->
        <section class="mb-8">
          <h2 class="text-lg font-bold mb-4 pb-2 border-b border-gray-200">Basic Information</h2>
          
          <!-- Recipe Title -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Recipe Title*</label>
            <input 
              name="title"
              v-model="values.title" 
              @input="setFieldValue('title', $event.target.value)"
              placeholder="e.g., Grandma's Chocolate Chip Cookies" 
              class="w-full px-3 py-2 border bg-white border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
            />
            <span class="text-red-500 text-xs mt-1" v-if="touched?.title && errors.title">{{ errors.title }}</span>
          </div>
          
          <!-- Description -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-1">Description*</label>
            <textarea
            name="description" 
              v-model="values.description"
              @input="setFieldValue('description', $event.target.value)" 
              placeholder="Tell us about your recipe..." 
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 bg-white rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
            ></textarea>
            <span class="text-red-500 text-xs mt-1" v-if="touched?.description && errors.description">{{ errors.description }}</span>
          </div>
          
          <!-- Category & Difficulty -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Category*</label>
              <select
              name="category" 
                v-model="values.category"
                @input="setFieldValue('category', $event.target.value)" 
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              >
                <option value="">Select category</option>
                <option value="breakfast">Breakfast</option>
                <option value="lunch">Lunch</option>
                <option value="dinner">Dinner</option>
                <option value="dessert">Dessert</option>
                <option value="snack">Snack</option>
              </select>
              <span class="text-red-500 text-xs mt-1" v-if="touched?.category && errors.category">{{ errors.category }}</span>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Difficulty*</label>
              <select
                name="difficulty"
                @input="setFieldValue('difficulty', $event.target.value)" 
                v-model="values.difficulty" 
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              >
                <option value="">Select difficulty</option>
                <option value="easy">Easy</option>
                <option value="medium">Medium</option>
                <option value="hard">Hard</option>
              </select>
              <span class="text-red-500 text-xs mt-1" v-if="touched?.difficulty && errors.difficulty">{{ errors.difficulty }}</span>
            </div>
          </div>
          
          <!-- Prep & Cook Time -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Prep Time (minutes)*</label>
              <input
                name="prepTime" 
                v-model="values.prepTime"
                @input="setFieldValue('prepTime', Number($event.target.value))" 
                type="number" 
                min="1"
                class="w-full px-3 py-2 border bg-white border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
              <span class="text-red-500 text-xs mt-1" v-if="touched?.prepTime && errors.prepTime">{{ errors.prepTime }}</span>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Cook Time (minutes)*</label>
              <input
              name="cookTime" 
                @input="setFieldValue('cookTime', Number($event.target.value))"
                v-model="values.cookTime" 
                type="number" 
                min="1"
                class="w-full px-3 py-2 border bg-white border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
              <span class="text-red-500 text-xs mt-1" v-if="touched?.cookTime && errors.cookTime">{{ errors.cookTime }}</span>
            </div>
          </div>
        </section>

        <!-- Recipe Images Section -->
        <section class="mb-8">
          <h2 class="text-lg font-bold mb-4 pb-2 border-b border-gray-200">Recipe Images</h2>
          <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
            <div class="flex flex-col items-center justify-center">
              <svg class="w-12 h-12 text-gray-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
              </svg>
              <p class="text-sm text-gray-500 mb-2">Drag & drop images or click to browse</p>
              <p class="text-xs text-gray-400 mb-4">Support: JPG, PNG (Max 5MB each)</p>
              <input type="file"
                multiple
                accept="image/jpeg,image/png"
                @change="handleImageUpload"
                class="hidden"
                id="image-upload"
                ref="fileInput"
              />
              <label for="image-upload" class="px-4 py-2 bg-orange-500 text-white rounded-md cursor-pointer hover:bg-orange-600 transition">
                Choose Files
              </label>
            </div>
            <span class="text-red-500 text-xs mt-2" v-if="touched?.images && errors.images">{{ errors.images }}</span>
            <!-- Image Preview Thumbnails -->
            <div v-if="imagePreviews.length" class="mt-4 grid grid-cols-2 sm:grid-cols-3 gap-4">
              <div
                v-for="(image, index) in imagePreviews"
                :key="index"
                class="relative border rounded overflow-hidden"
              >
                <img :src="image.url" :alt="image.name" class="w-full h-32 object-cover" />
                <button
                  type="button"
                  @click="removeImage(index)"
                  class="absolute top-1 right-1 bg-white rounded-full p-1 shadow hover:bg-gray-100"
                >
                  <svg class="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>

          </div>
        </section>

        <!-- Ingredients Section -->
        <section class="mb-8">
          <h2 class="text-lg font-bold mb-4 pb-2 border-b border-gray-200">Ingredients</h2>
          <div class="space-y-3">
            <div v-for="(ingredient, idx) in ingredientFields" :key="ingredient.key" class="flex items-center gap-3">
              <input
                :name="`ingredients_${idx}_name`"
                v-model="values.ingredients[idx].name"
                @input="setFieldValue(`ingredients[${idx}].name`, $event.target.value)"
                placeholder="Ingredient name"
                class="flex-1 px-3 py-2 bg-white border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
              <input
                :name="`ingredients_${idx}_quantity`"
                @input="setFieldValue(`ingredients[${idx}].quantity`, $event.target.value)"
                v-model="values.ingredients[idx].quantity"
                placeholder="Quantity"
                class="w-1/3 px-3 py-2 border bg-white border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
              <button
                type="button"
                @click="removeIngredient(idx)"
                class="text-red-500 hover:text-red-700"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
              </button>
            </div>
            <button
              type="button"
              @click="addIngredient({ name: '', quantity: '' })"
              class="flex items-center text-orange-500 hover:text-orange-700 mt-2"
            >
              <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
              Add ingredient
            </button>
            <span class="text-red-500 text-xs mt-1" v-if="touched?.ingredients && errors.ingredients">{{ errors.ingredients }}</span>
          </div>
        </section>

        <!-- Preparation Steps Section -->
        <section class="mb-8">
          <h2 class="text-lg font-bold mb-4 pb-2 border-b border-gray-200">Preparation Steps</h2>
          <div class="space-y-3">
            <div v-for="(step, idx) in stepFields" :key="step.key" class="flex gap-3">
              <div class="flex-1">
                <div class="flex items-start">
                  <span class="mt-2 mr-2 font-bold text-gray-500">{{ idx + 1 }}.</span>
                  <textarea
                    :name="`steps_${idx}_description`"
                    @input="setFieldValue(`steps[${idx}].description`, $event.target.value)"
                    v-model="values.steps[idx].description"
                    placeholder="Describe this step..."
                    rows="2"
                    class="flex-1 px-3 py-2 bg-white border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
                  ></textarea>
                </div>
              </div>
              <button
                type="button"
                @click="removeStep(idx)"
                class="text-red-500 hover:text-red-700 mt-2"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                </svg>
              </button>
            </div>
            <button
              type="button"
              @click="addStep({ description: '' })"
              class="flex items-center text-orange-500 hover:text-orange-700 mt-2"
            >
              <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
              Add Step
            </button>
            <span class="text-red-500 text-xs mt-1" v-if="touched?.steps && errors.steps">{{ errors.steps }}</span>
          </div>
        </section>

        <!-- Pricing Section -->
       
        <section class="mb-6">
          <h2 class="text-lg font-bold mb-4 pb-2 border-b border-gray-200">Pricing</h2>
          <div class="flex items-center mb-2">
            <input
            type="checkbox"
            id="premium-recipe"
            :checked="values.isPremium"
            @change="setFieldValue('isPremium', $event.target.checked)"
            class="w-4 h-4 text-orange-500 rounded focus:ring-orange-500"
            />
            <label for="premium-recipe" class="ml-2 text-sm font-medium text-gray-700">Premium Recipe</label>
          </div>
          <p class="text-xs text-gray-500 mb-3">Charge users to access this recipe</p>
          
          <!-- Price Input (Conditional) -->
          <div v-if="values.isPremium" class="mt-2 transition-all duration-200 ease-in-out">
            <label class="block text-sm font-medium text-gray-700 mb-1">Price (ETB)*</label>
            <div class="relative rounded-md shadow-sm">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <span class="text-gray-500">ETB</span>
              </div>
              <input
              name="price"
              :value="values.price"
              @input="handlePriceInput($event)"
              type="number"
              min="1"
              placeholder="Enter amount"
              class="block w-full pl-12 pr-3 py-2 border bg-white border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500"
              />
            </div>
            <span class="text-red-500 text-xs mt-1" v-if="touched?.price && errors.price">{{ errors.price }}</span>
          </div>
        </section>

        <!-- Submit Button -->
        <button
          type="submit"
          :disabled="isSubmitting"
          class="w-full cursor-pointer bg-orange-500 text-white py-3 px-4 rounded-md hover:bg-orange-600 transition font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
        >
          <svg v-if="isSubmitting" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isSubmitting ? 'Submitting Recipe...' : 'Submit Recipe' }}
        </button>
      </form>
      
    </div>
  </div>
  
</template>

