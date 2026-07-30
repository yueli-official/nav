<script setup lang="ts">
defineProps<{
  error: { statusCode?: number; statusMessage?: string; message?: string };
}>();
</script>

<template>
  <UApp>
    <main class="mx-auto flex min-h-screen max-w-xl flex-col items-center justify-center gap-5 px-6 text-center">
      <p class="text-sm font-semibold uppercase tracking-[0.2em] text-primary">
        {{ error.statusCode || 500 }}
      </p>
      <h1 class="text-3xl font-bold">页面暂时不可用</h1>
      <p class="text-muted">
        {{ error.statusMessage || error.message || "请稍后重试。" }}
      </p>
      <UButton to="/" icon="i-tabler-home">返回首页</UButton>
    </main>
  </UApp>
</template>
