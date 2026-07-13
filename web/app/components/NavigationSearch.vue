<script setup lang="ts">
const { placeholder, resultCount } = defineProps<{
  placeholder: string;
  resultCount?: number;
}>();

const query = defineModel<string>({ required: true });
const input = useTemplateRef<{ inputRef?: HTMLInputElement }>("search-input");

function focusSearch() {
  input.value?.inputRef?.focus();
}

defineShortcuts({
  "/": {
    usingInput: false,
    handler: focusSearch,
  },
  meta_k: focusSearch,
  ctrl_k: focusSearch,
});
</script>

<template>
  <div class="relative">
    <UInput
      ref="search-input"
      v-model="query"
      :placeholder="placeholder"
      icon="i-tabler-search"
      size="xl"
      autocomplete="off"
      aria-label="搜索导航站点"
      class="w-full"
      :ui="{
        base: 'h-14 rounded-xl bg-default/94 pl-12 pr-24 text-base shadow-sm ring-1 ring-default focus-visible:ring-2 focus-visible:ring-primary sm:h-16 sm:text-lg',
        leadingIcon: 'ml-1 size-5 text-primary sm:size-6',
      }"
    >
      <template #trailing>
        <UButton
          v-if="query"
          color="neutral"
          variant="ghost"
          size="sm"
          icon="i-tabler-x"
          aria-label="清空搜索"
          @click="
            () => {
              query = '';
            }
          "
        />
        <UKbd v-else value="/" class="hidden sm:inline-flex" />
      </template>
    </UInput>
    <p
      v-if="query && resultCount !== undefined"
      class="mt-2 text-sm text-muted"
      aria-live="polite"
    >
      找到 {{ resultCount }} 个匹配站点
    </p>
  </div>
</template>
