# 集群模板功能相关表

## cc_SetServiceTemplateRelation

#### 作用

存放集群模板与服务模板关联关系信息

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| bk_biz_id           | NumberLong | 业务id   |
| set_template_id     | NumberLong | 集群模板id |
| service_template_id | NumberLong | 服务模板id |
| bk_supplier_account | String     | 开发商ID  |

## cc_SetTemplate

#### 作用

集群模板信息表

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| id                  | NumberLong | 模板id   |
| name                | String     | 模板名称   |
| bk_biz_id           | NumberLong | 业务id   |
| creator             | String     | 创建者    |
| modifier            | String     | 更新者    |
| bk_supplier_account | String     | 开发商ID  |
| create_time         | ISODate    | 创建时间   |
| last_time           | ISODate    | 最后更新时间 |

## cc_SetTemplateAttr

#### 作用

存放集群模板属性配置信息

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| id                  | NumberLong | 自增id   |
| bk_biz_id           | NumberLong | 业务id   |
| set_template_id     | NumberLong | 集群模板id |
| bk_attribute_id     | NumberLong | 属性字段id |
| bk_property_value   | String     | 属性字段值  |
| creator             | String     | 创建者    |
| modifier            | String     | 更新者    |
| bk_supplier_account | String     | 开发商ID  |
| create_time         | ISODate    | 创建时间   |
| last_time           | ISODate    | 最后更新时间 |