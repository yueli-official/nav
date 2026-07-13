<script setup lang="ts">
import { PlatformUserMenu } from "@platform/ui/components";
import type { PlatformUserMenuAction } from "@platform/ui/components";

const { widthClass, brandName, brandTagline } = defineProps<{
  widthClass: string;
  brandName?: string;
  brandTagline?: string;
}>();
const { user, loggedIn, isAdmin, login, logout } = useAuth();
const accountUrl = computed(
  () => useRuntimeConfig().public.accountUrl || "http://localhost:3000",
);
const contextActions = computed<PlatformUserMenuAction[]>(() =>
  isAdmin.value
    ? [{ label: "控制台", icon: "i-tabler-layout-dashboard", to: "/manage" }]
    : [],
);
const utilityActions = computed<PlatformUserMenuAction[]>(() => [
  {
    label: "用户设置",
    icon: "i-tabler-user-cog",
    onSelect: async () => {
      await navigateTo(accountUrl.value, { external: true });
    },
  },
]);
async function handleLogin() {
  await login();
}
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
            >{{ brandName }}</span
          >
          <span class="hidden text-xs text-muted sm:block">{{
            brandTagline
          }}</span>
        </span>
      </NuxtLink>

      <nav
        class="ml-auto hidden items-center gap-1 md:flex"
        aria-label="页面导航"
      >
        <UButton to="#search" color="neutral" variant="ghost" label="搜索" />
        <UButton
          to="#catalog"
          color="neutral"
          variant="ghost"
          label="分类目录"
        />
      </nav>

      <div class="ml-auto flex items-center gap-1 md:ml-2">
        <UTooltip text="切换颜色模式">
          <UColorModeButton
            color="neutral"
            variant="ghost"
            aria-label="切换颜色模式"
          />
        </UTooltip>
        <PlatformUserMenu
          v-if="loggedIn"
          :name="user?.name"
          :email="user?.email"
          :context-actions="contextActions"
          :utility-actions="utilityActions"
          :logout
        />
        <UButton
          v-else
          color="neutral"
          variant="ghost"
          icon="i-tabler-login-2"
          label="登录"
          @click="handleLogin"
        />
      </div>
    </div>
  </header>
</template>
