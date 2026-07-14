<script setup lang="ts">
import {
  ManageClientBoundary,
  ManageSaveDock,
  ManageSettingCard,
  ManageSettingsLayout,
} from "@platform/manage/components";
import { useActionFeedback } from "@platform/manage/use-action-feedback";
import { useManageSettings } from "@platform/manage/use-manage-settings";
import { z } from "zod";
import type {
  NavigationSettingsResponse,
  NavigationSiteCopy,
} from "~/types/navigation";

definePageMeta({ layout: "manage", middleware: "auth" });
useSeoMeta({ title: "站点设置 · 月离导航" });

const { isAdmin } = useAuth();
const { call } = useApi();
const route = useRoute();
const router = useRouter();
type SectionKey = "site" | "featured" | "search" | "footer";

const section = ref<SectionKey>("site");
const sections = [
  {
    key: "site",
    label: "品牌",
    icon: "i-tabler-compass",
    description: "导航名称、品牌短句与分享摘要",
  },
  {
    key: "featured",
    label: "首页精选",
    icon: "i-tabler-sparkles",
    description: "管理“本周值得逛”的内容来源与排序规则",
  },
  {
    key: "search",
    label: "搜索",
    icon: "i-tabler-search",
    description: "全局搜索面板的引导文案与检索范围",
  },
  {
    key: "footer",
    label: "页脚",
    icon: "i-tabler-layout-bottombar",
    description: "公开页面底部的品牌说明",
  },
] as const;
const activeSection = computed(
  () => sections.find((item) => item.key === section.value) ?? sections[0],
);
const form = reactive<NavigationSiteCopy>({
  name: "",
  title: "",
  description: "",
  searchPlaceholder: "",
  footerTagline: "",
});
const schema = z.object({
  name: z.string().trim().min(1, "请输入站点名称").max(120),
  title: z.string().trim().min(1, "请输入首页标题").max(200),
  description: z.string().trim().min(1, "请输入站点介绍").max(500),
  searchPlaceholder: z.string().trim().min(1, "请输入搜索提示文案").max(160),
  footerTagline: z.string().trim().min(1, "请输入页脚标语").max(300),
});
const settingsForm = useTemplateRef<{ submit: () => Promise<void> }>(
  "settings-form",
);
const saveError = ref("");
const {
  status: saveStatus,
  pending: markSaving,
  success: markSaved,
  reset: resetSave,
} = useActionFeedback();
const settingsState = useManageSettings({
  snapshot: () => form,
  restore: (snapshot) => Object.assign(form, snapshot),
});

const { data, pending, error, refresh } = await useAsyncData(
  "nav-settings",
  () => call<NavigationSettingsResponse>("/api/v1/admin/nav/settings"),
  { server: false },
);
watch(
  data,
  (value) => {
    if (!value?.settings) return;
    Object.assign(form, value.settings);
    nextTick(settingsState.capture);
  },
  { immediate: true },
);
watch(
  () => route.query.section,
  (value) => {
    section.value = sections.some((item) => item.key === value)
      ? (value as SectionKey)
      : "site";
  },
  { immediate: true },
);
watch(section, (value) => {
  if (route.query.section !== value)
    void router.replace({ query: { ...route.query, section: value } });
});

async function save() {
  if (!isAdmin.value) return;
  markSaving();
  saveError.value = "";
  try {
    await call<NavigationSettingsResponse>("/api/v1/admin/nav/settings", {
      method: "PATCH",
      body: { ...form },
    });
    await refresh();
    settingsState.capture();
    markSaved();
  } catch (failure) {
    resetSave();
    const apiError = failure as { data?: { message?: string } };
    saveError.value = apiError.data?.message || "保存失败，请稍后重试。";
  }
}
function submitSettings() {
  void settingsForm.value?.submit();
}
function discard() {
  settingsState.discard();
  saveError.value = "";
  resetSave();
}
</script>

<template>
  <UForm
    ref="settings-form"
    :schema="schema"
    :state="form"
    class="contents"
    @submit="save"
  >
    <ManageClientBoundary :rows="4">
      <ManageSettingsLayout
        v-model:active-section="section"
        :title="activeSection.label"
        :description="activeSection.description"
        :sections="sections"
      >
        <template #actions>
          <UButton
            to="/"
            target="_blank"
            color="neutral"
            variant="soft"
            icon="i-tabler-external-link"
            label="查看公开站点"
          />
        </template>

        <template #notice>
          <UAlert
            v-if="!isAdmin"
            color="neutral"
            variant="subtle"
            icon="i-tabler-lock"
            title="只读设置"
            description="只有 Nav 管理员可以修改公开站点设置。"
          />
        </template>

        <ManageSettingCard
          v-if="pending"
          title="正在加载设置"
          description="读取当前站点的已保存配置。"
        >
          <div class="grid gap-4" role="status" aria-label="正在加载设置">
            <USkeleton class="h-9 w-full" />
            <USkeleton class="h-9 w-full" />
            <USkeleton class="h-24 w-full" />
          </div>
        </ManageSettingCard>

        <UAlert
          v-else-if="error"
          color="error"
          variant="subtle"
          icon="i-tabler-alert-circle"
          title="设置加载失败"
          description="站点配置尚未初始化或服务不可用。"
        >
          <template #actions>
            <UButton
              color="error"
              variant="soft"
              icon="i-tabler-refresh"
              label="重新加载"
              @click="() => refresh()"
            />
          </template>
        </UAlert>

        <ManageSettingCard
          v-else-if="section === 'site'"
          title="品牌与分享信息"
          description="这些内容用于公开页页头、浏览器标题和社交分享摘要。"
        >
          <div class="grid gap-5">
            <UFormField
              name="name"
              label="导航名称"
              description="显示在公开页左上角，并作为页面标题的品牌名称。"
              required
            >
              <UInput v-model="form.name" :disabled="!isAdmin" class="w-full" />
            </UFormField>
            <UFormField
              name="title"
              label="品牌短句"
              description="显示在导航名称下方，同时用于社交分享标题。"
              required
            >
              <UInput
                v-model="form.title"
                :disabled="!isAdmin"
                class="w-full"
              />
            </UFormField>
            <UFormField
              name="description"
              label="站点描述"
              description="用于搜索引擎和社交平台理解本站内容。"
              required
            >
              <UTextarea
                v-model="form.description"
                :disabled="!isAdmin"
                :rows="4"
                class="w-full"
              />
            </UFormField>
          </div>
        </ManageSettingCard>

        <ManageSettingCard
          v-else-if="section === 'featured'"
          title="本周值得逛"
          description="首页精选是内容治理能力，不在站点文案中维护。"
        >
          <div class="space-y-4">
            <div
              class="overflow-hidden rounded-xl border border-primary/20 bg-primary/5"
            >
              <div class="flex items-start gap-3 p-4 sm:p-5">
                <span
                  class="grid size-10 shrink-0 place-items-center rounded-lg bg-default text-primary shadow-sm ring-1 ring-primary/15"
                >
                  <UIcon
                    name="i-tabler-sparkles"
                    class="size-5"
                    aria-hidden="true"
                  />
                </span>
                <div class="min-w-0">
                  <h3 class="font-display font-semibold text-highlighted">
                    编辑精选优先，访问数据补位
                  </h3>
                  <p class="mt-1 text-sm leading-6 text-muted">
                    标记“加入首页精选”的站点会按目录顺序进入榜单；不足五项时，再用真实访问次数较高的站点补齐。
                  </p>
                </div>
              </div>
            </div>
            <div
              class="flex flex-col gap-3 rounded-xl border border-default bg-elevated/45 p-4 sm:flex-row sm:items-center sm:justify-between"
            >
              <div class="min-w-0">
                <p class="text-sm font-medium text-highlighted">
                  在站点条目编辑器中管理精选状态
                </p>
                <p class="mt-1 text-sm text-muted">
                  打开任一站点，切换“加入首页精选”即可更新首页榜单。
                </p>
              </div>
              <UButton
                to="/manage"
                color="neutral"
                variant="soft"
                trailing-icon="i-tabler-arrow-right"
                label="管理精选条目"
                class="shrink-0 self-start sm:self-auto"
              />
            </div>
          </div>
        </ManageSettingCard>

        <ManageSettingCard
          v-else-if="section === 'search'"
          title="全局搜索"
          description="搜索按钮、快捷键和首页精选区共用同一个搜索面板。"
        >
          <div class="grid gap-5">
            <UFormField
              name="searchPlaceholder"
              label="搜索提示文案"
              description="建议说明可搜索的内容，不要重复写“点击搜索”。"
              required
            >
              <UInput
                v-model="form.searchPlaceholder"
                :disabled="!isAdmin"
                icon="i-tabler-search"
                class="w-full"
              />
            </UFormField>
            <div
              class="grid gap-3 rounded-xl border border-default bg-elevated/45 p-4 text-sm sm:grid-cols-2"
            >
              <div>
                <p class="font-medium text-highlighted">检索范围</p>
                <p class="mt-1 leading-6 text-muted">
                  名称、描述、域名、分类、主题与标签
                </p>
              </div>
              <div>
                <p class="font-medium text-highlighted">快捷键</p>
                <p class="mt-1 leading-6 text-muted">/、Ctrl K 或 Command K</p>
              </div>
            </div>
          </div>
        </ManageSettingCard>

        <ManageSettingCard
          v-else
          title="页脚说明"
          description="保持简短，说明这个导航站为什么存在。"
        >
          <UFormField
            name="footerTagline"
            label="页脚标语"
            description="显示在所有公开页面底部，建议控制在一句话以内。"
            required
          >
            <UTextarea
              v-model="form.footerTagline"
              :disabled="!isAdmin"
              :rows="3"
              class="w-full"
            />
          </UFormField>
        </ManageSettingCard>

        <ManageSaveDock
          :dirty="settingsState.dirty.value"
          :status="saveStatus"
          :error="saveError"
          :disabled="!isAdmin"
          @discard="discard"
          @save="submitSettings"
        />
      </ManageSettingsLayout>
    </ManageClientBoundary>
  </UForm>
</template>
