## 目录结构

proto 目录中存储了 cmdb 的 protobuf 文件和生成的 grpc 代码，目录结构如下：

```text
proto
├── auth-server 包含鉴权服务的proto文件和生成的代码
├── Makefile
└── readme.md
```

## 代码生成

1. 安装 32.1 版本的 protoc
   1. 在 https://github.com/protocolbuffers/protobuf/releases 中根据版本和设备类型下载对应的 protoc 安装包
   2. 解压安装包
   3. 将 bin/protoc 复制到 $PATH 中并将 include 目录下的内容复制到系统的 include 目录中
2. 在本目录中使用 make 命令或者在根目录中使用 make proto 命令进行代码生成，该命令会自动清理存量的生成代码后安装依赖并生成新的 grpc 代码
3. 依赖的 golang 代码生成工具的版本如下，如需调整需要同步修改 Makefile 中的版本号，并执行 make clean-tools 命令清理依赖后再重新生成代码：
   - protoc-gen-go：1.36.10
   - protoc-gen-go-grpc：1.5.1
   - protoc-gen-grpc-gateway：2.27.3
