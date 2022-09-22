蓝鲸GSE-Data部署安装说明
========================

# 安装包结构

```shell
├── etc
│   └── gse_data.conf.template
├── gse_data
├── install
│   ├── generate.sh
│   ├── gse_data.env
│   └── install.sh
├── tools
└── README.md
```

目录说明，

- `etc`: 配置文件模板, 该文件会基于`gse_data.env`变量集进行渲染，生成最终运行配置;
- `gse_data`: 编译后的二进制;
- `install`: 安装工具, 包含配置生成脚本`generate.sh`，配置变量`gse_data.env`, 安装脚本`install.sh`等;
- `tools`: 日常运营维护工具, 用于问题排查，系统运行检查等;
- `README.md`: 安装说明文档;

# 安装教程

## 安装Zookeeper

[参见官方安装教程] <http://zookeeper.apache.org/doc/current/index.html>

## 安装运行GSE-Data

### 编辑配置变量

修改`install`目录下`gse_data.env`中的配置变量(详细参见注释说明)

### 生成配置

执行`install`目录下`generate.sh`脚本，根据`gse_data.env`渲染最终需要的`gse_data.conf`,

```shell
sh generate.sh -e ./gse_data.env -t ../etc/gse_data.conf.template > ../etc/gse_data.conf
```

### 安装GSE-Data

执行`install`目录下`install.sh`脚本，将在`gse_data.env`中`{HOME_DIR}`下安装块二进制文件和渲染后的配置文件

### 拉起进程

执行`./gse_data run`启动进程, 也可以根据自身运维环境配置systemd拉起模块进程。

至此已完成部署工作。

### 健康检查

安装完成之后，可利用健康检查接口确认服务是否正确部署:

```shell
curl -vv http://127.0.0.1:59402/healthz
```

### 内部信息检查

检查内部路由管理信息，以及该节点的其他详细运行时信息:

```shell
curl -vv http://127.0.0.1:59402/stack
```

# 其他

## 安装环境Bash

系统内提供的shell脚本工具均基于bash，若使用dash脚本的系统如Ubuntu 需修改为bash后使用:

- 1.执行`ls -l /bin/sh`命令，若得到结果`/bin/sh -> dash`，则说明shell的解释器为`dash`;
- 2.执行`dpkg-reconfigure dash`命令，然后选择no;
- 3.再次执行`ls -l /bin/sh`命令，若得到结果`/bin/sh -> bash`，则说明成功更改shell的解释器为`bash`;
