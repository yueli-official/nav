<script setup lang="ts">
import type { NavigationCategory } from "~/types/navigation";

const { categories } = defineProps<{ categories: NavigationCategory[] }>();
const selected = defineModel<string>({ required: true });

const selectItems = computed(() =>
  categories.map((category) => ({
    label: category.title,
    value: category.id,
    icon: category.icon,
  })),
);
</script>

<template>
  <div>
    <USelectMenu
      v-model="selected"
      :items="selectItems"
      value-key="value"
      label-key="label"
      :search-input="false"
      class="w-full md:hidden"
      aria-label="选择分类"
    />

    <nav
      class="hidden overflow-x-auto overflow-y-hidden pb-1 md:block"
      aria-label="导航分类"
    >
      <div
        class="grid min-w-max auto-cols-[minmax(132px,1fr)] grid-flow-col gap-2 xl:min-w-0"
      >
        <button
          v-for="category in categories"
          :key="category.id"
          type="button"
          class="group flex min-h-12 items-center gap-2 rounded-xl border px-3 py-2 text-left text-sm font-medium transition duration-200 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary"
          :class="
            selected === category.id
              ? 'border-primary/45 bg-primary/10 text-primary'
              : 'border-default bg-default text-muted hover:border-primary/30 hover:bg-elevated hover:text-highlighted'
          "
          :aria-current="selected === category.id ? 'true' : undefined"
          @click="
            () => {
              selected = category.id;
            }
          "
        >
          <UIcon :name="category.icon" class="size-5 shrink-0" />
          <span class="truncate">{{ category.title }}</span>
        </button>
      </div>
    </nav>
  </div>
</template>
