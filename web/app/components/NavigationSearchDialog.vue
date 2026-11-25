<script setup lang="ts">
import { portalRecommendations, searchPortalEntries } from "~/data/portal";

defineProps<{ placeholder?: string }>();

const { open, closeSearch } = useNavigationSearch();
const query = ref("");
const searchInput = useTemplateRef<{ focus: () => void }>("search-input");
const searching = computed(() => Boolean(query.value.trim()));
const results = computed(() => searchPortalEntries(query.value));
const displayedEntries = computed(() => searching.value ? results.value : portalRecommendations);

watch(open, async (isOpen) => {
  if (!isOpen) return;
  query.value = "";
  await nextTick();
  searchInput.value?.focus();
});

async function openFirstResult() {
  const first = displayedEntries.value[0];
  if (!first) return;
  closeSearch();
  await navigateTo(first.url, { external: true });
}
</script>

<template>
  <UModal
    v-model:open="open"
    title="搜索站群入口"
    :ui="{ content: 'overflow-hidden sm:max-w-lg', body: 'p-0 sm:p-0' }"
  >
    <template #body>
      <div class="border-b border-default p-4">
        <NavigationSearch
          ref="search-input"
          v-model="query"
          :placeholder="placeholder || '搜索工具、文档或站点名称'"
          :result-count="searching ? results.length : undefined"
          @submit="openFirstResult"
        />
      </div>

      <div class="max-h-[min(60dvh,32rem)] overflow-y-auto p-2.5">
        <ul v-if="displayedEntries.length" class="space-y-0.5">
          <li v-for="entry in displayedEntries" :key="entry.url">
            <PortalEntryLink :entry="entry" @click="closeSearch" />
          </li>
        </ul>
        <p v-else class="px-4 py-6 text-sm text-muted" role="status">
          没有匹配的入口，试试软件名或更短的关键词。
        </p>
      </div>

      <p class="border-t border-default px-4 py-2.5 text-xs text-muted">Enter 打开第一项</p>
    </template>
  </UModal>
</template>
