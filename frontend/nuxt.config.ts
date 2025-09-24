export default defineNuxtConfig({
  css: ['~/assets/css/main.css'],
  modules: [
    '@nuxt/image',
    '@nuxt/fonts',
    'nuxt-lucide-icons',
    '@nuxt/ui',
    '@pinia/nuxt'
  ],
  
  vite: {
    plugins: [],
  },
  lucide: {
    namePrefix: 'Icon'
  },
  routeRules: {
    '/': { redirect: '/register' },
  },
  compatibilityDate: '2025-07-21',
  devtools: { enabled: true },
  ssr: true,
  app: {
    baseURL: '/',
  },
  postcss: {
    plugins: {
      '@tailwindcss/postcss': {},
      autoprefixer: {},
    },
  },
})
