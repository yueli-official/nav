<script setup lang="ts">
import type { AccountMenuAction } from "@yueli/ui/account-menu/pattern";

const { widthClass, brandName, brandTagline } = defineProps<{
  widthClass: string;
  brandName?: string;
  brandTagline?: string;
}>();
const { openSearch } = useNavigationSearch();
const { user } = useAuth();
const {
  accessStatus,
  accessError,
  refresh: refreshMe,
  clear: clearMe,
} = useMe();
const accountActions = computed<readonly AccountMenuAction[]>(() => {
  if (!user.value) return [];
  switch (accessStatus.value) {
    case "ready_manage":
      return [
        {
          label: "控制台",
          icon: "i-tabler-layout-dashboard",
          to: "/manage",
        },
      ];
    case "ready_requestable":
      return [
        {
          label: "申请成为内容维护者",
          description: "申请本站当前开放的维护角色",
          icon: "i-tabler-user-edit",
          to: "/contribute",
        },
      ];
    case "ready_no_access":
      return [
        {
          label: "本站暂无维护权限",
          description: "当前没有向你开放的维护角色",
          icon: "i-tabler-user-off",
          disabled: true,
        },
      ];
    case "suspended":
      return [
        {
          label: "本站成员资格已暂停",
          description: "公开导航仍可浏览；已认证操作暂不可用",
          icon: "i-tabler-user-pause",
          disabled: true,
        },
      ];
    case "error":
      return [
        {
          label: "本站权限状态不可用 · 重试",
          description: accessError.value?.correlationId
            ? `读取失败，参考编号：${accessError.value.correlationId}`
            : "未能读取本站授权",
          icon: "i-tabler-refresh-alert",
          onSelect: async () => {
            await refreshMe();
          },
        },
      ];
    case "resolving":
    case "signed_out":
      return [
        {
          label: "正在确认本站权限",
          description: "正在读取本站授权",
          icon: "i-tabler-loader-2",
          disabled: true,
        },
      ];
  }
});
if (import.meta.client) {
  watch(
    () =>
      user.value
        ? `${user.value.sub}:${user.value.userKey || ""}`
        : null,
    (identity) => {
      if (!identity) clearMe();
      else void refreshMe();
    },
    { immediate: true },
  );
}
</script>

<template>
  <header
    class="yueli-topbar sticky top-0 z-30 border-b border-default backdrop-blur-xl"
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

      <div class="ml-auto flex items-center gap-1">
        <UButton
          color="neutral"
          variant="soft"
          icon="i-tabler-search"
          label="搜索"
          @click="openSearch"
        />
        <UTooltip text="切换颜色模式">
          <UColorModeButton
            color="neutral"
            variant="ghost"
            aria-label="切换颜色模式"
          />
        </UTooltip>
        <ConsumerAccountControl :context-actions="accountActions" />
      </div>
    </div>
  </header>
</template>
