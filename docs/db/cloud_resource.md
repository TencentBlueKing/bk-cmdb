# 云资源相关功能表

## cc_CloudAccount

#### 作用

存放云账户信息

#### 表结构

| 字段                    | 类型         | 描述     |
|-----------------------|------------|--------|
| _id                   | ObjectId   | 数据唯一ID |
| bk_account_name       | String     | 云账户名称  |
| bk_account_id         | NumberLong | 云账户id  |
| bk_cloud_vendor       | String     | 云厂商    |
| bk_description        | String     | 云账户描述  |
| bk_can_delete_account | Boolean    | 是否可删除  |
| bk_creator            | String     | 创建人    |
| bk_last_editor        | String     | 最后更新人  |
| create_time           | ISODate    | 创建时间   |
| last_time             | ISODate    | 最后更新时间 |

## cc_CloudSyncHistory

#### 作用

存放云同步任务历史信息

#### 表结构

| 字段                    | 类型         | 描述                                              |
|-----------------------|------------|-------------------------------------------------|
| _id                   | ObjectId   | 数据唯一ID                                          |
| bk_task_id            | String     | 任务ID，由taskserver生成的唯一ID                         |
| task_type             | String     | 任务标识，用于业务方识别任务，同时表示所在的任务队列，每个队列都配置了回调接口等信息      |
| bk_account_id         | NumberLong | 该任务关联的实例id，用于查询一个实例对应的所有任务，并且防止对同一个实例多次创建新的同步任务 |
| bk_sync_status        | String     | 任务执行状态                                          |
| bk_status_description | String     | 任务状态描述                                          |
| bk_detail             | Object     | 任务详情信息                                          |
| create_time           | ISODate    | 创建时间                                            |
| last_time             | ISODate    | 最后更新时间                                          |
| bk_supplier_account   | String     | 开发商ID，默认0                                       |

## cc_CloudSyncTask

#### 作用

存放云同步任务表信息

#### 表结构

| 字段                    | 类型         | 描述                                              |
|-----------------------|------------|-------------------------------------------------|
| _id                   | ObjectId   | 数据唯一ID                                          |
| bk_task_id            | String     | 任务ID，由taskserver生成的唯一ID                         |
| bk_task_name          | String     | 任务名称                                            |
| bk_resource_type      | String     | 任务标识，用于业务方识别任务，同时表示所在的任务队列，每个队列都配置了回调接口等信息      |
| bk_account_id         | NumberLong | 该任务关联的实例id，用于查询一个实例对应的所有任务，并且防止对同一个实例多次创建新的同步任务 |
| bk_sync_status        | String     | 任务执行状态                                          |
| bk_status_description | String     | 任务状态描述                                          |
| create_time           | ISODate    | 创建时间                                            |
| last_time             | ISODate    | 最后更新时间                                          |
| bk_supplier_account   | String     | 开发商ID，默认0                                       |