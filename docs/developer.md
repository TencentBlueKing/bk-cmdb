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
│           ├── router.go 路由
│           ├── service.go 服务定义
│           └── user.go 业务逻辑
├── docs 各类文档
│   └── developer.md
├── go.mod
├── go.sum
├── LICENSE.txt
├── Makefile
├── pkg 公共包
│   ├── healthz 自身服务healthz接口
│   ├── logger slog初始化
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
