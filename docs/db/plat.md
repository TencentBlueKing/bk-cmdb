# 管控区域相关表

## cc_PlatBase

#### 作用

管控区域信息表

#### 表结构

| 字段                  | 类型         | 描述      |
|---------------------|------------|---------|
| _id                 | ObjectId   | 数据唯一ID  |
| bk_cloud_name       | String     | 管控区域名   |
| bk_cloud_id         | NumberLong | 管控区域id  |
| bk_cloud_vendor     | String     | 管控区域类型  |
| bk_creator          | String     | 创建者     |
| bk_last_editor      | String     | 最后修改人   |
| bk_region           | String     | VPC所属地域 |
| bk_status           | String     | 状态      |
| bk_status_detail    | String     | 状态详情    |
| bk_vpc_id           | String     | VPC唯一标识 |
| bk_vpc_name         | String     | VPC名称   |
| bk_account_id       | NumberLong | 云账户ID   |
| bk_supplier_account | String     | 开发商ID   |
| create_time         | ISODate    | 创建时间    |
| last_time           | ISODate    | 最后更新时间  |