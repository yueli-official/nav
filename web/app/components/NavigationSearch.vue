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

defineExpose({ focus });
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
        base: 'h-12 rounded-xl bg-default pl-11 pr-12 text-base shadow-sm ring-1 ring-default focus-visible:ring-2 focus-visible:ring-primary',
        leadingIcon: 'ml-1 size-5 text-primary',
      }"
      @keydown.enter.prevent="emit('submit')"
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
