<script setup lang="ts">
import { ManageClientBoundary } from "~/utils/manageComponents";
import {
  CollectionPanel,
  type CollectionControl,
  type CollectionControlValue,
  type CollectionPanelMessages,
  type CollectionPanelState,
} from "@yueli/ui/collection/pattern";
import type {
  AdminNavigationLink,
  NavigationCheckExemptionResponse,
  NavigationCheckJob,
  NavigationCheckJobResponse,
  NavigationChecksResponse,
  NavigationHealthCounts,
  NavigationStartCheckJobResponse,
} from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "站点检查 · 月离导航" });

const ALL = "__all__";
const CHECK_JOB_STORAGE_KEY = "nav-check-job-id";
const CHECK_JOB_POLL_MS = 500;
const emptyCounts: NavigationHealthCounts = {
  all: 0,
  unchecked: 0,
  healthy: 0,
  redirected: 0,
  broken: 0,
  timeout: 0,
  error: 0,
  exempt: 0,
};
const { call } = useApi();
const { can } = useMe();
const canRunChecks = computed(() => can("nav.health_check.run"));
const search = ref("");
const q = ref("");
const health = ref("");
const page = ref(1);
const size = ref(15);
const selected = ref(new Set<string>());
const activeJob = ref<NavigationCheckJob>();
const startingJob = ref(false);
const resumingJob = ref(false);
const joinedExistingJob = ref(false);
const checkError = ref("");
const exemptionError = ref("");
const updatingExemptions = ref(new Set<string>());
let searchTimer: ReturnType<typeof setTimeout> | undefined;
let checkPollTimer: ReturnType<typeof setTimeout> | undefined;
let checkPollToken = 0;

watch(search, (value) => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    q.value = value.trim();
    page.value = 1;
  }, 250);
});
onScopeDispose(() => {
  clearTimeout(searchTimer);
  stopCheckJobPolling();
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
      checkableTotal: 0,
      page: 1,
      size: 15,
    }),
  },
);
const links = computed(() => data.value?.links ?? []);
const counts = computed(() => data.value?.counts ?? emptyCounts);
const total = computed(() => data.value?.total ?? 0);
const checkableTotal = computed(() => data.value?.checkableTotal ?? 0);
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
const healthOptions = computed(() => [
  { label: `全部 · ${counts.value.all}`, value: ALL },
  { label: `待处理 · ${issueCount.value}`, value: "issue" },
  { label: `正常 · ${counts.value.healthy}`, value: "healthy" },
  { label: `失效 · ${counts.value.broken}`, value: "broken" },
  { label: `超时 · ${counts.value.timeout}`, value: "timeout" },
  { label: `重定向 · ${counts.value.redirected}`, value: "redirected" },
  { label: `异常 · ${counts.value.error}`, value: "error" },
  { label: `未检查 · ${counts.value.unchecked}`, value: "unchecked" },
  { label: `免检 · ${counts.value.exempt}`, value: "exempt" },
]);
const controls = computed<CollectionControl[]>(() => [
  {
    kind: "select",
    id: "health",
    label: "检查状态",
    value: health.value || ALL,
    options: healthOptions.value,
    icon: "i-tabler-heartbeat",
    class: "w-40",
  },
]);
const activeFilterCount = computed(() => (health.value ? 1 : 0));
const pageSizes = [15, 30, 50] as const;
const panelState = computed<CollectionPanelState>(() =>
  error.value ? "error" : pending.value ? "loading" : "ready",
);
const messages = computed<CollectionPanelMessages>(() => ({
  searchPlaceholder: "搜索站点名称、地址或描述…",
  searchAction: "搜索",
  filtersAction: "筛选",
  activeFilters: (count) => `筛选（${count}）`,
  clearFilters: "清除筛选",
  selectPage: "选择当前页站点",
  selectItem: (label) => `选择站点：${label}`,
  bulkRegion: "链接检查批量操作",
  selected: (count) => `已选择 ${count} 个站点`,
  selectAllResults: "选择全部结果",
  clearSelection: "取消选择",
  emptyTitle:
    q.value || health.value ? "没有匹配的检查结果" : "还没有可检查的站点",
  emptyDescription:
    q.value || health.value
      ? "请调整搜索或检查状态后重试。"
      : "添加站点后即可在这里运行检查。",
  errorTitle: "检查列表加载失败",
  retry: "重新加载",
  showing: (first, last, count) => `显示 ${first}–${last}，共 ${count} 个`,
  pageSize: "每页",
  pageSizeControl: "每页检查数量",
  pageSizeOption: (value) => `${value} 个`,
}));
const pageIds = computed(() =>
  links.value.filter((link) => !link.healthCheckExempt).map((link) => link.id),
);
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
const checking = computed(
  () =>
    startingJob.value ||
    resumingJob.value ||
    activeJob.value?.status === "running",
);
const checkButtonLabel = computed(() => {
  if (activeJob.value?.status === "running") {
    return `正在检查 ${activeJob.value.completed}/${activeJob.value.total}`;
  }
  return q.value || health.value
    ? checkableTotal.value
      ? `检查筛选结果 ${checkableTotal.value} 项`
      : "当前筛选无可检查站点"
    : `检查全部 ${checkableTotal.value} 项`;
});
const checkAlertTitle = computed(() => {
  const job = activeJob.value;
  if (!job) return "";
  if (job.status === "failed")
    return `检查在 ${job.completed}/${job.total} 项后中断`;
  if (job.status === "completed") return `已完成 ${job.total} 个站点检查`;
  return `正在检查 ${job.completed}/${job.total} 个站点`;
});
const checkAlertDescription = computed(() => {
  const job = activeJob.value;
  if (!job) return "";
  if (job.status === "failed")
    return `${job.error || "检查任务未能完成，请稍后重试。"} 已完成的结果已经保留。`;
  if (job.status === "completed")
    return "结果已刷新，异常站点已集中到待处理筛选。";
  if (joinedExistingJob.value)
    return "已有检查正在运行，已接入当前进度。检查在后台并发进行，可以继续查看或筛选列表。";
  return "检查在后台并发进行，可以继续查看或筛选列表。";
});
const checkAlertColor = computed(() => {
  if (activeJob.value?.status === "failed") return "error" as const;
  if (activeJob.value?.status === "completed") return "success" as const;
  return "primary" as const;
});
const checkAlertIcon = computed(() => {
  if (activeJob.value?.status === "failed") return "i-tabler-alert-triangle";
  if (activeJob.value?.status === "completed") return "i-tabler-circle-check";
  return "i-tabler-heartbeat";
});
const linkKey = (link: NavigationChecksResponse["links"][number]) => link.id;
const linkLabel = (link: NavigationChecksResponse["links"][number]) =>
  link.title;
const isSelected = (id: string | number) => selected.value.has(String(id));
const isLinkSelectable = (link: AdminNavigationLink) =>
  canRunChecks.value && !link.healthCheckExempt;
const isUpdatingExemption = (id: string) =>
  updatingExemptions.value.has(id);

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
function submitSearch(value: string) {
  clearTimeout(searchTimer);
  search.value = value;
  q.value = value.trim();
  page.value = 1;
}
function changeControl(id: string, value: CollectionControlValue) {
  if (id === "health" && typeof value === "string") {
    health.value = value === ALL ? "" : value;
  }
}
function clearFilters() {
  health.value = "";
}
function changePageSize(value: number) {
  size.value = value;
  page.value = 1;
}

async function setCheckExemption(link: AdminNavigationLink, exempt: boolean) {
  if (
    !canRunChecks.value ||
    checking.value ||
    updatingExemptions.value.has(link.id)
  )
    return;

  exemptionError.value = "";
  const updating = new Set(updatingExemptions.value);
  updating.add(link.id);
  updatingExemptions.value = updating;
  try {
    await call<NavigationCheckExemptionResponse>(
      `/api/v1/admin/nav/checks/${link.id}/exemption`,
      { method: "PUT", body: { exempt } },
    );
    const nextSelection = new Set(selected.value);
    nextSelection.delete(link.id);
    selected.value = nextSelection;
    await refresh();
  } catch (failure) {
    exemptionError.value = apiFailureMessage(
      failure,
      exempt
        ? "无法将站点设为免检，请稍后重试。"
        : "无法恢复站点检查，请稍后重试。",
    );
  } finally {
    const next = new Set(updatingExemptions.value);
    next.delete(link.id);
    updatingExemptions.value = next;
  }
}

function stopCheckJobPolling() {
  checkPollToken++;
  clearTimeout(checkPollTimer);
  checkPollTimer = undefined;
}

function rememberCheckJob(id?: string) {
  if (!import.meta.client) return;
  try {
    if (id) sessionStorage.setItem(CHECK_JOB_STORAGE_KEY, id);
    else sessionStorage.removeItem(CHECK_JOB_STORAGE_KEY);
  } catch {
    // Session storage is optional; the in-page task still completes normally.
  }
}

function rememberedCheckJob() {
  if (!import.meta.client) return "";
  try {
    return sessionStorage.getItem(CHECK_JOB_STORAGE_KEY) || "";
  } catch {
    return "";
  }
}

function progressValueText(value: number | null | undefined, max: number) {
  return `已检查 ${value ?? 0}/${max} 个站点`;
}

function apiFailureMessage(failure: unknown, fallback: string) {
  const apiError = failure as { data?: { message?: string } };
  return apiError.data?.message || fallback;
}

async function acceptCheckJob(job: NavigationCheckJob, reused: boolean) {
  activeJob.value = job;
  joinedExistingJob.value = reused;
  checkError.value = "";
  if (job.status === "running") {
    rememberCheckJob(job.id);
    return;
  }
  rememberCheckJob();
  if (job.status === "completed") selected.value = new Set();
  await refresh();
}

async function pollCheckJob(jobId: string, token: number) {
  try {
    const result = await call<NavigationCheckJobResponse>(
      `/api/v1/admin/nav/checks/jobs/${jobId}`,
    );
    if (token !== checkPollToken) return;
    resumingJob.value = false;
    await acceptCheckJob(result.job, joinedExistingJob.value);
    if (result.job.status === "running" && token === checkPollToken) {
      checkPollTimer = setTimeout(
        () => void pollCheckJob(jobId, token),
        CHECK_JOB_POLL_MS,
      );
    }
  } catch (failure) {
    if (token !== checkPollToken) return;
    resumingJob.value = false;
    activeJob.value = undefined;
    joinedExistingJob.value = false;
    rememberCheckJob();
    checkError.value = apiFailureMessage(
      failure,
      "后台检查状态已失效，请重新运行检查。",
    );
    await refresh();
  }
}

function followCheckJob(job: NavigationCheckJob, reused: boolean) {
  stopCheckJobPolling();
  const token = checkPollToken;
  void acceptCheckJob(job, reused).then(() => {
    if (job.status === "running" && token === checkPollToken) {
      checkPollTimer = setTimeout(
        () => void pollCheckJob(job.id, token),
        CHECK_JOB_POLL_MS,
      );
    }
  });
}

onMounted(() => {
  const jobId = rememberedCheckJob();
  if (!jobId) return;
  stopCheckJobPolling();
  resumingJob.value = true;
  void pollCheckJob(jobId, checkPollToken);
});

async function runChecks(scope: "filtered" | "selected", ids: string[] = []) {
  const selectedIDs = ids.length ? ids : [...selected.value];
  if (
    !canRunChecks.value ||
    (scope === "filtered" && !checkableTotal.value) ||
    checking.value ||
    (scope === "selected" && !selectedIDs.length)
  )
    return;
  startingJob.value = true;
  checkError.value = "";
  try {
    const result = await call<NavigationStartCheckJobResponse>(
      "/api/v1/admin/nav/checks/run",
      {
        method: "POST",
        body:
          scope === "selected"
            ? { scope: "selected", ids: selectedIDs }
            : {
                scope: "filtered",
                q: q.value || undefined,
                health: health.value || undefined,
              },
      },
    );
    followCheckJob(result.job, result.reused);
  } catch (failure) {
    checkError.value = apiFailureMessage(
      failure,
      "无法启动检查，请稍后重试。",
    );
  } finally {
    startingJob.value = false;
  }
}
</script>

<template>
  <ManagePage
    id="checks"
    title="链接检查"
    description="按结果筛选站点可访问性；受地区或代理限制的站点可设为免检，最后一次结果仍会保留。"
    icon="i-tabler-heartbeat"
    main-id="manage-main"
    body-class="flex min-h-0 w-full flex-col gap-5"
  >
    <template #actions>
      <UButton
        icon="i-tabler-heartbeat"
        :label="checkButtonLabel"
        :loading="checking"
        :disabled="!canRunChecks || !checkableTotal || checking"
        @click="runChecks('filtered')"
      />
    </template>

    <ManageClientBoundary :rows="8">
      <div class="flex min-h-0 flex-1 flex-col gap-3">
        <UAlert
          v-if="activeJob"
          data-check-job-progress
          role="status"
          aria-live="polite"
          :color="checkAlertColor"
          variant="subtle"
          :icon="checkAlertIcon"
          :title="checkAlertTitle"
        >
          <template #description>
            <p>{{ checkAlertDescription }}</p>
            <UProgress
              v-if="activeJob.status === 'running'"
              class="mt-2"
              size="sm"
              color="primary"
              :model-value="activeJob.completed"
              :max="activeJob.total"
              :get-value-text="progressValueText"
            />
          </template>
          <template v-if="activeJob.status === 'failed'" #actions>
            <UButton
              :label="
                activeJob.scope === 'selected'
                  ? '重新检查已选站点'
                  : '重新检查'
              "
              icon="i-tabler-refresh"
              size="xs"
              color="error"
              variant="soft"
              :disabled="
                activeJob.scope === 'selected' && selectedCount === 0
              "
              @click="runChecks(activeJob.scope)"
            />
          </template>
        </UAlert>

        <UAlert
          v-else-if="checkError"
          color="error"
          variant="subtle"
          icon="i-tabler-alert-triangle"
          title="检查没有启动"
          :description="checkError"
        />

        <UAlert
          v-if="exemptionError"
          color="error"
          variant="subtle"
          icon="i-tabler-alert-triangle"
          title="免检状态没有保存"
          :description="exemptionError"
        />

        <CollectionPanel
          data-check-list-panel
          class="flex min-h-0 flex-1 flex-col [&>[aria-live=polite]]:min-h-0 [&>[aria-live=polite]]:flex-1 [&>[aria-live=polite]]:overflow-y-auto [&>[aria-live=polite]]:overscroll-contain [&>footer]:shrink-0"
          v-model:search="search"
          :items="links"
          :item-key="linkKey"
          :item-label="linkLabel"
          :controls="controls"
          :messages="messages"
          :state="panelState"
          error-message="请检查 Nav API、数据库和登录状态。"
          :total="total"
          :page="page"
          :page-size="size"
          :page-sizes="pageSizes"
          :active-filter-count="activeFilterCount"
          :selection-count="selectedCount"
          :page-selected="isPageSelected"
          :page-indeterminate="isPageIndeterminate"
          :is-selected="isSelected"
          :is-item-selectable="isLinkSelectable"
          label="链接检查列表"
          selectable
          @search="submitSearch"
          @control-change="changeControl"
          @clear-filters="clearFilters"
          @retry="refresh"
          @toggle-page="togglePage"
          @toggle-item="toggle"
          @clear-selection="selected = new Set()"
          @page-change="page = $event"
          @page-size-change="changePageSize"
        >
          <template #view>
            <UButton
              icon="i-tabler-refresh"
              color="neutral"
              variant="outline"
              size="xs"
              aria-label="刷新检查结果"
              @click="refresh()"
            />
          </template>

          <template #columns>
            <div class="grid grid-cols-[minmax(0,1fr)_auto] gap-3">
              <span>站点与地址</span>
              <span class="hidden w-60 text-right md:block"
                >检查结果与操作</span
              >
            </div>
          </template>

          <template #bulk-actions>
            <UButton
              :label="`检查已选 ${selectedCount} 项`"
              icon="i-tabler-heartbeat"
              size="xs"
              color="primary"
              variant="soft"
              :loading="checking"
              :disabled="!canRunChecks || checking"
              @click="runChecks('selected')"
            />
          </template>

          <template #item="{ item: link }">
            <div
              class="grid min-w-0 gap-3 md:grid-cols-[minmax(0,1fr)_minmax(12rem,18rem)_auto] md:items-center"
            >
              <div class="min-w-0">
                <div class="flex min-w-0 items-center gap-2">
                  <span
                    class="grid size-7 shrink-0 place-items-center overflow-hidden rounded-md bg-primary/10 text-primary"
                  >
                    <NavigationFavicon
                      :id="link.id"
                      :title="link.title"
                      :revision="link.faviconRevision"
                      image-class="size-5"
                    />
                  </span>
                  <p class="truncate text-sm font-medium text-highlighted">
                    {{ link.title }}
                  </p>
                </div>
                <p class="mt-1 truncate text-xs text-muted">{{ link.url }}</p>
                <p
                  v-if="link.healthError && !link.healthCheckExempt"
                  class="mt-1 line-clamp-1 text-xs text-error"
                >
                  {{ link.healthError }}
                </p>
              </div>
              <div
                class="flex flex-wrap items-center gap-2 text-xs text-muted md:justify-end"
              >
                <template v-if="link.healthCheckExempt">
                  <span
                    v-if="
                      link.healthStatus && link.healthStatus !== 'unchecked'
                    "
                  >
                    上次：{{ healthMeta[link.healthStatus || "unchecked"].label }}
                  </span>
                  <UBadge
                    label="免检"
                    color="neutral"
                    icon="i-tabler-shield-off"
                    variant="subtle"
                    size="sm"
                  />
                </template>
                <span v-else-if="link.healthHttpStatus" class="tabular-nums"
                  >HTTP {{ link.healthHttpStatus }}</span
                >
                <span
                  v-if="!link.healthCheckExempt && link.healthLatencyMs"
                  class="tabular-nums"
                  >{{ link.healthLatencyMs }} ms</span
                >
                <UBadge
                  v-if="!link.healthCheckExempt"
                  v-bind="healthMeta[link.healthStatus || 'unchecked']"
                  variant="subtle"
                  size="sm"
                />
              </div>
              <div class="flex flex-wrap items-center justify-end gap-1">
                <UButton
                  :label="link.healthCheckExempt ? '恢复检查' : '设为免检'"
                  :icon="
                    link.healthCheckExempt
                      ? 'i-tabler-shield-check'
                      : 'i-tabler-shield-off'
                  "
                  color="neutral"
                  :variant="link.healthCheckExempt ? 'soft' : 'ghost'"
                  size="xs"
                  :loading="isUpdatingExemption(link.id)"
                  :disabled="checking || !canRunChecks"
                  @click="setCheckExemption(link, !link.healthCheckExempt)"
                />
                <UButton
                  v-if="!link.healthCheckExempt"
                  icon="i-tabler-refresh"
                  color="neutral"
                  variant="ghost"
                  size="sm"
                  square
                  :aria-label="`重新检查 ${link.title}`"
                  :disabled="checking || !canRunChecks"
                  @click="runChecks('selected', [link.id])"
                />
              </div>
            </div>
          </template>
        </CollectionPanel>
      </div>
    </ManageClientBoundary>
  </ManagePage>
</template>
