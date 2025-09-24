<script setup>

</script>


<template>
    <div class="bg-white  w-full h-[70vh]">
        <div class="mx-[30vw] pt-[70px]">
            <h1 class="mb-5 text-4xl ml-[7vw] font-serif  font-bold">Top Recipe Creators</h1>

        </div>
        <div class="flex gap-9 mt-8 justify-center">
            <div class="w-[240px] p-4 bg-[#fae3cd] rounded-xl flex flex-col items-center shadow-sm">
                <!-- Icon in circle -->
                <div class="w-16 h-16 bg-white hover:shadow-slate-300 hover:scale-110 rounded-full flex items-center justify-center shadow-md">
                  <IconUsers class="text-orange-500 w-8 h-8" />
                </div>
            
                <!-- User info -->
                <div class="mt-4 text-center">
                  <p class="font-bold text-sm text-gray-800">Alana Germy</p>
                  <p class="text-xs text-gray-600 mt-1">35 Recipes Submitted</p>
                
                  <!-- View Button -->
                  <button
                    class="mt-3 text-sm bg-orange-300 hover:bg-orange-500 text-gray-800 hover:text-white font-medium py-1 px-3 rounded-md transition-all duration-200"
                  >
                    View Recipes
                  </button>
                </div>
            </div>
            <div class="w-[240px] p-4 bg-[#fae3cd] rounded-xl flex flex-col items-center shadow-sm">
                <!-- Icon in circle -->
                <div class="w-16 h-16 hover:shadow-slate-300 hover:scale-110 bg-white rounded-full flex items-center justify-center shadow-md">
                  <IconUsers class="text-orange-500 w-8 h-8" />
                </div>
            
                <!-- User info -->
                <div class="mt-4 text-center">
                  <p class="font-bold text-sm text-gray-800">Mike Rodrigaze</p>
                  <p class="text-xs text-gray-600 mt-1">32 Recipes Submitted</p>
                
                  <!-- View Button -->
                  <button
                    class="mt-3 text-sm bg-orange-300 hover:bg-orange-500 text-gray-800 hover:text-white font-medium py-1 px-3 rounded-md transition-all duration-200"
                  >
                    View Recipes
                  </button>
                </div>
            </div>
            <div class="w-[240px] p-4 bg-[#fae3cd] rounded-xl flex flex-col items-center shadow-sm">
                <!-- Icon in circle -->
                <div class="w-16 h-16 bg-white rounded-full hover:shadow-slate-300 hover:scale-110 flex items-center justify-center shadow-md">
                  <IconUsers class="text-orange-500 w-8 h-8" />
                </div>
            
                <!-- User info -->
                <div class="mt-4 text-center">
                  <p class="font-bold text-sm text-gray-800">Emma Wilson</p>
                  <p class="text-xs text-gray-600 mt-1">25 Recipes Submitted</p>
                
                  <!-- View Button -->
                  <button
                    class="mt-3 text-sm bg-orange-300 hover:bg-orange-500 text-gray-800 hover:text-white font-medium py-1 px-3 rounded-md transition-all duration-200"
                  >
                    View Recipes
                  </button>
                </div>
            </div>
            <div class="w-[240px] p-4 bg-[#fae3cd] rounded-xl flex flex-col items-center shadow-sm">
                <!-- Icon in circle -->
                <div class="w-16 h-16 bg-white hover:shadow-slate-300 hover:scale-110 rounded-full flex items-center justify-center shadow-md">
                  <IconUsers class="text-orange-500 w-8 h-8" />
                </div>
            
                <!-- User info -->
                <div class="mt-4 text-center">
                  <p class="font-bold text-sm text-gray-800">David Kim</p>
                  <p class="text-xs text-gray-600 mt-1">25 Recipes Submitted</p>
                
                  <!-- View Button -->
                  <button
                    class="mt-3 text-sm bg-orange-300 hover:bg-orange-500 text-gray-800 hover:text-white font-medium py-1 px-3 rounded-md transition-all duration-200"
                  >
                    View Recipes
                  </button>
                </div>
            </div>
        </div>
    </div>
</template>