<script setup lang="ts">
import { BackToTop } from '@yuelioi/ui';
import { ToastContainer } from '@yuelioi/toast';

import "@yuelioi/toast/dist/toast.css"

</script>

<template>
  <ToastContainer></ToastContainer>
  <RouterView v-slot="{ Component }">
    <keep-alive>
      <component :is="Component" />
    </keep-alive>
  </RouterView>
  <BackToTop class="right-3 bottom-6"></BackToTop>
</template>
