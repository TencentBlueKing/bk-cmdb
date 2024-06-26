# 资源目录用户自定义配置相关表

## cc_UserCustom

#### 作用

存放用户自定义配置信息，前端使用

#### 表结构

| 字段                          | 类型           | 描述                                                |
|-----------------------------|--------------|---------------------------------------------------|
| _id                         | ObjectId     | 数据唯一ID                                            |
| id                          | NumberLong   | 自增id                                              |
| bk_user                     | String       | 创建人                                               |
| menu_resource_collection    | String Array | 自定义资源导航栏配置，模型唯一标识列表，即此列表中模型唯一标识对应的模型会追加到资源页面的导航栏中 |
| resource_host_common_filter | Array        | 自定义主机高级筛选配置                                       |
| bk_supplier_account         | String       | 开发商ID                                             |

#### resource_host_common_filter.[x] 字段结构示例

如下数据示例则代表在进行主机高级筛选时会默认显示出集群分组下的集群名称作为主机筛选条件

| 字段          | 类型     | 描述   |
|-------------|--------|------|
| bk_set_name | String | 集群名称 |
| set         | String | 条件分组 |