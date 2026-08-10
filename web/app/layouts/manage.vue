<script setup lang="ts">
import type {
  AdminNavigationItem,
  AdminSearchGroup,
  AdminShellMessages,
} from "@yueli/ui/admin";

const route = useRoute();
const { brand } = useSiteRuntime();
const sidebarOpen = ref(false);
const { can, isAdministrator } = useMe();
const sidebarAppearance = computed(() =>
  route.query.sidebar === "baseline" ? "framed" : "commercial",
);

const messages: AdminShellMessages = {
  skipToContent: "跳到主要内容",
  search: "搜索后台",
  searchPlaceholder: "搜索页面与常用操作",
};

function closeSidebar() {
  sidebarOpen.value = false;
}

function active(path: string, exact = false) {
  return exact ? route.path === path : route.path.startsWith(path);
}

const navigation = computed<readonly AdminNavigationItem[]>(() => [
  ...(can("nav.link.update") || can("nav.link.submit")
    ? [
        {
          label: "站点链接",
          icon: "i-tabler-world-www",
          to: "/manage",
          active: active("/manage", true),
          onSelect: closeSidebar,
        },
      ]
    : []),
  ...(can("nav.structure.manage")
    ? [
        {
          label: "分类与主题",
          icon: "i-tabler-folders",
          to: "/manage/categories",
          active: active("/manage/categories"),
          onSelect: closeSidebar,
        },
        {
          label: "标签治理",
          icon: "i-tabler-hash",
          to: "/manage/tags",
          active: active("/manage/tags"),
          onSelect: closeSidebar,
        },
      ]
    : []),
  ...(can("nav.health_check.run")
    ? [
        {
          label: "链接检查",
          icon: "i-tabler-heartbeat",
          to: "/manage/checks",
          active: active("/manage/checks"),
          onSelect: closeSidebar,
        },
      ]
    : []),
  ...(can("nav.settings.manage")
    ? [
        {
          label: "站点设置",
          icon: "i-tabler-settings",
          to: "/manage/settings",
          active: active("/manage/settings"),
          onSelect: closeSidebar,
        },
      ]
    : []),
  ...(isAdministrator.value
    ? [
        {
          label: "成员",
          icon: "i-tabler-users",
          to: "/manage/members",
          active: active("/manage/members"),
          onSelect: closeSidebar,
        },
        {
          label: "权限策略",
          icon: "i-tabler-shield-lock",
          to: "/manage/authorization",
          active: active("/manage/authorization"),
          onSelect: closeSidebar,
        },
      ]
    : []),
]);

const searchGroups = computed<readonly AdminSearchGroup[]>(() => {
  const pageItems = navigation.value.map((item, index) => ({
    id: `nav-page-${index}`,
    label: item.label,
    icon: item.icon,
    to: item.to,
  }));
  const actionItems = [
    ...(can("nav.link.submit")
      ? [
          {
            id: "new-link",
            label: "添加站点链接",
            icon: "i-tabler-world-plus",
            to: "/manage?action=create",
          },
        ]
      : []),
    ...(can("nav.health_check.run")
      ? [
          {
            id: "run-checks",
            label: "运行链接检查",
            icon: "i-tabler-heartbeat",
            to: "/manage/checks",
          },
        ]
      : []),
    ...(isAdministrator.value
      ? [
          {
            id: "authorization",
            label: "权限策略",
            icon: "i-tabler-shield-lock",
            to: "/manage/authorization",
          },
        ]
      : []),
  ];
  return [
    { id: "manage-pages", label: "管理页面", items: pageItems },
    ...(actionItems.length
      ? [{ id: "manage-actions", label: "常用操作", items: actionItems }]
      : []),
  ];
});

</script>

<template>
  <ClientOnly>
    <YAdminShell
      v-model:open="sidebarOpen"
      :navigation="navigation"
      :search-groups="searchGroups"
      :messages="messages"
      :sidebar-appearance="sidebarAppearance"
      storage-key="nav-manage"
      main-id="manage-main"
      :default-size="16"
      :min-size="14"
      :max-size="20"
    >
      <template #brand="{ collapsed }">
        <UButton
          to="/"
          color="neutral"
          variant="ghost"
          :block="!collapsed"
          :square="collapsed"
          :aria-label="`返回${brand}首页`"
          :class="[
            'min-h-11 gap-2 px-1.5',
            !collapsed && 'w-full justify-start',
            collapsed && 'aspect-square justify-center px-0',
          ]"
          @click="closeSidebar"
        >
          <span
            class="grid size-7 shrink-0 place-items-center rounded-md bg-primary/10 text-primary"
          >
            <UIcon name="i-tabler-compass" class="size-4" />
          </span>
          <span
            v-if="!collapsed"
            class="min-w-0 truncate text-sm font-semibold text-highlighted"
          >
            {{ brand }}
          </span>
        </UButton>
      </template>

      <template #sidebar-footer="{ collapsed }">
        <ConsumerManageAccountControl
          home-to=""
          show-appearance
          :trigger-mode="collapsed ? 'collapsed' : 'sidebar'"
        />
      </template>

      <slot />
      <YBackToTop
        target-id="manage-main"
        scroll-container-id="manage-main"
        avoid-selector="[data-manage-dock], [data-back-to-top-avoid]"
        label="返回顶部"
      />
      <ManageSidebarPrototypeSwitcher />
    </YAdminShell>

    <template #fallback>
      <div
        class="fixed inset-0 grid place-items-center bg-default text-sm text-muted"
        role="status"
      >
        正在打开控制台
      </div>
    </template>
  </ClientOnly>
</template>
