<script setup lang="ts">
import {
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
const section = ref<"site" | "search" | "footer">("site");
const sections = [
  {
    key: "site",
    label: "站点基础",
    icon: "i-tabler-adjustments-horizontal",
    description: "站点名称、首页标题与说明",
  },
  {
    key: "search",
    label: "搜索",
    icon: "i-tabler-search",
    description: "首页搜索框的引导文案",
  },
  {
    key: "footer",
    label: "页脚",
    icon: "i-tabler-layout-bottombar",
    description: "公开页面底部的品牌说明",
  },
] as const;
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
const settingsForm = ref<{ submit: () => Promise<void> }>();
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
    section.value = value === "search" || value === "footer" ? value : "site";
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
    ref="settingsForm"
    :schema="schema"
    :state="form"
    class="contents"
    @submit="save"
  >
    <ManageSettingsLayout
      v-model:active-section="section"
      :title="
        sections.find((item) => item.key === section)?.label || '站点设置'
      "
      :description="sections.find((item) => item.key === section)?.description"
      :sections="sections"
    >
      <template #notice
        ><UAlert
          v-if="!isAdmin"
          color="neutral"
          variant="subtle"
          icon="i-tabler-lock"
          title="只读设置"
          description="只有 Nav 管理员可以修改公开站点设置。"
      /></template>

      <ManageSettingCard
        v-if="pending"
        title="正在加载设置"
        description="读取当前站点的已保存配置。"
        ><div class="grid gap-4">
          <USkeleton class="h-9 w-full" /><USkeleton
            class="h-9 w-full"
          /><USkeleton class="h-24 w-full" /></div
      ></ManageSettingCard>
      <UAlert
        v-else-if="error"
        color="error"
        variant="subtle"
        icon="i-tabler-alert-circle"
        title="设置加载失败"
        description="站点配置尚未初始化或服务不可用。"
      />

      <ManageSettingCard
        v-else-if="section === 'site'"
        title="站点基础"
        description="这些内容用于公开页标题、导航品牌和搜索引擎描述。"
      >
        <div class="grid gap-4">
          <UFormField name="name" label="站点名称" required
            ><UInput v-model="form.name" :disabled="!isAdmin" class="w-full"
          /></UFormField>
          <UFormField name="title" label="首页标题" required
            ><UInput v-model="form.title" :disabled="!isAdmin" class="w-full"
          /></UFormField>
          <UFormField name="description" label="站点介绍" required
            ><UTextarea
              v-model="form.description"
              :disabled="!isAdmin"
              :rows="4"
              class="w-full"
          /></UFormField>
        </div>
      </ManageSettingCard>

      <ManageSettingCard
        v-else-if="section === 'search'"
        title="搜索体验"
        description="提示用户可检索名称、域名、描述、分类与标签。"
      >
        <UFormField name="searchPlaceholder" label="搜索框占位文案" required
          ><UInput
            v-model="form.searchPlaceholder"
            :disabled="!isAdmin"
            class="w-full"
        /></UFormField>
      </ManageSettingCard>

      <ManageSettingCard
        v-else
        title="页脚内容"
        description="保持简短，说明站点的整理目的。"
      >
        <UFormField name="footerTagline" label="页脚标语" required
          ><UTextarea
            v-model="form.footerTagline"
            :disabled="!isAdmin"
            :rows="3"
            class="w-full"
        /></UFormField>
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
  </UForm>
</template>
