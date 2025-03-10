![](docs/resource/img/bk-cmdb.png)
---
[![license](https://img.shields.io/badge/license-mit-brightgreen.svg?style=flat)](https://github.com/TencentBlueKing/bk-cmdb/blob/master/LICENSE.txt)
[![Release Version](https://img.shields.io/badge/release-3.2.19-brightgreen.svg)](https://github.com/TencentBlueKing/bk-cmdb/releases)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/TencentBlueKing/bk-cmdb/pulls)
[![BK Pipelines Status](https://api.bkdevops.qq.com/process/api/external/pipelines/projects/cc/p-c02db56ac633447eb2e740b3fd0b6d2b/badge?X-DEVOPS-PROJECT-ID=cc)](http://api.bkdevops.qq.com/process/api-html/user/builds/projects/cc/pipelines/p-c02db56ac633447eb2e740b3fd0b6d2b/latestFinished?X-DEVOPS-PROJECT-ID=cc)


[(English Documents Available)](readme_en.md)

> **重要提示**: `master` 分支在开发过程中可能处于 *不稳定或者不可用状态* 。
请通过[releases](https://github.com/TencentBlueKing/bk-cmdb/releases) 而非 `master` 去获取稳定的二进制文件。

蓝鲸配置平台（蓝鲸CMDB）是一个面向资产及应用的企业级配置管理平台。

蓝鲸配置平台提供了全新自定义模型管理，用户不仅可以方便地实现内置模型属性的拓展，同时也能够根据不同的企业需求随时新增模型和关联关系，把网络、中间件、虚拟资源等纳入到CMDB的管理中。除此之外还增加了更多符合场景需要的新功能：机器数据快照、数据自动发现、变更事件主动推送、更加精细的权限管理、可拓展的业务拓扑等功能。

在技术构建上，架构的核心聚焦于资源，我们把CMDB管理的原子资源分为主机、进程和通用对象三种类型，并构建了对这些资源的原子操作层。在这些原子操作之上，我们构建了更贴近用户操作的场景层，场景层通过对不同资源的组合操作来完成用户的请求。

## Overview
* [设计理念](docs/overview/design.md)
* [架构设计](docs/overview/architecture.md)
* [代码目录](docs/overview/code_framework.md)
* [数据库表结构设计](docs/db/README.md)

![front-page](docs/resource/img/frontpage.png)

## Features
* 拓扑化的主机管理：主机基础属性、主机快照数据、主机归属关系管理
* 组织架构管理：可扩展的基于业务的组织架构管理
* 模型管理：既能管理业务、集群、主机等内置模型，也能自定义模型
* 进程管理：基于模块的主机进程管理
* 事件注册与推送：提供基于回调方式的事件注册与推送
* 通用权限管理：灵活的基于用户组的权限管理
* 操作审计：用户操作行为的审计与回溯

如果想了解以上功能的详细说明，请参考[功能说明](https://bk.tencent.com/docs/markdown/CMDB/UserGuide/Introduce/Overview.md)

## Experience
* [在线体验蓝鲸CMDB](https://cmdb-exp.bktencent.com/start) 🔥 **用户名密码：admin:admin**
* [极速体验容器化部署蓝鲸CMDB](docs/wiki/container-support.md)

## Getting started
* [下载与编译](docs/overview/source_compile.md)
* [安装部署](docs/overview/installation.md)
* [版本升级说明](docs/wiki/db_upgrade.md)
* [API使用说明见这里](docs/apidoc/)
* [使用CMDB开源版替换社区版](docs/overview/upgrade-from-ce.md)
* [使用Helm部署CMDB到K8S环境](docs/support-file/helm/backend/README.md)

## Roadmap
* [版本日志](docs/support-file/changelog/release.md)

## Support
- [wiki](https://github.com/TencentBlueKing/bk-cmdb/wiki)
- [白皮书](https://docs.bk.tencent.com/cmdb/)
- [蓝鲸论坛](https://bk.tencent.com/s-mart/community)
- [蓝鲸 DevOps 在线视频教程](https://bk.tencent.com/s-mart/video/)
- [蓝鲸社区版交流1群](https://jq.qq.com/?_wv=1027&k=5zk8F7G)
- 技术交流QQ群(305496802), 扫码入群戳[这里](docs/resource/img/qq.png)

## BlueKing Community
- [BK-CI](https://github.com/Tencent/bk-ci)：蓝鲸持续集成平台是一个开源的持续集成和持续交付系统，可以轻松将你的研发流程呈现到你面前。
- [BK-BCS](https://github.com/Tencent/bk-bcs)：蓝鲸容器管理平台是以容器技术为基础，为微服务业务提供编排管理的基础服务平台。
- [BK-PaaS](https://github.com/Tencent/bk-PaaS)：蓝鲸PaaS平台是一个开放式的开发平台，让开发者可以方便快捷地创建、开发、部署和管理SaaS应用。
- [BK-SOPS](https://github.com/Tencent/bk-sops)：标准运维（SOPS）是通过可视化的图形界面进行任务流程编排和执行的系统，是蓝鲸体系中一款轻量级的调度编排类SaaS产品。


## Contributing
如果你有好的意见或建议，欢迎给我们提 Issues 或 Pull Requests，为蓝鲸开源社区贡献力量。关于bk-cmdb分支管理、Issue 以及 PR 规范，
请阅读 [Contributing Guide](docs/CONTRIBUTING.md)。

[腾讯开源激励计划](https://opensource.tencent.com/contribution) 鼓励开发者的参与和贡献，期待你的加入。

## Partners

![jiawei](docs/resource/img/jiawei-logo.png)　　![youole](docs/resource/img/youole-logo.png)　　![honops](docs/resource/img/honops-logo.png)

## License
项目基于 MIT 协议，详细请参考 [LICENSE](LICENSE.txt)。

我们承诺未来不会更改适用于交付给任何人的当前项目版本的开源许可证（MIT 协议）。

