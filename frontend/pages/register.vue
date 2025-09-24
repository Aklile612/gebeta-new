<script setup>
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import { useRouter } from "#app"
definePageMeta({
  layout:false
})
const router = useRouter()
// Validation schema
const schema = yup.object({
  name: yup.string().required('Name is required').min(2, 'Name must be at least 2 characters'),
  email: yup.string().required('Email is required').email('Please enter a valid email'),
  password: yup.string().required('Password is required').min(6, 'Password must be at least 6 characters')
})

// Form setup with vee-validate
const { defineField, handleSubmit, errors, isSubmitting } = useForm({
  validationSchema: schema
})

// Define form fields
const [name, nameAttrs] = defineField('name')
const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

// Form submission handler
const onSubmit = handleSubmit(async (values) => {
  try {
    const response = await $fetch('http://localhost:8081/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: {
        name: values.name,
        email: values.email,
        password: values.password
      }
    })
    
    // Handle successful registration
    console.log('Registration successful:', response)
    
    // You can add navigation or success message here
    router.push('/')
  } catch (error) {
    console.error('Registration failed:', error)
    // Handle error - you can add error message display here
  }
})
</script>


<template>
  <div class="flex min-h-screen justify-center">
    <!-- image background -->
    <div class="hero w-1/2 min-h-screen bg-no-repeat bg-center bg-cover" style="background-image: url('/image/register.png');">
      <div class="flex flex-col justify-end items-center text-neutral-content w-full h-full pb-10 text-center">
        <div class="font-bold max-w-fit justify-center mb-9">
          <div class=" w-8 h-8 "><IconUtensils/></div>
          <p class="mt-3 mb-7 font-bold text-2xl">Every Recipe Has a Story</p>
          <p class="mb-7  text-sm ">Join Thousends of food Lovers creating and sharing there favorite recipe</p>
        </div>
      </div>
    </div>

  
    <!-- Optional Other Side -->
        <!-- Right Side (Form) -->
        <div class="w-1/2 bg-[#fae3cd] flex items-center justify-center">
  <div class="flex flex-col items-center  -mt-4"> 
    <!-- Icon + Gebeta Name -->
    <div class="flex flex-col items-center">
      <div class="w-16 h-16 bg-orange-500 rounded-full flex items-center justify-center shadow-md">
        <IconUtensils class="text-white  w-8 h-8" />
      </div>
      <p class="text-2xl mt-4 font-mono font-bold">Gebeta</p>
      <p class="text-sm text-gray-600">Recipe Sharing Community</p>
    </div>

    <!-- Welcome Text -->
    <div class="text-center mb-3">
      <p class="text-3xl font-bold">Create Your Account</p>
      <p class=" text-sm">Start Your Curlinary Journey today</p>
    </div>

    <!-- Form -->
    <form @submit.prevent="onSubmit" class="w-full -mt-3">
      <fieldset class="bg-[#fae3cd]   w-[370px] p-6 mx-5">
        <!-- Removed fieldset style classes that gave dark bg -->
        <label class="label text-lg" for="name"></label>
        <div class="relative mb-4">
          <input
            id="name"
            type="text"
            class="input input-bordered bg-white  h-10 pl-10  w-[300px] rounded-[6px]   text-black placeholder-text-gray-500 placeholder:text-xs "
            placeholder="Full Name"
            v-model="name"
            v-bind="nameAttrs"
          />
          
          <IconUser class="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-500" />
          <div v-if="errors.name" class="text-red-500 text-xs mt-1">{{ errors.name }}</div>
        </div>
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
        <label class="flex items-center gap-2 mb-2 whitespace-nowrap text-sm">
          <input type="checkbox" class="checkbox checked:bg-orange-500 bg-white" />
          <span>
            I agree to the
            <span class="text-orange-500 hover:underline">Terms of service</span>
            and
            <span class="text-orange-500 hover:underline">privacy policy</span>
          </span>
        </label>
        <button
          type="submit"
          :disabled="isSubmitting"
          class=" bg-orange-500 hover:bg-orange-600 text-white flex justify-center items-center group hover:scale-110 ease-out transition-all w-[300px] h-12  text-center  text-sm  font-semibold rounded-[6px] disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {{ isSubmitting ? 'Signing Up...' : 'SignUp' }}
        </button>
      </fieldset>
    </form>

    <!-- Link to Register -->
    <div>
      <NuxtLink to="/login">
        Already Have an account
        <span class="text-orange-500 hover:underline font-semibold">Log In</span>
      </NuxtLink>
    </div>
  </div>
</div>
  </div>
  
  </template>
  