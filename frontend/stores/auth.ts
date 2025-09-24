import { defineStore } from "pinia"

export const useAuthStore = defineStore("auth", {
  state: () => ({
    token: null as string | null
  }),
  getters: {
    isAuthenticated: (state) => !!state.token
  },
  actions: {
    setToken(newToken: string) {
      this.token = newToken
      localStorage.setItem("token", newToken)
    },
    loadToken() {
      const savedToken = localStorage.getItem("token")
      if (savedToken) {
        this.token = savedToken
      }
    },
    clearToken() {
      this.token = null
      localStorage.removeItem("token")
    }
  }
})
