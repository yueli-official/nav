<template>
  <!-- Mobile Header -->
  <header class="sticky top-0 z-50 bg-card/50 backdrop-blur-lg border-b border-border">
    <div class="flex items-center justify-between px-4 py-3">
      <div class="flex items-center gap-2">
        <div
          class="w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-chart-2 flex items-center justify-center">
          <span class="text-white font-bold text-sm">Y</span>
        </div>
        <h1 class="text-lg font-bold">月离导航</h1>
      </div>

      <div class="flex items-center gap-2">
        <ThemeToggle class="btn size-9 btn-outline"></ThemeToggle>

        <button class="btn btn-outline btn-icon" @click="settingOpen = true">
          <!-- 未登录 -->
          <svg
            v-if="!store.isLogin"
            class="size-5"
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24">
            <g
              fill="none"
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2">
              <path
                d="M9.671 4.136a2.34 2.34 0 0 1 4.659 0a2.34 2.34 0 0 0 3.319 1.915a2.34 2.34 0 0 1 2.33 4.033a2.34 2.34 0 0 0 0 3.831a2.34 2.34 0 0 1-2.33 4.033a2.34 2.34 0 0 0-3.319 1.915a2.34 2.34 0 0 1-4.659 0a2.34 2.34 0 0 0-3.32-1.915a2.34 2.34 0 0 1-2.33-4.033a2.34 2.34 0 0 0 0-3.831A2.34 2.34 0 0 1 6.35 6.051a2.34 2.34 0 0 0 3.319-1.915" />
              <circle cx="12" cy="12" r="3" />
            </g>
          </svg>
          <!-- 已登录 -->
          <svg
            v-else
            class="size-5"
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24">
            <g
              fill="none"
              stroke="currentColor"
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2">
              <circle cx="12" cy="8" r="5" />
              <path d="M20 21a8 8 0 0 0-16 0" />
            </g>
          </svg>
        </button>

        <button
          @click="mobileMenuOpen = !mobileMenuOpen"
          class="btn lg:hidden size-9 btn-square btn-outline">
          <svg
            v-if="!mobileMenuOpen"
            class="size-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24">
            <path d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Mobile Menu -->
    <div v-show="mobileMenuOpen" class="border-t border-border px-4 py-3 space-y-2">
      <button
        v-for="cat in categories"
        :key="cat.cid"
        @click="selectCategory(cat)"
        :class="[
          'w-full text-left px-4 py-2.5 rounded-lg font-medium transition-all',
          selectedCategory.cid === cat.cid
            ? 'bg-primary text-primary-foreground shadow-sm'
            : 'bg-secondary text-secondary-foreground hover:bg-secondary/80',
        ]">
        {{ cat.title }} {{ cat.cid }}
      </button>
    </div>
  </header>

  <YamiDialog
    v-model="settingOpen"
    @confirm="handleConfirm"
    title="设置密钥"
    :closeOnClickOutside="false">
    <div class="p-2 flex flex-col gap-3">
      <label for="pwd" class="text-sm">请输入管理员密码</label>
      <input id="pwd" type="text" class="input input-primary w-full" v-model="key" />
    </div>
  </YamiDialog>
</template>

<script setup lang="ts">
import type { Category } from '@/models'
import { ThemeToggle } from '@yuelioi/ui'

const store = useAuthStore()

const key = ref('')

function handleConfirm() {
  store.setApiKey(key.value)
  settingOpen.value = false

  toast.success('设置密钥成功')
  location.reload()
}

const settingOpen = ref(false)

const mobileMenuOpen = ref(false)
const selectedCategory = defineModel<Category>('selectedCategory', { required: true })
defineProps<{
  categories: Category[]
}>()

const selectCategory = (cat: Category) => {
  selectedCategory.value = cat
  mobileMenuOpen.value = false
}
</script>
