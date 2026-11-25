interface PortalQuickLink {
  title: string;
  path: string;
  action: string;
  icon: string;
  keywords?: string;
}

interface PortalSite {
  id: string;
  name: string;
  url: string;
  logo: string;
  description: string;
  links: PortalQuickLink[];
}

export interface PortalEntry {
  title: string;
  site: string;
  url: string;
  action: string;
  icon: string;
}

const siteDefinitions: PortalSite[] = [
  {
    id: "apps",
    name: "月离软件",
    url: "https://apps.yuelili.com/",
    logo: "https://apps.yuelili.com/favicon.svg",
    description: "查找创作和办公软件、插件与实用工具。每款软件都有独立的发布页，可查看版本、适用平台和下载选项。",
    links: [
      { title: "软件目录", path: "/software", action: "浏览下载", icon: "i-tabler-package", keywords: "找软件和插件 应用 工具 下载 安装" },
      { title: "剪辑与特效工具", path: "/software/effectcraft", action: "下载软件", icon: "i-tabler-wand", keywords: "EffectCraft 音视频 后期 Windows" },
      { title: "PDF 工具", path: "/software/pdfcraft", action: "下载软件", icon: "i-tabler-file-type-pdf", keywords: "PdfCraft 办公 PDF Windows" },
    ],
  },
  {
    id: "blog",
    name: "月离博客",
    url: "https://cg.yuelili.com/",
    logo: "https://cg.yuelili.com/favicon.ico?v=20260929",
    description: "阅读面向后期创作者的 AE 教程、效果拆解和工作流经验。可以跟着系列从头学，也可以直接查某个效果的做法。",
    links: [
      { title: "找 AE 教程", path: "/search?q=AE", action: "阅读教程", icon: "i-tabler-sparkles", keywords: "After Effects AfterEffects Adobe 动效 合成" },
      { title: "教程系列", path: "/series", action: "浏览系列", icon: "i-tabler-books" },
      { title: "文章归档", path: "/archive", action: "浏览文章", icon: "i-tabler-archive" },
    ],
  },
  {
    id: "docs",
    name: "月离文档",
    url: "https://docs.yuelili.com/",
    logo: "https://docs.yuelili.com/favicon.svg",
    description: "查阅软件手册、插件说明和开发参考。文档按项目成套整理，适合查具体参数、API 或脚本写法。",
    links: [
      { title: "文档集", path: "/collections", action: "浏览文档", icon: "i-tabler-books", keywords: "说明手册 Houdini VEX 插件" },
      { title: "Premiere 脚本开发", path: "/pr-scripting", action: "阅读文档", icon: "i-tabler-file-code", keywords: "Premiere Pro PR Adobe 脚本 ExtendScript 自动化 API" },
      { title: "JavaScript Tools 文档", path: "/javascript-tools", action: "阅读文档", icon: "i-tabler-code", keywords: "javascript-tools Adobe 脚本 JavaScript" },
    ],
  },
  {
    id: "paste",
    name: "月离代码分享",
    url: "https://paste.yuelili.com/",
    logo: "https://paste.yuelili.com/favicon.svg",
    description: "把代码、脚本或文本打包成一个链接。一次分享可包含多个文件，适合发复现示例、配置片段和临时笔记。",
    links: [
      { title: "分享代码和文本", path: "/", action: "在线使用", icon: "i-tabler-code", keywords: "Paste 粘贴 代码片段 多文件 脚本 新建分享" },
      { title: "我的分享", path: "/mine", action: "登录后管理", icon: "i-tabler-folder" },
    ],
  },
  {
    id: "forum",
    name: "月离社区",
    url: "https://forum.yuelili.com/",
    logo: "https://forum.yuelili.com/favicon.svg",
    description: "向其他创作者提问、分享经验和交流作品。遇到具体问题时，可以在帖子里附上版本、截图或代码链接，方便大家一起排查。",
    links: [
      { title: "浏览讨论", path: "/", action: "进入社区", icon: "i-tabler-messages", keywords: "论坛 经验交流" },
      { title: "发帖提问", path: "/new", action: "登录后发帖", icon: "i-tabler-pencil", keywords: "求助 发帖" },
      { title: "社区指南", path: "/guidelines", action: "阅读指南", icon: "i-tabler-book" },
    ],
  },
  {
    id: "gallery",
    name: "月离图库",
    url: "https://gallery.yuelili.com/",
    logo: "https://gallery.yuelili.com/favicon.svg?v=20261024-brand",
    description: "浏览图片作品和主题合集，查找视觉灵感与创作参考。也可以通过榜单发现近期受到关注的作品。",
    links: [
      { title: "浏览图片", path: "/images", action: "查看作品", icon: "i-tabler-photo", keywords: "图库 灵感 创作参考" },
      { title: "主题合集", path: "/collections", action: "浏览合集", icon: "i-tabler-folders" },
      { title: "图片榜单", path: "/rankings", action: "查看榜单", icon: "i-tabler-chart-bar" },
    ],
  },
  {
    id: "tv",
    name: "月离影视",
    url: "https://tv.yuelili.com/",
    logo: "https://tv.yuelili.com/favicon.svg",
    description: "浏览和搜索电影、剧集，进入详情页查看简介、集数和播放入口。还可以从观看记录继续上次的进度。",
    links: [
      { title: "浏览影视", path: "/", action: "查看影片", icon: "i-tabler-device-tv", keywords: "电影 剧集 电视剧" },
      { title: "搜索片名", path: "/search", action: "搜索影视", icon: "i-tabler-search" },
      { title: "观看记录", path: "/history", action: "继续观看", icon: "i-tabler-history" },
    ],
  },
  {
    id: "video",
    name: "fun",
    url: "https://fun.yuelili.com/",
    logo: "https://fun.yuelili.com/favicon.svg",
    description: "发现社区分享的 B 站视频和创作者。可以按内容浏览，找到喜欢的创作者，也可以提交自己的推荐。",
    links: [
      { title: "浏览视频", path: "/browse", action: "查看视频", icon: "i-tabler-player-play", keywords: "Bilibili B站 哔哩哔哩 沙雕" },
      { title: "创作者目录", path: "/creators", action: "浏览创作者", icon: "i-tabler-users" },
      { title: "投稿视频", path: "/submit", action: "登录后投稿", icon: "i-tabler-upload" },
    ],
  },
  {
    id: "www",
    name: "月离万事屋官网",
    url: "https://www.yuelili.com/",
    logo: "https://www.yuelili.com/favicon.svg?v=20261124-brand",
    description: "了解月离在做的产品和项目，查看各站的介绍与地址。需要反馈、合作或联系月离时，也可以从这里找到联系方式。",
    links: [
      { title: "了解月离", path: "/about", action: "阅读介绍", icon: "i-tabler-info-circle" },
      { title: "站点目录", path: "/sites", action: "查看站点", icon: "i-tabler-world", keywords: "所有站点 站群" },
      { title: "联系月离", path: "/contact", action: "打开联系页", icon: "i-tabler-mail", keywords: "反馈 问题 建议 合作 QQ" },
    ],
  },
  {
    id: "account",
    name: "月离用户中心",
    url: "https://account.yuelili.com/",
    logo: "https://account.yuelili.com/favicon.svg?v=20260929",
    description: "管理月离站群共用的账号。登录后可以更新个人资料、设置安全选项，以及查看和管理已登录的设备。",
    links: [
      { title: "个人资料", path: "/me", action: "登录后管理", icon: "i-tabler-user-circle", keywords: "账号 账户 管理月离账号" },
      { title: "账号安全", path: "/security", action: "登录后管理", icon: "i-tabler-shield-lock", keywords: "密码 安全" },
      { title: "登录设备", path: "/sessions", action: "登录后管理", icon: "i-tabler-devices", keywords: "会话 设备" },
    ],
  },
];

export const portalSites = siteDefinitions.map((site) => ({
  ...site,
  hostname: new URL(site.url).hostname,
  links: site.links.map((link) => ({ ...link, url: new URL(link.path, site.url).href })),
}));

const siteEntries = portalSites.map((site) =>
  site.links.map(({ title, url, action, icon, keywords }) => ({
    entry: { title, site: site.name, url, action, icon },
    text: [title, site.name, action, site.url, keywords].join(" ").normalize("NFKC").toLowerCase().replace(/\s+/g, ""),
  })),
);

const searchableEntries = siteEntries.flat();
export const portalRecommendations = siteEntries.flatMap(([first]) => first ? [first.entry] : []);

export function searchPortalEntries(query: string) {
  const terms = query.normalize("NFKC").trim().toLowerCase().split(/\s+/);
  return searchableEntries
    .filter(({ text }) => terms.every((term) => text.includes(term)))
    .map(({ entry }) => entry);
}
