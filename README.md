# 项目初始化脚手架

## 环境要求
    0.go env -w CGO_ENABLED="1" oracle的交叉编译
    1.github.com/cengsin/oracle v1.0.0
    2.gorm.io/driver/mysql v1.1.3
	3.gorm.io/driver/postgres v1.2.0
    4.gorm.io/driver/sqlserver v1.1.1
	5.gorm.io/gorm v1.21.16

### 功能

- [x] 日志模块
- [x] 多通道日志
    - 参考 `config/logging.go`
- [x] 连接`MySQL` `PostgreSQL` `Oracle` 数据库
- [x] 缓存`Redis`连接
- [x] 消息队列`rabbitMQ`连接
- [x] 数据库迁移
    - 操作流程
        1. `db/migrations` 目录下新建迁移文件, 例如 `2020_5_7_17_59_create_users_table.go`
        2. 编辑新建的文件, 实现 `Key(), Up(), Down()`
            - `Key()` 该迁移文件的唯一标识, 推荐使用文件名
            - `Up()` 执行 `migrator` 操作时会调用
            - `Down()` 执行 `migrator rollback` 操作调用
        3. 注册迁移文件到 `pkg/migrate/migrations.go -> MigrateFiles` 变量中
        4. 如果 go run main.go server 启动项目：
            - 执行迁移: `go run main.go migrator`,
            - 执行回滚: `go run main.go migrator rollback`
        5. 如果 `make run` 启动项目：
            - 执行迁移: `./uims migrator`,
            - 执行回滚: `./uims migrator rollback`
    - 手动实现的, 按步数迁移回滚操作后续补充上
    - 注意!
        - 迁移文件中不要使用 `model` 来创建表, 目的是维持迁移文件的版本性不随着 `model` 的变更而变动
- [x] 配置模块
- [ ] 储存模块(储存接口)
- [x] 命令行工具
- [x] 路由结构
    - 路由定义参考 `routes/api/api.go`
- [x] ORM
- [x] Swagger
    - 使用 `go build -o bin/swag vendor/github.com/swaggo/swag/cmd/swag/main.go` 生成 `bin/swag` 可执行文件
    - 使用 `./bin/swag init` 生成新的 swagger 文档
    - 通过 `http://localhost:8080/swagger/index.html` 访问 swagger 页面
- [x] 中间件
    - 参考 `routes/api/api.go` 中使用了 api 中间件组