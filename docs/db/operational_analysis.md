# 运营分析功能相关表

## cc_AuditLog

#### 作用

存放审计日志

#### 表结构

| 字段                   | 类型         | 描述                                    |
|----------------------|------------|---------------------------------------|
| _id                  | ObjectId   | 数据唯一ID                                |
| id                   | NumberLong | 自增ID                                  |
| audit_type           | String     | 操作的资源类型大类                             |
| user                 | String     | 操作人                                   |
| resource_type        | String     | 操作的具体资源类型                             |
| action               | String     | 操作类型如create(新增)/update(更新)/delete(删除) |
| operate_from         | String     | 操作审计的来源平台                             |
| operation_time       | NumberLong | 操作时间                                  |
| operation_detail     | Object     | 操作审计详情，字段示例如下                         |
| bk_biz_id            | NumberLong | 业务id                                  |
| resource_id          | NumberLong | 资源id                                  |
| resource_name        | String     | 资源名称                                  |
| code                 | String     | 状态码                                   |
| rid                  | String     | 请求id                                  |
| extend_resource_name | String     | 扩展资源名                                 |
| bk_supplier_account  | String     | 开发商ID                                 |

#### operation_detail 字段结构示例

##### 转移主机

| 字段          | 类型         | 描述                          |
|-------------|------------|-----------------------------|
| bk_biz_id   | NumberLong | 业务id归还主机为归还前业务ID，否则为转移后业务ID |
| bk_biz_name | String     | 业务名称                        |
| pre_data    | Object     | 主机转移前业务、集群、模块信息             |
| cur_data    | Object     | 主机转移后业务、集群、模块信息             |

##### 实例关联操作

| 字段             | 类型         | 描述     |
|----------------|------------|--------|
| asst_id        | String     | 关联关系id |
| asst_kind      | String     | 关联类型   |
| src_obj_id     | String     | 源模型    |
| dest_obj_id    | String     | 目标模型   |
| dest_inst_id   | NumberLong | 目标实例id |
| dest_inst_name | String     | 目标实例名称 |

##### 模型实例操作

| 字段            | 类型     | 描述      |
|---------------|--------|---------|
| pre_data      | Object | 实例变更前数据 |
| cur_data      | Object | 实例变更后数据 |
| update_fields | Object | 更新的字段   |
| bk_obj_id     | String | 模型唯一标识  |

##### 模型属性字段操作

| 字段            | 类型     | 描述      |
|---------------|--------|---------|
| pre_data      | Object | 实例变更前数据 |
| cur_data      | Object | 实例变更后数据 |
| update_fields | Object | 更新的字段   |
| bk_obj_id     | String | 模型唯一标识  |
| bk_obj_name   | String | 模型名称    |

##### 引用类型实例操作

| 字段             | 类型     | 描述           |
|----------------|--------|--------------|
| pre_data       | Object | 实例变更前数据      |
| cur_data       | Object | 实例变更后数据      |
| update_fields  | Object | 更新的字段        |
| bk_obj_id      | String | 引用模型唯一标识     |
| src_obj_id     | String | 引用模型所属模型唯一标识 |
| bk_property_id | String | 属性字段id       |

##### 服务实例操作

| 字段            | 类型           | 描述      |
|---------------|--------------|---------|
| pre_data      | Object       | 实例变更前数据 |
| cur_data      | Object       | 实例变更后数据 |
| update_fields | Object       | 更新的字段   |
| bk_host_id    | NumberLong   | 主机id    |
| processes     | Object Array | 进程实例操作  |

##### 进程实例操作（聚合在服务实例操作中）

| 字段               | 类型         | 描述      |
|------------------|------------|---------|
| pre_data         | Object     | 实例变更前数据 |
| cur_data         | Object     | 实例变更后数据 |
| update_fields    | Object     | 更新的字段   |
| action           | String     | 操作类型    |
| bk_process_ids   | NumberLong | 进程id    |
| bk_process_names | String     | 进程名称    |

## cc_ChartConfig

#### 作用

存放运营统计图表配置数据信息

#### 表结构

| 字段                  | 类型         | 描述       |
|---------------------|------------|----------|
| _id                 | ObjectId   | 数据唯一ID   |
| config_id           | NumberLong | 图表配置数据ID |
| report_type         | String     | 统计类型     |
| name                | String     | 图表名称     |
| bk_obj_id           | String     | 统计实例类型   |
| width               | String     | 图表宽度     |
| chart_type          | String     | 图表类型     |
| field               | String     | 分类字段     |
| x_axis_count        | NumberLong | x轴数量     |
| create_time         | ISODate    | 创建时间     |
| bk_supplier_account | String     | 开发商ID    |

## cc_ChartData

#### 作用

存放运营统计数据

#### 表结构

| 字段                  | 类型       | 描述     |
|---------------------|----------|--------|
| _id                 | ObjectId | 数据唯一ID |
| report_type         | String   | 统计类型   |
| data                | Object   | 统计数据   |
| create_time         | ISODate  | 创建时间   |
| bk_supplier_account | String   | 开发商ID  |

## cc_ChartPosition

#### 作用

存放运营统计中主机与其它模型实例的统计图表顺序数据

#### 表结构

| 字段                  | 类型         | 描述                  |
|---------------------|------------|---------------------|
| _id                 | ObjectId   | 数据唯一ID              |
| bk_biz_id           | NumberLong | 业务id                |
| position            | Object     | 图表顺序数据，数字小的图表对应位置靠前 |
| bk_supplier_account | String     | 开发商ID               |

#### position 字段结构

| 字段   | 类型    | 描述       |
|------|-------|----------|
| host | Array | 主机图表顺序   |
| inst | Array | 模型实例图表顺序 |