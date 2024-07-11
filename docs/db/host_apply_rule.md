# 主机属性自动应用功能相关表

## cc_HostApplyRule

#### 作用

保存主机属性自动应用信息

#### 表结构

| 字段                  | 类型         | 描述                              |
|---------------------|------------|---------------------------------|
| _id                 | ObjectId   | 数据唯一ID                          |
| id                  | NumberLong | 自增id                            |
| bk_biz_id           | NumberLong | 业务id                            |
| bk_module_id        | NumberLong | 模块id，当模块id为0代表主机属性自动应用作用于服务模板   |
| service_template_id | NumberLong | 服务模板id，当服务模板id为0代表主机属性自动应用作用于模块 |
| bk_attribute_id     | NumberLong | 属性id                            |
| bk_property_value   | String     | 自动应用的属性值                        |
| creator             | String     | 创建人                             |
| modifier            | String     | 修改人                             |
| bk_supplier_account | String     | 开发商ID                           |
| create_time         | ISODate    | 创建时间                            |
| last_time           | ISODate    | 最后更新时间                          |