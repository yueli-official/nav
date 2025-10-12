<template>
  <div class="h-screen flex flex-col bg-background text-foreground">
    <MobileHeader
      :categories="categories"
      v-model:selectedCategory="selectedCategory"></MobileHeader>

    <div class="flex overflow-hidden h-full">
      <!-- Desktop Sidebar -->
      <aside class="hidden lg:flex flex-col w-48 border-r border-border bg-card z-10">
        <nav class="flex-1 p-4 space-y-2 flex flex-col">
          <div
            v-for="cat in categories"
            :key="cat.cid"
            @click="selectCategory(cat)"
            :class="[
              'text-left btn-lg btn btn-wide group',
              selectedCategory === cat
                ? 'bg-primary text-primary-foreground shadow-sm'
                : 'hover:bg-secondary/80',
            ]">
            <span class="w-full"> {{ cat.title }}</span>

            <div
              v-if="isLogin"
              class="hidden group-hover:block"
              @click="modifyGroupRef?.openCreate(selectedCategory.cid)">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                class="size-4 hidden group-hover:block"
                width="24"
                height="24"
                viewBox="0 0 24 24">
                <path
                  fill="none"
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M5 12h14m-7-7v14" />
              </svg>
            </div>

            <svg
              v-if="isLogin"
              xmlns="http://www.w3.org/2000/svg"
              @click.stop="modifyCateRef?.openEdit(cat)"
              class="size-4 hidden group-hover:block"
              width="24"
              height="24"
              viewBox="0 0 24 24">
              <path
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497zM15 5l4 4" />
            </svg>
          </div>

          <button v-if="isLogin" class="mt-auto btn btn-secondary btn-lg btn-wide" @click="addCate">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="size-4 mr-3"
              width="24"
              height="24"
              viewBox="0 0 24 24">
              <path
                fill="none"
                stroke="currentColor"
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M5 12h14m-7-7v14" /></svg
            ><span>添加分类</span>
          </button>
        </nav>
      </aside>

      <!-- Main Content -->
      <main class="flex-1 overflow-y-auto">
        <div class="max-w-6xl mx-auto p-4 lg:p-8">
          <!-- Search Bar -->
          <SearchInput
            v-model:search="searchQuery"
            class="my-12 bg-background/50 mx-auto max-w-2xl"></SearchInput>

          <div v-if="searchQuery" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <template v-for="c in filteredCollections" :key="c.link">
              <CollectionCard :item="c"></CollectionCard>
            </template>
          </div>

          <div class="flex items-center justify-between">
            <!-- Tags Filter -->
            <div v-if="allGroups.length > 0 && !searchQuery" class="flex flex-wrap gap-2 mb-6">
              <button
                v-for="group in allGroups"
                :key="group.gid"
                @click="jumpGroup(group)"
                class="badge badge-primary badge-lg font-semibold rounded-full transition-all">
                {{ group.title }}
              </button>
            </div>
          </div>

          <!-- Groups -->
          <div v-if="filteredGroups.length > 0 && !searchQuery" class="space-y-8">
            <section v-for="groupData in filteredGroups" :key="groupData.group.gid" class="group">
              <!-- Group Title -->

              <div class="title-deco relative mb-4 flex" :id="groupData.group.gid">
                <h2 class="">{{ groupData.group.title }}</h2>

                <div class="flex items-center gap-2 absolute right-0">
                  <button
                    v-if="isLogin"
                    @click="modifyGroupRef?.openEdit(selectedCategory.cid, groupData.group)"
                    class="btn btn-sm btn-icon btn-primary opacity-0 group-hover:opacity-100 transition-opacity">
                    <svg
                      xmlns="http://www.w3.org/2000/svg"
                      width="24"
                      height="24"
                      viewBox="0 0 24 24">
                      <path
                        fill="none"
                        stroke="currentColor"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        stroke-width="2"
                        d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497zM15 5l4 4" />
                    </svg>
                  </button>
                  <button
                    @click="
                      modifyCollectionRef?.openCreate(selectedCategory.cid, groupData.group.gid)
                    "
                    class="btn btn-sm btn-icon btn-primary opacity-0 group-hover:opacity-100 transition-opacity">
                    +
                  </button>
                </div>
              </div>

              <!-- Collections Grid -->
              <div class="grid grid-cols-1 md:grid-cols-2 group lg:grid-cols-3 gap-4">
                <template v-for="(item, idx) in groupData.collections" :key="idx">
                  <CollectionCard
                    :item="item"
                    @open-edit="openCollectionEdit(item, groupData.group.gid)"></CollectionCard>
                </template>
              </div>
            </section>
          </div>

          <!-- Empty State -->
          <div v-if="filteredGroups.length == 0" class="text-center py-16">
            <p class="text-lg text-muted-foreground">未找到相关资源</p>
          </div>
        </div>
        <TheFooter />
      </main>
    </div>
  </div>
  <ModifyCategory ref="modifyCateRef"></ModifyCategory>
  <ModifyCollection ref="modifyCollectionRef"></ModifyCollection>
  <ModifyGroup ref="modifyGroupRef"></ModifyGroup>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { SearchInput } from '@yuelioi/ui'
import type ModifyCollection from '@/components/ModifyCollection.vue'
import type ModifyCategory from '@/components/ModifyCategory.vue'
import type ModifyGroup from '@/components/ModifyGroup.vue'
import type { Collection } from '@/models'

const store = useAuthStore()
const { isLogin } = store

const modifyCollectionRef = ref<InstanceType<typeof ModifyCollection> | null>()
const modifyCateRef = ref<InstanceType<typeof ModifyCategory> | null>()
const modifyGroupRef = ref<InstanceType<typeof ModifyGroup> | null>()

const searchQuery = ref('')

const apiData = ref<BookmarkData>([])
const searchTerm = ref('')
const selectedCategory = ref<Category>({
  cid: '',
  title: '',
  order: 0,
})

function addCate() {
  modifyCateRef.value?.openCreate()
}

function jumpGroup(group: Group) {
  const el = document.getElementById(group.gid)
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

function openCollectionEdit(item: Collection, gid: string) {
  modifyCollectionRef.value?.openEdit(item, selectedCategory.value.cid, gid, item.link)
}

// 获取数据
onMounted(async () => {
  try {
    const resp = await navApi.getNavs()
    apiData.value = resp.data.data || []
    if (apiData.value.length > 0) {
      selectedCategory.value = apiData.value[0].category
    }
  } catch (error) {
    console.error('Failed to fetch data:', error)
  }
})

// 获取所有分类
const categories = computed(() => {
  return apiData.value.map((d) => d.category).sort((a, b) => a.order - b.order)
})

// 获取当前分类的数据
const currentCategoryData = computed(() => {
  return apiData.value.find((d) => d.category.cid === selectedCategory.value.cid)
})

// 获取当前分类下的所有标签
const allGroups = computed(() => {
  if (!currentCategoryData.value) return []
  return currentCategoryData.value.groups.map((group) => group.group)
})

// 获取搜索过滤
const filteredCollections = computed(() => {
  const q = (searchQuery.value || '').trim().toLowerCase()
  const all = apiData.value.flatMap((cat) => (cat.groups || []).flatMap((g) => g.collections || []))

  if (!q) return all

  return apiData.value
    .flatMap((cat) => (cat.groups || []).flatMap((g) => g.collections || []))
    .filter(
      (c) =>
        (c.title || '').toLowerCase().includes(q) ||
        (c.description || '').toLowerCase().includes(q) ||
        (c.tags || []).some((t) => (t || '').toLowerCase().includes(q)),
    )
})

// 过滤后的分组数据
const filteredGroups = computed(() => {
  if (!currentCategoryData.value) return []

  return currentCategoryData.value.groups.map((group) => ({
    ...group,
    collections: group.collections.filter((item) => {
      const matchSearch =
        !searchTerm.value ||
        item.title.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
        (item.description &&
          item.description.toLowerCase().includes(searchTerm.value.toLowerCase()))
      return matchSearch
    }),
  }))
})

const selectCategory = (cat: Category) => {
  selectedCategory.value = cat
}
</script>
