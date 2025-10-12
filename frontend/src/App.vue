<script setup lang="ts">
import { ToastContainer } from '@yuelioi/toast'

import '@yuelioi/toast/dist/toast.css'
</script>

<template>
  <ToastContainer></ToastContainer>
  <RouterView v-slot="{ Component }">
    <keep-alive>
      <component :is="Component" />
    </keep-alive>
  </RouterView>
</template>
