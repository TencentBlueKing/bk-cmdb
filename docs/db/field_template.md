# 字段组合模板功能相关表

## cc_FieldTemplate

#### 作用

存放字段组合模版信息

#### 表结构

| 字段                  | 类型         | 描述     |
|---------------------|------------|--------|
| _id                 | ObjectId   | 数据唯一ID |
| id                  | NumberLong | 自增id   |
| name                | String     | 模版名称   |
| description         | String     | 模版描述   |
| creator             | String     | 创建人    |
| modifier            | String     | 修改人    |
| bk_supplier_account | String     | 开发商ID  |
| create_time         | ISODate    | 创建时间   |
| last_time           | ISODate    | 最后更新时间 |

## cc_ObjAttDesTemplate

#### 作用

字段组合模版属性表

#### 表结构

| 字段                  | 类型         | 描述                        |
|---------------------|------------|---------------------------|
| _id                 | ObjectId   | 数据唯一ID                    |
| bk_template_id      | NumberLong | 字段组合模版id                  |
| bk_property_id      | String     | 属性ID                      |
| editable            | Boolean    | 是否可编辑                     |
| isrequired          | Boolean    | 是否必填                      |
| option              | —          | 用户自定义内容，存储的内容及数据格式由字段类型决定 |
| unit                | String     | 单位                        |
| placeholder         | String     | 用户提示                      |
| bk_property_name    | String     | 属性名，用于展示                  |
| bk_property_type    | String     | 属性字段数据类型                  |
| ismultiple          | Boolean    | 字段是否支持可多选                 |
| default             | —          | 字段默认值，存储的内容及数据格式由字段类型决定   |
| creator             | String     | 创建人                       |
| modifier            | String     | 修改人                       |
| bk_supplier_account | String     | 开发商ID                     |
| create_time         | ISODate    | 创建时间                      |
| last_time           | ISODate    | 最后更新时间                    |

## cc_ObjectUniqueTemplate

#### 作用

字段组合模板字段唯一校验信息表

#### 表结构

| 字段                  | 类型         | 描述       |
|---------------------|------------|----------|
| _id                 | ObjectId   | 数据唯一ID   |
| id                  | NumberLong | 模型id     |
| bk_template_id      | NumberLong | 字段组合模板id |
| creator             | String     | 创建人      |
| modifier            | String     | 更新人      |
| keys                | Array      | 属性字段id列表 |
| bk_supplier_account | String     | 开发商ID    |
| create_time         | ISODate    | 创建时间     |
| last_time           | ISODate    | 最后更新时间   |

## cc_ObjFieldTemplateRelation

#### 作用

字段组合模板与模型关联关系信息表

#### 表结构

| 字段                  | 类型         | 描述       |
|---------------------|------------|----------|
| _id                 | ObjectId   | 数据唯一ID   |
| object_id           | NumberLong | 模型id     |
| bk_template_id      | NumberLong | 字段组合模板id |
| bk_supplier_account | String     | 开发商ID    |