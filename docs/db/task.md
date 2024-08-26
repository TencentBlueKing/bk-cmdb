# 异步任务相关表

## cc_APITask

#### 作用

存储具体执行的异步任务信息（会定时清理过期任务）

#### 表结构

| 字段                  | 类型         | 描述                                              |
|---------------------|------------|-------------------------------------------------|
| _id                 | ObjectId   | 数据唯一ID                                          |
| task_id             | String     | 任务ID，由taskserver生成的唯一ID                         |
| task_type           | String     | 任务标识，用于业务方识别任务，同时表示所在的任务队列，每个队列都配置了回调接口等信息      |
| bk_inst_id          | NumberLong | 该任务关联的实例id，用于查询一个实例对应的所有任务，并且防止对同一个实例多次创建新的同步任务 |
| user                | String     | 任务创建者                                           |
| header              | Object     | 请求的http header，包括rid等信息，可以用于问题定位                |
| status              | String     | 任务执行状态                                          |
| detail              | Object     | 子任务详情列表                                         |
| bk_supplier_account | String     | 开发商ID                                           |
| create_time         | ISODate    | 创建时间                                            |
| last_time           | ISODate    | 最后更新时间                                          |

## cc_APITaskSyncHistory

#### 作用

存储执行的异步任务历史

#### 表结构

| 字段                  | 类型         | 描述                                              |
|---------------------|------------|-------------------------------------------------|
| _id                 | ObjectId   | 数据唯一ID                                          |
| task_id             | String     | 任务ID，由taskserver生成的唯一ID                         |
| task_type           | String     | 任务标识，用于业务方识别任务，同时表示所在的任务队列，每个队列都配置了回调接口等信息      |
| bk_inst_id          | NumberLong | 该任务关联的实例id，用于查询一个实例对应的所有任务，并且防止对同一个实例多次创建新的同步任务 |
| status              | String     | 任务执行状态                                          |
| creator             | String     | 任务创建者                                           |
| bk_supplier_account | String     | 开发商ID                                           |
| create_time         | ISODate    | 创建时间                                            |
| last_time           | ISODate    | 最后更新时间                                          |