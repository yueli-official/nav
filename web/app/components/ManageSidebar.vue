<script setup lang="ts">
const route = useRoute();
const { brand } = useSiteRuntime();
const nav = [
  { label: "站点", icon: "i-tabler-world-www", to: "/manage" },
  { label: "分类与主题", icon: "i-tabler-folders", to: "/manage/categories" },
  { label: "标签", icon: "i-tabler-hash", to: "/manage/tags" },
  { label: "设置", icon: "i-tabler-settings", to: "/manage/settings" },
];

function isActive(to: string) {
  return to === "/manage" ? route.path === to : route.path.startsWith(to);
}
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
        v-for="item in nav"
        :key="item.to"
        :to="item.to"
        class="flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition"
        :class="
          isActive(item.to)
            ? 'bg-primary/10 text-primary'
            : 'text-muted hover:bg-elevated hover:text-default'
        "
      >
        <UIcon :name="item.icon" class="size-5 shrink-0" />
        <span>{{ item.label }}</span>
      </NuxtLink>
    </nav>
  </div>
</template>
