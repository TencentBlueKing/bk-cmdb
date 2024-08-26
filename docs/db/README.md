# 数据库表结构设计

*说明：蓝鲸配置平台（蓝鲸CMDB）项目底层使用MongoDB进行数据存储，MongoDB是一个面向文档的数据库，因此它没有传统的表结构，而是使用集合（Collection）和文档（Document）来组织数据，本目录下文档中涉及的表结构和字段解释仅供参考，请以实际环境中的结构为准。*

## 表结构文档分类及作用

| 作用             | 分类文档                                              |
|----------------|---------------------------------------------------|
| 内置模型相关表        | [built-in_model](built-in_model.md)               |
| 业务下相关资源表       | [business](business.md)                           |
| 云资源相关功能表       | [cloud_resource](cloud_resource.md)               |
| 容器数据纳管功能相关表    | [container_data_manage](container_data_manage.md) |
| 字段组合模板功能相关表    | [field_template](field_template.md)               |
| 主机属性自动应用功能相关表  | [host_apply_rule](host_apply_rule.md)             |
| 实例相关表          | [instance](instance.md)                           |
| 主线模型相关表        | [mainline_model](mainline_model.md)               |
| 模型相关表          | [model](model.md)                                 |
| 运营分析功能相关表      | [operational_analysis](operational_analysis.md)   |
| 管控区域相关表        | [plat](plat.md)                                   |
| 服务模板功能相关表      | [service_template](service_template.md)           |
| 集群模板功能相关表      | [set_template](set_template.md)                   |
| 异步任务相关表        | [task](task.md)                                   |
| 资源目录用户自定义配置相关表 | [user_custom](user_custom.md)                     |
| 资源变更事件相关表      | [watch](watch.md)                                 |
| 不归属任何一种分类的表    | [other](other.md)                                 |

