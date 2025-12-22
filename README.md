# Vision Rainbow Zijin

## 🌈 项目概述
**项目口号**：一站式现代化Web应用平台，融合高效管理与实时交互体验

**详细描述**：
Vision Rainbow Zijin 是一个基于前后端分离架构的综合性Web应用，旨在为用户提供流畅、直观的在线体验。该项目整合了用户管理、商品展示、订单处理和实时通讯等核心功能，满足现代Web应用的多元化需求。

通过采用Vue 3与Go语言的高效组合，项目实现了前端界面的精美呈现与后端服务的稳定运行，同时保证了系统的可扩展性和可维护性，适用于各类商业展示与交易场景。

**项目类型**：Web应用

## ✨ 核心特性与功能
1. **完整用户系统**：支持用户注册、登录、个人信息管理及权限控制，确保账户安全与个性化体验
2. **商品与订单管理**：提供商品浏览、分类筛选、购物车操作及订单处理全流程，满足电商场景需求
3. **实时聊天功能**：集成WebSocket实现用户间实时通讯，支持消息记录与表情互动
4. **响应式设计**：适配各种设备屏幕尺寸，提供一致的用户体验
5. **高效后台管理**：为管理员提供用户管理、商品审核、数据统计等功能模块

## 🛠️ 技术栈与依赖
### 前端技术
- 核心框架：Vue 3、Vite
- UI组件库：Element Plus
- 状态管理：Pinia
- 路由管理：Vue Router
- HTTP客户端：Axios
- 样式处理：Sass
- 其他依赖：bootstrap、vue3-emoji-picker、@element-plus/icons-vue

### 后端技术
- 编程语言：Go
- Web框架：Gin
- 实时通信：WebSocket
- 数据库：支持主流关系型数据库（配置文件可指定）

### 环境要求
- 前端：Node.js (v20.19.0 或 v22.12.0+)
- 后端：Go 1.18+
- 浏览器：推荐Chrome/Edge/Brave等Chromium内核浏览器

## 🚀 快速开始
### 1. 克隆项目
```bash
git clone <项目仓库地址>
cd vision-rainbow-zijin
2. 前端安装与运行
   1. 克隆项目
```bash
 # 进入前端目录
cd frontend

# 安装依赖
npm install

# 开发环境运行（热重载）
npm run dev

# 构建生产版本
npm run build
```
3. 后端安装与运行
```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 运行服务
go run main.go
```
服务启动后，可访问 http://localhost:8081 查看应用（后端默认端口为 8081）
基础使用示例
```javascript
 // 前端API调用示例（获取商品列表）
import axios from 'axios';

const fetchProducts = async () => {
  try {
    const response = await axios.get('/api/products');
    return response.data;
  } catch (error) {
    console.error('获取商品列表失败:', error);
    return [];
  }
};
```
### 配置说明
- 前端配置：修改 `frontend/.env` 文件配置 API 地址等环境变量
- 后端配置：修改 `backend/config` 目录下的相关文件配置数据库连接、端口等信息

### 📚 进阶指南
#### 详细使用方法
##### 用户模块
- 注册：通过首页注册表单提交用户信息
- 登录：使用账号密码或第三方登录方式
- 个人中心：修改个人信息、查看订单历史、管理收货地址

##### 商品与订单模块
- 浏览商品：通过分类导航或搜索功能查找商品
- 购物车：添加商品、修改数量、删除商品
- 下单流程：从购物车选择商品，填写收货信息，完成支付

##### 实时聊天
- 查看联系人列表，选择对话对象
- 发送文本消息与表情
- 查看聊天历史记录

#### 常见用例
**用例 1：普通用户购物流程**
1. 注册并登录系统
2. 浏览商品并将心仪商品加入购物车
3. 在购物车确认商品及数量，提交订单
4. 选择支付方式完成支付
5. 在个人中心查看订单状态

**用例 2：管理员商品管理**
1. 使用管理员账号登录
2. 进入后台管理界面
3. 上传新商品信息或编辑现有商品
4. 审核用户提交的订单
5. 查看销售数据统计

### 📂 项目结构
```plaintext
 vision-rainbow-zijin/
├── frontend/                 # 前端项目目录
│   ├── src/
│   │   ├── components/       # 通用组件
│   │   ├── router/           # 路由配置
│   │   ├── store/            # Pinia状态管理
│   │   ├── utils/            # 工具函数
│   │   ├── App.vue           # 根组件
│   │   └── main.js           # 入口文件
│   ├── public/               # 静态资源
│   └── package.json          # 依赖配置
│
├── backend/                  # 后端项目目录
│   ├── internal/
│   │   ├── controllers/      # 控制器（处理请求）
│   │   ├── middleware/       # 中间件（身份验证等）
│   │   └── models/           # 数据模型
│   ├── pkg/
│   │   └── ws/               # WebSocket相关实现
│   ├── config/               # 配置文件
│   └── main.go               # 程序入口
│
└── README.md                 # 项目说明文档
```
