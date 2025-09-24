<script setup>
import { ref, computed } from 'vue'
import {
  Combobox,
  ComboboxInput,
  ComboboxButton,
  ComboboxOptions,
  ComboboxOption,
  TransitionRoot,
} from '@headlessui/vue'
import difficulty from './difficulty.vue'
import { ChevronUpDownIcon, CheckIcon } from '@heroicons/vue/20/solid' // or lucide icons

// Example categories, replace with real data from your backend or API
const categories = ref([
  { id: 1, name: 'Breakfast' },
  { id: 2, name: 'Lunch' },
  { id: 3, name: 'Dinner' },
  { id: 4, name: 'Dessert' },
  { id: 5, name: 'Snacks' },
])

const selectedCategory = ref(null) // ✅ correct use of ref
const query = ref('')

const filteredCategories = computed(() => {
  if (!query.value) return categories.value
  return categories.value.filter(category =>
    category.name.toLowerCase().includes(query.value.toLowerCase())
  )
})
watch(query, (val) => console.log("🔍 query updated:", val))
watch(filteredCategories, (val) => console.log("📦 filtered:", val))
watch(selectedCategory, (val) => console.log("✅ selected:", val))

</script>
<template>
<div class="flex justify-center gap-10">
    <div class="my-10">
      <label class="input  bg-white border-black">
      <svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
        <g stroke-linejoin="round"
            stroke-linecap="round"
            stroke-width="2.5"
            fill="none"
            stroke="currentColor"
        >
            <circle cx="11" cy="11" r="8"></circle>
            <path d="m21 21-4.3-4.3"></path>
          </g>
      </svg>
      <input type="search" class="grow bg-[#FFFEFC]" placeholder="Search Recipes..." />
      <kbd class="kbd bg-[#FFFEFC] kbd-sm">⌘</kbd>
      <kbd class="kbd bg-[#FFFEFC] kbd-sm">K</kbd>
      </label>
    </div>
    
    <div class=" flex mt-2 gap-6 ">
        <!-- <p class="my-3 ml-20 font-serif">All catagories</p> -->
        <div class=" bg-white mt-8 w-72">
            <Combobox v-model="selectedCategory">
            <div class="relative mt-1">
            <div class="relative w-full cursor-default overflow-hidden rounded-lg bg-white text-left shadow-md focus:outline-none focus-visible:ring-2 focus-visible:ring-white/75 focus-visible:ring-offset-2 focus-visible:ring-offset-teal-300 sm:text-sm"
            >
            <ComboboxInput
              placeholder="Select a category"
                class="w-full border-none py-2 pl-3 pr-10 text-sm leading-5 text-gray-900 bg-slate-300 focus:ring-0"
                :displayValue="(category) => category?.name || ''"
                 @input="event => query.value = event.target.value"
                
            />
            <ComboboxButton
                class="absolute inset-y-0 right-0 flex items-center pr-2"
            >
                <ChevronUpDownIcon
                class="h-5 w-5 text-gray-400"
                aria-hidden="true"
                />
            </ComboboxButton>
            </div>
            <TransitionRoot
            leave="transition ease-in duration-100"
            leave-from="opacity-100"
            leave-to="opacity-0"
            @after-leave="query = ''"
            >
            <ComboboxOptions class="absolute mt-1 max-h-60 w-full overflow-auto rounded-md bg-white py-1 text-base shadow-lg ring-1 ring-black/5 focus:outline-none sm:text-sm"
            >
                <div
                v-if="filteredCategories.length === 0 && query !== ''"
                class="relative cursor-default select-none px-4 py-2 text-gray-700"
                >
                Nothing found.
                </div>

                <ComboboxOption
                v-for="category in filteredCategories"
                as="template"
                :key="category.id"
                :value="category"
                v-slot="{ selected, active }"
                >
                <li
                    class="relative cursor-default select-none py-2 pl-10 pr-4"
                    :class="{
                    'bg-orange-300 text-white': active,
                    'text-gray-900': !active,
                    }"
                >
                    <span
                    class="block truncate"
                    :class="{ 'font-medium': selected, 'font-normal': !selected }"
                    >
                    {{ category.name }}
                    </span>
                    <span
                    v-if="selected"
                    class="absolute inset-y-0 left-0 flex items-center pl-3"
                    :class="{ 'text-white': active, 'text-teal-600': !active }"
                    >
                    <CheckIcon class="h-5 w-5" aria-hidden="true" />
                    </span>
                </li>
                </ComboboxOption>
            </ComboboxOptions>
            </TransitionRoot>
        </div>
        </Combobox>
        </div>
        <div class="mt-7">
            <difficulty/>
        </div>
    </div>
</div>

</template>