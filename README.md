闲趣--
一个基于前后端分离架构的电商平台，支持商品浏览、订单管理、实时聊天及管理员后台等完整功能，采用现代化技术栈实现高效交互与数据处理。
项目概述
本项目采用前后端分离架构，后端以 Go 语言为核心提供 RESTful API 及 WebSocket 实时通信服务，前端基于 Vue 3 构建响应式用户界面。平台涵盖用户系统、商品管理、订单处理、购物车及实时聊天等核心模块，同时支持管理员对平台数据的全面管控，满足电商场景的多样化需求。
技术栈
后端
核心框架：Go 1.24 + Gin 1.11.0（高性能 Web 框架）
数据存储：SQLite（轻量嵌入式数据库，通过glebarez/sqlite驱动连接）
ORM 工具：GORM 1.31.1（数据库操作封装）
身份认证：JWT（基于golang-jwt/jwt/v5，支持用户 ID 与角色绑定，令牌有效期 24 小时）
实时通信：WebSocket（基于gorilla/websocket，通过 Hub-Client 模式管理连接，支持点对点消息）
文件处理：支持图片上传，自动创建存储目录并生成动态访问 URL
安全特性：密码加密（bcrypt 算法，加密难度 14）、权限校验
前端
框架与构建工具：Vue 3 + Vite 7.1.11
UI 组件库：Element Plus 2.11.8
状态管理：Pinia 3.0.3
路由管理：Vue Router 4.6.3
网络请求：Axios 1.13.2
样式方案：Sass（支持变量与嵌套语法，实现玻璃态、渐变等视觉效果）
开发工具：推荐 VS Code + Vue 官方插件（禁用 Vetur），浏览器建议使用 Chromium 内核（搭配 Vue DevTools）
其他组件：表情选择器（vue3-emoji-picker）、自定义图标组件
功能特点
用户模块
认证功能：注册 / 登录（JWT 令牌生成与验证，包含用户 ID 和角色信息）
信息管理：支持修改昵称、头像、手机号等个人资料，密码修改需验证旧密码（bcrypt 加密存储）
权限控制：基于角色的访问限制（如商品修改需验证所有权）
商品模块
商品操作：发布 / 编辑商品（支持封面图上传，限制 JPG/PNG 格式，最大 10MB）、分类浏览、详情查询
状态管理：商品发布后需管理员审核，支持修改商品信息（仅商品所有者可操作）
视觉体验：商品卡片 hover 效果、发布页玻璃态设计、表单动画过渡
订单与购物车
购物车管理：添加、删除、查询商品
订单处理：创建订单（支持批量操作）、更新支付状态，管理员可查询所有订单
订单展示：包含商品缩略图、名称、价格等信息，支持分页浏览
实时聊天
通信机制：基于 WebSocket 实现实时消息收发，通过 Hub 管理客户端连接，支持用户索引快速查找
消息功能：支持文本 / 图片消息（类型区分），消息自动存入 SQLite 数据库（包含发送者、接收者、时间等信息）
界面设计：气泡式聊天界面，区分自己与对方消息，支持图片预览
管理员功能
后台管理：管理员登录入口，专用管理界面（如订单管理、用户管理）
数据管控：商品审核、用户状态更新、订单查询与统计
操作体验：搜索筛选、状态切换按钮、分页控制等便捷操作
快速开始
环境要求
后端：Go 1.24+
前端：Node.js 20.19.0+ 或 22.12.0+
后端启动步骤
进入后端目录
bash
运行
cd vision-rainbow-zijin/backend
初始化依赖
bash
运行
go mod tidy
启动服务（默认端口 8081，自动创建uploads目录用于文件存储）
bash
运行
go run main.go
前端启动步骤
进入前端目录
bash
运行
cd vision-rainbow-zijin/frontend
安装依赖
bash
运行
npm install
# 或使用yarn
yarn install
启动开发服务器（默认地址：http://localhost:5173）
bash
运行
npm run dev
# 或
yarn dev
构建生产版本（用于部署）
bash
运行
npm run build
项目结构
plaintext
vision-rainbow-zijin/
├── backend/                # 后端Go代码
│   ├── main.go             # 入口文件（路由配置、服务启动）
│   ├── go.mod              # 依赖管理
│   ├── xianqu.db           # SQLite数据库文件
│   ├── uploads/            # 上传文件存储目录（自动创建）
│   ├── internal/
│   │   ├── controllers/    # 控制器（用户、商品、订单、文件等）
│   │   ├── services/       # 业务逻辑（如密码修改）
│   │   ├── models/         # 数据模型（用户、商品、消息等）
│   │   ├── utils/          # 工具函数（JWT、加密等）
│   │   └── middleware/     # 中间件（认证等）
│   └── pkg/ws/             # WebSocket相关（Hub、Client）
└── frontend/               # 前端Vue代码
    ├── src/
    │   ├── main.js         # 入口文件（Vue初始化、插件配置）
    │   ├── router/         # 路由配置
    │   ├── components/     # 通用组件（图标等）
    │   ├── views/          # 页面组件（发布、聊天、管理员页面等）
    │   ├── assets/         # 静态资源（样式、图片等）
    │   └── stores/         # Pinia状态管理
    ├── package.json        # 依赖管理
    └── README.md           # 前端开发说明
核心接口说明
用户相关：
POST /api/register：用户注册
POST /api/login：用户登录（返回 JWT 令牌）
PUT /api/user/profile：更新个人资料
PUT /api/user/password：修改密码
商品相关：
GET /api/products：获取商品列表
GET /api/products/:id：获取商品详情
POST /api/products：发布商品
PUT /api/products/:id：更新商品（需验证所有权）
POST /api/upload：上传图片（返回图片 URL）
订单相关：
GET /api/orders：获取订单列表
POST /api/orders：创建订单
PUT /api/orders/:id/status：更新订单状态
实时聊天：
GET /api/ws：WebSocket 连接（用于实时消息）
GET /api/chat/history：获取聊天历史记录
管理员接口：
GET /api/admin/orders：管理员查询所有订单
PUT /api/admin/products/:id/audit：审核商品状态
