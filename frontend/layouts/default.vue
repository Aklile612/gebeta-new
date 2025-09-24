<template>
    <div class="layout">
      <Header />
  
      <main class="min-h-screen">
        <slot />
      </main>
  
      <Footer />
    </div>
  </template>
  
  <script setup>
  import Header from '~/components/Header.vue'
  import Footer from '~/components/Footer.vue'
  </script>
  