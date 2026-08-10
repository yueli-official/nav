<script setup lang="ts">
// PROTOTYPE — compares the current framed sidebar with the proposed commercial direction.
const route = useRoute();
const router = useRouter();
const visible = import.meta.dev;

const variants = [
  { key: "commercial", label: "推荐方案" },
  { key: "baseline", label: "当前方案" },
] as const;

const currentIndex = computed(() =>
  route.query.sidebar === "baseline" ? 1 : 0,
);
const current = computed(() => variants[currentIndex.value]);

async function cycle(direction: -1 | 1) {
  const nextIndex =
    (currentIndex.value + direction + variants.length) % variants.length;
  const next = variants[nextIndex]!;
  const query = { ...route.query };
  if (next.key === "commercial") delete query.sidebar;
  else query.sidebar = next.key;
  await router.replace({ query });
}

function handleKeydown(event: KeyboardEvent) {
  const target = event.target as HTMLElement | null;
  if (
    target?.matches("input, textarea, select, [contenteditable='true']") ||
    (event.key !== "ArrowLeft" && event.key !== "ArrowRight")
  ) {
    return;
  }
  event.preventDefault();
  void cycle(event.key === "ArrowLeft" ? -1 : 1);
}

onMounted(() => window.addEventListener("keydown", handleKeydown));
onBeforeUnmount(() => window.removeEventListener("keydown", handleKeydown));
</script>

<template>
  <div
    v-if="visible"
    class="fixed bottom-16 left-1/2 z-[90] flex -translate-x-1/2 items-center gap-1 rounded-full bg-inverted p-1 text-inverted shadow-lg"
    data-sidebar-prototype-switcher
  >
    <UButton
      type="button"
      color="neutral"
      variant="ghost"
      icon="i-tabler-chevron-left"
      square
      size="sm"
      aria-label="查看上一个侧栏方案"
      class="rounded-full text-inverted hover:bg-default/15"
      @click="cycle(-1)"
    />
    <span class="min-w-24 px-2 text-center text-xs font-medium">
      {{ current.label }}
    </span>
    <UButton
      type="button"
      color="neutral"
      variant="ghost"
      icon="i-tabler-chevron-right"
      square
      size="sm"
      aria-label="查看下一个侧栏方案"
      class="rounded-full text-inverted hover:bg-default/15"
      @click="cycle(1)"
    />
  </div>
</template>
