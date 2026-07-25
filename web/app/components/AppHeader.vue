<script setup lang="ts">
import type { AccountMenuAction } from "@yueli/ui/account-menu/pattern";

const { widthClass, brandName, brandTagline } = defineProps<{
  widthClass: string;
  brandName?: string;
  brandTagline?: string;
}>();
const { openSearch } = useNavigationSearch();
const { user } = useAuth();
const { canManage, refresh: refreshMe } = useMe();
const accountActions = computed<readonly AccountMenuAction[]>(() =>
  canManage.value
    ? [
        {
          label: "控制台",
          icon: "i-tabler-layout-dashboard",
          to: "/manage",
        },
      ]
    : user.value
      ? [
          {
            label: "申请维护导航",
            icon: "i-tabler-user-edit",
            to: "/contribute",
          },
        ]
      : [],
);
onMounted(() => {
  if (user.value) void refreshMe();
});
watch(user, (value) => {
  if (value) void refreshMe();
});
</script>

<template>
  <header
    class="yueli-topbar sticky top-0 z-30 border-b border-default bg-default/88 backdrop-blur-xl"
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
