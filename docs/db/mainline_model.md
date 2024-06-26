# 主线模型相关表

## cc_ApplicationBase

#### 作用

存放业务信息

#### 表结构

| 字段                  | 类型         | 描述                    |
|---------------------|------------|-----------------------|
| _id                 | ObjectId   | 数据唯一ID                |
| bk_biz_name         | String     | 业务名                   |
| bk_biz_tester       | String     | 测试人员                  |
| operator            | String     | 操作人员                  |
| time_zone           | String     | 时区                    |
| bk_biz_maintainer   | String     | 运维人员                  |
| bk_biz_productor    | String     | 产品人员                  |
| default             | NumberLong | 是否是默认业务，1代表是，0代表否     |
| language            | String     | 语言：1-中文，2-英文          |
| life_cycle          | String     | 生命周期：1-测试中，2-已上线，3-停运 |
| bk_biz_id           | NumberLong | 业务ID                  |
| bk_biz_developer    | String     | 开发人员                  |
| bk_data_status      | String     | 值为 disable 时代表该业务已归档  |
| bk_supplier_account | String     | 开发商ID                 |
| create_time         | ISODate    | 创建时间                  |
| last_time           | ISODate    | 最后更新时间                |

**注意**：此处仅对业务内置的模型字段做说明，业务表结构字段取决于用户在业务模型中定义的属性字段

## cc_ObjectBase_0_pub_{模型id}

#### 作用

存放用户自定义拓扑节点信息

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| bk_inst_name        | String     | 实例名称   |
| bk_inst_id          | NumberLong | 实例id   |
| bk_obj_id           | String     | 模型id   |
| bk_supplier_account | String     | 开发商ID  |
| create_time         | ISODate    | 创建时间   |
| last_time           | ISODate    | 最后更新时间 |
| bk_biz_id           | NumberLong | 业务ID   |
| bk_parent_id        | NumberLong | 父节点id  |

**注意**：此处仅对内置的模型字段做说明，自定义拓扑表结构字段取决于用户在自定义拓扑模型中定义的属性字段

## cc_SetBase

#### 作用

存放集群相关信息

#### 表结构

| 字段                   | 类型         | 描述                   |
|----------------------|------------|----------------------|
| _id                  | ObjectId   | 数据唯一ID               |
| default              | NumberLong | 是否为默认集群              |
| bk_set_env           | String     | 环境类型                 |
| set_template_version | String     | 集群模板                 |
| bk_biz_id            | NumberLong | 业务ID                 |
| bk_set_name          | String     | 集群名称                 |
| set_template_id      | NumberLong | 集群模板id，为0代表不通过集群模板创建 |
| bk_capacity          | NumberLong | 设计容量                 |
| bk_set_id            | NumberLong | 集群id                 |
| bk_parent_id         | NumberLong | 父节点id                |
| bk_set_desc          | String     | 集群描述                 |
| bk_service_status    | String     | 服务状态                 |
| description          | String     | 备注                   |
| bk_supplier_account  | String     | 开发商ID                |
| create_time          | ISODate    | 创建时间                 |
| last_time            | ISODate    | 最后更新时间               |

**注意**：此处仅对内置的模型字段做说明，集群表结构字段取决于用户在集群模型中定义的属性字段

## cc_ModuleBase

#### 作用

模块信息表

#### 表结构

| 字段                  | 类型                   | 描述               |
|---------------------|----------------------|------------------|
| _id                 | ObjectId             | 数据唯一ID           |
| default             | NumberLong           | 是否是默模块，1代表是，0代表否 |
| operator            | String               | 维护人              |
| service_template_id | NumberLong           | 服务模板id           |
| bk_module_type      | String               | 模块类型模块类型         |
| bk_bak_operator     | String               | 备份维护人            |
| bk_module_id        | NumberLong           | 模块id             |
| bk_biz_id           | NumberLong           | 业务id             |
| bk_parent_id        | NumberLong           | 父节点id            |
| host_apply_enabled  | Boolean              | 是否开启主机属性自动应用     |
| bk_module_name      | String               | 模块名              |
| bk_set_id           | NumberLongNumberLong | 集群id             |
| service_category_id | NumberLong           | 服务类型id           |
| set_template_id     | NumberLong           | 集群模板id           |
| bk_supplier_account | String               | 开发商ID            |
| create_time         | ISODate              | 创建时间             |
| last_time           | ISODate              | 最后更新时间           |

**注意**：此处仅对内置的模型字段做说明，模块表结构字段取决于用户在模块模型中定义的属性字段

## cc_ModuleHostConfig

#### 作用

模块主机关联关系信息表

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| bk_biz_id           | NumberLong | 业务id   |
| bk_host_id          | NumberLong | 主机id   |
| bk_module_id        | NumberLong | 模块id   |
| bk_set_id           | NumberLong | 集群id   |
| bk_supplier_account | String     | 开发商ID  |