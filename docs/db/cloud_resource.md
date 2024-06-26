# 云资源相关功能表

## cc_CloudAccount

#### 作用

存放云账户信息

#### 表结构

| 字段                    | 类型         | 描述      |
|-----------------------|------------|---------|
| _id                   | ObjectId   | 数据唯一ID  |
| bk_account_name       | String     | 云账户名称   |
| bk_account_id         | NumberLong | 云账户id   |
| bk_secret_id          | String     | 云账户密钥id |
| bk_secret_key         | String     | 云账户密钥   |
| bk_cloud_vendor       | String     | 云厂商     |
| bk_description        | String     | 云账户描述   |
| bk_can_delete_account | Boolean    | 是否可删除   |
| bk_creator            | String     | 创建人     |
| bk_last_editor        | String     | 最后更新人   |
| create_time           | ISODate    | 创建时间    |
| last_time             | ISODate    | 最后更新时间  |

## cc_CloudSyncHistory

#### 作用

存放云同步任务历史信息

#### 表结构

| 字段                    | 类型       | 描述        |
|-----------------------|----------|-----------|
| _id                   | ObjectId | 数据唯一ID    |
| bk_task_id            | String   | 任务ID      |
| bk_history_id         | String   | 云同步任务历史ID |
| bk_sync_status        | String   | 任务执行状态    |
| bk_status_description | Object   | 任务状态描述    |
| bk_detail             | Object   | 任务详情信息    |
| create_time           | ISODate  | 创建时间      |
| bk_supplier_account   | String   | 开发商ID     |

#### bk_status_description 字段结构示例

| 字段         | 类型     | 描述        |
|------------|--------|-----------|
| cost_time  | Float  | 云同步任务花费时间 |
| error_info | String | 错误信息      |

#### bk_detail 字段结构示例

| 字段      | 类型     | 描述     |
|---------|--------|--------|
| update  | Object | 更新的云实例 |
| new_add | Object | 新增的云实例 |

#### update 和 new_add 字段结构示例

| 字段    | 类型           | 描述      |
|-------|--------------|---------|
| count | NumberLong   | 云实例数量   |
| ips   | String Array | 云实例IP数组 |

## cc_CloudSyncTask

#### 作用

存放云同步任务表信息

#### 表结构

| 字段                    | 类型           | 描述       |
|-----------------------|--------------|----------|
| _id                   | ObjectId     | 数据唯一ID   |
| bk_task_id            | String       | 任务ID     |
| bk_task_name          | String       | 任务名称     |
| bk_resource_type      | String       | 资源类型     |
| bk_account_id         | NumberLong   | 云账户id    |
| bk_cloud_vendor       | String       | 云厂商      |
| bk_sync_status        | String       | 任务执行状态   |
| bk_status_description | Object       | 任务状态描述   |
| bk_last_sync_time     | ISODate      | 任务最后同步时间 |
| bk_sync_all           | Boolean      | 是否同步所有实例 |
| bk_sync_all_dir       | NumberLong   | 同步目录     |
| bk_sync_vpcs          | Object Array | VPC详情    |
| bk_creator            | String       | 创建人      |
| bk_last_editor        | String       | 最后更新人    |
| create_time           | ISODate      | 创建时间     |
| last_time             | ISODate      | 最后更新时间   |
| bk_supplier_account   | String       | 开发商ID    |

#### bk_status_description 字段结构示例

| 字段         | 类型     | 描述        |
|------------|--------|-----------|
| cost_time  | Float  | 云同步任务花费时间 |
| error_info | String | 错误信息      |

#### bk_sync_vpcs 字段结构示例

| 字段            | 类型         | 描述        |
|---------------|------------|-----------|
| bk_vpc_id     | String     | vpc id    |
| bk_vpc_name   | String     | vpc 名称    |
| bk_region     | String     | 地域        |
| bk_host_count | NumberLong | 主机数量      |
| bk_sync_dir   | NumberLong | 同步到主机池的目录 |
| bk_cloud_id   | NumberLong | 管控区域      |
| destroyed     | Boolean    | 云实例是否已释放  |