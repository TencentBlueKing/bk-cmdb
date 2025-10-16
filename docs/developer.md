# CMDB 开发快速指南

## 代码结构

```text
bk-cmdb
├── bin 工具目录
│   └── lint
├── build 项目构建产物
│   └── cmdb_apiserver
├── cmd 服务入口&自身服务业务逻辑
│   └── apiserver
│       ├── apiserver.go main入口
│       ├── app 命令行主逻辑
│       ├── etc 命令行配置模版
│       ├── options 命令行参数
│       └── service 服务业务逻辑
├── docs 各类文档
│   └── developer.md
├── go.mod
├── go.sum
├── LICENSE.txt
├── Makefile
├── pkg 公共包
│   ├── healthz 自身服务healthz接口
│   ├── logger 自定义Contextual&Structured Logger
│   ├── rest http服务框架
│   ├── dal 数据访问层
│   │   └── dao 数据访问对象
│   ├── runtime
│   │   └── cli 命令行入口封装
│   ├── validator struct参数校验
│   └── version 版本号
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
