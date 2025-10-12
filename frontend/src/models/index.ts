// 收藏项
export interface Collection {
  title: string
  link: string
  description?: string
  tags?: string[]
}

// 分组信息
export interface Group {
  gid: string
  title: string
  order: number
}

// 带收藏项的分组
export interface GroupWithCollections {
  group: Group
  collections: Collection[]
}

// 分类信息
export interface Category {
  cid: string
  title: string
  order: number
}

// 完整的分类数据
export interface CategoryData {
  category: Category
  groups: GroupWithCollections[]
}

// 整个数据结构（数组）
export type BookmarkData = CategoryData[]
