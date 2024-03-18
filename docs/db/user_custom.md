# 资源目录用户自定义配置相关表

## cc_UserCustom

#### 作用

存放用户自定义配置信息

#### 表结构

| 字段                       | 类型         | 描述             |
|--------------------------|------------|----------------|
| _id                      | ObjectId   | 数据唯一ID         |
| id                       | NumberLong | 自增id           |
| bk_user                  | String     | 创建人            |
| menu_resource_collection | Array      | 自定义配置的模型唯一标识列表 |
| bk_supplier_account      | String     | 开发商ID，默认0      |