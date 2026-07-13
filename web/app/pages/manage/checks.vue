<script setup lang="ts">
import {
  ManageCollectionFooter,
  ManageCollectionToolbar,
  ManageEmpty,
  ManageHeader,
  ManagePageSelection,
  ManageTabs,
  SkeletonList,
} from "@platform/manage/components";
import type {
  NavigationChecksResponse,
  NavigationHealthCounts,
} from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "站点检查 · 月离导航" });

const emptyCounts: NavigationHealthCounts = {
  all: 0,
  unchecked: 0,
  healthy: 0,
  redirected: 0,
  broken: 0,
  timeout: 0,
  error: 0,
};
const { call } = useApi();
const { isAdmin } = useAuth();
const search = ref("");
const q = ref("");
const health = ref("");
const page = ref(1);
const size = ref(15);
const selected = ref(new Set<string>());
const checking = ref(false);
const checkMessage = ref("");
let searchTimer: ReturnType<typeof setTimeout> | undefined;

watch(search, (value) => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    q.value = value.trim();
    page.value = 1;
  }, 250);
});
watch(health, () => {
  page.value = 1;
  selected.value = new Set();
});
watch([q, page, size], () => {
  selected.value = new Set();
});

const { data, pending, error, refresh } = await useAsyncData(
  "nav-health-checks",
  () =>
    call<NavigationChecksResponse>("/api/v1/admin/nav/checks", {
      query: {
        q: q.value || undefined,
        health: health.value || undefined,
        page: page.value,
        size: size.value,
      },
    }),
  {
    server: false,
    watch: [q, health, page, size],
    default: () => ({
      links: [],
      counts: emptyCounts,
      total: 0,
      page: 1,
      size: 15,
    }),
  },
);
const links = computed(() => data.value?.links ?? []);
const counts = computed(() => data.value?.counts ?? emptyCounts);
const total = computed(() => data.value?.total ?? 0);
const totalPages = computed(() =>
  Math.max(1, Math.ceil(total.value / size.value)),
);
const issueCount = computed(
  () =>
    counts.value.redirected +
    counts.value.broken +
    counts.value.timeout +
    counts.value.error,
);
const tabs = computed(() => [
  { key: "", label: "全部", count: counts.value.all },
  { key: "issue", label: "待处理", count: issueCount.value },
  { key: "unchecked", label: "未检查", count: counts.value.unchecked },
  { key: "healthy", label: "正常", count: counts.value.healthy },
  { key: "redirected", label: "重定向", count: counts.value.redirected },
  { key: "broken", label: "失效", count: counts.value.broken },
]);
const pageIds = computed(() => links.value.map((link) => link.id));
const isPageSelected = computed(
  () =>
    pageIds.value.length > 0 &&
    pageIds.value.every((id) => selected.value.has(id)),
);
const isPageIndeterminate = computed(
  () =>
    !isPageSelected.value && pageIds.value.some((id) => selected.value.has(id)),
);
const selectedCount = computed(() => selected.value.size);

const healthMeta = {
  unchecked: {
    label: "未检查",
    color: "neutral",
    icon: "i-tabler-clock-question",
  },
  healthy: { label: "正常", color: "success", icon: "i-tabler-circle-check" },
  redirected: { label: "重定向", color: "warning", icon: "i-tabler-route" },
  broken: { label: "失效", color: "error", icon: "i-tabler-link-off" },
  timeout: {
    label: "超时",
    color: "warning",
    icon: "i-tabler-clock-exclamation",
  },
  error: { label: "异常", color: "error", icon: "i-tabler-alert-triangle" },
} as const;

function togglePage(value: boolean) {
  const next = new Set(selected.value);
  pageIds.value.forEach((id) => {
    if (value) next.add(id);
    else next.delete(id);
  });
  selected.value = next;
}
function toggle(id: string, value: boolean) {
  const next = new Set(selected.value);
  if (value) next.add(id);
  else next.delete(id);
  selected.value = next;
}
async function runChecks() {
  const ids = selectedCount.value ? [...selected.value] : pageIds.value;
  if (!isAdmin.value || !ids.length || checking.value) return;
  checking.value = true;
  checkMessage.value = "";
  try {
    const result = await call<{ checked: number }>(
      "/api/v1/admin/nav/checks/run",
      { method: "POST", body: { ids } },
    );
    checkMessage.value = `已完成 ${result.checked} 个站点检查`;
    selected.value = new Set();
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    checkMessage.value = apiError.data?.message || "检查失败，请稍后重试。";
  } finally {
    checking.value = false;
  }
}
</script>

<template>
  <div class="space-y-6">
    <ManageHeader title="站点检查">
      <template #subtitle
        >定期验证站点可访问性；异常结果集中在“待处理”，便于修正或归档。</template
      >
      <template #actions>
        <UButton
          icon="i-tabler-heartbeat"
          :label="selectedCount ? `检查已选 ${selectedCount} 项` : '检查当前页'"
          :loading="checking"
          :disabled="!isAdmin || !links.length"
          @click="runChecks"
        />
      </template>
    </ManageHeader>

    <ManageTabs v-model="health" :items="tabs" />
    <ManageCollectionToolbar
      v-model:search="search"
      search-placeholder="搜索站点名称、地址或描述…"
    >
      <template #actions
        ><UButton
          icon="i-tabler-refresh"
          color="neutral"
          variant="outline"
          aria-label="刷新检查结果"
          @click="refresh()"
      /></template>
    </ManageCollectionToolbar>
    <UAlert
      v-if="checkMessage"
      color="neutral"
      variant="subtle"
      icon="i-tabler-info-circle"
      :description="checkMessage"
    />
    <UAlert
      v-if="error"
      color="error"
      icon="i-tabler-alert-circle"
      title="检查列表加载失败"
      description="请检查 Nav API 与数据库状态。"
    />
    <SkeletonList v-else-if="pending" :rows="8" />
    <ManageEmpty
      v-else-if="!links.length"
      icon="i-tabler-heartbeat"
      :text="q || health ? '没有匹配的检查结果' : '还没有可检查的站点'"
    />

    <div
      v-else
      class="overflow-hidden rounded-xl border border-default bg-default divide-y divide-default"
    >
      <article
        v-for="link in links"
        :key="link.id"
        class="grid gap-3 px-4 py-3 sm:grid-cols-[auto_minmax(0,1fr)_auto] sm:items-center"
      >
        <UCheckbox
          :model-value="selected.has(link.id)"
          :aria-label="`选择 ${link.title}`"
          @update:model-value="toggle(link.id, Boolean($event))"
        />
        <div class="min-w-0">
          <div class="flex min-w-0 items-center gap-2">
            <span
              class="grid size-6 shrink-0 place-items-center overflow-hidden rounded bg-primary/10 text-primary"
            >
              <NavigationFavicon
                :id="link.id"
                :title="link.title"
                image-class="size-5"
              />
            </span>
            <p class="truncate text-sm font-medium text-highlighted">
              {{ link.title }}
            </p>
          </div>
          <p class="mt-1 truncate text-xs text-muted">{{ link.url }}</p>
          <p
            v-if="link.healthError"
            class="mt-1 line-clamp-1 text-xs text-error"
          >
            {{ link.healthError }}
          </p>
        </div>
        <div
          class="flex flex-wrap items-center justify-end gap-2 text-xs text-muted"
        >
          <span v-if="link.healthHttpStatus" class="tabular-nums"
            >HTTP {{ link.healthHttpStatus }}</span
          >
          <span v-if="link.healthLatencyMs" class="tabular-nums"
            >{{ link.healthLatencyMs }} ms</span
          >
          <UBadge
            v-bind="healthMeta[link.healthStatus || 'unchecked']"
            variant="subtle"
            size="sm"
          />
          <UButton
            icon="i-tabler-refresh"
            color="neutral"
            variant="ghost"
            size="sm"
            square
            :aria-label="`重新检查 ${link.title}`"
            :disabled="checking || !isAdmin"
            @click="
              selected = new Set([link.id]);
              runChecks();
            "
          />
        </div>
      </article>
    </div>

    <ManageCollectionFooter
      v-if="total"
      v-model:page="page"
      v-model:size="size"
      :total="total"
      :total-pages="totalPages"
      label="检查选择与分页"
    >
      <template #selection>
        <ManagePageSelection
          :model-value="isPageSelected"
          :indeterminate="isPageIndeterminate"
          label="选择当前页站点"
          @update:model-value="togglePage"
        />
        <span class="text-xs">{{
          selectedCount ? `已选 ${selectedCount} 项` : `共 ${total} 项`
        }}</span>
      </template>
    </ManageCollectionFooter>
  </div>
</template>
