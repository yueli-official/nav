<script setup lang="ts">
const { placeholder, resultCount } = defineProps<{
  placeholder: string;
  resultCount?: number;
}>();

const query = defineModel<string>({ required: true });
const input = useTemplateRef<{ inputRef?: HTMLInputElement }>("search-input");
const emit = defineEmits<{ submit: [] }>();

function focus() {
  input.value?.inputRef?.focus();
}

function clearSearch() {
  query.value = "";
  focus();
}

function submitSearch(event: KeyboardEvent) {
  if (event.isComposing) return;
  event.preventDefault();
  emit("submit");
}

defineExpose({ focus });
</script>

<template>
  <div class="nav-search">
    <UInput
      ref="search-input"
      v-model="query"
      :placeholder="placeholder"
      leading-icon="i-tabler-search"
      size="xl"
      autocomplete="off"
      aria-label="搜索站群入口"
      class="w-full"
      :ui="{
        base: 'h-12 rounded-xl bg-default pl-11 pr-12 text-sm shadow-sm ring-1 ring-default focus-visible:ring-1 focus-visible:ring-primary',
        leadingIcon: 'ml-1 size-5 text-primary',
      }"
      @keydown.enter="submitSearch"
    >
      <template #trailing>
        <UButton
          v-if="query"
          color="neutral"
          variant="ghost"
          size="sm"
          icon="i-tabler-x"
          aria-label="清空搜索"
          @click="clearSearch"
        />
      </template>
    </UInput>
    <p
      v-if="query.trim() && resultCount"
      class="mt-2 text-sm text-muted"
      aria-live="polite"
    >
      找到 {{ resultCount }} 个入口
    </p>
  </div>
</template>
