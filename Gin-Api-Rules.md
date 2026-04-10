---
trigger: always_on
---
# Gin API Base 项目开发规则

## 项目概述

这是一个基于 Gin 框架的现代化 Go API 基础模板，采用**四层架构**设计，技术栈精简、结构清晰。

## 目录结构与分层职责

```
.
├── config/              # 配置管理（Viper + JSON）
├── handler/             # 处理器层：接收请求、参数校验、调用 service、返回响应
├── service/             # 业务逻辑层：核心业务、数据组装、请求/响应结构体定义
├── model/               # 数据访问层：纯 CRUD 操作函数（不含模型定义）
├── middleware/           # 中间件：CORS、JWT 认证、请求日志
├── pkg/                 # 公共工具包
│   ├── resp/            #   统一 JSON 响应
│   └── jwtutil/         #   JWT Token 工具
├── router/              # 路由注册（统一入口）
├── server/              # 基础设施
│   ├── mysql/           #   MySQL 连接 + 模型定义 + 自动建表（package mysqlServer）
│   │   ├── mysql.go     #     连接初始化、连接池、RegisterModel、AutoMigrate、GetDB()
│   │   └── model.go     #     所有 GORM 模型定义 + init() 注册
│   └── redis/           #   Redis 连接管理（package redisServer）
│       └── redis.go     #     连接初始化、GetRedis()
└── main.go              # 入口：配置 → 数据库(含建表) → 路由 → 启动
```

## 核心开发规范

### 1. 分层调用规则（严格遵守）

```
Handler → Service → Model → mysqlServer.GetDB()
```

- **Handler 层** (`handler/`)：只做参数绑定、调用 service、返回响应。禁止直接操作数据库。
- **Service 层** (`service/`)：编写业务逻辑，定义请求/响应结构体（XxxReq / XxxResp），调用 model 层。构造模型对象时使用 `mysqlServer.Xxx`。
- **Model 层** (`model/`)：纯数据访问层，只编写 CRUD 函数。**不包含模型定义**，模型类型通过 `mysqlServer.Xxx` 引用。
- **Server 层** (`server/mysql/`、`server/redis/`)：负责数据库连接、模型定义、自动建表，不含业务逻辑。

### 2. 新增业务模块步骤

当新增一个业务模块（如 "文章管理"），按以下顺序操作：

1. `server/mysql/model.go` — 定义 Article 结构体 + TableName()，在 `init()` 中追加 `&Article{}` 注册
2. `model/article.go` — 编写 CRUD 函数（使用 `mysqlServer.Article` 类型）
3. `service/article.go` — 定义请求/响应结构体 + 业务逻辑函数
4. `handler/article.go` — 编写 HTTP 处理函数
5. `router/router.go` — 注册新路由

### 3. 命名约定

| 类别 | 规则 | 示例 |
|------|------|------|
| 文件名 | 小写单数 | `user.go`, `article.go` |
| 包名 | 小写单词或驼峰 | `handler`, `service`, `model`, `mysqlServer`, `redisServer` |
| Handler 函数 | 动词 + 名词 | `GetUser`, `CreateArticle`, `ListUsers` |
| Service 函数 | 同 Handler | `Register`, `Login`, `GetUser` |
| Model 函数 | 动词 + 模型名 | `CreateUser`, `GetUserByID`, `ListUsers` |
| 请求结构体 | Xxx + Req | `RegisterReq`, `UpdateUserReq` |
| 响应结构体 | Xxx + Resp / Info | `LoginResp`, `UserInfo` |
| GORM 模型 | 单数名词 | `User`, `Article` |
| 表名 | 小写单数 | `user`, `article` |

### 4. 响应格式（统一使用 `pkg/resp`）

所有接口必须使用 `pkg/resp` 包返回响应，禁止直接 `c.JSON()`。只有两个方法：

```go
resp.OK(c, data)             // 成功（data 可以是任意类型：结构体、字符串、nil）
resp.Fail(c, code, msg)      // 失败（code < 1000 作为 HTTP 状态码，code >= 1000 为业务错误码 HTTP 200）
```

常见用法：
```go
resp.OK(c, user)                           // 返回数据
resp.OK(c, "操作成功")                       // 返回消息
resp.OK(c, resp.Page{...})                 // 返回分页
resp.Fail(c, 400, "参数错误")               // 参数错误
resp.Fail(c, 401, "未登录")                 // 未授权
resp.Fail(c, 404, "不存在")                 // 资源不存在
resp.Fail(c, 500, "服务器错误")             // 服务器错误
resp.Fail(c, 1001, "用户名已存在")           // 业务错误（HTTP 200）
```

响应 JSON 格式：`{"code": 0, "msg": "success", "data": {...}}`

### 5. 路由规范

- 所有路由在 `router/router.go` 中统一注册
- API 前缀统一为 `/api`
- 公开接口放在 `api` 组下，需要认证的接口放在 `auth` 子组下（使用 `middleware.JWTAuth()`）
- RESTful 风格：GET 查询、POST 创建、PUT 更新、DELETE 删除

### 6. 中间件规范

- 全局中间件在 `router.Setup()` 中注册
- 中间件文件放在 `middleware/` 目录，每个中间件一个文件
- 中间件函数签名：`func XxxMiddleware() gin.HandlerFunc`

### 7. 模型定义规范（`server/mysql/model.go`）

所有 GORM 模型统一定义在 `server/mysql/model.go` 中，通过 `init()` 自动注册，启动时自动建表：

```go
package mysqlServer

// 模型定义
type Article struct {
    gorm.Model
    Title   string `gorm:"column:title;type:varchar(200);not null" json:"title"`
    Content string `gorm:"column:content;type:text"                json:"content"`
}

func (Article) TableName() string { return "article" }

// 在 init() 中追加注册
func init() {
    RegisterModel(
        &User{},
        &Article{}, // 新增
    )
}
```

- 使用 `gorm.Model` 内嵌（自带 ID/CreatedAt/UpdatedAt/DeletedAt）
- 必须显式定义 `TableName()` 方法
- 必须写完整的 gorm tag（column、type、约束）
- 密码等敏感字段用 `json:"-"` 隐藏
- 新增模型在 `init()` 中追加 `RegisterModel()`，无需改其他文件

### 8. Model 数据访问层规范（`model/`）

`model/` 只包含纯 CRUD 函数，不含模型定义。通过 `mysqlServer.Xxx` 引用模型类型：

```go
package model

import mysqlServer "gin-api/server/mysql"

func CreateArticle(article *mysqlServer.Article) error {
    return mysqlServer.GetDB().Create(article).Error
}

func GetArticleByID(id uint) (*mysqlServer.Article, error) {
    var article mysqlServer.Article
    err := mysqlServer.GetDB().First(&article, id).Error
    return &article, err
}
```

### 9. Handler 层标准模板

```go
// XxxHandler 功能说明
// HTTP方法 /api/path
func XxxHandler(c *gin.Context) {
    // 1. 参数绑定与校验
    var req service.XxxReq
    if err := c.ShouldBindJSON(&req); err != nil {
        resp.Fail(c, 400, "参数错误: "+err.Error())
        return
    }

    // 2. 调用 service
    data, err := service.XxxFunc(&req)
    if err != nil {
        resp.Fail(c, 1001, err.Error())
        return
    }

    // 3. 返回响应
    resp.OK(c, data)
}
```

### 10. 错误处理规范

- Handler 层：参数错误用 `resp.Fail(c, 400, ...)`，业务错误用 `resp.Fail(c, 1xxx, ...)`
- Service 层：返回 `error`，使用 `errors.New()` 或 `fmt.Errorf("xxx: %w", err)`
- Model 层：直接返回 GORM 的 error
- 检查记录不存在：`errors.Is(err, gorm.ErrRecordNotFound)`

### 11. 包引用规范

由于 `server/mysql` 目录的包名为 `mysqlServer`，`server/redis` 目录的包名为 `redisServer`，import 时需使用别名：

```go
import (
    mysqlServer "gin-api/server/mysql"
    redisServer "gin-api/server/redis"
)
```

## 技术栈

| 组件 | 技术 | 用途 |
|------|------|------|
| Web 框架 | Gin v1.10 | HTTP 路由与中间件 |
| ORM | GORM v1.31 | MySQL 数据库操作 |
| 缓存 | go-redis v9 | Redis 缓存 |
| 认证 | golang-jwt/jwt v5 | JWT Token |
| 配置 | Viper v1.19 | JSON 配置管理 |
| 日志 | Go 标准库 slog | 结构化请求日志 |
| 链路追踪 | google/uuid | 请求 traceID |

## 配置管理

- 配置文件：`config/config.json`（开发）、`config/config.prod.json`（生产）
- 通过环境变量 `GO_ENV=prod` 切换环境
- 使用 `viper.GetString("key.subkey")` 读取配置
- 生产配置文件已在 `.gitignore` 中忽略

## 编码风格

- 注释使用中文，简洁明了
- import 分组：标准库 → 项目内部包 → 第三方包
- 函数注释格式：`// FuncName 功能描述`
- 每个文件只处理一个业务实体
