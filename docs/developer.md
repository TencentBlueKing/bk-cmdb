# CMDB 开发快速指南

## 代码结构

```text
bk-cmdb
├── bin 工具目录
│   └── lint
├── build 项目构建产物
│   └── cmdb_apiserver
├── cmd 服务入口&自身服务业务逻辑
│   └── api-server
│   │   ├── api_server.go main入口
│   │   ├── app 命令行主逻辑
│   │   ├── options 命令行参数
│   │   └── service 服务业务逻辑
│   └── auth-server 鉴权服务
│       ├── auth_server.go main入口
│       ├── app 命令行主逻辑
│       ├── options 命令行参数
│       └── service 服务业务逻辑
├── docs 各类文档
│   └── developer.md
├── go.mod
├── go.sum
├── LICENSE.txt
├── Makefile
├── pkg 公共包
│   ├── config-center 配置中心，封装了配置文件读取、配置文件变更监听、配置注册与发现逻辑
│       ├── config 系统使用的配置类型定义
│   │   └── etcd 使用etcd实现配置注册与发现
│   ├── etcd 封装了etcd相关的配置和操作
│   ├── healthz 自身服务healthz接口
│   ├── logger 自定义Contextual&Structured Logger
│   ├── proto protobuf文件和生成的grpc代码
│   ├── rest http服务框架
│   ├── i18n 国际化处理
│   ├── dal 数据访问层
│   │   ├── dao 数据访问对象
│   │   ├── gen 静态结构生成工具
│   │   ├── gendo 生成的静态模型数据
│   │   ├── orm 对象关系映射能力封装
│   │   ├── table 模型结构定义
│   │   └── types 类型定义
│   ├── runtime
│       ├── cli 命令行入口封装
│   │   └── server 创建和运行通用服务的逻辑封装
│   ├── service-discovery 服务注册与发现逻辑封装
│   │   └── etcd 使用etcd实现服务注册与发现
│   ├── validator struct参数校验
│   └── version 版本号
├── resource 资源文件
│   └── translations 对应国际化翻译文件
│       └── en/zh等 语言类型定义
│           ├── error-code 错误码翻译
│           └── text 内置及其他文本信息翻译
├── readme_en.md
└── readme.md
```

## 编辑器配置

golang 1.25开启jsonv2需要添加环境变量, 在`.vscode/settings.json`添加如下配置, goroot请跟进实际路径修改

```json
{
    "go.goroot": "/opt/go/sdk/go1",
    "go.toolsEnvVars": {
        "GOEXPERIMENT": "jsonv2"
    }
}
```

## 代码检查&编译

- make lint， 检测golang代码规范
- make test, 执行单元测试
- make build, 编译命令行

运行apiserver服务

```bash
./build/bk-cmdb-apiserver
```

## 开发API服务

cmdb api服务默认是unary函数，是以请求 -> 响应模式

实例如下

```go
// HelloReq Hello请求参数
type HelloReq struct {
	Name string `json:"name", validator:"required"`
}

// Validate Hello参数校验
func (h *HelloReq) Validate() error {
	if req.Name != "hello" {
		return fmt.Errorf("name not equal hello")
	}

    return nil
}

// HelloResp Hello响应
type HelloResp struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type service struct{}

// Hello 接口实现
func (s *service) Hello(ctx context.Context, req *HelloReq) (*HelloResp, error) {
	// 业务逻辑
	name := req.Name + " world"

	// 响应
	resp := &HelloResp{Name: name, Age: 18}
	return resp, nil
}
```

## req decode 规则
rest框架可以对request请求自动解析到业务自定义结构体，定义使用`req`做tag, 使用`in`表示来源，规则如下
- 语法习惯同社区, 可参考https://pkg.go.dev/encoding/json/v2#example-package-FormatFlags
- 不能和json同时使用
- `in`目前支持query/path，不合法的会提示错误
- time类型支持`format`参数, 格式如`format:2006-01-02`,表示按日期解析

示例如下
```go
// UserInfoReq 个人信息Req
type UserInfoReq struct {
	Username string     `json:"name" req:"-,in:query"`
	Age      int        `req:"age,in:query" validate:"required"`
	Games    *[]*string `json:"games" req:"-"`
	BirthDay time.Time  `req:"birthday,in:query,format:2006-01-02"`
	Ko       []byte     `json:"-" req:"ko,in:query"`
}
```
