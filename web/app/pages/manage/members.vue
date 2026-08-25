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
  NavigationMember,
  NavigationMemberCounts,
  NavigationMembersResponse,
} from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "成员 · 月离导航" });

const ALL = "__all__";
const emptyCounts: NavigationMemberCounts = {
  all: 0,
  active: 0,
  suspended: 0,
};
const { call } = useApi();
const { me } = useMe();
const toast = createNavNotifier(useToast());
const runtime = useRuntimeConfig();
const search = ref("");
const q = ref("");
const status = ref("");
const role = ref("");
const page = ref(1);
const size = ref(15);
const selected = ref<NavigationMember>();
const detailOpen = ref(false);
const suspendOpen = ref(false);
const suspensionReason = ref("");
const suspensionReasonTouched = ref(false);
const saving = ref(false);
const avatarFailures = ref(new Set<string>());
let searchTimer: ReturnType<typeof setTimeout> | undefined;

watch(search, (value) => {
  clearTimeout(searchTimer);
  searchTimer = setTimeout(() => {
    q.value = value.trim();
    page.value = 1;
  }, 250);
});
onScopeDispose(() => clearTimeout(searchTimer));
watch([status, role], () => {
  page.value = 1;
});

const { data, pending, error, refresh } = await useAsyncData(
  "nav-members",
  () =>
    call<NavigationMembersResponse>("/api/v1/admin/nav/members", {
      query: {
        q: q.value || undefined,
        status: status.value || undefined,
        role: role.value || undefined,
        page: page.value,
        size: size.value,
      },
    }),
  {
    server: false,
    watch: [q, status, role, page, size],
    default: () => ({
      members: [],
      counts: emptyCounts,
      roles: [],
      total: 0,
      page: 1,
      size: 15,
    }),
  },
);
const members = computed(() => data.value?.members ?? []);
const counts = computed(() => data.value?.counts ?? emptyCounts);
const total = computed(() => data.value?.total ?? 0);
const roleOptions = computed(() => [
  { label: "全部角色", value: ALL },
  ...(data.value?.roles ?? []).map((item) => ({
    label: item.displayName,
    value: item.key,
  })),
]);
const statusOptions = computed(() => [
  { label: `全部成员 · ${counts.value.all}`, value: ALL },
  { label: `正常 · ${counts.value.active}`, value: "active" },
  { label: `已暂停 · ${counts.value.suspended}`, value: "suspended" },
]);
const controls = computed<CollectionControl[]>(() => [
  {
    kind: "select",
    id: "status",
    label: "成员状态",
    value: status.value || ALL,
    options: statusOptions.value,
    icon: "i-tabler-user-check",
    class: "w-40",
  },
  {
    kind: "select",
    id: "role",
    label: "本站角色",
    value: role.value || ALL,
    options: roleOptions.value,
    icon: "i-tabler-shield-check",
    class: "w-44",
  },
]);
const activeFilterCount = computed(
  () => Number(Boolean(status.value)) + Number(Boolean(role.value)),
);
const panelState = computed<CollectionPanelState>(() =>
  error.value ? "error" : pending.value ? "loading" : "ready",
);
const messages = computed<CollectionPanelMessages>(() => ({
  searchPlaceholder: "搜索显示名、账号或用户键…",
  searchAction: "搜索",
  filtersAction: "筛选",
  activeFilters: (count) => `筛选（${count}）`,
  clearFilters: "清除筛选",
  selectPage: "选择当前页成员",
  selectItem: (label) => `选择成员：${label}`,
  bulkRegion: "成员批量操作",
  selected: (count) => `已选择 ${count} 位成员`,
  selectAllResults: "选择全部结果",
  clearSelection: "取消选择",
  emptyTitle:
    q.value || status.value || role.value
      ? "没有匹配的导航成员"
      : "还没有导航成员",
  emptyDescription:
    q.value || status.value || role.value
      ? "请调整搜索、成员状态或本站角色后重试。"
      : "用户完成登录并首次以已认证身份进入 Nav 后，会自动出现在这里。",
  errorTitle: "成员列表加载失败",
  retry: "重新加载",
  showing: (first, last, count) => `显示 ${first}–${last}，共 ${count} 位`,
  pageSize: "每页",
  pageSizeControl: "每页成员数量",
  pageSizeOption: (value) => `${value} 位`,
}));
const pageSizes = [15, 30, 50] as const;
const exactTimeFormatter = new Intl.DateTimeFormat("zh-CN", {
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
  hour: "2-digit",
  minute: "2-digit",
  second: "2-digit",
  hour12: false,
});
const suspensionReasonError = computed(() =>
  suspensionReasonTouched.value && !suspensionReason.value.trim()
    ? "请填写暂停原因。"
    : undefined,
);

function submitSearch(value: string) {
  clearTimeout(searchTimer);
  search.value = value;
  q.value = value.trim();
  page.value = 1;
}

function changeControl(id: string, value: CollectionControlValue) {
  if (typeof value !== "string") return;
  if (id === "status") status.value = value === ALL ? "" : value;
  if (id === "role") role.value = value === ALL ? "" : value;
}

function clearFilters() {
  status.value = "";
  role.value = "";
}

function changePageSize(value: number) {
  size.value = value;
  page.value = 1;
}

function openMember(member: NavigationMember) {
  selected.value = member;
  detailOpen.value = true;
}

function displayName(member: NavigationMember) {
  return member.displayName || member.handle || "导航成员";
}

function initials(member: NavigationMember) {
  return Array.from(displayName(member).trim()).slice(0, 1).join("").toUpperCase();
}

function memberLabel(member: NavigationMember) {
  return displayName(member);
}

function memberKey(member: NavigationMember) {
  return member.userKey;
}

function avatarURL(member: NavigationMember) {
  if (!member.avatarMediaKey || avatarFailures.value.has(member.userKey)) return "";
  const account = String(runtime.public.accountUrl || "").replace(/\/$/, "");
  return `${account}/media/${encodeURIComponent(member.avatarMediaKey)}?format=webp&name=thumbnail`;
}

function markAvatarFailed(userKey: string) {
  avatarFailures.value = new Set([...avatarFailures.value, userKey]);
}

function identityURL(member: NavigationMember) {
  const account = String(runtime.public.accountUrl || "").replace(/\/$/, "");
  return member.handle ? `${account}/@${encodeURIComponent(member.handle)}` : account;
}

function roleSourceLabel(source: string) {
  return (
    {
      bootstrap: "初始管理员",
      direct: "直接授予",
      application: "申请获批",
      invitation: "邀请加入",
      automatic: "自动规则",
      group: "成员组",
    }[source] || source
  );
}

function exactTime(value?: string) {
  if (!value) return "—";
  const instant = new Date(value);
  return Number.isNaN(instant.getTime()) ? "—" : exactTimeFormatter.format(instant);
}

function openSuspend() {
  suspensionReason.value = "";
  suspensionReasonTouched.value = false;
  suspendOpen.value = true;
}

async function setMemberStatus(nextStatus: "active" | "suspended") {
  const member = selected.value;
  if (!member || saving.value) return;
  if (nextStatus === "suspended" && !suspensionReason.value.trim()) {
    suspensionReasonTouched.value = true;
    return;
  }
  saving.value = true;
  try {
    const result = await call<{ member: NavigationMember }>(
      `/api/v1/admin/nav/members/${encodeURIComponent(member.userKey)}/status`,
      {
        method: "PATCH",
        body: {
          status: nextStatus,
          reason: nextStatus === "suspended" ? suspensionReason.value.trim() : "",
        },
      },
    );
    selected.value = result.member;
    suspendOpen.value = false;
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    toast.add({
      title: "成员状态更新失败",
      description: apiError.data?.message || "请稍后重试。",
      color: "error",
      icon: "i-tabler-alert-circle",
    });
  } finally {
    saving.value = false;
  }
}
</script>

<template>
  <!--
  THESIS: 成员是身份、本站关系、权限三条可追溯链，而不是一张账号副本表。
  OWN-WORLD: 继承 Nav 管理台的中性 CollectionPanel、主色状态提示与紧凑行式信息密度。
  STORY: 管理员先识别人，再确认成员状态，最后查看权限与活动；高影响操作进入详情层。
  FIRST VIEWPORT: 固定搜索筛选、三段关系列头、内部滚动成员行与常驻分页器。
  FORM: grounded candidate 3，relationship ledger；seed c73a0b46。
  FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, and DESIGN.md
  -->
  <ManagePage
    id="members"
    title="成员"
    description="查看谁已加入本站、当前成员状态及其权限与活动；账号资料仍由用户中心管理。"
    icon="i-tabler-users"
    main-id="manage-main"
    body-class="flex min-h-0 w-full flex-col gap-5"
  >
    <ManageClientBoundary :rows="8">
      <div class="flex min-h-0 flex-1 flex-col">
        <CollectionPanel
          data-member-list-panel
          class="flex min-h-0 flex-1 flex-col [&>[aria-live=polite]]:min-h-0 [&>[aria-live=polite]]:flex-1 [&>[aria-live=polite]]:overflow-y-auto [&>[aria-live=polite]]:overscroll-contain [&>footer]:shrink-0"
          v-model:search="search"
          :items="members"
          :item-key="memberKey"
          :item-label="memberLabel"
          :controls="controls"
          :messages="messages"
          :state="panelState"
          error-message="请检查 Nav API、Membership 数据库和 Identity 公开资料服务。"
          :total="total"
          :page="page"
          :page-size="size"
          :page-sizes="pageSizes"
          :active-filter-count="activeFilterCount"
          label="导航成员列表"
          @search="submitSearch"
          @control-change="changeControl"
          @clear-filters="clearFilters"
          @retry="refresh"
          @page-change="page = $event"
          @page-size-change="changePageSize"
        >
          <template #view>
            <UButton
              icon="i-tabler-refresh"
              color="neutral"
              variant="outline"
              size="xs"
              aria-label="刷新成员列表"
              @click="refresh()"
            />
          </template>

          <template #columns>
            <div
              class="hidden grid-cols-[minmax(15rem,1.2fr)_minmax(12rem,0.8fr)_minmax(16rem,1fr)_auto] gap-4 md:grid"
            >
              <span>身份</span>
              <span>本站关系</span>
              <span>权限与活动</span>
              <span class="w-10 text-right">操作</span>
            </div>
          </template>

          <template #item="{ item: member }">
            <div
              class="grid min-w-0 gap-3 md:grid-cols-[minmax(15rem,1.2fr)_minmax(12rem,0.8fr)_minmax(16rem,1fr)_auto] md:items-center md:gap-4"
            >
              <button
                type="button"
                class="flex min-w-0 items-center gap-3 rounded-lg text-left focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
                @click="openMember(member)"
              >
                <span
                  class="relative grid size-10 shrink-0 place-items-center overflow-hidden rounded-full bg-primary/10 text-sm font-semibold text-primary"
                >
                  {{ initials(member) }}
                  <img
                    v-if="avatarURL(member)"
                    :src="avatarURL(member)"
                    alt=""
                    class="absolute inset-0 size-full object-cover"
                    @error="markAvatarFailed(member.userKey)"
                  />
                </span>
                <span class="min-w-0">
                  <span class="block truncate text-sm font-medium text-highlighted">
                    {{ displayName(member) }}
                  </span>
                  <span class="mt-0.5 block truncate text-xs text-muted">
                    {{ member.handle ? `@${member.handle}` : "未设置用户句柄" }}
                  </span>
                  <code class="mt-1 block truncate text-[11px] text-dimmed">{{ member.userKey }}</code>
                </span>
              </button>

              <div class="min-w-0">
                <p class="mb-1.5 text-[11px] font-medium uppercase tracking-wide text-dimmed md:hidden">
                  本站关系
                </p>
                <div class="flex flex-wrap items-center gap-2">
                  <UBadge
                    :label="member.status === 'active' ? '正常' : '已暂停'"
                    :color="member.status === 'active' ? 'success' : 'warning'"
                    :icon="member.status === 'active' ? 'i-tabler-user-check' : 'i-tabler-user-pause'"
                    variant="subtle"
                  />
                  <span class="text-xs text-muted">加入 {{ rel(member.joinedAt) }}</span>
                </div>
                <p class="mt-2 text-xs text-muted">最近进入 {{ rel(member.lastSeenAt) }}</p>
              </div>

              <div class="min-w-0">
                <p class="mb-1.5 text-[11px] font-medium uppercase tracking-wide text-dimmed md:hidden">
                  权限与活动
                </p>
                <div class="flex flex-wrap gap-1.5">
                  <UBadge
                    v-if="!member.roles.length"
                    label="普通成员"
                    color="neutral"
                    variant="soft"
                    size="sm"
                  />
                  <UBadge
                    v-for="memberRole in member.roles"
                    :key="`${memberRole.key}:${memberRole.source}`"
                    :label="memberRole.displayName || memberRole.key"
                    color="primary"
                    variant="soft"
                    size="sm"
                  />
                </div>
                <p class="mt-2 text-xs text-muted">
                  投稿 {{ member.submissionCount }} · 待处理申请 {{ member.pendingApplications }}
                </p>
              </div>

              <UButton
                icon="i-tabler-chevron-right"
                color="neutral"
                variant="ghost"
                size="sm"
                square
                class="justify-self-end"
                :aria-label="`查看成员：${displayName(member)}`"
                @click="openMember(member)"
              />
            </div>
          </template>
        </CollectionPanel>
      </div>
    </ManageClientBoundary>

    <USlideover
      v-model:open="detailOpen"
      :title="selected ? displayName(selected) : '成员详情'"
      description="身份资料、本站成员关系和权限分别由各自系统管理。"
      :ui="{ content: 'w-full sm:max-w-lg' }"
    >
      <template #body>
        <div v-if="selected" class="space-y-6">
          <section aria-labelledby="member-identity-title">
            <div class="flex items-center gap-3">
              <span
                class="relative grid size-12 shrink-0 place-items-center overflow-hidden rounded-full bg-primary/10 text-base font-semibold text-primary"
              >
                {{ initials(selected) }}
                <img
                  v-if="avatarURL(selected)"
                  :src="avatarURL(selected)"
                  alt=""
                  class="absolute inset-0 size-full object-cover"
                  @error="markAvatarFailed(selected.userKey)"
                />
              </span>
              <div class="min-w-0">
                <h2 id="member-identity-title" class="truncate font-semibold text-highlighted">
                  {{ displayName(selected) }}
                </h2>
                <p class="mt-0.5 truncate text-sm text-muted">
                  {{ selected.handle ? `@${selected.handle}` : selected.userKey }}
                </p>
              </div>
            </div>
            <div class="mt-4 rounded-lg border border-default bg-elevated/50 p-3">
              <h3 class="text-xs font-medium text-toned">Identity 用户</h3>
              <code class="mt-1 block break-all text-xs text-muted">{{ selected.userKey }}</code>
              <UButton
                class="mt-3"
                :to="identityURL(selected)"
                external
                target="_blank"
                label="在用户中心查看资料"
                icon="i-tabler-external-link"
                color="neutral"
                variant="outline"
                size="xs"
              />
            </div>
          </section>

          <section class="border-t border-default pt-5" aria-labelledby="member-relation-title">
            <div class="flex items-center justify-between gap-3">
              <h2 id="member-relation-title" class="text-sm font-semibold text-highlighted">Nav 成员关系</h2>
              <UBadge
                :label="selected.status === 'active' ? '正常' : '已暂停'"
                :color="selected.status === 'active' ? 'success' : 'warning'"
                variant="subtle"
              />
            </div>
            <dl class="mt-3 grid grid-cols-1 gap-3 text-sm sm:grid-cols-2">
              <div class="rounded-lg bg-elevated/50 p-3">
                <dt class="text-xs text-muted">加入本站</dt>
                <dd class="mt-1 font-medium text-highlighted">
                  <time :datetime="selected.joinedAt">{{ exactTime(selected.joinedAt) }}</time>
                </dd>
              </div>
              <div class="rounded-lg bg-elevated/50 p-3">
                <dt class="text-xs text-muted">最近进入</dt>
                <dd class="mt-1 font-medium text-highlighted">
                  <time :datetime="selected.lastSeenAt">{{ exactTime(selected.lastSeenAt) }}</time>
                </dd>
              </div>
            </dl>
            <UAlert
              v-if="selected.status === 'suspended'"
              class="mt-3"
              color="warning"
              variant="subtle"
              icon="i-tabler-user-pause"
              title="已认证操作已暂停"
              description="公开内容仍可浏览，本站已认证操作已被阻止。"
            />
            <dl
              v-if="selected.status === 'suspended'"
              class="mt-3 grid grid-cols-1 gap-3 text-sm sm:grid-cols-2"
            >
              <div class="rounded-lg bg-elevated/50 p-3">
                <dt class="text-xs text-muted">暂停时间</dt>
                <dd class="mt-1 font-medium text-highlighted">
                  <time v-if="selected.suspendedAt" :datetime="selected.suspendedAt">
                    {{ exactTime(selected.suspendedAt) }}
                  </time>
                  <span v-else>—</span>
                </dd>
              </div>
              <div class="rounded-lg bg-elevated/50 p-3">
                <dt class="text-xs text-muted">操作人</dt>
                <dd class="mt-1 break-all font-medium text-highlighted">
                  <code>{{ selected.suspendedBy || "—" }}</code>
                </dd>
              </div>
              <div class="rounded-lg bg-elevated/50 p-3 sm:col-span-2">
                <dt class="text-xs text-muted">暂停原因</dt>
                <dd class="mt-1 whitespace-pre-wrap text-highlighted">
                  {{ selected.suspensionReason || "—" }}
                </dd>
              </div>
            </dl>
          </section>

          <section class="border-t border-default pt-5" aria-labelledby="member-access-title">
            <h2 id="member-access-title" class="text-sm font-semibold text-highlighted">本站权限与活动</h2>
            <p class="mt-1 text-xs text-muted">成员存在不代表拥有维护权限。</p>
            <div class="mt-3 flex flex-wrap gap-2">
              <UTooltip
                v-for="memberRole in selected.roles"
                :key="`${memberRole.key}:${memberRole.source}`"
                :text="roleSourceLabel(memberRole.source)"
              >
                <UBadge
                  :label="memberRole.displayName || memberRole.key"
                  color="primary"
                  variant="soft"
                />
              </UTooltip>
              <UBadge
                v-if="!selected.roles.length"
                label="无维护角色"
                color="neutral"
                variant="soft"
              />
            </div>
            <dl class="mt-3 grid grid-cols-2 gap-3 text-sm">
              <div class="rounded-lg bg-elevated/50 p-3">
                <dt class="text-xs text-muted">提交站点</dt>
                <dd class="mt-1 text-lg font-semibold tabular-nums text-highlighted">{{ selected.submissionCount }}</dd>
              </div>
              <div class="rounded-lg bg-elevated/50 p-3">
                <dt class="text-xs text-muted">待处理申请</dt>
                <dd class="mt-1 text-lg font-semibold tabular-nums text-highlighted">{{ selected.pendingApplications }}</dd>
              </div>
            </dl>
          </section>

          <section class="border-t border-default pt-5" aria-labelledby="member-actions-title">
            <h2 id="member-actions-title" class="text-sm font-semibold text-highlighted">成员操作</h2>
            <p class="mt-1 text-xs leading-5 text-muted">
              暂停只影响本站已认证操作，不修改用户中心账号，也不会删除授权和审计历史。
            </p>
            <UButton
              v-if="selected.status === 'active'"
              class="mt-3"
              label="暂停成员资格"
              icon="i-tabler-user-pause"
              color="warning"
              variant="soft"
              :disabled="selected.userKey === me?.userKey"
              @click="openSuspend"
            />
            <UButton
              v-else
              class="mt-3"
              label="恢复成员资格"
              icon="i-tabler-user-check"
              :loading="saving"
              @click="setMemberStatus('active')"
            />
            <p
              v-if="selected.status === 'active' && selected.userKey === me?.userKey"
              class="mt-2 text-xs text-muted"
            >
              不能暂停当前管理员自己的成员资格。
            </p>
          </section>
        </div>
      </template>
    </USlideover>

    <UModal
      v-model:open="suspendOpen"
      title="暂停成员资格"
      description="该成员仍可浏览公开导航，但所有本站已认证操作会立即被阻止。"
    >
      <template #body>
        <div class="space-y-4">
          <UAlert
            color="warning"
            variant="subtle"
            icon="i-tabler-alert-triangle"
            :title="selected ? `即将暂停 ${displayName(selected)}` : '即将暂停成员'"
            description="已有角色、申请和活动历史会保留，恢复后可继续使用。"
          />
          <UFormField
            label="暂停原因"
            hint="必填；会显示在成员详情中"
            required
            :error="suspensionReasonError"
          >
            <UTextarea
              v-model="suspensionReason"
              :rows="4"
              maxlength="500"
              class="w-full"
              placeholder="例如：账号归属待确认"
              @blur="suspensionReasonTouched = true"
            />
          </UFormField>
        </div>
      </template>
      <template #footer>
        <div class="flex w-full justify-end gap-2">
          <UButton
            label="取消"
            color="neutral"
            variant="outline"
            @click="() => { suspendOpen = false }"
          />
          <UButton
            label="确认暂停"
            icon="i-tabler-user-pause"
            color="warning"
            :loading="saving"
            :disabled="!suspensionReason.trim()"
            @click="setMemberStatus('suspended')"
          />
        </div>
      </template>
    </UModal>
  </ManagePage>
</template>
