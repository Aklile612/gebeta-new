<script setup>
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import { useRouter } from "#app"
import { useAuthStore } from '../stores/auth'


definePageMeta({
  layout:false
})

const router = useRouter()
const authStore = useAuthStore()
// Validation schema
const schema = yup.object({
  email: yup.string().required('Email is required').email('Please enter a valid email'),
  password: yup.string().required('Password is required').min(6, 'Password must be at least 6 characters')
})

// Form setup with vee-validate
const { defineField, handleSubmit, errors, isSubmitting } = useForm({
  validationSchema: schema
})

// Define form fields
const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

// Form submission handler
const onSubmit = handleSubmit(async (values) => {
  try {
    const response = await $fetch('http://localhost:8081/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: {
        email: values.email,
        password: values.password
      }
    })
    
    // Handle successful login
    console.log('Login successful:', response)
    console.log('Response token:', response.Token)
    
    // Store the token and navigate to home page
    if (response.Token) {
      authStore.setToken(response.Token)
      console.log('Token stored in store:', authStore.token)
      console.log('Token in localStorage:', localStorage.getItem('token'))
    } else {
      console.error('No token in response:', response)
    }
    router.push('/')
  } catch (error) {
    console.error('Login failed:', error)
    // Handle error - you can add error message display here
  }
})
</script>


<template>
  <div class="flex min-h-screen justify-center">
    <!-- image background -->
    <div class="hero w-1/2 min-h-screen bg-no-repeat bg-center bg-cover" style="background-image: url('/image/login.png');">
      <div class="flex flex-col justify-end items-center text-neutral-content w-full h-full pb-10 text-center">
        <div class="max-w-md">
          <p class="mt-7 mr-4 font-bold text-2xl">Discover Amazing Recipes</p>
          <p class="mt-7 mr-4 text-lg">Join thousands of food lovers sharing their favorite recipes</p>
        </div>
      </div>
    </div>

  
    <!-- Optional Other Side -->
        <!-- Right Side (Form) -->
        <div class="w-1/2 bg-[#fae3cd] flex items-center justify-center">
  <div class="flex flex-col items-center space-y-6 -mt-6"> <!-- reduced vertical space with -mt -->
    <!-- Icon + Gebeta Name -->
    <div class="flex flex-col items-center">
      <div class="w-16 h-16 bg-orange-500 rounded-full flex items-center justify-center shadow-md">
        <IconUtensils class="text-white w-8 h-8" />
      </div>
      <p class="text-2xl mt-2 font-mono font-bold">Gebeta</p>
      <p class="text-sm text-gray-600">Recipe Sharing Community</p>
    </div>

    <!-- Welcome Text -->
    <div class="text-center">
      <p class="text-3xl font-bold">Welcome Back</p>
      <p class="font-bold text-sm">Login to share your next delicious recipe!</p>
    </div>

    <!-- Form -->
    <form @submit.prevent="onSubmit" class="w-full">
      <fieldset class="bg-[#fae3cd]   w-[370px] mx-5 p-6 ">
        <!-- Removed fieldset style classes that gave dark bg -->
        <label class="label text-lg" for="email"></label>
        <div class="relative mb-4">
          <input
            id="email"
            type="email"
            class="input input-bordered bg-white  h-10 pl-10  w-[300px] rounded-[6px]   text-black placeholder-text-gray-500 placeholder:text-xs"
            placeholder="Email"
            v-model="email"
            v-bind="emailAttrs"
          />
          <IconMail class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
          <div v-if="errors.email" class="text-red-500 text-xs mt-1">{{ errors.email }}</div>
        </div>

        <label class="label text-lg" for="password"></label>
        <div class="relative mb-4">
          <input
            id="password"
            type="password"
            class="input input-bordered bg-white  h-10 pl-10  w-[300px] rounded-[6px]   text-black placeholder-text-gray-500 placeholder:text-xs"
            placeholder="Password"
            v-model="password"
            v-bind="passwordAttrs"
          />
          <IconLock class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
          <div v-if="errors.password" class="text-red-500 text-xs mt-1">{{ errors.password }}</div>
        </div>

        <button
          type="submit"
          :disabled="isSubmitting"
          class=" bg-orange-500 hover:bg-orange-600 text-white flex justify-center items-center group hover:scale-110 ease-out transition-all w-[300px] h-12  text-center  text-sm  font-semibold rounded-[6px] disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isSubmitting ? 'Logging in...' : 'Login' }}
        </button>
      </fieldset>
    </form>

    <!-- Link to Register -->
    <div>
      <NuxtLink to="/register">
        <span class="font-semibold ">New Here?</span>
        <span class="text-orange-500 font-semibold hover:underline">Create an account</span>
      </NuxtLink>
    </div>
  </div>
</div>
  </div>
  
  </template>
  