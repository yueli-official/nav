CREATE TABLE nav_categories (
    id          TEXT PRIMARY KEY,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    icon        TEXT NOT NULL,
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE nav_groups (
    id          TEXT PRIMARY KEY,
    category_id TEXT NOT NULL REFERENCES nav_categories(id) ON DELETE RESTRICT,
    title       TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT uq_nav_groups_category UNIQUE (id, category_id)
);

CREATE TABLE nav_links (
    id          TEXT PRIMARY KEY,
    category_id TEXT NOT NULL REFERENCES nav_categories(id) ON DELETE RESTRICT,
    group_id    TEXT NOT NULL,
    title       TEXT NOT NULL,
    url         TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    tags        JSONB NOT NULL DEFAULT '[]'::jsonb,
    keywords    JSONB NOT NULL DEFAULT '[]'::jsonb,
    kind        TEXT NOT NULL,
    featured    BOOLEAN NOT NULL DEFAULT FALSE,
    status      TEXT NOT NULL DEFAULT 'draft',
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT fk_nav_links_group_category FOREIGN KEY (group_id, category_id)
        REFERENCES nav_groups(id, category_id) ON DELETE RESTRICT,
    CONSTRAINT ck_nav_links_status CHECK (status IN ('published', 'draft', 'archived')),
    CONSTRAINT ck_nav_links_kind CHECK (kind IN ('official', 'tool', 'community', 'learning', 'resource', 'reference', 'research')),
    CONSTRAINT ck_nav_links_url CHECK (url ~ '^https?://')
);

CREATE INDEX ix_nav_groups_category ON nav_groups(category_id, sort_order);
CREATE INDEX ix_nav_links_browse ON nav_links(status, category_id, group_id, sort_order);
CREATE INDEX ix_nav_links_featured ON nav_links(featured, sort_order) WHERE status = 'published';

INSERT INTO nav_categories (id, title, description, icon, sort_order) VALUES ('create', '创意制作', '设计灵感、视觉创作、动效与三维工作流。', 'i-tabler-palette', 0);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('inspiration', 'create', '灵感与作品', '发现优秀作品、视觉趋势和创作者。', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('behance', 'create', 'inspiration', 'Behance', 'https://www.behance.net/', 'Adobe 旗下创意作品展示与发现平台。', '["作品集","设计","灵感"]'::jsonb, '[]'::jsonb, 'community', TRUE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('dribbble', 'create', 'inspiration', 'Dribbble', 'https://dribbble.com/', 'UI、品牌与插画设计作品社区。', '["UI","品牌","插画"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('artstation', 'create', 'inspiration', 'ArtStation', 'https://www.artstation.com/', '游戏、影视与数字艺术作品平台。', '["概念设计","游戏美术","3D"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('pinterest', 'create', 'inspiration', 'Pinterest', 'https://www.pinterest.com/', '用视觉收藏板组织灵感和参考资料。', '["灵感","收藏","图片"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('design-tools', 'create', '设计工具', '从界面设计到配色与交互练习。', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('figma-community', 'create', 'design-tools', 'Figma Community', 'https://www.figma.com/community', '组件、模板、插件与公开设计文件。', '["Figma","UI","组件"]'::jsonb, '[]'::jsonb, 'resource', TRUE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('adobe-color', 'create', 'design-tools', 'Adobe Color', 'https://color.adobe.com/', '创建、提取和探索配色方案。', '["配色","色彩","Adobe"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('bezier-game', 'create', 'design-tools', 'The Bézier Game', 'https://bezier.method.ac/', '通过关卡练习钢笔工具与贝塞尔曲线。', '["贝塞尔","钢笔工具","练习"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('uigradients', 'create', 'design-tools', 'uiGradients', 'https://uigradients.com/', '浏览并复制适合界面使用的渐变组合。', '["渐变","配色","CSS"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('motion-3d', 'create', '动效与三维', '动效、程序化三维与实时图形入口。', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('entagma', 'create', 'motion-3d', 'Entagma', 'https://entagma.com/', 'Houdini、程序化设计和高级 CG 教程。', '["Houdini","程序化","CG"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('blender', 'create', 'motion-3d', 'Blender', 'https://www.blender.org/', '开源三维创作套件及官方资源入口。', '["Blender","3D","开源"]'::jsonb, '[]'::jsonb, 'official', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('sidefx', 'create', 'motion-3d', 'SideFX Houdini', 'https://www.sidefx.com/', 'Houdini 产品、学习资料与社区入口。', '["Houdini","VFX","3D"]'::jsonb, '[]'::jsonb, 'official', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('unreal-engine', 'create', 'motion-3d', 'Unreal Engine', 'https://www.unrealengine.com/', '实时三维引擎、样例与学习资源。', '["实时渲染","游戏引擎","3D"]'::jsonb, '[]'::jsonb, 'official', FALSE, 'published', 3);
INSERT INTO nav_categories (id, title, description, icon, sort_order) VALUES ('develop', '开发工程', '权威文档、开发实验、代码社区与规范。', 'i-tabler-code', 1);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('references', 'develop', '文档与规范', '直接进入权威技术资料。', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('mdn', 'develop', 'references', 'MDN Web Docs', 'https://developer.mozilla.org/', 'Web 平台技术与浏览器 API 文档。', '["Web","JavaScript","CSS"]'::jsonb, '[]'::jsonb, 'reference', TRUE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('typescript', 'develop', 'references', 'TypeScript Documentation', 'https://www.typescriptlang.org/docs/', 'TypeScript 官方手册与参考文档。', '["TypeScript","文档","类型"]'::jsonb, '[]'::jsonb, 'reference', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('go-docs', 'develop', 'references', 'Go Documentation', 'https://go.dev/doc/', 'Go 官方文档、教程与语言规范。', '["Go","后端","文档"]'::jsonb, '[]'::jsonb, 'reference', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('rust-book', 'develop', 'references', 'The Rust Programming Language', 'https://doc.rust-lang.org/book/', 'Rust 官方入门与语言实践手册。', '["Rust","系统编程","文档"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('playgrounds', 'develop', '实验与工具', '快速验证代码、格式和前端方案。', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('codepen', 'develop', 'playgrounds', 'CodePen', 'https://codepen.io/', '前端代码实验、作品与社区平台。', '["前端","Playground","CSS"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('transform-tools', 'develop', 'playgrounds', 'Transform', 'https://transform.tools/', '在多种代码与数据格式之间快速转换。', '["格式转换","JSON","TypeScript"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('regex101', 'develop', 'playgrounds', 'regex101', 'https://regex101.com/', '带解释、测试和调试能力的正则工具。', '["正则","调试","工具"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('diagrams', 'develop', 'playgrounds', 'diagrams.net', 'https://app.diagrams.net/', '绘制流程图、架构图和技术草图。', '["流程图","架构","协作"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('communities', 'develop', '社区与练习', '讨论技术、阅读源码并持续练习。', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('github', 'develop', 'communities', 'GitHub', 'https://github.com/', '开源代码托管、协作与项目发现平台。', '["开源","Git","协作"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('linux-do', 'develop', 'communities', 'LINUX DO', 'https://linux.do/', '面向开发者的中文技术社区。', '["中文社区","Linux","开发"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('leetcode-cn', 'develop', 'communities', '力扣', 'https://leetcode.cn/', '算法题、面试练习与竞赛社区。', '["算法","面试","练习"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('codeforces', 'develop', 'communities', 'Codeforces', 'https://codeforces.com/', '算法竞赛、题库与全球排名平台。', '["竞赛","算法","题库"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 3);
INSERT INTO nav_categories (id, title, description, icon, sort_order) VALUES ('ai', 'AI 与研究', '通用模型、开发平台、开源生态与研究资料。', 'i-tabler-sparkles', 2);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('assistants', 'ai', 'AI 产品', '面向写作、研究、编程和多模态创作。', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('chatgpt', 'ai', 'assistants', 'ChatGPT', 'https://chatgpt.com/', 'OpenAI 的通用 AI 助手。', '["对话","写作","多模态"]'::jsonb, '[]'::jsonb, 'tool', TRUE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('qwen', 'ai', 'assistants', '通义千问', 'https://tongyi.aliyun.com/', '阿里云推出的通用 AI 助手。', '["中文","对话","多模态"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('runway', 'ai', 'assistants', 'Runway', 'https://runwayml.com/', '面向视频与视觉创作的生成式 AI 工具。', '["视频","生成式 AI","创作"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('deepmind', 'ai', 'assistants', 'Google DeepMind', 'https://deepmind.google/', 'AI 研究成果、项目与科学进展。', '["研究","科学","Google"]'::jsonb, '[]'::jsonb, 'research', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('ai-developers', 'ai', '开发与模型', '构建、评估和部署 AI 应用。', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('openai-platform', 'ai', 'ai-developers', 'OpenAI Platform', 'https://platform.openai.com/', 'OpenAI API、控制台与开发文档入口。', '["API","模型","开发"]'::jsonb, '[]'::jsonb, 'official', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('hugging-face', 'ai', 'ai-developers', 'Hugging Face', 'https://huggingface.co/', '开源模型、数据集与机器学习应用社区。', '["开源模型","数据集","ML"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('kaggle', 'ai', 'ai-developers', 'Kaggle', 'https://www.kaggle.com/', '数据集、Notebook、竞赛与机器学习社区。', '["数据科学","竞赛","Notebook"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('deeplearning-ai', 'ai', 'ai-developers', 'DeepLearning.AI', 'https://www.deeplearning.ai/', '机器学习与生成式 AI 课程和资讯。', '["课程","深度学习","AI"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 3);
INSERT INTO nav_categories (id, title, description, icon, sort_order) VALUES ('assets', '创作素材', '图片、字体、图标、三维、材质与音视频资源。', 'i-tabler-photo', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('images', 'assets', '图片与插画', '可用于设计参考和项目制作的视觉资源。', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('pexels', 'assets', 'images', 'Pexels', 'https://www.pexels.com/', '免费图片与视频素材平台。', '["图片","视频","免费素材"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('pixabay', 'assets', 'images', 'Pixabay', 'https://pixabay.com/', '图片、插画、视频、音乐与音效资源。', '["图片","音频","视频"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('undraw', 'assets', 'images', 'unDraw', 'https://undraw.co/illustrations', '可调整主题色的开源插画集合。', '["插画","SVG","开源"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('nasa-images', 'assets', 'images', 'NASA Image Library', 'https://images.nasa.gov/', 'NASA 官方图片、视频与音频资料库。', '["太空","科学","影像"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('icons-fonts', 'assets', '图标与字体', '建立稳定视觉语言的基础素材。', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('iconify', 'assets', 'icons-fonts', 'Iconify', 'https://iconify.design/', '统一检索并使用多套开源图标集合。', '["图标","Iconify","开源"]'::jsonb, '[]'::jsonb, 'resource', TRUE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('tabler-icons', 'assets', 'icons-fonts', 'Tabler Icons', 'https://tabler.io/icons', '风格统一的开源线性图标库。', '["图标","Tabler","SVG"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('font-awesome', 'assets', 'icons-fonts', 'Font Awesome', 'https://fontawesome.com/', '常用图标工具包与品牌图标集合。', '["图标","Web","品牌"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('fontspace', 'assets', 'icons-fonts', 'FontSpace', 'https://www.fontspace.com/', '按风格与授权方式浏览字体资源。', '["字体","排版","资源"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('three-d', 'assets', '三维与材质', '模型、贴图、HDRI 与动作素材。', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('poly-haven', 'assets', 'three-d', 'Poly Haven', 'https://polyhaven.com/', '开放授权的 HDRI、材质与三维模型。', '["HDRI","材质","3D"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('ambientcg', 'assets', 'three-d', 'ambientCG', 'https://ambientcg.com/', 'CC0 授权的 PBR 材质与环境资源。', '["PBR","材质","CC0"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('sketchfab', 'assets', 'three-d', 'Sketchfab', 'https://sketchfab.com/', '在线查看、发布与发现三维模型。', '["3D","模型","预览"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('mixamo', 'assets', 'three-d', 'Mixamo', 'https://www.mixamo.com/', '角色自动绑定与动作资源平台。', '["动画","角色","3D"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 3);
INSERT INTO nav_categories (id, title, description, icon, sort_order) VALUES ('learn', '学习成长', '系统课程、公开课、语言与知识探索。', 'i-tabler-school', 4);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('courses', 'learn', '课程平台', '从大学课程到职业技能的系统学习。', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('coursera', 'learn', 'courses', 'Coursera', 'https://www.coursera.org/', '大学与企业提供的在线课程和证书项目。', '["在线课程","大学","证书"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('edx', 'learn', 'courses', 'edX', 'https://www.edx.org/', '全球高校与机构的在线课程平台。', '["公开课","大学","学习"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('khan-academy', 'learn', 'courses', 'Khan Academy', 'https://www.khanacademy.org/', '免费学习数学、科学、经济与编程。', '["免费课程","基础教育","科学"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('xuetangx', 'learn', 'courses', '学堂在线', 'https://www.xuetangx.com/', '面向中文用户的在线课程与认证平台。', '["中文课程","大学","认证"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('language', 'learn', '语言学习', '听说训练、词典与真实语料入口。', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('ted', 'learn', 'language', 'TED', 'https://www.ted.com/', '演讲、字幕和多语言知识内容。', '["演讲","英语","字幕"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('bbc-learning-english', 'learn', 'language', 'BBC Learning English', 'https://www.bbc.co.uk/learningenglish/', '英语听力、词汇、语法与新闻课程。', '["英语","听力","BBC"]'::jsonb, '[]'::jsonb, 'learning', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('cambridge-dictionary', 'learn', 'language', 'Cambridge Dictionary', 'https://dictionary.cambridge.org/', '英语释义、例句、发音与语法参考。', '["词典","英语","发音"]'::jsonb, '[]'::jsonb, 'reference', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('ojad', 'learn', 'language', 'OJAD', 'https://www.gavo.t.u-tokyo.ac.jp/ojad/', '日语单词与句子的重音、韵律查询工具。', '["日语","发音","重音"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 3);
INSERT INTO nav_categories (id, title, description, icon, sort_order) VALUES ('utilities', '实用工具', '文件处理、效率协作、搜索识别与快速转换。', 'i-tabler-tool', 5);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('file-media', 'utilities', '文件与图像', '常见格式转换、压缩和图像处理。', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('convertio', 'utilities', 'file-media', 'Convertio', 'https://convertio.co/zh/', '在线转换文档、图片、音视频等格式。', '["格式转换","文件","在线工具"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('tinypng', 'utilities', 'file-media', 'TinyPNG', 'https://tinypng.com/', '压缩 WebP、PNG 与 JPEG 图片。', '["图片压缩","PNG","WebP"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('remove-bg', 'utilities', 'file-media', 'remove.bg', 'https://www.remove.bg/zh', '自动移除图片背景并导出透明图像。', '["抠图","背景移除","图片"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('bigjpg', 'utilities', 'file-media', 'Bigjpg', 'https://bigjpg.com/', '面向插画和照片的 AI 放大工具。', '["图片放大","AI","降噪"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('productivity', 'utilities', '效率与协作', '组织知识、项目与创作资料。', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('processon', 'utilities', 'productivity', 'ProcessOn', 'https://www.processon.com/', '在线流程图、思维导图与协作白板。', '["流程图","思维导图","协作"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('overleaf', 'utilities', 'productivity', 'Overleaf', 'https://www.overleaf.com/', '在线 LaTeX 编辑、模板与协作平台。', '["LaTeX","写作","协作"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('eagle', 'utilities', 'productivity', 'Eagle', 'https://eagle.cool/', '面向设计师的本地素材收集与管理工具。', '["素材管理","设计","效率"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('geogebra', 'utilities', 'productivity', 'GeoGebra', 'https://www.geogebra.org/calculator', '在线数学绘图、几何与计算工具。', '["数学","绘图","教学"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('discovery', 'utilities', '搜索与识别', '从图片、场景和站点线索反向查找。', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('trace-moe', 'utilities', 'discovery', 'trace.moe', 'https://trace.moe/', '通过动画截图查找作品、集数与时间点。', '["动画","以图搜图","识别"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('tineye', 'utilities', 'discovery', 'TinEye', 'https://tineye.com/', '反向图片搜索与来源追踪工具。', '["图片搜索","来源","识别"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('google-images', 'utilities', 'discovery', 'Google Images', 'https://images.google.com/', '图片搜索与反向图片查找入口。', '["图片搜索","Google","搜索"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('miku-tools', 'utilities', 'discovery', 'MikuTools', 'https://tools.miku.ac/', '集合多种轻量在线工具的中文工具箱。', '["工具箱","在线工具","中文"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 3);
INSERT INTO nav_categories (id, title, description, icon, sort_order) VALUES ('life', '生活与兴趣', '阅读、文化、游戏与日常探索入口。', 'i-tabler-coffee', 6);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('culture', 'life', '阅读与文化', '书影音资料、作品社区与文化内容。', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('bangumi', 'life', 'culture', 'Bangumi 番组计划', 'https://bgm.tv/', '动画、书籍、音乐与游戏条目社区。', '["ACG","条目","社区"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('mikan', 'life', 'culture', '蜜柑计划', 'https://mikanani.me/', '动画信息与订阅聚合入口。', '["动画","订阅","聚合"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('vimeo', 'life', 'culture', 'Vimeo', 'https://vimeo.com/', '高质量视频作品发布与发现平台。', '["视频","作品","影视"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('internet-archive', 'life', 'culture', 'Internet Archive', 'https://archive.org/', '网站、书籍、音视频与软件数字档案馆。', '["档案","图书","历史"]'::jsonb, '[]'::jsonb, 'resource', FALSE, 'published', 3);
INSERT INTO nav_groups (id, category_id, title, description, sort_order) VALUES ('games', 'life', '游戏与互动', '游戏资料、社区与互动地图。', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('gamekee', 'life', 'games', 'GameKee', 'https://www.gamekee.com/', '游戏百科、攻略与玩家共建内容。', '["游戏","Wiki","攻略"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 0);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('opgg', 'life', 'games', 'OP.GG', 'https://www.op.gg/', '英雄联盟等游戏的数据与对局分析。', '["游戏数据","英雄联盟","分析"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 1);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('blueprintue', 'life', 'games', 'BlueprintUE', 'https://blueprintue.com/', '在线分享和查看 Unreal Engine 蓝图。', '["Unreal","蓝图","游戏开发"]'::jsonb, '[]'::jsonb, 'community', FALSE, 'published', 2);
INSERT INTO nav_links (id, category_id, group_id, title, url, description, tags, keywords, kind, featured, status, sort_order) VALUES ('autopiano', 'life', 'games', '自由钢琴', 'https://www.autopiano.cn/', '可用键盘演奏的在线钢琴工具。', '["音乐","钢琴","互动"]'::jsonb, '[]'::jsonb, 'tool', FALSE, 'published', 3);
