<script setup lang="ts">
import { ManageShell, ManageUserMenu } from "@platform/manage/components";

const { user, logout } = useAuth();
const { brand } = useSiteRuntime();
const accountUrl = computed(
  () => useRuntimeConfig().public.accountUrl || "http://localhost:3000",
);
</script>

<template>
  <ManageShell
    :site-name="brand"
    context-label="站点管理"
    home-to="/manage"
    storage-key="nav-manage"
    shell-class="platform-app-shell"
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
