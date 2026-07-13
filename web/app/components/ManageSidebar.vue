<script setup lang="ts">
const route = useRoute();
const { brand } = useSiteRuntime();
</script>

<template>
  <div class="flex h-full flex-col bg-elevated/30">
    <NuxtLink
      to="/manage"
      class="font-display flex h-16 items-center gap-2 border-b border-default px-5 font-semibold text-highlighted"
    >
      <span
        class="grid size-8 place-items-center rounded-lg bg-primary/10 text-primary"
      >
        <UIcon name="i-tabler-compass" class="size-5" />
      </span>
      {{ brand }}
    </NuxtLink>

    <nav aria-label="导航站后台" class="flex-1 space-y-1 p-3">
      <NuxtLink
        to="/manage"
        class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition"
        :class="
          route.path === '/manage'
            ? 'bg-primary/10 text-primary'
            : 'text-muted hover:bg-elevated hover:text-default'
        "
      >
        <UIcon name="i-tabler-world-www" class="size-5 shrink-0" />
        <span>站点管理</span>
      </NuxtLink>
    </nav>
  </div>
</template>
