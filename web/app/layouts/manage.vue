<script setup lang="ts">
import { ManageShell, ManageUserMenu } from "@platform/manage/components";

const { user, logout } = useAuth();
const { brand } = useSiteRuntime();
const accountUrl = computed(
  () => useRuntimeConfig().public.accountUrl || "http://localhost:3000",
);
const route = useRoute();
const showBackToTop = computed(() => route.path === "/manage/settings");
const contextLabel = computed(() => {
  if (route.path.startsWith("/manage/categories")) return "分类与主题";
  if (route.path.startsWith("/manage/tags")) return "标签管理";
  if (route.path.startsWith("/manage/checks")) return "站点检查";
  if (route.path.startsWith("/manage/settings")) return "站点设置";
  return "站点管理";
});
</script>

<template>
  <ManageShell
    :site-name="brand"
    :context-label="contextLabel"
    home-to="/manage"
    storage-key="nav-manage"
    shell-class="platform-app-shell"
    :show-back-to-top="showBackToTop"
  >
    <template #sidebar>
      <ManageSidebar />
    </template>
    <template #user>
      <ManageUserMenu
        :name="user?.name"
        :email="user?.email"
        :settings-to="accountUrl"
        :logout
      />
    </template>
    <slot />
  </ManageShell>
</template>
