<script setup lang="ts">
import type {
  AdminNavigationItem,
  AdminSearchGroup,
  AdminShellMessages,
} from "@yueli/ui/admin";

const route = useRoute();
const { brand } = useSiteRuntime();
const { can, isAdministrator } = useMe();

const currentLabel = computed(() => {
  if (route.path === "/manage") return "站点链接";
  if (route.path.startsWith("/manage/categories")) return "分类与主题";
  if (route.path.startsWith("/manage/tags")) return "标签治理";
  if (route.path.startsWith("/manage/checks")) return "链接检查";
  if (route.path.startsWith("/manage/settings")) return "站点设置";
  if (route.path.startsWith("/manage/members")) return "成员";
  if (route.path.startsWith("/manage/authorization")) return "权限策略";
  return "站点链接";
});

const messages: AdminShellMessages = {
  skipToContent: "跳到主要内容",
  search: "搜索控制台",
  searchPlaceholder: "搜索页面与常用操作",
  currentLocation: "当前位置",
};

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
        },
        {
          label: "标签治理",
          icon: "i-tabler-hash",
          to: "/manage/tags",
          active: active("/manage/tags"),
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
        },
        {
          label: "权限策略",
          icon: "i-tabler-shield-lock",
          to: "/manage/authorization",
          active: active("/manage/authorization"),
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
  <YAdminConsoleLayout
    :navigation="navigation"
    :search-groups="searchGroups"
    :messages="messages"
    storage-key="nav-manage"
    main-id="manage-main"
    :brand-label="brand"
    brand-icon="i-tabler-compass"
    brand-to="/"
    :context-label="brand"
    :current-label="currentLabel"
    back-to-top-label="返回顶部"
    data-nav-manage-shell
  >
    <template #account="{ collapsed }">
      <ConsumerManageAccountControl
        home-to=""
        show-appearance
        :trigger-mode="collapsed ? 'collapsed' : 'sidebar'"
      />
    </template>
    <slot />
  </YAdminConsoleLayout>
</template>

<style scoped>
@media (max-width: 640px) {
  [data-nav-manage-shell] :deep(button),
  [data-nav-manage-shell] :deep(a[href]),
  [data-nav-manage-shell] :deep(summary) {
    min-height: 44px;
  }

  [data-nav-manage-shell] :deep(button[aria-label]),
  [data-nav-manage-shell] :deep(a[aria-label]) {
    min-width: 44px;
  }
}
</style>
