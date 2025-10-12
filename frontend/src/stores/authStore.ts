import { defineStore } from 'pinia'
import { useStorage } from '@vueuse/core'

import { computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  const apiKey = useStorage('apikey', '')
  const isLogin = computed(() => apiKey.value != '')

  function setApiKey(key: string) {
    apiKey.value = key
  }

  return {
    isLogin,
    apiKey,
    setApiKey,
  }
})
