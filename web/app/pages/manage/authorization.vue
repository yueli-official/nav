<script setup lang="ts">
interface RoleView {
  key: string;
  displayName: string;
  kind: string;
  status: string;
  protected: boolean;
  capabilities: string[];
  assignmentSources: string[];
}
interface ApplicationView {
  id: string;
  subject: string;
  role: string;
  reason: string;
  createdAt: string;
}
interface ConsoleView {
  activeRevision: number;
  policy: { number: number; state: string };
  roles: RoleView[];
  automaticRules: { key: string; enabled: boolean }[];
  applications: ApplicationView[];
  capabilities: { key: string; displayName: string }[];
}

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "权限策略 · 月离导航" });

const { call } = useApi();
const { isAdministrator } = useMe();
const toast = useToast();
const busy = ref(false);
const createRoleOpen = ref(false);
const roleForm = reactive({
  key: "",
  displayName: "",
  capabilities: [] as string[],
  assignmentSources: ["application", "invitation", "direct"] as string[],
});
const { data, pending, error, refresh } = await useAsyncData(
  "nav-authorization-console",
  () => call<ConsoleView>("/api/v1/authorization/manage/console"),
  { server: false },
);
const consoleState = computed(() => data.value);
const draft = computed(() => consoleState.value?.policy.state === "draft");

async function mutate(task: () => Promise<unknown>, success: string) {
  if (busy.value) return;
  busy.value = true;
  try {
    const result = await task();
    if (result === false) return;
    toast.add({ title: success, color: "success", icon: "i-tabler-check" });
    await refresh();
  } catch (failure) {
    const apiError = failure as { data?: { message?: string } };
    toast.add({
      title: "操作失败",
      description: apiError.data?.message || "请刷新后重试。",
      color: "error",
      icon: "i-tabler-alert-circle",
    });
  } finally {
    busy.value = false;
  }
}

function createDraft() {
  const active = consoleState.value?.activeRevision;
  if (!active) return;
  return mutate(
    () =>
      call("/api/v1/authorization/manage/policies/drafts", {
        method: "POST",
        body: { expectedActiveRevision: active },
      }),
    "策略草稿已创建",
  );
}

function toggleRoleCapability(role: RoleView, capability: string) {
  if (!draft.value || role.protected) return;
  const capabilities = role.capabilities.includes(capability)
    ? role.capabilities.filter((item) => item !== capability)
    : [...role.capabilities, capability];
  return mutate(
    () =>
      call(
        `/api/v1/authorization/manage/policies/${consoleState.value?.policy.number}/roles/${role.key}/capabilities`,
        { method: "PUT", body: { capabilities } },
      ),
    "角色能力已更新到草稿",
  );
}

function toggleAutomatic(enabled: boolean) {
  const rule = consoleState.value?.automaticRules[0];
  if (!draft.value || !rule) return;
  return mutate(
    () =>
      call(
        `/api/v1/authorization/manage/policies/${consoleState.value?.policy.number}/automatic/${rule.key}`,
        { method: "PUT", body: { enabled } },
      ),
    enabled ? "已在草稿中启用新成员自动授权" : "已在草稿中关闭新成员自动授权",
  );
}

async function validateAndActivate() {
  const state = consoleState.value;
  if (!draft.value || !state) return;
  await mutate(async () => {
    const validation = await call<{ valid: boolean; violations: string[] }>(
      `/api/v1/authorization/manage/policies/${state.policy.number}/validate`,
      { method: "POST" },
    );
    if (!validation.valid) throw new Error(validation.violations.join("；"));
    const impact = await call<{ addedBindings: number; removedBindings: number }>(
      `/api/v1/authorization/manage/policies/${state.policy.number}/preview`,
      { method: "POST" },
    );
    if (impact.removedBindings > 0 && !window.confirm(
      `本次发布会移除 ${impact.removedBindings} 项能力绑定，是否继续？`,
    )) return false;
    await call(
      `/api/v1/authorization/manage/policies/${state.policy.number}/activate`,
      {
        method: "POST",
        body: { expectedActiveRevision: state.activeRevision },
      },
    );
  }, "权限策略已发布");
}

function review(application: ApplicationView, decision: "approve" | "reject") {
  return mutate(
    () =>
      call(
        `/api/v1/authorization/manage/applications/${application.id}/review`,
        {
          method: "POST",
          body: {
            decision,
            reason: decision === "approve" ? "管理员批准" : "管理员拒绝",
          },
        },
      ),
    decision === "approve" ? "申请已批准" : "申请已拒绝",
  );
}

function createRole() {
  const revision = consoleState.value?.policy.number;
  if (!draft.value || !revision) return;
  return mutate(async () => {
    await call(
      `/api/v1/authorization/manage/policies/${revision}/roles`,
      { method: "POST", body: roleForm },
    );
    createRoleOpen.value = false;
    Object.assign(roleForm, {
      key: "",
      displayName: "",
      capabilities: [],
      assignmentSources: ["application", "invitation", "direct"],
    });
  }, "自定义角色已创建");
}
</script>

<template>
  <ManagePage
    id="authorization"
    title="权限策略"
    description="配置本站角色、能力、申请流程与新成员自动授权；成员目录在独立页面管理。"
    icon="i-tabler-shield-lock"
    main-id="manage-main"
    body-class="w-full space-y-5"
  >
    <template #actions>
      <UButton
        v-if="consoleState && !draft"
        label="创建策略草稿"
        icon="i-tabler-file-plus"
        :loading="busy"
        @click="createDraft"
      />
      <UButton
        v-else-if="draft"
        label="验证并发布"
        icon="i-tabler-rocket"
        :loading="busy"
        @click="validateAndActivate"
      />
    </template>

    <ManageClientBoundary :rows="8">
      <UAlert
        v-if="!isAdministrator"
        color="error"
        icon="i-tabler-lock"
        title="只有管理员可以管理本站权限"
        description="本站没有超级管理员，也不会继承 Identity 或其他站点的管理员身份。"
      />
      <SkeletonList v-else-if="pending" :rows="8" />
      <UAlert
        v-else-if="error || !consoleState"
        color="error"
        icon="i-tabler-alert-circle"
        title="权限配置加载失败"
        description="Authorization 实例暂时不可用，请检查 Nav API 与数据库状态。"
      />

      <template v-else>
      <section class="rounded-xl border border-default bg-default p-4 sm:p-5">
        <div class="flex flex-wrap items-start justify-between gap-3">
          <div>
            <p class="text-sm font-semibold text-highlighted">策略发布</p>
            <p class="mt-1 text-sm text-muted">
              当前生效修订 {{ consoleState.activeRevision }}；正在查看修订
              {{ consoleState.policy.number }}。
            </p>
          </div>
          <UBadge
            :label="draft ? '草稿' : '已生效'"
            :color="draft ? 'warning' : 'success'"
            variant="subtle"
          />
        </div>
        <p class="mt-3 text-xs leading-5 text-muted">
          所有角色和自动规则修改先进入草稿，验证影响后一次发布，不会逐项改变线上权限。
        </p>
      </section>

      <section class="space-y-3" aria-labelledby="roles-title">
        <div class="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h2 id="roles-title" class="text-sm font-semibold text-highlighted">
              角色与能力
            </h2>
            <p class="mt-1 text-xs text-muted">
              管理员角色受保护；内容维护者和自定义角色只使用明确勾选的能力。
            </p>
          </div>
          <UButton
            v-if="draft"
            label="新建自定义角色"
            icon="i-tabler-user-plus"
            color="neutral"
            variant="soft"
            @click="() => { createRoleOpen = true }"
          />
        </div>
        <div class="grid gap-3 lg:grid-cols-2">
          <article
            v-for="role in consoleState.roles"
            :key="role.key"
            class="rounded-xl border border-default bg-default p-4"
          >
            <div class="flex items-center justify-between gap-3">
              <div>
                <p class="font-medium text-highlighted">{{ role.displayName }}</p>
                <p class="mt-0.5 text-xs text-muted">{{ role.key }}</p>
              </div>
              <UBadge
                :label="role.protected ? '受保护' : role.kind === 'custom' ? '自定义' : '内置'"
                color="neutral"
                variant="soft"
              />
            </div>
            <div class="mt-4 space-y-2">
              <UCheckbox
                v-for="capability in consoleState.capabilities"
                :key="capability.key"
                :model-value="role.capabilities.includes(capability.key)"
                :label="capability.displayName"
                :disabled="role.protected || !draft || busy"
                @update:model-value="toggleRoleCapability(role, capability.key)"
              />
            </div>
          </article>
        </div>
      </section>

      <section class="rounded-xl border border-default bg-default p-4 sm:p-5">
        <div class="flex items-start justify-between gap-4">
          <div>
            <h2 class="text-sm font-semibold text-highlighted">
              新导航成员自动获得内容维护者权限
            </h2>
            <p class="mt-1 text-xs leading-5 text-muted">
              默认关闭。启用后，用户首次以已登录身份加入导航时授予本站维护权限；不修改用户中心身份或全局角色。
            </p>
          </div>
          <USwitch
            :model-value="consoleState.automaticRules[0]?.enabled ?? false"
            :disabled="!draft || busy"
            @update:model-value="toggleAutomatic(Boolean($event))"
          />
        </div>
      </section>

      <section class="space-y-3" aria-labelledby="applications-title">
        <div>
          <h2 id="applications-title" class="text-sm font-semibold text-highlighted">
            待处理申请
          </h2>
          <p class="mt-1 text-xs text-muted">批准后授权立即生效，不需要修改 Identity。</p>
        </div>
        <ManageEmpty
          v-if="!consoleState.applications.length"
          icon="i-tabler-user-check"
          text="当前没有待处理申请"
        />
        <div v-else class="divide-y divide-default overflow-hidden rounded-xl border border-default bg-default">
          <article
            v-for="application in consoleState.applications"
            :key="application.id"
            class="grid gap-3 p-4 sm:grid-cols-[minmax(0,1fr)_auto] sm:items-center"
          >
            <div class="min-w-0">
              <p class="truncate text-sm font-medium text-highlighted">
                {{ application.subject }}
              </p>
              <p class="mt-1 text-xs text-muted">
                申请 {{ application.role }} · {{ application.reason || "未填写原因" }}
              </p>
            </div>
            <div class="flex gap-2">
              <UButton
                label="拒绝"
                color="neutral"
                variant="outline"
                :disabled="busy"
                @click="review(application, 'reject')"
              />
              <UButton
                label="批准"
                :disabled="busy"
                @click="review(application, 'approve')"
              />
            </div>
          </article>
        </div>
      </section>
      </template>
    </ManageClientBoundary>

    <UModal
      v-model:open="createRoleOpen"
      title="新建自定义角色"
      description="角色能力写入当前策略草稿，发布后才生效。"
    >
      <template #body>
        <div class="space-y-4">
          <UFormField label="角色标识" required>
            <UInput v-model="roleForm.key" placeholder="editor" class="w-full" />
          </UFormField>
          <UFormField label="显示名称" required>
            <UInput v-model="roleForm.displayName" placeholder="编辑" class="w-full" />
          </UFormField>
          <UFormField label="能力">
            <div class="space-y-2">
              <UCheckbox
                v-for="capability in consoleState?.capabilities ?? []"
                :key="capability.key"
                :model-value="roleForm.capabilities.includes(capability.key)"
                :label="capability.displayName"
                @update:model-value="
                  roleForm.capabilities = $event
                    ? [...roleForm.capabilities, capability.key]
                    : roleForm.capabilities.filter((item) => item !== capability.key)
                "
              />
            </div>
          </UFormField>
        </div>
      </template>
      <template #footer>
        <div class="flex w-full justify-end gap-2">
          <UButton
            label="取消"
            color="neutral"
            variant="outline"
            @click="() => { createRoleOpen = false }"
          />
          <UButton
            label="创建角色"
            :disabled="!roleForm.key.trim() || !roleForm.displayName.trim()"
            :loading="busy"
            @click="createRole"
          />
        </div>
      </template>
    </UModal>
  </ManagePage>
</template>
