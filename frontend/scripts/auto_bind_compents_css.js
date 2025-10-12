import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

// 从 auto components 提取

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const srcDir = path.join(__dirname, '../src')
const pluginsCssPath = path.join(srcDir, 'assets/plugins.css')
const dtsPath = path.join(srcDir, 'components.d.ts')
const uiPackageName = '@yuelioi/ui'

// 🔹 1️⃣ 确保 plugins.css 存在
if (!fs.existsSync(pluginsCssPath)) {
  fs.writeFileSync(pluginsCssPath, "@reference '@yuelioi/yami/src/index.css';\n\n", 'utf-8')
  console.log('🪄 已创建 plugins.css')
}

// 🔹 2️⃣ 读取 components.d.ts，解析组件名
if (!fs.existsSync(dtsPath)) {
  console.error('❌ 未找到 src/components.d.ts，请先运行 Vite 生成自动导入声明文件')
  process.exit(1)
}

const dtsContent = fs.readFileSync(dtsPath, 'utf-8')

// 匹配像 “BackToTop: typeof import('@yuelioi/ui')['BackToTop']” 的部分
const regex = /:\s*typeof import\(['"]@yuelioi\/ui['"]\)\[['"]([A-Za-z0-9_]+)['"]\]/g
const components = new Set()
let match

while ((match = regex.exec(dtsContent)) !== null) {
  components.add(match[1])
}

// 🔹 3️⃣ 生成 @source 路径
const lines = Array.from(components).map(
  (c) => `@source "../../node_modules/${uiPackageName}/src/components/${c}.vue";`,
)

// 🔹 4️⃣ 避免重复追加
const originalContent = fs.readFileSync(pluginsCssPath, 'utf-8')
const newLines = lines.filter((line) => !originalContent.includes(line))

if (newLines.length > 0) {
  fs.appendFileSync(pluginsCssPath, newLines.join('\n') + '\n', 'utf-8')
  console.log(`✅ plugins.css 已更新，追加组件:`, newLines)
} else {
  console.log('ℹ️ 没有新的组件需要追加')
}
