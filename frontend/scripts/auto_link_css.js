import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

// 从src vue里提取

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)

const srcDir = path.join(__dirname, '../src')
const pluginsCssPath = path.join(srcDir, 'assets/plugins.css')
const uiPackageName = '@yuelioi/ui'

// 确保 plugins.css 存在
if (!fs.existsSync(pluginsCssPath)) {
  fs.writeFileSync(pluginsCssPath, "@reference '@yuelioi/yami/src/index.css';\n\n", 'utf-8')
}

// 递归查找 .vue 文件
function findVueFiles(dir) {
  let results = []
  for (const file of fs.readdirSync(dir)) {
    const fullPath = path.join(dir, file)
    const stat = fs.statSync(fullPath)
    if (stat.isDirectory()) {
      results = results.concat(findVueFiles(fullPath))
    } else if (file.endsWith('.vue')) {
      results.push(fullPath)
    }
  }
  return results
}

const vueFiles = findVueFiles(srcDir)
const componentSet = new Set()

// 遍历每个 .vue 文件
for (const filePath of vueFiles) {
  const content = fs.readFileSync(filePath, 'utf-8')
  // 匹配 import { Xxx } from '@yuelioi/ui'
  const importRegex = new RegExp(`import\\s+\\{([^}]+)\\}\\s+from\\s+['"]${uiPackageName}['"]`, 'g')
  let match
  while ((match = importRegex.exec(content)) !== null) {
    const components = match[1]
      .split(',')
      .map((c) => c.trim())
      .filter(Boolean)
    components.forEach((c) => componentSet.add(c))
  }
}

// 生成 @source 路径
const lines = Array.from(componentSet).map(
  (c) => `@source "../../node_modules/${uiPackageName}/src/components/${c}.vue";`,
)

// 避免重复追加
const originalContent = fs.readFileSync(pluginsCssPath, 'utf-8')
const newLines = lines.filter((line) => !originalContent.includes(line))

if (newLines.length > 0) {
  fs.appendFileSync(pluginsCssPath, newLines.join('\n') + '\n', 'utf-8')
  console.log('✅ plugins.css 已更新，追加组件:', newLines)
} else {
  console.log('ℹ️ 没有新的组件需要追加')
}
