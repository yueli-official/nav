<template>
  <div class="relative group/card">
    <!-- 链接提示框 -->
    <div
      class="absolute z-10 bg-popover/95 glass text-popover-foreground px-3 py-2 rounded-md shadow-lg border border-border/50 -top-14 left-1/2 -translate-x-1/2 opacity-0 group-hover/card:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap max-w-xs truncate text-sm">
      <div class="flex items-center gap-2">
        <svg
          class="w-3.5 h-3.5 flex-shrink-0 text-muted-foreground"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
        </svg>
        <span class="truncate">{{ item.link }}</span>
      </div>
      <!-- 小三角形 -->
      <div
        class="absolute -bottom-1 left-1/2 -translate-x-1/2 w-2 h-2 bg-popover/95 border-r border-b border-border/50 rotate-45"></div>
    </div>

    <!-- 卡片主体 -->
    <a
      :href="item.link"
      target="_blank"
      rel="noopener noreferrer"
      class="card-light-sweep block relative h-48 overflow-hidden p-6 bg-card/80 border border-border/60 rounded-xl shadow-sm group-hover/card:shadow-xl hover:border-foreground/15 hover:-translate-y-[3px] transition-all duration-350 group"
      style="transition-timing-function: cubic-bezier(0.23, 1, 0.32, 1);">
      <!-- 背景装饰 -->
      <div
        class="absolute inset-0 bg-gradient-to-br from-foreground/[0.02] via-transparent to-transparent group-hover/card:from-foreground/[0.05] transition-opacity duration-300"></div>
      <!-- 右上角装饰圆弧 -->
      <div
        class="absolute -top-16 -right-16 w-40 h-40 rounded-full border border-foreground/[0.06] group-hover/card:border-foreground/[0.12] transition-all duration-500"></div>
      <div
        class="absolute -top-20 -right-20 w-52 h-52 rounded-full border border-foreground/[0.04] group-hover/card:border-foreground/[0.08] transition-all duration-500"></div>

      <!-- 内容区域 -->
      <div class="relative z-10 h-full flex flex-col">
        <!-- 标题和图标 -->
        <div class="flex items-start gap-3 mb-3">
          <h3
            class="font-semibold text-lg line-clamp-1 text-card-foreground group-hover/card:text-foreground transition-colors duration-200">
            {{ item.title }}
          </h3>
          <div class="ml-auto"></div>
          <div
            v-if="store.isLogin"
            @click.stop.prevent="openCollectionEdit"
            class="flex-shrink-0 w-8 h-8 rounded-full bg-accent/50 flex items-center justify-center opacity-0 group-hover/card:opacity-100 transition-all duration-200 hover:bg-accent">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="size-4"
              width="24"
              height="24"
              viewBox="0 0 24 24">
              <g
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2">
                <path d="M12 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                <path
                  d="M18.375 2.625a1 1 0 0 1 3 3l-9.013 9.014a2 2 0 0 1-.853.505l-2.873.84a.5.5 0 0 1-.62-.62l.84-2.873a2 2 0 0 1 .506-.852z" />
              </g>
            </svg>
          </div>
          <div
            class="flex-shrink-0 w-8 h-8 rounded-full bg-accent/50 flex items-center justify-center opacity-0 group-hover/card:opacity-100 transition-all duration-200 group-hover/card:rotate-45 hover:bg-accent">
            <svg class="w-4 h-4 text-foreground" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
            </svg>
          </div>
        </div>

        <!-- 描述 -->
        <p
          v-if="item.description"
          class="text-sm text-muted-foreground mb-4 line-clamp-2 leading-relaxed">
          {{ item.description }}
        </p>

        <!-- 标签 -->
        <div v-if="item.tags && item.tags.length > 0" class="flex flex-wrap gap-2 mt-auto">
          <span
            v-for="tag in item.tags"
            :key="tag"
            class="inline-flex items-center gap-1.5 px-2.5 py-1 text-xs font-medium bg-accent/80 group-hover/card:bg-accent text-accent-foreground rounded-full transition-colors duration-200">
            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
            </svg>
            {{ tag }}
          </span>
        </div>
      </div>

      <!-- 装饰性光效 -->
      <div
        class="absolute -right-8 -top-8 w-32 h-32 bg-foreground/[0.04] rounded-full blur-3xl opacity-0 group-hover/card:opacity-100 transition-opacity duration-500"></div>
    </a>
  </div>
</template>

<script setup lang="ts">
const store = useAuthStore()

interface Collection {
  title: string
  description?: string
  link: string
  tags?: string[]
}

const props = defineProps<{ item: Collection }>()
const emit = defineEmits(['openEdit'])

function openCollectionEdit() {
  emit('openEdit', props.item)
}
</script>
