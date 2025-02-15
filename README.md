# gRPC-todoList
***
grpc + grpc-gateway + GORM + MySQL 实现的简易备忘录

***
# 项目结构
1. 架构图
![image](https://github.com/user-attachments/assets/3b483955-f762-419f-9cb3-a2fc9bc6a72a)
2. grpc-todolist 项目目录
```
. grpc-todolist
├── app            // 各个微服务
│    ├── gateway   // grpc 网关
│    ├── task      // 任务模块
│    └── user      // 用户模块
└── idl
     ├── google    // grpc-gateway idl
     └── todolist  // protoc 接口定义
```
3. gateway 模块
```
. gateway
├── api 
│    └── main.go  // grpc gateway 启动入口
├── go.mod
├── go.sum
└── mws
    ├── cors.go   // 跨域中间件 
    └── jwt.go    // jwt 中间件
```
4. task 模块
```
. task
├── biz
│    ├── repository  // 数据访问层
│    └── service     // 业务逻辑
│── rpc              // rpc 调用
├── config
│    ├── config.go   // 配置读取入口
│    └── test        // test 环境下的配置文件
│       └── config.yaml
├── go.mod
├── go.sum
├── ioc
│    ├── wire.go     
│    └── wire_gen.go // 依赖注入
├── main.go          // 入口文件
├── .env             // 项目启动环境变量
└── rpc_gen 
     └── task        // grpc 生成代码
```
3. user 模块
```
. user
├── biz
│    ├── repository  // 数据访问层
│    └── service     // 业务逻辑
│── rpc              // rpc 调用
├── config
│    ├── config.go   // 配置读取入口
│    └── test        // test 环境下的配置文件
│       └── config.yaml
├── go.mod
├── go.sum
├── ioc
│    ├── wire.go     
│    └── wire_gen.go // 依赖注入
├── main.go          // 入口文件
├── .env             // 项目启动环境变量
└── rpc_gen 
     └── user        // grpc 生成代码
```
***
# 项目前端
项目前端仓库在 https://github.com/crazyfrankie/grpc-todolist-web

# 项目配置文件
`.env` : 采用以下格式，且必须放置在 task & user 包下
```
MYSQL_USER=your_user
MYSQL_PASSWORD=your_password
MYSQL_HOST=your_host
MYSQL_PORT=your_port
MYSQL_DB=your_db
```
`config.yaml` : 大致采用以下格式，根据不同模块 `config.go` 的需求进行更改
```
server:
  addr: "your_addr"

mysql:
  dsn: "%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"

etcd:
  addr: "your_addr"
```
# 项目启动
项目环境: etcd + MySQL
## Docker
1. 修改 `docker-compose.yaml`
2. 执行 `docker-compose up -d`

然后执行 
```
make run-user 
make run-task
make run-api
```
## 本机环境
修改配置文件后执行
```
make run-user 
make run-task
make run-api
```
