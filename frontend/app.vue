<template>
  <NuxtLayout>
    <div class="bg-white min-h-screen">
      <NuxtPage />
    </div>
    <UNotifications :ui="{ position: 'top-0 right-0' }" />
  </NuxtLayout>
</template>
