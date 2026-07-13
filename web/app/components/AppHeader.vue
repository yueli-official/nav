<script setup lang="ts">
const { widthClass } = defineProps<{ widthClass: string }>();
const { brand } = useSiteRuntime();
</script>

<template>
  <header
    class="platform-topbar sticky top-0 z-30 border-b border-default bg-default/88 backdrop-blur-xl"
  >
    <div
      class="mx-auto flex h-16 w-full items-center gap-3 px-4 sm:px-6 lg:px-8"
      :class="widthClass"
    >
      <NuxtLink
        to="/"
        class="group flex min-w-0 items-center gap-3 rounded-lg focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-primary"
      >
        <span
          class="grid size-9 shrink-0 place-items-center rounded-xl bg-primary text-inverted shadow-sm"
        >
          <UIcon name="i-tabler-compass" class="size-5" />
        </span>
        <span class="min-w-0">
          <span
            class="font-display block truncate text-base font-semibold tracking-tight text-highlighted"
            >{{ brand }}</span
          >
          <span class="hidden text-xs text-muted sm:block"
            >精选互联网工作台</span
          >
        </span>
      </NuxtLink>

      <nav
        class="ml-auto hidden items-center gap-1 md:flex"
        aria-label="页面导航"
      >
        <UButton to="#featured" color="neutral" variant="ghost" label="精选" />
        <UButton to="#catalog" color="neutral" variant="ghost" label="分类" />
      </nav>

      <div class="ml-auto flex items-center gap-1 md:ml-2">
        <UButton
          to="/manage"
          color="neutral"
          variant="ghost"
          icon="i-tabler-settings"
          label="管理"
        />
        <UTooltip text="切换颜色模式">
          <UColorModeButton
            color="neutral"
            variant="ghost"
            aria-label="切换颜色模式"
          />
        </UTooltip>
      </div>
    </div>
  </header>
</template>
