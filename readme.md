![](docs/resource/img/BKCN.jpg)
---
[![license](https://img.shields.io/badge/license-mit-brightgreen.svg?style=flat)](https://github.com/Tencent/bk-cmdb/blob/master/LICENSE)
[![Release Version](https://img.shields.io/badge/release-3.1.0-brightgreen.svg)](https://github.com/Tencent/bk-cmdb/releases)
[![Build Status](https://travis-ci.org/Tencent/bk-cmdb.svg?branch=master)](https://travis-ci.org/Tencent/bk-cmdb)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/Tencent/bk-cmdb/pulls)

[(English Documents Available)](readme_en.md)

> **重要提示**: `master` 分支在开发过程中可能处于 *不稳定或者不可用状态* 。
请通过[releases](https://github.com/tencent/bk-cmdb/releases) 而非 `master` 去获取稳定的二进制文件。

蓝鲸配置平台（蓝鲸CMDB）是一个基于运维场景设计的企业配置管理服务。

本次开源的是全新重构的3.0版本。相对于2.0提供了全新自定义模型管理，用户不仅可以方便地实现内置模型属性的拓展，同时也能够根据不同的企业需求随时新增模型和关联关系，把网络、中间件、虚拟资源等纳入到CMDB的管理中。除此之外还增加了更多符合场景需要的新功能：机器数据快照、数据自动发现、变更事件主动推送、更加精细的权限管理、可拓展的业务拓扑等功能。

在技术构建上，新版本核心聚焦于资源，我们把CMDB管理的原子资源分为主机、进程和通用对象三种类型，并构建了对这些资源的原子操作层。在这些原子操作之上，我们构建了更贴近用户操作的场景层，场景层通过对不同资源的组合操作来完成用户的请求。

此次重构使用golang作为开发语言，相比于2.0版本，系统的运行效率得到较大提升。此外采用了微服务架构设计，系统的部署发布可以支持传统方式和容器方式。

开源的版本会与蓝鲸社区版、企业版中内置的蓝鲸配置平台版本保持一致并且同步更新。我们的目标是打造能够兼容不同行业、不同架构的统一配置管理平台，成为业界领先的通用性强、易用性好的免费开源CMDB，欢迎对此感兴趣的同仁能够参与其中。



## Overview
* [架构设计](docs/overview/architecture.md)
* [代码目录](docs/overview/code_framework.md)
* [设计理念](docs/overview/design.md)
* [使用场景](docs/overview/usecase.md)

## Features
* 拓扑化的主机管理：主机基础属性、主机快照数据、主机归属关系管理
* 组织架构管理：可扩展的基于业务的组织架构管理
* 模型管理：既能管理业务、集群、主机等内置模型，也能自定义模型
* 进程管理：基于模块的主机进程管理
* 事件注册与推送：提供基于回调方式的事件注册与推送
* 通用权限管理：灵活的基于用户组的权限管理
* 操作审计：用户操作行为的审计与回溯

如果想了解以上功能的详细说明，请参考[功能说明](http://bk.tencent.com/document/bkprod/000120.html)

## Getting started
* [下载与编译](docs/overview/source_compile.md)
* [安装部署](docs/overview/installation.md)
* [API使用说明](docs/apidoc/readme.md)

## Version plan
* [版本迭代规则](docs/VERSION.md)

## Support
1. 参考bk-cmdb安装文档 [安装文档](docs/overview/installation.md)
2. 阅读 [源码](https://github.com/Tencent/bk-cmdb/tree/master)
3. 阅读 [wiki](https://github.com/Tencent/bk-cmdb/wiki/cmdb-3.0) 或者寻求帮助
4. 了解蓝鲸社区相关信息：[蓝鲸社区版交流1群](https://jq.qq.com/?_wv=1027&k=5zk8F7G)
5. 联系我们，技术交流QQ群：

![qq](docs/resource/img/qq.png)

## Contributing
关于 bk-cmdb 分支管理、issue 以及 pr 规范，请阅读 [bk-cmdb Contributing Guide](docs/CONTRIBUTING.md)。

## FAQ

https://github.com/Tencent/bk-cmdb/wiki/FAQ

## License
bk-cmdb 是基于 MIT 协议， 详细请参考 [LICENSE](LICENSE) 。
